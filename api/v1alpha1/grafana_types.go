package v1alpha1

// Grafana defines the desired state of Grafana
type Grafana struct {
	Host     string `json:"host,omitempty"`
	NodePort int32  `json:"nodePort,omitempty"`
}
