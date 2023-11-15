package v1alpha1

// Zookeeper defines the desired state of Zookeeper
type Zookeeper struct {
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
}

func (z *Zookeeper) SetDefault(c *Pulsar) bool {
	changed := false

	if z.Image.SetDefault(c, ZookeeperComponent) {
		changed = true
	}

	if z.Size == 0 {
		z.Size = ZookeeperClusterDefaultNodeNum
		changed = true
	}

	if z.Pod.SetDefault(c, ZookeeperComponent) {
		changed = true
	}
	return changed
}
