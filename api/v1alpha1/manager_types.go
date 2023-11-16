package v1alpha1

// Manager defines the desired state of Manager
type Manager struct {
	ConfigMap map[string]string `json:"configMap,omitempty"`
	// Image is the  container image. default is apachepulsar/pulsar-all:latest
	Image ContainerImage `json:"image,omitempty"`

	// Labels specifies the labels to attach to pods the operator creates for
	// the broker cluster.
	Labels map[string]string `json:"labels,omitempty"`

	// Size (DEPRECATED) is the expected size of the broker cluster.
	Size int32 `json:"size,omitempty"`

	// Pod defines the policy to create pod for the broker cluster.
	//
	// Updating the pod does not take effect on any existing pods.
	Pod PodPolicy `json:"pod,omitempty"`

	StorageClassName string `json:"storageClassName,omitempty"`
	StorageCapacity  int32  `json:"storageCapacity,omitempty"`
	FrontendNodePort int32  `json:"frontendNodePort,omitempty"`
	BackendNodePort  int32  `json:"backendNodePort,omitempty"`
}

func (m *Manager) SetDefault(c *Pulsar) bool {
	changed := false

	if m.Image.SetDefault(c, ManagerComponent) {
		changed = true
	}

	if m.Size == 0 {
		m.Size = ManagerDefaultNodeNum
		changed = true
	}

	if m.Pod.SetDefault(c, ManagerComponent) {
		changed = true
	}
	return changed
}
