package v1alpha1

// +operator-sdk:csv:customresourcedefinitions:type=spec
type Provider string

const (
	JWT Provider = "jwt"
)

type Authentication struct {
	// Authentication is the authentication policy for the broker cluster.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Enabled bool `json:"enabled,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Provider Provider `json:"provider,omitempty"`
}
