package v1alpha1

// Broker defines the desired state of Broker
type Broker struct {
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

	Authentication Authentication `json:"auth,omitempty"`
}

type Authentication struct {
	// Authentication is the authentication policy for the broker cluster.
	Enabled bool `json:"enabled,omitempty"`
}

func (b *Broker) SetDefault(cluster *Pulsar) bool {
	changed := false

	if b.Image.SetDefault(cluster, BrokerComponent) {
		changed = true
	}

	if b.Size == 0 {
		b.Size = BrokerClusterDefaultNodeNum
		changed = true
	}

	if b.Pod.SetDefault(cluster, BrokerComponent) {
		changed = true
	}
	return changed
}
