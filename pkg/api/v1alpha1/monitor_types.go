package v1alpha1

type Monitor struct {
	Prometheus Prometheus `json:"prometheus,omitempty"`
	Grafana    Grafana    `json:"grafana,omitempty"`
	Enable     bool       `json:"enable,omitempty"`
}

func (m *Monitor) SetDefault(c *PulsarCluster) bool {
	changed := false
	if m.Prometheus.SetDefault(c) {
		changed = true
	}
	if m.Grafana.SetDefault(c) {
		changed = true
	}
	return changed
}
