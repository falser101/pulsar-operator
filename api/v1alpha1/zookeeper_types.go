package v1alpha1

// Zookeeper defines the desired state of Zookeeper
type Zookeeper struct {
	// Image is the  container image. default is apachepulsar/pulsar-all:latest
	Image ContainerImage `json:"image,omitempty"`

	// Labels specifies the labels to attach to pods the operator creates for
	// the broker cluster.
	Labels map[string]string `json:"labels,omitempty"`

	// Replicas is the expected size of the broker cluster.
	Replicas int32 `json:"replicas,omitempty"`

	// Pod defines the policy to create pod for the broker cluster.
	//
	// Updating the pod does not take effect on any existing pods.
	Pod PodPolicy `json:"pod,omitempty"`

	// Storage class name
	//
	// PVC of storage class name
	StorageClassName string `json:"storageClassName,omitempty"`

	StorageCapacity int32 `json:"storageCapacity,omitempty"`
}

func (z *Zookeeper) SetDefault(c *Pulsar) bool {
	changed := false

	if z.Image.SetDefault(c, ZookeeperComponent) {
		changed = true
	}

	if z.Replicas == 0 {
		z.Replicas = ZookeeperClusterDefaultNodeNum
		changed = true
	}

	if z.StorageClassName != "" && z.StorageCapacity == 0 {
		z.StorageCapacity = ZookeeperClusterDefaultStorageCapacity
		changed = true
	}

	if z.Pod.SetDefault(c, ZookeeperComponent) {
		changed = true
	}
	return changed
}
