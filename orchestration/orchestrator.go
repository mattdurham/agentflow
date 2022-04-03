package orchestration

import (
	"agentflow/components"
	"agentflow/components/integrations"
	"agentflow/components/logs"
	"agentflow/components/remotewrites"
	"agentflow/config"
	"agentflow/types"
	"agentflow/types/actorstate"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/go-kit/kit/log"
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

func (u *Orchestrator) StartActorSystem(as *actor.ActorSystem, root *actor.RootContext) error {
	u.actorSystem = as
	u.rootContext = root
	var agentLog *logs.Agent
	// Find if they have defined the agent logger
	for _, nodeCfg := range u.cfg.Nodes {
		if nodeCfg.AgentLogs != nil {
			no, err := logs.NewAgent(nodeCfg.Name, root)
			if err != nil {
				return err
			}
			agentLog = no.(*logs.Agent)
			u.addPID(no)
			break
		}
	}
	// If they have not defined one, then create an internal one
	if agentLog == nil {
		no, err := logs.NewAgent("__agent_log", root)
		if err != nil {
			return err
		}
		agentLog = no.(*logs.Agent)
	}
	logger := log.NewLogfmtLogger(agentLog)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	global := &types.Global{
		Log: logger,
	}
	// Generate the Nodes
	for _, nodeCfg := range u.cfg.Nodes {
		if nodeCfg.MetricGenerator != nil {
			no, err := components.NewMetricGenerator(nodeCfg.Name, *nodeCfg.MetricGenerator, global)
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
		} else if nodeCfg.LogFileWriter != nil {
			no, err := logs.NewFileWriter(nodeCfg.Name, *nodeCfg.LogFileWriter)
			if err != nil {
				return err
			}
			u.addPID(no)
		} else if nodeCfg.Github != nil {
			no, err := integrations.NewGithub(nodeCfg.Name, nodeCfg.Github, *global)
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
