package v1alpha1

type Monitor struct {
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Prometheus Prometheus `json:"prometheus,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Grafana Grafana `json:"grafana,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Enable bool `json:"enable,omitempty"`
}

func (m *Monitor) SetDefault(c *Pulsar) bool {
	changed := false
	if m.Prometheus.SetDefault(c) {
		changed = true
	}
	if m.Grafana.SetDefault(c) {
		changed = true
	}
	return changed
}
