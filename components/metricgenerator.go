package components

import (
	"agentflow/config"
	"agentflow/types"
	"agentflow/types/actorstate"
	"agentflow/types/pogo"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/scheduler"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"math/rand"
	"time"
)

type MetricGenerator struct {
	config config.MetricGenerator
	cancel scheduler.CancelFunc
	out    []*actor.PID
	self   *actor.PID
	name   string
	log log.Logger
	index int
}

func (mg *MetricGenerator) Output() actorstate.InOutType {
	return actorstate.Metrics
}

func (mg *MetricGenerator) AllowableInputs() []actorstate.InOutType {
	return []actorstate.InOutType{}
}

func NewMetricGenerator(name string, cfg config.MetricGenerator, global *types.Global) (actorstate.FlowActor, error) {
	return &MetricGenerator{
		config: cfg,
		name:   name,
		log: global.Log,
	}, nil
}

func (mg *MetricGenerator) Name() string {
	return mg.name
}

func (mg *MetricGenerator) PID() *actor.PID {
	return mg.self
}

func (mg *MetricGenerator) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case actorstate.Init:
		mg.self = ctx.Self()
		mg.out = msg.Children
	case actorstate.Start:
		sched := scheduler.NewTimerScheduler(ctx)
		mg.cancel = sched.SendRepeatedly(1*time.Millisecond, mg.config.SpawnInterval, ctx.Self(), "SendMore")
	case actorstate.Done:
		mg.cancel()
	case string:
		if msg != "SendMore" {
			return
		}
		metrics := mg.makeMetrics()
		for _, o := range mg.out {
			cpy := make([]pogo.Metric, len(metrics))
			copy(cpy, metrics)
			ctx.Send(o, cpy)
		}
		_ = level.Info(mg.log).Log("msg","creating logs","length",len(metrics), "index", mg.index)
		mg.index++
	}
}

func (mg *MetricGenerator) makeMetrics() []pogo.Metric {
	metrics := make([]pogo.Metric, 0)
	for i := 0; i < 100; i++ {
		metrics = append(metrics, pogo.NewMetric(fmt.Sprintf("gen_%d",i), rand.Float64(), time.Now(), nil, nil))
	}
	return metrics
}
