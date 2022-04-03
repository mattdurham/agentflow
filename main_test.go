package main

import (
	"agentflow/config"
	"agentflow/orchestration"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"testing"
	"time"
)

func TestLots(t *testing.T) {
	cfg := &config.Config{}
	cfgStr, err := ioutil.ReadFile("/Users/matt/Utils/agent_flow_configs/agent_flow_stress.yml")
	err = yaml.Unmarshal(cfgStr, cfg)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 30; i ++ {
		flt := config.MetricFilter{Filters: []config.MetricFilterFilter{
			{
				Action:   "add_label",
				AddValue: fmt.Sprintf("filter_%d", i),
				AddLabel: fmt.Sprintf("filter_%d",i),
			},
		}}
		if i != 29 {
			cfg.Nodes = append(cfg.Nodes, config.Node{
				Name:    fmt.Sprintf("filter_%d", i),
				Outputs: []string{fmt.Sprintf("filter_%d", i+1)},
				MetricFilter: &flt,
			})
		} else {
			cfg.Nodes = append(cfg.Nodes, config.Node{
				Name:    fmt.Sprintf("filter_%d", i),
				Outputs: []string{"rw"},
				MetricFilter: &flt,
			})
		}
	}
	orch := orchestration.NewOrchestrator(*cfg)
	err = orch.StartActorSystem()
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Minute)
}
