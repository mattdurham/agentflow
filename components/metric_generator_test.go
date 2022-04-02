package components

import (
	"agentflow/config"
	"agentflow/types/actorstate"
	"agentflow/types/pogo"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMetricGenerator(t *testing.T) {
	as := actor.NewActorSystem()
	root := actor.NewRootContext(as, nil)
	no, err := NewMetricGenerator("test", config.MetricGenerator{ SpawnInterval: 1 * time.Second})
	require.NoError(t, err)
	props := actor.PropsFromProducer(func() actor.Actor { return no })
	pid := root.Spawn(props)
	ta := &TestActor{}
	testActorProps := actor.PropsFromProducer(func() actor.Actor { return ta })
	testPID := root.Spawn(testActorProps)
	root.Send(pid, actorstate.Init{Children: []*actor.PID{testPID}})
	root.Send(pid, actorstate.Start{})
	time.Sleep(3 * time.Second)
	require.True(t, ta.Found)
}

type TestActor struct {
	Found bool
}

func (t *TestActor) Receive(c actor.Context) {
	switch c.Message().(type) {
	case []pogo.Metric:
		t.Found = true
		c.Send(c.Sender(), actorstate.Done{})
	}
}
