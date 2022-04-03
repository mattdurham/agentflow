package main

import (
	"agentflow/config"
	"agentflow/orchestration"
	"github.com/AsynkronIT/protoactor-go/actor"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

func main() {
	as := actor.NewActorSystem()
	root := actor.NewRootContext(as, nil)
	cfgStr, err := ioutil.ReadFile("/Users/matt/Utils/agent_flow_configs/agent_flow_simply.yml")
	if err != nil {
		panic(err)
	}
	cfg := &config.Config{}
	err = yaml.Unmarshal(cfgStr, cfg)
	if err != nil {
		panic(err)
	}
	orch := orchestration.NewOrchestrator(*cfg)
	err = orch.StartActorSystem(as, root)
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Minute)
}
