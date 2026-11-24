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

	// ImpersonateServiceAccountID is the ID of a service account to
	// impersonate. The provider authenticates with the identity selected by
	// Type, then exchanges its token for a token of this service account and
	// manages resources as it. The authenticating identity must be allowed to
	// impersonate the service account.
	// +optional
	ImpersonateServiceAccountID string `json:"impersonateServiceAccountID,omitempty"`
}
