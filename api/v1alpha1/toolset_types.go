package v1alpha1

type Toolset struct {
	// Image is the  container image. default is apachepulsar/pulsar-all:latest
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Image ContainerImage `json:"image,omitempty"`

	// Labels specifies the labels to attach to pods the operator creates for
	// the broker cluster.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Labels map[string]string `json:"labels,omitempty"`

	// Size (DEPRECATED) is the expected size of the toolset cluster.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Replicas int32 `json:"replicas,omitempty"`

	// Pod defines the policy to create pod for the toolset cluster.
	//
	// Updating the pod does not take effect on any existing pods.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Pod PodPolicy `json:"pod,omitempty"`
}

func (t *Toolset) SetDefault(c *PulsarCluster) bool {
	changed := false

	if t.Image.SetDefault(c, ToolsetComponent) {
		changed = true
	}

	if t.Pod.SetDefault(c, ToolsetComponent) {
		changed = true
	}
	return changed
}
