package gnmi

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"math"
	"regexp"
	"strings"
	"time"

	//go_metrics "github.com/kentik/go-metrics"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/kt/counters"

	auth_pb "github.com/nileshsimaria/jtimon/authentication"
	"github.com/nileshsimaria/jtimon/telemetry"
	"github.com/openconfig/gnmi/client"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type (
	sensorUpdate struct {
		systemId string
		key      *gnmi.Path
		val      interface{}
	}

	updateState struct {
		log  logger.ContextL
		conf *kt.SnmpDeviceConfig

		// deltaBasis, lastVal, and lastDelta are all keyed on interface name
		// (ifName) (e.g. "ge-0/0/1" or "ge-0/0/1.10").  We GET readings from the
		// switch, and SEND values as flow, on different schedules.  It's
		// possible that we can GET new values from the switch both MORE often
		// than we send flow, and LESS often.
		//
		// The deltas that we send should always be by comparison to the value
		// from the last time we SENT data.  (NOT by comparison to the last time
		// we GOT data.)
		//
		// So when we GET data, we compare it to deltaBasis and store the delta
		// in lastDelta, and the value itself in lastVal.
		//
		// When we SEND data, we send the deltas, copy values from lastVal to
		// deltaBasis (leaving alone any in deltaBasis that've not been updated),
		// and clear lastVal and lastDelta.
		deltaBasis map[string]*counters.CounterSet
		lastVal    map[string]*counters.CounterSet
		lastDelta  map[string]map[string]uint64

		// Key: ifName.
		name2idx     map[string]uint32
		updatedIfIdx int
		sensorName   string
		metrics      *kt.SnmpDeviceMetric
	}
)

var (
	ST_ifConnectorPresent = "ifConnectorPresent"

	ST_IfIndex      = "ifindex" // interface index
	ST_HighSpeed    = "high-speed"
	ST_Description  = "description"
	ST_IP           = "ip"
	ST_PrefixLength = "prefix-length"

	ST_InOctets     = "in-octets"
	ST_InUcastPkts  = "in-unicast-pkts"
	ST_OutOctets    = "out-octets"
	ST_OutUcastPkts = "out-unicast-pkts"
	ST_InErrors     = "in-errors"
	ST_OutErrors    = "out-errors"

	// all under .../state/counters
	ST_InDiscards       = "in-discards"
	ST_OutDiscards      = "out-discards"
	ST_OutMulticastPkts = "out-multicast-pkts"
	ST_OutBroadcastPkts = "out-broadcast-pkts"
	ST_InMulticastPkts  = "in-multicast-pkts"
	ST_InBroadcastPkts  = "in-broadcast-pkts"

	// Top-level nodes
	// "/components",
	// "/junos",
	// "/network-instances",
	// "/interfaces",
	// "/lacp",
	// "/lldp",
	// "/mpls",
	// "/arp-information",
	// "/nd6-information",
	// "/ipv6-ra",

	// // From the Juniper Telemetry Explorer
	// "/telemetry-system",
	// "/debug",
	// "/vlans", // internal error in na-grpcd log on pe3.nyc1
	// // "/bgp-rib",
	// "/local-routes",

	interfaces    = "/interfaces/interface"
	iState        = interfaces + "/state"
	subInterfaces = interfaces + "/subinterfaces/subinterface"
	iSubState     = subInterfaces + "/state"

	// Paths we subscribe to & sample interval
	paths = map[string]time.Duration{
		// Can't just subscribe to .../state, as that has a bunch of other stuff
		// (including all the counters, mentioned explicitly below) that we don't
		// want (in this particular subscription).
		iState + "/" + ST_IfIndex:     300 * time.Second,
		iState + "/" + ST_Description: 300 * time.Second,

		iSubState + "/" + ST_IfIndex:     300 * time.Second,
		iSubState + "/" + ST_Description: 300 * time.Second,
		// IP addresses are only on subinterfaces.  An interface can have 0, 1,
		// or more IPs.  We store v4 IPs first and then v6.  The first IP goes in
		// InterfaceData.IPAddr, the rest are in InterfaceData.AliasAddr.
		subInterfaces + "/ipv4/addresses/address/state": 300 * time.Second,
		subInterfaces + "/ipv6/addresses/address/state": 300 * time.Second,

		// For unknown reasons this also gets us iState/parent_ae_name and
		// iState/high-speed.  Subscribing to iState/high-speed directly gets us
		// an "Authorization failed" error.
		iState + "/counters": 30 * time.Second,

		// This does *not* get us iSubState/high-speed (or parent_ae_name).
		iSubState + "/counters": 30 * time.Second,
	}

	// Could be modified during testing
	initialSleep = time.Minute
)

