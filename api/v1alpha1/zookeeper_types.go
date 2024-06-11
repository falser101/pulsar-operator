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

	// zk ports
	Ports ZKPorts `json:"ports,omitempty"`

	ConfigData map[string]string `json:"configData,omitempty"`

	Volumes ZKVolumes `json:"volumes,omitempty"`
}

type ZKPorts struct {
	Http           int32 `json:"http,omitempty"`
	Client         int32 `json:"client,omitempty"`
	ClientTls      int32 `json:"clientTls,omitempty"`
	Follower       int32 `json:"follower,omitempty"`
	LeaderElection int32 `json:"leaderElection,omitempty"`
}

type ZKVolumes struct {
	Data PVC `json:"data,omitempty"`
}

func (z *Zookeeper) SetDefault(c *PulsarCluster) bool {
	changed := false

	if z.Image.SetDefault(c, ZookeeperComponent) {
		changed = true
	}

	if z.Replicas == 0 {
		z.Replicas = ZookeeperClusterDefaultNodeNum
		changed = true
	}

	if z.Pod.SetDefault(c, ZookeeperComponent) {
		changed = true
	}

	if z.setPortsDefault(c) {
		changed = true
	}

	return changed
}

func (z *Zookeeper) setPortsDefault(c *PulsarCluster) (changed bool) {
	if z.Ports.Http == 0 {
		z.Ports.Http = 8000
		changed = true
	}

	if z.Ports.Client == 0 {
		z.Ports.Client = 2181
		changed = true
	}

	if z.Ports.Follower == 0 {
		z.Ports.Follower = 2888
		changed = true
	}

	if z.Ports.LeaderElection == 0 {
		z.Ports.LeaderElection = 3888
		changed = true
	}
	return
}
