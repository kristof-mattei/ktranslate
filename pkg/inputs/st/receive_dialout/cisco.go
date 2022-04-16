package receive_dialout

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/kt/counters"

	telemetry "github.com/kentik/bigmuddy-network-telemetry-proto/proto_go"
	cisco_generic_counters "github.com/kentik/bigmuddy-network-telemetry-proto/proto_go/cisco_ios_xr_infra_statsd_oper/infra_statistics/interfaces/interface/latest/generic_counters"

	"github.com/golang/protobuf/proto"
)

// CiscoTelemetryMDT plugin for IOS XR, IOS XE and NXOS platforms
type CiscoTelemetryMDT struct {
	logger.ContextL

	raw          chan []byte
	jchfChan     chan []*kt.JCHF
	metrics      *STMetric
	do           DialoutServer
	inputs       chan *kt.JCHF
	maxBatchSize int
	dev          *kt.Device
	devName      string
	counterSet   map[string]*counters.CounterSet
}

func NewCiscoTelemetryMDT(devName string, dev *kt.Device, maxBatchSize int, log logger.Underlying, jchfChan chan []*kt.JCHF, metrics *STMetric) (*CiscoTelemetryMDT, error) {
	c := &CiscoTelemetryMDT{
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "st:cisco:" + devName}, log),
		raw:          make(chan []byte, maxBatchSize),
		jchfChan:     jchfChan,
		metrics:      metrics,
		maxBatchSize: maxBatchSize,
		inputs:       make(chan *kt.JCHF, maxBatchSize),
		dev:          dev,
		devName:      devName,
		counterSet:   map[string]*counters.CounterSet{},
	}

	c.Infof("Cisco MDI setup for %s", dev.Name)

	return c, nil
}

func (c *CiscoTelemetryMDT) Process(pkt []byte) {
	select {
	case c.raw <- pkt:
	default:
		c.Warnf("ST chan full")
	}
}

func (c *CiscoTelemetryMDT) Run(ctx context.Context) {
	c.Infof("Cisco running")
	sendTicker := time.NewTicker(kt.SendBatchDuration)
	defer sendTicker.Stop()
	batch := make([]*kt.JCHF, 0, c.maxBatchSize)

	for {
		select {
		case local := <-c.inputs:
			batch = append(batch, local)
			if len(batch) >= c.maxBatchSize {
				c.jchfChan <- batch
				batch = make([]*kt.JCHF, 0, c.maxBatchSize)
			}
		case <-sendTicker.C:
			if len(batch) > 0 {
				c.jchfChan <- batch
				batch = make([]*kt.JCHF, 0, c.maxBatchSize)
			}
		case raw := <-c.raw:
			err := c.handleTelemetry(raw)
			if err != nil {
				c.Errorf("Invalid Input: %v", err)
			}

		case <-ctx.Done():
			c.Infof("Cisco ST Done")
			sendTicker.Stop()
			return
		}
	}
}

var wantedType = reflect.TypeOf((*cisco_generic_counters.IfstatsbagGeneric)(nil))

