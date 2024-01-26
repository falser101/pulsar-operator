package v1alpha1

// Proxy defines the desired state of Proxy
type Proxy struct {
	// Image is the  container image. default is apachepulsar/pulsar-all:latest
	Image ContainerImage `json:"image,omitempty"`

	// Labels specifies the labels to attach to pods the operator creates for
	// the proxy cluster.
	Labels map[string]string `json:"labels,omitempty"`

	// Size (DEPRECATED) is the expected size of the proxy cluster.
	Size int32 `json:"size,omitempty"`

	// Pod defines the policy to create pod for the proxy cluster.
	//
	// Updating the pod does not take effect on any existing pods.
	Pod PodPolicy `json:"pod,omitempty"`

	HttpServerPort   int32 `json:"httpServerPort,omitempty"`
	PulsarServerPort int32 `json:"pulsarServerPort,omitempty"`
}

func (p *Proxy) SetDefault(cluster *Pulsar) bool {
	changed := false

	if p.Image.SetDefault(cluster, ProxyComponent) {
		changed = true
	}

	if p.Size == 0 {
		p.Size = ProxyClusterDefaultNodeNum
		changed = true
	}

	if p.Pod.SetDefault(cluster, ProxyComponent) {
		changed = true
	}
	return changed
}