type GNMI struct {
	log      logger.ContextL
	jchfChan chan []*kt.JCHF
	conf     *kt.SnmpDeviceConfig
	metrics  *kt.SnmpDeviceMetric
	ctx      context.Context
	cancel   context.CancelFunc
	recvCh   chan *telemetry.OpenConfigData
	updateCh chan *sensorUpdate
	subReq   telemetry.SubscriptionRequest
}

func NewGNMIClient(jchfChan chan []*kt.JCHF, conf *kt.SnmpDeviceConfig, metrics *kt.SnmpDeviceMetric, log logger.ContextL) (*GNMI, error) {
	kt := GNMI{
		log:      log,
		jchfChan: jchfChan,
		conf:     conf,
		metrics:  metrics,
		// We don't want lots of these to be outstanding.
		recvCh: make(chan *telemetry.OpenConfigData, 1),
		// But there are several of these for each of the above.
		updateCh: make(chan *sensorUpdate, 1000),
	}

	err := kt.openST()
	if err != nil {
		return nil, err

	}

	return &kt, nil
}

func (g *GNMI) GetName() string {
	return "gNMI"
}

func (g *GNMI) Run(ctx context.Context, dur time.Duration) {
	ctxC, cancel := context.WithCancel(ctx)
	g.ctx = ctxC
	g.cancel = cancel

	go g.processUpdates()
}

func (g *GNMI) openST() error {
	d := &client.Destination{Addrs: []string{g.conf.DeviceIP}}
	if g.conf.Ext.GNMIConfig.CaCert != "" && g.conf.Ext.GNMIConfig.ClientCert != "" && g.conf.Ext.GNMIConfig.ClientKey != "" {
		certPool := x509.NewCertPool()
		caCert, err := mkPEM("CERTIFICATE", g.conf.Ext.GNMIConfig.CaCert)
		if err != nil {
			return err
		}
		ok := certPool.AppendCertsFromPEM([]byte(caCert))
		if !ok {
			return fmt.Errorf("Could not append the caCert")
		}

		clientCert, err := mkPEM("CERTIFICATE", g.conf.Ext.GNMIConfig.ClientCert)
		if err != nil {
			return err
		}
		clientKey, err := mkPEM("RSA PRIVATE KEY", g.conf.Ext.GNMIConfig.ClientKey)
		if err != nil {
			return err
		}
		certificate, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil {
			return err
		}
		d.TLS = &tls.Config{
			RootCAs:      certPool,
			Certificates: []tls.Certificate{certificate},
		}
	}
	if g.conf.Ext.GNMIConfig.Username != "" && g.conf.Ext.GNMIConfig.Password != "" {
		d.Credentials = &client.Credentials{
			Username: g.conf.Ext.GNMIConfig.Username,
			Password: g.conf.Ext.GNMIConfig.Password,
		}
	}

	//g.metrics.STSeen.Update(0)

	for path, sampleInterval := range paths {
		var pathM telemetry.Path
		pathM.Path = path
		pathM.SampleFrequency = uint32(sampleInterval / time.Millisecond)
		g.subReq.PathList = append(g.subReq.PathList, &pathM)
	}

	go g.subscribe(g.conf.DeviceName, d)
	go g.splitUpdates()

	//g.metrics.STSeen.Update(1)

	return nil
}

func mkPEM(name, cert string) (string, error) {
	re := regexp.MustCompile(
		fmt.Sprintf("^(-----BEGIN %s-----)?\n?"+
			"((?s).*?)\n?"+
			"(-----END %s-----)?\n?$",
			name, name))
	matches := re.FindStringSubmatch(cert)
	if matches == nil {
		return "", fmt.Errorf("Not a valid %s certificate: %s", name, cert)
	}
	return fmt.Sprintf("-----BEGIN %s-----\n%s\n-----END %s-----", name, matches[2], name), nil
}

