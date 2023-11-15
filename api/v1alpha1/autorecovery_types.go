package v1alpha1

// AutoRecovery defines the desired state of AutoRecovery
// +k8s:openapi-gen=true
type AutoRecovery struct {
	// Image is the  container image. default is apachepulsar/pulsar-all:latest
	Image ContainerImage `json:"image,omitempty"`

	// Labels specifies the labels to attach to pods the operator creates for
	// the bookie cluster.
	Labels map[string]string `json:"labels,omitempty"`

	// Size (DEPRECATED) is the expected size of the bookie cluster.
	Size int32 `json:"size,omitempty"`

	// Pod defines the policy to create pod for the bookie cluster.
	//
	// Updating the pod does not take effect on any existing pods.
	Pod PodPolicy `json:"pod,omitempty"`
}

func (b *AutoRecovery) SetDefault(cluster *Pulsar) bool {
	changed := false

	if b.Image.SetDefault(cluster, AutoRecoveryComponent) {
		changed = true
	}

	if b.Size == 0 {
		b.Size = BookieClusterDefaultNodeNum
		changed = true
	}

	if b.Pod.SetDefault(cluster, AutoRecoveryComponent) {
		changed = true
	}

	return changed
}
