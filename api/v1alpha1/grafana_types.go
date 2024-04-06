package v1alpha1

type Security struct {
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	AdminUser string `json:"adminUser,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	AdminPassword string `json:"adminPassword,omitempty"`
}

// Grafana defines the desired state of Grafana
type Grafana struct {
	// Image is the  container image. default is apachepulsar/pulsar-all:latest
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Image ContainerImage `json:"image,omitempty"`

	// Labels specifies the labels to attach to pods the operator creates for
	// the broker cluster.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Labels map[string]string `json:"labels,omitempty"`

	// Replicas is the expected size of the broker cluster.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Replicas int32 `json:"replicas,omitempty"`

	// Pod defines the policy to create pod for the broker cluster.
	//
	// Updating the pod does not take effect on any existing pods.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Pod PodPolicy `json:"pod,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	NodePort int32 `json:"nodePort,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	StorageCapacity int32 `json:"storageCapacity,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	StorageClassName string `json:"storageClassName,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Security Security `json:"security,omitempty"`
}

func (g *Grafana) SetDefault(c *PulsarCluster) bool {
	changed := false

	if g.Image.SetDefault(c, MonitorGrafanaComponent) {
		changed = true
	}

	if g.Replicas == 0 {
		g.Replicas = 1
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
