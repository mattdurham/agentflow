package remotewrites

import (
	"agentflow/config"
	"agentflow/types/actorstate"
	"agentflow/types/pogo"
	"context"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/castai/promwrite"
)

type SimpleMetric struct {
	name   string
	self   *actor.PID
	client *promwrite.Client
}

func (f *SimpleMetric) AllowableInputs() []actorstate.InOutType {
	return []actorstate.InOutType{actorstate.Metrics}
}

func (f *SimpleMetric) Output() actorstate.InOutType {
	return actorstate.None
}

func NewSimpleMetric(name string, cfg config.SimpleRemoteWrite ) (actorstate.FlowActor, error) {
	cl := promwrite.NewClient(cfg.URL)
	return &SimpleMetric{
		name:   name,
		client: cl,
	}, nil
}

func (f *SimpleMetric) Receive(c actor.Context) {
	switch msg := c.Message().(type) {
	case actorstate.Start:
		f.self = c.Self()
	case []pogo.Metric:
		writable := make([]promwrite.TimeSeries, len(msg))
		for i, m := range msg {
			lblSet := m.Labels()
			lbls := make([]promwrite.Label, 0)

			for k, v := range lblSet {
				lbls = append(lbls, promwrite.Label{
					Name:  k,
					Value: v,
				})
			}

			ts := promwrite.TimeSeries{
				Labels: lbls,
				Sample: promwrite.Sample{
					Time:  m.Timestamp(),
					Value: m.Value(),
				},
			}
			writable[i] = ts
		}

		_, err := f.client.Write(context.Background(), &promwrite.WriteRequest{
			TimeSeries: writable,
		})
		if err != nil {
			fmt.Printf("error writing %s \n", err)
		}
	}
}

func (f *SimpleMetric) Name() string {
	return f.name
}

func (f SimpleMetric) PID() *actor.PID {
	return f.PID()
}