// Subscribe dials, subscribes, and runs the receive loop.  Some errors (for
// example, authentication errors in Dial) can permanently terminate all ST
// efforts (since, for example, if we have the wrong TLS certs or the wrong
// username/password, nothing will ever work until those are fixed, so might as
// well just abort everything).
func (g *GNMI) subscribe(deviceId string, d *client.Destination) {
	defer g.cancel()

	sleepTime := initialSleep
	connErrCnt := 0

	for g.ctx.Err() == nil {
		g.log.Debugf("%s dialing %s", deviceId, d.Addrs[0])
		conn, retry, err := Dial(g.ctx, deviceId, d)
		if err != nil {
			if !retry {
				g.log.Errorf("%v", err)
				break
			}

			connErrCnt++
			if connErrCnt > 5 {
				sleepTime = 5 * initialSleep
			}
			g.log.Warnf("subscribe: Dial error, try #%d; sleeping %v: error: %v", connErrCnt, sleepTime, err)
			time.Sleep(sleepTime)
			continue
		}
		g.log.Infof("Connected %s to %s", deviceId, d.Addrs[0])
		sleepTime = initialSleep
		connErrCnt = 0

		err = g.recvLoop(conn)
		if err != nil {
			g.log.Errorf("%v", err)
			break
		}
	}
	close(g.recvCh)
	g.log.Infof("%s: Subscribe loop exited", deviceId)
}

func (g *GNMI) recvLoop(conn *grpc.ClientConn) error {
	octc := telemetry.NewOpenConfigTelemetryClient(conn)
	subscrCtx, subscrCancel := context.WithCancel(g.ctx)
	defer subscrCancel()

	stream, err := octc.TelemetrySubscribe(subscrCtx, &g.subReq)
	if err != nil {
		//g.metrics.STErrors.Mark(1)
		return fmt.Errorf("Error subscribing to ST: %v", err)
	}

	// hdr, err := stream.Header()
	// if err != nil {
	// 	g.log.Warnf(st.logPrefix, "Failed to get header for stream: %v", err)
	// 	break
	// }
	// st.log.Debugf(st.logPrefix, "gRPC headers from host %s", server)
	// for k, v := range hdr {
	// 	st.log.Debugf(st.logPrefix, "  %s: %s", k, v)
	// }

	for subscrCtx.Err() == nil {
		ocData, err := stream.Recv()
		if err == io.EOF {
			g.log.Debugf("recvLoop: EOF")
			break
		}
		if err != nil {
			if p, ok := status.FromError(err); ok {
				// FIXME: This text could probably vary from vendor to vendor
				// or switch OS version to version.
				if p.Message() == "Authorization failed" {
					// Auth is wrong; just bail on all ST operations.
					g.log.Debugf("recvLoop: Authorization failed")
					//st.metrics.STErrors.Mark(1)
					return err
				}

				if p.Message() == "grpc: the client connection is closing" ||
					p.Code() == codes.Canceled {
					g.log.Debugf("recvLoop: connection closed or cancelled")
					break
				}
			}

			//st.metrics.STErrors.Mark(1)
			g.log.Infof("TelemetrySubscribe Recv error: %v, %T", err, err)
			break
		}

		if ocData.SyncResponse {
			ocData = nil
		}

		// Try to send, but if can't, drop it.
		select {
		case g.recvCh <- ocData:
		default:
			g.log.Warnf("Dropping streamed data")
		}
	}
	return nil
}

// Read from recvCh, send individual updates to updateCh
func (g *GNMI) splitUpdates() {
OUTER:
	for {
		select {
		case u, ok := <-g.recvCh:
			if !ok {
				break OUTER
			}
			// Send the systemId only once per update.
			g.updateCh <- &sensorUpdate{
				systemId: u.SystemId,
			}

			// st.log.Infof(st.logPrefix, "ocData: %s, %d, %d, %s, %d, %d, %t",
			// 	u.SystemId, u.ComponentId, u.SubComponentId,
			// 	u.Path, u.SequenceNumber, u.Timestamp,
			// 	u.SyncResponse)
			var prefix string
			for _, kv := range u.Kv {
				if kv.Key == "__prefix__" {
					prefix = kv.GetStrValue()
					// st.log.Infof(st.logPrefix, "Prefix: %s", prefix)
					continue
				}
				// st.log.Infof(st.logPrefix, "kv: prefix: %s, key: %s, value: %+v", prefix, kv.GetKey(), kv.GetValue())
				fullKey := prefix + kv.GetKey()
				structuredKey, err := ygot.StringToStructuredPath(fullKey)
				if err != nil {
					g.log.Warnf("StringToStructuredPath error on %s: %v", fullKey, err)
					continue
				}
				g.updateCh <- &sensorUpdate{
					key: structuredKey,
					val: kv.GetValue(),
				}
			}

		case <-g.ctx.Done():
			break OUTER
		}
	}
	close(g.updateCh)
}

