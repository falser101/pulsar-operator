package v1alpha1

// +operator-sdk:csv:customresourcedefinitions:type=spec
type Provider string

const (
	JWT Provider = "jwt"
)

type Auth struct {
	// Authentication is the authentication policy for the broker cluster.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	AuthenticationEnabled bool `json:"authenticationEnabled,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	AuthorizationEnabled bool `json:"authorizationEnabled,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	AuthenticationProvider Provider `json:"authenticationProvider,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	SuperUserRoles []string `json:"superUserRoles,omitempty"`
	ProxyRoles     []string `json:"proxyRoles,omitempty"`
}