// Handle telemetry packet from any transport, decode and add as measurement
func (c *CiscoTelemetryMDT) handleTelemetry(data []byte) error {
	itelem := &telemetry.Telemetry{}
	l := len(data)
	err := proto.Unmarshal(data, itelem)
	if err != nil {
		if l > 1024 {
			data = data[:1024]
		}
		return err
	}
	c.Debugf("Cisco MDT received %d bytes of %s", l, itelem.EncodingPath)

	protoMsg := telemetry.EncodingPathToMessageReflectionSet(&telemetry.ProtoKey{
		EncodingPath: itelem.EncodingPath,
		// version ignored (for now?)
	})
	if protoMsg == nil {
		return fmt.Errorf("Unknown Cisco MDT EncodingPath: %s", itelem.EncodingPath)
	}

	// These are the only types I get from the Cogent Cisco router
	//   ImDescInfo
	//   IfstatsbagProto
	//   IfstatsbagGeneric
	//   StatsdbagDatarate
	// Ignore everything but IfstatsbagGeneric
	contentType := protoMsg.MessageReflection(telemetry.PROTO_CONTENT_MSG)
	keysType := protoMsg.MessageReflection(telemetry.PROTO_KEYS_MSG)
	// c.log.Infof(c.logPrefix, "protoMsg content: %s, keys: %s", contentType.String(), keysType.String())
	if contentType != wantedType { // reflect.Type is comparable
		c.Debugf("Got %v, expected %v; skipping it", contentType, wantedType)
		return nil
	}

	numRows, numFlow := 0, 0
	for i, gpb := range itelem.DataGpb.GetRow() {
		numRows++
		c.Debugf("itelem.DataGpb.GetRow: %d", i)

		keysValue := reflect.New(keysType.Elem())
		err := proto.Unmarshal(gpb.Keys, keysValue.Interface().(proto.Message))
		if err != nil {
			c.Infof("Cisco failed to decode keys: %s", keysType.String())
			continue
		}

		contentValue := reflect.New(contentType.Elem())
		err = proto.Unmarshal(gpb.Content, contentValue.Interface().(proto.Message))
		if err != nil {
			c.Infof("Cisco failed to decode: %s", contentType.String())
			continue
		}

		stats, ok := contentValue.Interface().(*cisco_generic_counters.IfstatsbagGeneric)
		if !ok {
			// shouldn't happen, but just in case
			c.Debugf("Could not convert interface to cisco_generic_counters.IfstatsbagGeneric")
			continue
		}
		keys, ok := keysValue.Interface().(*cisco_generic_counters.IfstatsbagGeneric_KEYS)
		if !ok {
			// shouldn't happen, but just in case
			c.Debugf("Could not convert interface to cisco_generic_counters.IfstatsbagGeneric_KEYS")
			continue
		}

		flow := c.convertToCHF(keys.InterfaceName, stats, nil)
		if flow == nil {
			c.Debugf("Not sending flow for: %s", keys.InterfaceName)
		} else {
			numFlow++
			c.Debugf("Sending flow for: %s", keys.InterfaceName)
			c.metrics.STIntTotal.Mark(int64(1))
			c.inputs <- flow
		}
	}
	c.Debugf("Processed %d rows of DataGpb; sent %d flow records", numRows, numFlow)

	numRows, numFlow = 0, 0
	for _, kv := range itelem.GetDataGpbkv() {
		// c.log.Debugf(c.logPrefix, "DataGpbkv %d: %v", i, kv.String())

		// If this ever breaks, for a more general approach see
		// https://github.com/influxdata/telegraf/blob/6665b48008f38fc77118503be3ce4dc3f469273e/plugins/inputs/cisco_telemetry_mdt/cisco_telemetry_mdt.go.

		// Let's be as specific as possible here.
		if kv.Name == "" &&
			len(kv.Fields) == 2 &&
			(kv.Fields[0].Name == "keys" &&
				len(kv.Fields[0].Fields) == 1 &&
				kv.Fields[0].Fields[0].Name == "interface-name") &&
			(kv.Fields[1].Name == "content" &&
				len(kv.Fields[1].Fields) > 0) {

			numRows++

			interfaceName := kv.Fields[0].Fields[0].GetStringValue()
			flow := c.convertToCHF(interfaceName, nil, kv.Fields[1].Fields)
			if flow == nil {
				c.Debugf("Not sending flow for: %s", interfaceName)
			} else {
				numFlow++
				c.Debugf("Sending flow for: %s", interfaceName)
				c.metrics.STIntTotal.Mark(int64(1))
				c.inputs <- flow
			}
		}
	}
	c.Debugf("Processed %d rows of DataGpbkv; sent %d flow records", numRows, numFlow)

	return nil
}

// Stop listener and cleanup
func (c *CiscoTelemetryMDT) Close() {
	// ???
}

