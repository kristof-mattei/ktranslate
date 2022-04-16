package x

import (
	"context"
	"time"

	go_metrics "github.com/kentik/go-metrics"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/x/arista"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/x/gnmi"
	"github.com/kentik/ktranslate/pkg/kt"
)

// Code to handle various vendor extensions to snmp.
type Extension interface {
	Run(context.Context, time.Duration)
	GetName() string
}

func NewExtension(jchfChan chan []*kt.JCHF, conf *kt.SnmpDeviceConfig, metrics *kt.SnmpDeviceMetric, registry go_metrics.Registry, log logger.ContextL) (Extension, error) {
	if conf.Ext == nil { // No extensions set.
		return nil, nil
	}

	if conf.Ext.EAPIConfig != nil {
		return arista.NewEAPIClient(jchfChan, conf, metrics, log)
	} else if conf.Ext.GNMIConfig != nil {
		return gnmi.NewGNMIClient(jchfChan, conf, registry, log)
	}

	return nil, nil
}
