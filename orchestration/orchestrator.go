package orchestration

import (
	"agentflow/components"
	"agentflow/components/remotewrites"
	"agentflow/config"
	"agentflow/types/actorstate"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
)

type Orchestrator struct {
	cfg config.Config

	actorSystem *actor.ActorSystem
	rootContext *actor.RootContext
	nameToPID   map[string]*actor.PID
	pidToName   map[*actor.PID]string
	nameToActor map[string]actorstate.FlowActor
}

func NewOrchestrator(cfg config.Config) *Orchestrator {
	return &Orchestrator{
		cfg:         cfg,
		nameToPID:   map[string]*actor.PID{},
		pidToName:   map[*actor.PID]string{},
		nameToActor: map[string]actorstate.FlowActor{},
	}
}

func (u *Orchestrator) StartActorSystem() error {
	u.actorSystem = actor.NewActorSystem()
	u.rootContext = actor.NewRootContext(u.actorSystem, nil)
	// Generate the Nodes
	for _, nodeCfg := range u.cfg.Nodes {
		if nodeCfg.MetricGenerator != nil {
			no, err := components.NewMetricGenerator(nodeCfg.Name, *nodeCfg.MetricGenerator)
			if err != nil {
				return err
			}
			u.addPID(no)

		} else if nodeCfg.MetricFilter != nil {
			no, err := components.NewMetricFilter(nodeCfg.Name, *nodeCfg.MetricFilter)
			if err != nil {
				return err
			}
			u.addPID(no)
		} else if nodeCfg.FakeMetricRemoteWrite != nil {
			no, err := remotewrites.NewFakeMetricRemoteWrite(nodeCfg.Name)
			if err != nil {
				return err
			}
			u.addPID(no)
		} else if nodeCfg.SimpleRemoteWrite != nil {
			no, err := remotewrites.NewSimpleMetric(nodeCfg.Name, *nodeCfg.SimpleRemoteWrite)
			if err != nil {
				return err
			}
			u.addPID(no)
		}
	}
	// Assign all the outputs
	for _, nodeCfg := range u.cfg.Nodes {
		outs := make([]*actor.PID, 0)
		for _, out := range nodeCfg.Outputs {
			pid, found := u.nameToPID[out]
			if !found {
				return fmt.Errorf("unable to find output %s on node named %s", out, nodeCfg.Name)
			}
			outs = append(outs, pid)
		}
		n, found := u.nameToPID[nodeCfg.Name]
		if !found {
			return fmt.Errorf("unable to find %s in name to pid", nodeCfg.Name)
		}
		u.rootContext.Send(n, actorstate.Init{Children: outs})
	}
	// Start the system
	for _, v := range u.nameToPID {
		u.rootContext.Send(v, actorstate.Start{})
	}
	return nil
}

func (u *Orchestrator) addPID(no actorstate.FlowActor) {
	props := actor.PropsFromProducer(func() actor.Actor { return no })
	pid := u.rootContext.Spawn(props)
	u.nameToPID[no.Name()] = pid
	u.pidToName[pid] = no.Name()
	u.nameToActor[no.Name()] = no
}
