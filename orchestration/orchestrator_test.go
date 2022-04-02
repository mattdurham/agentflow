package orchestration

import (
	"agentflow/config"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestNewOrchestrator(t *testing.T) {
	cfgStr := `
nodes:
- name: generator
  outputs: 
  - filter
  metric_generator:
    spawn_interval: 1m
- name: filter
  outputs:
  - rw
  metric_filter:
    filters:
    - action: add_label
      add_label: test_label
      add_value: test
- name: rw
  fake_metric_remote_write: {}
`
	cfg := &config.Config{}
	err := yaml.Unmarshal([]byte(cfgStr), cfg)
	require.NoError(t, err)
	orch := NewOrchestrator(*cfg)
	err = orch.StartActorSystem()
	require.NoError(t, err)
	no, found := orch.nameToPID["filter"]
	require.True(t, found)
	require.NotNil(t, no)
	no, found = orch.nameToPID["rw"]
	require.True(t, found)
	require.NotNil(t, no)
	no, found = orch.nameToPID["generator"]
	require.True(t, found)
	require.NotNil(t, no)
}
