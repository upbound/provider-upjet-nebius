package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	xpv2 "github.com/crossplane/crossplane-runtime/v2/apis/common/v2"
	"github.com/crossplane/upjet/v2/apis/configuration/v1alpha1"

	internalconfig "github.com/upbound/provider-nebius/internal/config"
)

// A ProviderConfigSpec defines the desired state of a cluster-scoped
// ClusterProviderConfig. Its credential secret references carry an explicit
// namespace, because a cluster-scoped config can reference secrets in any
// namespace. The namespaced ProviderConfig uses NamespacedProviderConfigSpec
// instead, which omits the namespace.
type ProviderConfigSpec struct {
	// +optional
	// +kubebuilder:validation:XValidation:rule="!has(self.exponentialFailureRateLimiter) || !has(self.exponentialFailureRateLimiter.baseDelay) || has(self.exponentialFailureRateLimiter.maxDelay) || duration(self.exponentialFailureRateLimiter.baseDelay) <= duration('60s')",message="when maxDelay is omitted it defaults to 60s; baseDelay must be <= 60s"
	ReconciliationPolicy *v1alpha1.ReconciliationPolicy `json:"reconciliationPolicy,omitempty"`

	// Credentials required to authenticate to this provider.
	Credentials ProviderCredentials `json:"credentials"`

	// ProjectID is the Nebius project ID used as the default parent for
	// project-parented resources. Individual resources may override it via
	// spec.forProvider.parentId. Optional.
	// +kubebuilder:validation:Optional
	ProjectID *string `json:"projectID,omitempty"`

	// Identity specifies the authentication identity configuration.
	// +kubebuilder:validation:Required
	Identity *internalconfig.Identity `json:"identity"`
}

// ProviderCredentials required to authenticate. The secret reference carries an
// explicit namespace, as used by the cluster-scoped ClusterProviderConfig.
type ProviderCredentials struct {
	// Source of the provider credentials.
	// +kubebuilder:validation:Enum=None;Secret;InjectedIdentity;Environment;Filesystem
	Source xpv1.CredentialsSource `json:"source"`

	xpv1.CommonCredentialSelectors `json:",inline"`
}

// A NamespacedProviderConfigSpec defines the desired state of a namespaced
// ProviderConfig. It mirrors ProviderConfigSpec but its credential secret
// references omit the namespace: they implicitly resolve to the namespace of
// the referencing managed resource.
type NamespacedProviderConfigSpec struct {
	// +optional
	// +kubebuilder:validation:XValidation:rule="!has(self.exponentialFailureRateLimiter) || !has(self.exponentialFailureRateLimiter.baseDelay) || has(self.exponentialFailureRateLimiter.maxDelay) || duration(self.exponentialFailureRateLimiter.baseDelay) <= duration('60s')",message="when maxDelay is omitted it defaults to 60s; baseDelay must be <= 60s"
	ReconciliationPolicy *v1alpha1.ReconciliationPolicy `json:"reconciliationPolicy,omitempty"`

	// Credentials required to authenticate to this provider.
	Credentials NamespacedProviderCredentials `json:"credentials"`

	// ProjectID is the Nebius project ID used as the default parent for
	// project-parented resources. Individual resources may override it via
	// spec.forProvider.parentId. Optional.
	// +kubebuilder:validation:Optional
	ProjectID *string `json:"projectID,omitempty"`

	// Identity specifies the authentication identity configuration.
	// +kubebuilder:validation:Required
	Identity *internalconfig.Identity `json:"identity"`
}

// NamespacedProviderCredentials required to authenticate. The secret reference
// omits the namespace and resolves to the namespace of the referencing managed
// resource. It mirrors xpv1.CommonCredentialSelectors but uses a
// LocalSecretKeySelector for the secret reference.
type NamespacedProviderCredentials struct {
	// Source of the provider credentials.
	// +kubebuilder:validation:Enum=None;Secret;InjectedIdentity;Environment;Filesystem
	Source xpv1.CredentialsSource `json:"source"`

	// Fs is a reference to a filesystem location that contains credentials that
	// must be used to connect to the provider.
	// +optional
	Fs *xpv1.FsSelector `json:"fs,omitempty"`

	// Env is a reference to an environment variable that contains credentials
	// that must be used to connect to the provider.
	// +optional
	Env *xpv1.EnvSelector `json:"env,omitempty"`

	// A SecretRef is a reference to a secret key in the same namespace as the
	// referencing managed resource that contains the credentials that must be
	// used to connect to the provider.
	// +optional
	SecretRef *xpv1.LocalSecretKeySelector `json:"secretRef,omitempty"`
}

// A ProviderConfigStatus reflects the observed state of a ProviderConfig.
type ProviderConfigStatus struct {
	xpv1.ProviderConfigStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// A ProviderConfig configures a Nebius provider.
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="SECRET-NAME",type="string",JSONPath=".spec.credentials.secretRef.name",priority=1
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,provider,nebius}
type ProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NamespacedProviderConfigSpec `json:"spec"`
	Status ProviderConfigStatus         `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProviderConfigList contains a list of ProviderConfig.
type ProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderConfig `json:"items"`
}

// +kubebuilder:object:root=true

// A ProviderConfigUsage indicates that a resource is using a ProviderConfig.
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="CONFIG-NAME",type="string",JSONPath=".providerConfigRef.name"
// +kubebuilder:printcolumn:name="RESOURCE-KIND",type="string",JSONPath=".resourceRef.kind"
// +kubebuilder:printcolumn:name="RESOURCE-NAME",type="string",JSONPath=".resourceRef.name"
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,provider,nebius}
type ProviderConfigUsage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	xpv2.TypedProviderConfigUsage `json:",inline"`
}

// +kubebuilder:object:root=true

// ProviderConfigUsageList contains a list of ProviderConfigUsage
type ProviderConfigUsageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderConfigUsage `json:"items"`
}

// +kubebuilder:object:root=true

// A ClusterProviderConfig configures a Nebius provider.
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="SECRET-NAME",type="string",JSONPath=".spec.credentials.secretRef.name",priority=1
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:resource:scope=Cluster,categories={crossplane,provider,nebius}
type ClusterProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderConfigSpec   `json:"spec"`
	Status ProviderConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterProviderConfigList contains a list of ProviderConfig.
type ClusterProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterProviderConfig `json:"items"`
}
