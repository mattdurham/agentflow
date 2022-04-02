package main

import (
	"agentflow/config"
	"agentflow/orchestration"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

func main() {
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
	err = orch.StartActorSystem()
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Minute)
}
