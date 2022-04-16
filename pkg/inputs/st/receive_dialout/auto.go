package receive_dialout

import (
	"context"
	"fmt"
	"net"
	"strings"
	"syscall"
	"time"

	go_metrics "github.com/kentik/go-metrics"

	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

	tt "github.com/Juniper/jtimon/protos/telemetry_top"
	ct "github.com/kentik/bigmuddy-network-telemetry-proto/proto_go"

	"github.com/golang/protobuf/proto"
)

const (
	MAX_UDP_PACKET = 4096

	TYPE_JUNOS   = "junos"
	TYPE_IOS     = "ios"
	TYPE_AUTO    = "auto"
	TYPE_UNKNOWN = "unknown"

	ST_InOctets         = "in-octets"
	ST_InUcastPkts      = "in-unicast-pkts"
	ST_OutOctets        = "out-octets"
	ST_OutUcastPkts     = "out-unicast-pkts"
	ST_InErrors         = "in-errors"
	ST_OutErrors        = "out-errors"
	ST_InDiscards       = "in-discards"
	ST_OutDiscards      = "out-discards"
	ST_OutMulticastPkts = "out-multicast-pkts"
	ST_OutBroadcastPkts = "out-broadcast-pkts"
	ST_InMulticastPkts  = "in-multicast-pkts"
	ST_InBroadcastPkts  = "in-broadcast-pkts"

	DeviceUpdateDuration = 1 * time.Hour
)

type STMetric struct {
	STIntTotal go_metrics.Meter
}

type AutoMetric struct {
	STErrors go_metrics.Meter
	Pkts     go_metrics.Meter
}

type DialoutServer interface {
	Run(context.Context)
	Close()
	Process([]byte)
}

type input struct {
	pkt  []byte
	addr net.Addr
}

// Figures out what kinda of packets are coming in and then spins up the appopriate handler for this.
type Auto struct {
	logger.ContextL
	jchfChan     chan []*kt.JCHF
	apic         *api.KentikApi
	registry     go_metrics.Registry
	inputs       chan *kt.JCHF
	devices      map[string]*kt.Device
	maxBatchSize int
	rawChan      chan *input
	sock         *net.UDPConn
	streamers    map[string]DialoutServer
	metrics      *AutoMetric
}

// Noop here because we don't know what to do.
func NewAuto(ctx context.Context, listenAddr string, maxBatchSize int, log logger.Underlying, registry go_metrics.Registry, jchfChan chan []*kt.JCHF, apic *api.KentikApi) (*Auto, error) {
	a := &Auto{
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "st:auto"}, log),
		jchfChan:     jchfChan,
		apic:         apic,
		devices:      apic.GetDevicesAsMap(0),
		maxBatchSize: maxBatchSize,
		rawChan:      make(chan *input, maxBatchSize),
		inputs:       make(chan *kt.JCHF, maxBatchSize),
		registry:     registry,
		streamers:    map[string]DialoutServer{},
		metrics: &AutoMetric{
			STErrors: go_metrics.GetOrRegisterMeter("st.errors", registry),
			Pkts:     go_metrics.GetOrRegisterMeter("st.pkts", registry),
		},
	}

	addr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		return nil, err
	}

	sock, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	a.Infof("Listening for ST on %s", addr.String())
	a.sock = sock
	go a.run(ctx)
	return a, nil
}

// Block until we get a packet and then figure out what to do with it.
func (a *Auto) run(ctx context.Context) {
	go a.listen()
	deviceTicker := time.NewTicker(DeviceUpdateDuration)
	defer deviceTicker.Stop()

	for {
		select {
		case raw := <-a.rawChan:
			if s, ok := a.streamers[raw.addr.String()]; ok {
				s.Process(raw.pkt) // And we're done.
				continue
			}

			// Here, We don't know about this sender.
			ptype, err := a.getType(raw.pkt)
			if err != nil {
				a.Errorf("Cannot guess proto type: %v", err)
				a.metrics.STErrors.Mark(1)
			} else {
				switch ptype {
				case TYPE_JUNOS:
					a.Warnf("JUNOS proto type not supported: %s", ptype)
					a.metrics.STErrors.Mark(1)

				case TYPE_IOS:
					a.Infof("Running as Cisco ST Listener for %s", raw.addr.String())
					pkts := strings.Split(raw.addr.String(), ":")
					var dev *kt.Device
					devName := pkts[0]
					if d, ok := a.devices[pkts[0]]; ok {
						dev = d
						devName = dev.Name
					}

					do, err := NewCiscoTelemetryMDT(devName, dev, a.maxBatchSize, a.GetLogger().GetUnderlyingLogger(), a.jchfChan, newSTMetric(devName, a.registry))
					if err != nil {
						a.Errorf("Cannot launch cisco ST: %v", err)
					} else {
						go do.Run(ctx)
						do.Process(raw.pkt) // Remember this packet.
						a.streamers[devName] = do
					}

				case TYPE_UNKNOWN:
					a.Warnf("Invalid proto type: %s", ptype)
					a.metrics.STErrors.Mark(1)

				default:
					a.Errorf("No proto type: %s -> %v", ptype, raw)
					a.metrics.STErrors.Mark(1)
				}
			}
		case <-deviceTicker.C:
			go func() {
				a.Infof("Updating the network flow device list.")
				a.devices = a.apic.GetDevicesAsMap(0)
			}()
		case <-ctx.Done():
			return
		}
	}
}

// Pass through to our underlaying impl.
func (a *Auto) Close() {
	for _, do := range a.streamers {
		do.Close()
	}
}

func (a *Auto) listen() {
	for {
		buf := make([]byte, MAX_UDP_PACKET)
		rlen, _, flags, addr, err := a.sock.ReadMsgUDP(buf, nil)
		if err != nil {
			a.Errorf("Error reading from UDP: %v", err)
			continue
		} else if (flags & syscall.MSG_TRUNC) != 0 {
			a.Errorf("Error reading from UDP: packet from %v exceeds payload length %v", addr, len(buf))
			continue
		}

		a.rawChan <- &input{pkt: buf[:rlen], addr: addr}
		a.metrics.Pkts.Mark(1)
	}
}

// Figures out if raw is a cisco or junos ST packet
func (a *Auto) getType(raw []byte) (string, error) {

	// Try junos first.
	ts := &tt.TelemetryStream{}
	if err := proto.Unmarshal(raw, ts); err == nil {
		if proto.HasExtension(ts.Enterprise, tt.E_JuniperNetworks) {
			return TYPE_JUNOS, nil
		}
	}

	// Now try cisco
	ctt := &ct.Telemetry{}
	err := proto.Unmarshal(raw, ctt)
	if err == nil {
		return TYPE_IOS, nil
	}

	// Now we don't know
	return TYPE_UNKNOWN, nil
}

func newSTMetric(devName string, registry go_metrics.Registry) *STMetric {
	return &STMetric{
		STIntTotal: go_metrics.GetOrRegisterMeter(fmt.Sprintf("st.flows^device_name=%s^force=true", devName), registry),
	}
}
