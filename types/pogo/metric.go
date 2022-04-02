package pogo

import "time"

type Metric struct {
	name     string
	value    float64
	ts       time.Time
	labels   map[string]string
	metadata map[string]string
}

func NewMetric(name string, value float64, ts time.Time, labels map[string]string, metadata map[string]string) Metric {
	m := Metric{
		name:     name,
		value:    value,
		ts:       ts,
		labels:   labels,
		metadata: metadata,
	}
	if m.labels == nil {
		m.labels = make(map[string]string)
	}
	if m.metadata == nil {
		m.labels = make(map[string]string)
	}
	m.labels["__name__"] = name
	return m
}

func CopyMetric(in Metric) Metric {
	return Metric{
		name:     in.Name(),
		value:    in.Value(),
		ts:       in.Timestamp(),
		labels:   in.Labels(),
		metadata: in.Metadata(),
	}
}

func (m *Metric) Name() string {
	return m.name
}

func (m *Metric) Value() float64 {
	return m.value
}

func (m *Metric) Timestamp() time.Time {
	return m.ts
}

func (m *Metric) Labels() map[string]string {
	return copyMap(m.labels)
}

func (m *Metric) Metadata() map[string]string {
	return copyMap(m.metadata)
}

func copyMap(in map[string]string) map[string]string {
	newMap := make(map[string]string)
	for k, v := range in {
		newMap[k] = v
	}
	return newMap
}
