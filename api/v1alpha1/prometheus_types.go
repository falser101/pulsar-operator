package v1alpha1

// Prometheus defines the desired state of Prometheus
type Prometheus struct {
	// Image is the  container image. default is apachepulsar/pulsar-all:latest
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Image ContainerImage `json:"image,omitempty"`

	// Labels specifies the labels to attach to pods the operator creates for
	// the broker cluster.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Labels map[string]string `json:"labels,omitempty"`

	// Replicas is the expected size of the broker cluster.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Replicas int32 `json:"replicas,omitempty"`

	// Pod defines the policy to create pod for the broker cluster.
	//
	// Updating the pod does not take effect on any existing pods.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Pod PodPolicy `json:"pod,omitempty"`

	// +operator-sdk:csv:customresourcedefinitions:type=spec
	StorageCapacity int32 `json:"storageCapacity,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	StorageClassName string `json:"storageClassName,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	NodePort int32 `json:"nodePort,omitempty"`
}

func (p *Prometheus) SetDefault(c *PulsarCluster) bool {
	changed := false

	if p.Image.SetDefault(c, MonitorPrometheusComponent) {
		changed = true
	}

	if p.Replicas == 0 {
		p.Replicas = 1
		changed = true
	}

	if p.Pod.SetDefault(c, MonitorPrometheusComponent) {
		changed = true
	}

	if p.StorageClassName != "" && p.StorageCapacity == 0 {
		p.StorageCapacity = PrometheusStorageDefaultCapacity
		changed = true
	}
	return changed
}