func (g *GNMI) processUpdates() {
	counterTimeSec := 60
	if g.conf.PollTimeSec > 0 {
		counterTimeSec = g.conf.PollTimeSec
	}

	sendFlow := time.NewTicker(time.Duration(counterTimeSec) * time.Second)
	defer sendFlow.Stop()
	us := newUpdateState(g.log, g.metrics, g.conf)

	for {
		select {
		case su, ok := <-g.updateCh:
			if !ok {
				return
			}
			// This could be better
			if su.systemId != "" {
				continue
			}
			err := us.processSingleUpdate(su)
			if err != nil {
				//st.metrics.STErrors.Mark(1)
				g.log.Errorf("%v", err)
			}

		case <-sendFlow.C:
			res := us.sendAsCHF()
			if len(res) > 0 {
				g.jchfChan <- res
			}

		case <-g.ctx.Done():
			return
		}
	}
}

func (g *GNMI) Close() {
	g.cancel()
}

func newUpdateState(log logger.ContextL, metrics *kt.SnmpDeviceMetric, conf *kt.SnmpDeviceConfig) *updateState {
	return &updateState{
		log:  log,
		conf: conf,
		// TODO: Refactor this set of map[string]<stuff> into a single map[string]SomeType.
		deltaBasis: map[string]*counters.CounterSet{},
		lastVal:    map[string]*counters.CounterSet{},
		lastDelta:  map[string]map[string]uint64{},
		name2idx:   map[string]uint32{},
		metrics:    metrics,
	}
}

func (us *updateState) processSingleUpdate(update *sensorUpdate) error {
	// us.log.Infof(st.logPrefix, "updateState.processSingleUpdate: update: %v", update)

	if update == nil || update.key == nil {
		if update != nil {
			return fmt.Errorf("nil key received: %v", update)
		}
		return nil
	}

	elem := update.key.Elem
	l := len(elem)
	if l == 0 {
		return fmt.Errorf("0-length path received: %v", update)
	}

	pathName := elem[l-1].Name
	switch pathName {
	case ST_IfIndex, ST_Description:
		if !verifyStatePath(update.key, pathName) {
			return fmt.Errorf("'%s' recieved; verifyStatePath failed; update: %v", pathName, update)
		}
	case ST_HighSpeed:
		if !verifyShortStatePath(update.key, pathName) {
			return fmt.Errorf("'%s' recieved; verifyShortStatePath failed; update: %v", pathName, update)
		}
	case ST_InOctets, ST_OutOctets,
		ST_InUcastPkts, ST_OutUcastPkts,
		ST_InErrors, ST_OutErrors,
		ST_InDiscards, ST_OutDiscards,
		ST_OutMulticastPkts, ST_OutBroadcastPkts,
		ST_InMulticastPkts, ST_InBroadcastPkts:
		if !verifyCounter(update.key, pathName) {
			return fmt.Errorf("Known counter recieved; verifyCounter failed; update: %v", update)
		}
	case ST_PrefixLength:
		if !verifyIPPath(update.key, pathName) {
			return fmt.Errorf("'%s' recieved; verifyIPPath failed; update: %v", pathName, update)
		}
	default:
		return nil
	}

	ifName := strings.Trim(elem[1].Key["name"], "'")
	if l == 6 || l == 7 || // long ifindex paths are len==6; long counter paths are len==7
		l == 9 { // ST_IP paths are len==9, but their prefix is the same
		ifName += "." + strings.Trim(elem[3].Key["index"], "'")
	}

	// If we got this far, it's a counter

	ifcDeltaBasis := us.deltaBasis[ifName]
	ifcLastVal, ok := us.lastVal[ifName]
	if !ok {
		ifcLastVal = counters.NewCounterSet()
		us.lastVal[ifName] = ifcLastVal
	}
	ifcDelta, ok := us.lastDelta[ifName]
	if !ok {
		ifcDelta = map[string]uint64{}
		us.lastDelta[ifName] = ifcDelta
	}

	val := update.val.(*telemetry.KeyValue_UintValue).UintValue
	// We only use ifIndex for logging here, so it's allowed to not exist yet.
	// If we get to sendAsCHF and it's still not known, throw a warning there.
	ifIndex := us.name2idx[ifName]
	us.log.Debugf("Interface %s(%d): counter %s: value: %d",
		ifName, ifIndex, pathName, val)
	ifcDelta[pathName] = ifcDeltaBasis.GetDelta(pathName, val)
	ifcLastVal.Values[pathName] = val

	return nil
}

