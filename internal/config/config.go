package config

// IdentityType describes which authentication method the provider client uses.
type IdentityType string

const (
	// IdentityTypeToken authenticates using a static IAM token from the credentials secret.
	IdentityTypeToken IdentityType = "Token"
	// IdentityTypeServiceAccount authenticates using a service-account key
	// (account_id, public_key_id, private_key) from the credentials secret.
	IdentityTypeServiceAccount IdentityType = "ServiceAccount"
)

// Identity specifies the authentication identity configuration.
type Identity struct {
	// Type specifies the authentication method.
	// Token: authenticate using a static IAM token from the credentials secret key "token".
	// ServiceAccount: authenticate using a service-account key (account_id, public_key_id,
	// private_key) from the credentials secret.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=Token;ServiceAccount
	Type IdentityType `json:"type"`
}
