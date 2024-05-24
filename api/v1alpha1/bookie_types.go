package v1alpha1

// Bookie defines the desired state of Bookie
// +k8s:openapi-gen=true
type Bookie struct {
	// Image is the  container image. default is apachepulsar/pulsar-all:latest
	Image ContainerImage `json:"image,omitempty"`

	// Labels specifies the labels to attach to pods the operator creates for
	// the bookie cluster.
	Labels map[string]string `json:"labels,omitempty"`

	// Replicas is the expected size of the bookie cluster.
	Replicas int32 `json:"replicas,omitempty"`

	// Pod defines the policy to create pod for the bookie cluster.
	//
	// Updating the pod does not take effect on any existing pods.
	Pod PodPolicy `json:"pod,omitempty"`

	// ConfigData is the configuration data for bookie
	ConfigData map[string]string `json:"configData,omitempty"`

	// Storage class name
	//
	// PVC of storage class name
	JournalStorageClassName string `json:"journalStorageClassName,omitempty"`

	// Storage request capacity(Gi unit) for journal
	JournalStorageCapacity string `json:"journalStorageCapacity,omitempty"`

	LedgersStorageClassName string `json:"ledgersStorageClassName,omitempty"`

	// Storage request capacity(Gi unit) for ledgers
	LedgersStorageCapacity string `json:"ledgersStorageCapacity,omitempty"`
}

func (b *Bookie) SetDefault(cluster *PulsarCluster) bool {
	changed := false

	if b.Image.SetDefault(cluster, BookieComponent) {
		changed = true
	}

	if b.Replicas == 0 {
		b.Replicas = BookieClusterDefaultNodeNum
		changed = true
	}

	if b.Pod.SetDefault(cluster, BookieComponent) {
		changed = true
	}

	if b.JournalStorageClassName != "" && b.JournalStorageCapacity == "" {
		b.JournalStorageCapacity = JournalStorageDefaultCapacity
		changed = true
	}

	if b.LedgersStorageClassName != "" && b.LedgersStorageCapacity == "" {
		b.LedgersStorageCapacity = LedgersStorageDefaultCapacity
		changed = true
	}

	return changed
}