// path must not be nil
func verifyStatePath(path *gnmi.Path, finalComponent string) bool {
	return (verifyShortStatePath(path, finalComponent) ||
		verifyLongStatePath(path, finalComponent))
}

func verifyShortStatePath(path *gnmi.Path, finalComponent string) bool {
	return len(path.Elem) == 4 &&
		verifyShortPrefix(path) &&
		path.Elem[3].Name == finalComponent
}

func verifyLongStatePath(path *gnmi.Path, finalComponent string) bool {
	return len(path.Elem) == 6 &&
		verifyLongPrefix(path) &&
		path.Elem[5].Name == finalComponent
}

// path must not be nil
func verifyCounter(path *gnmi.Path, counter string) bool {
	return (verifyShortCounter(path, counter) ||
		verifyLongCounter(path, counter))
}

func verifyShortCounter(path *gnmi.Path, counter string) bool {
	return len(path.Elem) == 5 &&
		verifyShortPrefix(path) &&
		path.Elem[3].Name == "counters" &&
		path.Elem[4].Name == counter
}

func verifyLongCounter(path *gnmi.Path, counter string) bool {
	return len(path.Elem) == 7 &&
		verifyLongPrefix(path) &&
		path.Elem[5].Name == "counters" &&
		path.Elem[6].Name == counter
}

func verifyShortPrefix(path *gnmi.Path) bool {
	return path.Elem[0].Name == "interfaces" &&
		path.Elem[1].Name == "interface" &&
		len(path.Elem[1].Key) > 0 &&
		path.Elem[1].Key["name"] != "" &&
		path.Elem[2].Name == "state"
}

func verifyLongPrefix(path *gnmi.Path) bool {
	return path.Elem[0].Name == "interfaces" &&
		path.Elem[1].Name == "interface" &&
		len(path.Elem[1].Key) > 0 &&
		path.Elem[1].Key["name"] != "" &&
		path.Elem[2].Name == "subinterfaces" &&
		path.Elem[3].Name == "subinterface" &&
		len(path.Elem[3].Key) > 0 &&
		path.Elem[3].Key["index"] != "" &&
		path.Elem[4].Name == "state"
}

func verifyIPPath(path *gnmi.Path, finalComponent string) bool {
	return path.Elem[0].Name == "interfaces" &&
		path.Elem[1].Name == "interface" &&
		len(path.Elem[1].Key) > 0 &&
		path.Elem[1].Key["name"] != "" &&
		path.Elem[2].Name == "subinterfaces" &&
		path.Elem[3].Name == "subinterface" &&
		len(path.Elem[3].Key) > 0 &&
		path.Elem[3].Key["index"] != "" &&
		(path.Elem[4].Name == "ipv4" || path.Elem[4].Name == "ipv6") &&
		path.Elem[5].Name == "addresses" &&
		path.Elem[6].Name == "address" &&
		len(path.Elem[6].Key) > 0 &&
		path.Elem[6].Key["ip"] != "" &&
		path.Elem[7].Name == "state" &&
		path.Elem[8].Name == finalComponent
}