func (c *CiscoTelemetryMDT) convertToCHF(ifcName string,
	stats *cisco_generic_counters.IfstatsbagGeneric,
	fields []*telemetry.TelemetryField) *kt.JCHF {

	c.Debugf("Converting to CHF, ifc: %s, fields: %v", ifcName, fields)

	in := kt.NewJCHF()
	in.CustomStr = make(map[string]string)
	in.CustomInt = make(map[string]int32)
	in.CustomBigInt = make(map[string]int64)
	in.EventType = kt.KENTIK_EVENT_TYPE
	in.Provider = kt.ProviderFlowDevice
	if c.dev != nil {
		in.DeviceName = c.dev.Name
		in.DeviceId = c.dev.ID
		in.CompanyId = c.dev.CompanyID
		c.dev.SetUserTags(in.CustomStr)
	} else {
		in.DeviceName = c.devName
	}

	var inPkts, inBytes, outPkts, outBytes, inErr, outErr, inDis, outDis, inPktsMulti, outPktsMulti, inPktsBroad, outPktsBroad uint64
	switch {
	case stats != nil:
		inPkts = c.counterSet[ST_InUcastPkts].SetValueAndReturnDelta(ifcName, stats.PacketsReceived)
		inBytes = c.counterSet[ST_InOctets].SetValueAndReturnDelta(ifcName, stats.BytesReceived)
		outPkts = c.counterSet[ST_OutUcastPkts].SetValueAndReturnDelta(ifcName, stats.PacketsSent)
		outBytes = c.counterSet[ST_OutOctets].SetValueAndReturnDelta(ifcName, stats.BytesSent)
		inErr = c.counterSet[ST_InErrors].SetValueAndReturnDelta(ifcName, uint64(stats.InputErrors))
		outErr = c.counterSet[ST_OutErrors].SetValueAndReturnDelta(ifcName, uint64(stats.OutputErrors))
		inDis = c.counterSet[ST_InDiscards].SetValueAndReturnDelta(ifcName, uint64(stats.InputQueueDrops))
		outDis = c.counterSet[ST_OutDiscards].SetValueAndReturnDelta(ifcName, uint64(stats.OutputQueueDrops))
		inPktsMulti = c.counterSet[ST_InMulticastPkts].SetValueAndReturnDelta(ifcName, stats.MulticastPacketsReceived)
		outPktsMulti = c.counterSet[ST_OutMulticastPkts].SetValueAndReturnDelta(ifcName, stats.MulticastPacketsSent)
		inPktsBroad = c.counterSet[ST_InBroadcastPkts].SetValueAndReturnDelta(ifcName, stats.BroadcastPacketsReceived)
		outPktsBroad = c.counterSet[ST_OutBroadcastPkts].SetValueAndReturnDelta(ifcName, stats.BroadcastPacketsSent)
	case fields != nil:
		for _, field := range fields {
			switch field.Name {
			case "packets-received":
				inPkts = c.counterSet[ST_InUcastPkts].SetValueAndReturnDelta(ifcName, field.GetUint64Value())
			case "bytes-received":
				inBytes = c.counterSet[ST_InOctets].SetValueAndReturnDelta(ifcName, field.GetUint64Value())
			case "packets-sent":
				outPkts = c.counterSet[ST_OutUcastPkts].SetValueAndReturnDelta(ifcName, field.GetUint64Value())
			case "bytes-sent":
				outBytes = c.counterSet[ST_OutOctets].SetValueAndReturnDelta(ifcName, field.GetUint64Value())
			case "input-errors":
				inErr = c.counterSet[ST_InErrors].SetValueAndReturnDelta(ifcName, uint64(field.GetUint32Value()))
			case "output-errors":
				outErr = c.counterSet[ST_OutErrors].SetValueAndReturnDelta(ifcName, uint64(field.GetUint32Value()))
			case "input-queue-drops":
				inDis = c.counterSet[ST_InDiscards].SetValueAndReturnDelta(ifcName, uint64(field.GetUint32Value()))
			case "output-queue-drops":
				outDis = c.counterSet[ST_OutDiscards].SetValueAndReturnDelta(ifcName, uint64(field.GetUint32Value()))
			case "multicast-packets-received":
				inPktsMulti = c.counterSet[ST_InMulticastPkts].SetValueAndReturnDelta(ifcName, field.GetUint64Value())
			case "multicast-packets-sent":
				outPktsMulti = c.counterSet[ST_OutMulticastPkts].SetValueAndReturnDelta(ifcName, field.GetUint64Value())
			case "broadcast-packets-received":
				inPktsBroad = c.counterSet[ST_InBroadcastPkts].SetValueAndReturnDelta(ifcName, field.GetUint64Value())
			case "broadcast-packets-sent":
				outPktsBroad = c.counterSet[ST_OutBroadcastPkts].SetValueAndReturnDelta(ifcName, field.GetUint64Value())
			}
		}
	}

	in.CustomStr["if_interface_name"] = ifcName
	in.CustomBigInt[ST_InUcastPkts] = int64(inPkts)
	in.CustomBigInt[ST_InOctets] = int64(inBytes)
	in.CustomBigInt[ST_OutUcastPkts] = int64(outPkts)
	in.CustomBigInt[ST_OutOctets] = int64(outBytes)
	in.CustomBigInt[ST_InErrors] = int64(inErr)
	in.CustomBigInt[ST_OutErrors] = int64(outErr)
	in.CustomBigInt[ST_InDiscards] = int64(inDis)
	in.CustomBigInt[ST_OutDiscards] = int64(outDis)
	in.CustomBigInt[ST_InMulticastPkts] = int64(inPktsMulti)
	in.CustomBigInt[ST_OutMulticastPkts] = int64(outPktsMulti)
	in.CustomBigInt[ST_InBroadcastPkts] = int64(inPktsBroad)
	in.CustomBigInt[ST_OutBroadcastPkts] = int64(outPktsBroad)

	return in
}
