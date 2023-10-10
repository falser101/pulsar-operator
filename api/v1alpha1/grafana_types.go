package v1alpha1

type Security struct {
	AdminUser     string `json:"adminUser,omitempty"`
	AdminPassword string `json:"adminPassword,omitempty"`
}

// Grafana defines the desired state of Grafana
type Grafana struct {
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
	Pod              PodPolicy `json:"pod,omitempty"`
	NodePort         int32     `json:"nodePort,omitempty"`
	StorageCapacity  int32     `json:"storageCapacity,omitempty"`
	StorageClassName string    `json:"storageClassName,omitempty"`
	Security         Security  `json:"security,omitempty"`
}

func (g *Grafana) SetDefault(c *PulsarCluster) bool {
	changed := false

	if g.Image.SetDefault(c, MonitorGrafanaComponent) {
		changed = true
	}

	if g.Size == 0 {
		g.Size = 1
		changed = true
	}

	if g.Pod.SetDefault(c, MonitorGrafanaComponent) {
		changed = true
	}

	if g.StorageClassName != "" && g.StorageCapacity == 0 {
		g.StorageCapacity = GrafanaStorageDefaultCapacity
		changed = true
	}

	if g.Security.AdminUser == "" {
		g.Security.AdminUser = GrafanaDefaultAdminUser
		changed = true
	}

	if g.Security.AdminPassword == "" {
		g.Security.AdminPassword = GrafanaDefaultAdminPassword
		changed = true
	}
	return changed
}