func (us *updateState) sendAsCHF() []*kt.JCHF {

	//us.metrics.STIntTotal.Update(int64(len(us.deltaBasis)))
	//us.metrics.STIntChanged.Update(int64(us.updatedIfIdx))
	us.updatedIfIdx = 0

	// us.log.Infof(us.logPrefix, "len us.deltaBasis: %d", len(us.deltaBasis))
	res := make([]*kt.JCHF, 0)
	for ifName, vals := range us.lastDelta {
		// Lookup ifindex
		if ifIndex, ok := us.name2idx[ifName]; ok {
			dst := kt.NewJCHF()
			dst.CustomStr = map[string]string{}
			dst.CustomInt = map[string]int32{}
			dst.CustomBigInt = map[string]int64{}
			dst.EventType = kt.KENTIK_EVENT_SNMP_DEV_METRIC
			dst.Provider = us.conf.Provider
			dst.DeviceName = us.conf.DeviceName
			dst.SrcAddr = us.conf.DeviceIP
			dst.Timestamp = time.Now().Unix()
			dst.CustomMetrics = map[string]kt.MetricInfo{}
			dst.InputPort = kt.IfaceID(ifIndex)
			dst.OutputPort = kt.IfaceID(ifIndex)

			for _, ctrName := range []string{
				ST_InUcastPkts, ST_OutUcastPkts, ST_InOctets, ST_OutOctets,
				ST_InErrors, ST_OutErrors,
				ST_InDiscards, ST_OutDiscards,
				ST_OutMulticastPkts, ST_OutBroadcastPkts,
				ST_InMulticastPkts, ST_InBroadcastPkts,
			} {
				if val, ok := vals[ctrName]; ok {
					dst.CustomBigInt[ctrName] = int64(val)
					dst.CustomMetrics[ctrName] = kt.MetricInfo{Oid: "gnmi", Mib: "gnmi", Profile: "gnmi", Type: "gnmi"}
				}
			}

			us.conf.SetUserTags(dst.CustomStr)
			res = append(res, dst)
		}
	}

	// Update us.deltaBasis with new values from lastVal; leave unchanged
	// anything that doesn't appear in lastVal.
	for ifName, cSet := range us.lastVal {
		if ifcDeltaBasis, ok := us.deltaBasis[ifName]; ok {
			for ctr, val := range cSet.Values {
				ifcDeltaBasis.Values[ctr] = val
			}
		} else {
			us.deltaBasis[ifName] = cSet
		}
	}

	us.lastVal = map[string]*counters.CounterSet{}
	us.lastDelta = map[string]map[string]uint64{}

	return res
}

///////////////////////////////////////////////////////////////////////////////
// Originally copied from vendor/github.com/openconfig/gnmi/client/gnmi/client.go

// Dial dials the server and logs in. If error is nil, returned conn has
// established a connection to d.
func Dial(ctx context.Context, deviceId string, d *client.Destination) (conn *grpc.ClientConn, retry bool, dErr error) {
	if len(d.Addrs) != 1 { // can't happen
		return nil, false, fmt.Errorf("Dial: d.Addrs must contain exactly one entry; got %+v", d.Addrs)
	}
	opts := []grpc.DialOption{
		grpc.WithTimeout(d.Timeout),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)),
	}

	switch d.TLS {
	case nil:
		opts = append(opts, grpc.WithInsecure())
	default:
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(d.TLS)))
	}

	conn, err := grpc.DialContext(ctx, d.Addrs[0], opts...)
	if err != nil {
		return nil, true, fmt.Errorf("%s dialing %s, timeout %v: error %v", deviceId, d.Addrs[0], d.Timeout, err)
	}

	if d.Credentials != nil {
		err := loginCheckJunos(deviceId, d.Credentials, conn)
		if err != nil {
			return nil, false, err // auth errors are non-retryable.
		}
	}
	return conn, true, nil
}

func loginCheckJunos(deviceId string, cred *client.Credentials, conn *grpc.ClientConn) error {
	if cred.Username != "" && cred.Password != "" {
		user := cred.Username
		pass := cred.Password

		lc := auth_pb.NewLoginClient(conn)
		dat, err := lc.LoginCheck(context.Background(),
			&auth_pb.LoginRequest{UserName: user,
				Password: pass,
				ClientId: deviceId,
			})
		if err != nil {
			return fmt.Errorf("Could not login: %v", err)
		}
		if !dat.Result {
			return fmt.Errorf("loginCheckJunos failed")
		}
	}
	return nil
}
