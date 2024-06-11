package v1alpha1

type PVC struct {
	StorageClassName string `json:"storageClassName,omitempty"`

	Capacity string `json:"capacity,omitempty"`
}
