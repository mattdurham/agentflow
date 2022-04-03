Agent Flow is prototype for a component based Grafana Agent built around Actors and Messages.

Example config

```
nodes:
- name: generator
  outputs: 
  - filter
  metric_generator:
    spawn_interval: 10s
- name: filter
  outputs:
  - filter2
  metric_filter:
    filters:
    - action: add_label
      add_label: test_label
      add_value: test
- name: filter2
  outputs:
  - rw
  metric_filter:
    filters:
    - action: add_label
      add_label: filter2_label
      add_value: "this is from filter 2"
- name: rw
  simple_metric_remote_write:
    url: https://user:pass@example.com
- name: agent_logs
  agent_logs: {}
  outputs:
  - filewriter
- name: filewriter
  log_file_writer:
    path: /Users/matt/Utils/agent_flow_configs/logs.txt
- name: github
  github:
    repositories:
    - grafana/agent
  outputs:
  - rw
```

The above generates the follow PlantUML

```
@startuml 
[agent_logs] 
[generator] 
[filter] 
[filter2] 
[rw] 
[filewriter] 
[github] 
[filter] -> [filter2] 
[filter2] -> [rw] 
[github] -> [rw] 
[agent_logs] -> [filewriter] 
[generator] -> [filter] 
@enduml 
```

![](./assets/graph.png)

It also generates the following Mermaid

```mermaid
graph LR 
	FILTER[filter - *components.MetricFilter] 
	FILTER2[filter2 - *components.MetricFilter] 
	RW[rw - *remotewrites.SimpleMetric] 
	FILEWRITER[filewriter - *logs.FileWriter] 
	GITHUB[github - *integrations.Github] 
	AGENT_LOGS[agent_logs - *logs.Agent] 
	GENERATOR[generator - *components.MetricGenerator] 
	AGENT_LOGS --> FILEWRITER 
	GENERATOR --> FILTER 
	FILTER --> FILTER2 
	FILTER2 --> RW 
	GITHUB --> RW 
```
