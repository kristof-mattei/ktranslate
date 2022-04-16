package st

import (
	"context"

	go_metrics "github.com/kentik/go-metrics"

	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/st/receive_dialout"
	"github.com/kentik/ktranslate/pkg/kt"
)

const ()

var ()

type KentikStreamer interface {
	Close()
}

func NewStreamer(ctx context.Context, proto string, maxBatchSize int, log logger.Underlying, registry go_metrics.Registry, jchfChan chan []*kt.JCHF, apic *api.KentikApi) (KentikStreamer, error) {
	return receive_dialout.NewAuto(ctx, proto, maxBatchSize, log, registry, jchfChan, apic)
}
