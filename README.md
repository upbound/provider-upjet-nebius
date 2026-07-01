# Provider Nebius

`provider-nebius` is a [Crossplane](https://crossplane.io/) provider for
[Nebius AI Cloud](https://nebius.com/), built with
[Upjet](https://github.com/crossplane/upjet) and backed by the
[`nebius/nebius`](https://github.com/nebius/terraform-provider-nebius)
Terraform provider.

It exposes XRM-conformant managed resources that let you manage Nebius
infrastructure — VPCs, compute instances, managed Kubernetes clusters, IAM,
storage, and more — directly from Kubernetes or Upbound. Every resource is
available in two flavors: cluster-scoped (`*.nebius.upbound.io`) and
namespaced (`*.nebius.m.upbound.io`).

## Authentication

The provider authenticates using an `identity.type` selected on the
`ProviderConfig`, combined with credentials read from the referenced `Secret`.
Two modes are supported:

### Mode 1 — Token

Authenticate using a static IAM token.

```json
{
  "token": "<IAM token>"
}
```

### Mode 2 — ServiceAccount

Authenticate using a service-account key.

```json
{
  "account_id": "<service account ID>",
  "public_key_id": "<public key ID>",
  "private_key": "<private key>"
}
```

## Getting Started

### 1. Install the provider

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-nebius
spec:
  package: xpkg.upbound.io/upbound/provider-nebius:v1.0.0
```

### 2. Create a credentials Secret

Replace `<IAM token>` with your actual value, or use the `ServiceAccount`
fields shown above instead.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: nebius-creds
  namespace: upbound-system
type: Opaque
stringData:
  creds: |
    {
      "token": "<IAM token>"
    }
```

### 3. Create a ProviderConfig

```yaml
apiVersion: nebius.upbound.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  projectID: project-e00example
  identity:
    type: Token
  credentials:
    source: Secret
    secretRef:
      name: nebius-creds
      namespace: upbound-system
      key: creds
```

### 4. Create a managed resource

```yaml
apiVersion: vpc.nebius.upbound.io/v1beta1
kind: Network
metadata:
  name: example
spec:
  forProvider:
    name: example
    parentId: project-e00example
  providerConfigRef:
    name: default
```

More examples for each resource are available under [`examples/cluster/`](examples/cluster/)
and [`examples/namespaced/`](examples/namespaced/).

## ProviderConfig fields

| Field | Required | Description |
|---|---|---|
| `spec.identity.type` | Yes | Authentication mode: `Token` or `ServiceAccount` |
| `spec.projectID` | No | Default Nebius project ID for project-parented resources; can be overridden per-resource via `spec.forProvider.parentId` |
| `spec.credentials.source` | Yes | One of `Secret`, `InjectedIdentity`, `Environment`, `Filesystem` |
| `spec.credentials.secretRef` | When source=Secret | Reference to the credentials Secret |
| `spec.reconciliationPolicy` | No | Rate-limiting policy for reconciliation |

## Developing

### Code generation

```console
make generate
```

This runs the Upjet code generator against the pinned `nebius/nebius`
Terraform provider schema and docs, and writes the generated APIs,
controllers, and CRDs for both the cluster-scoped and namespaced variants.

### Run locally against a cluster

```console
make run
```

### Run end-to-end tests

```console
make e2e
```

## Reporting issues

Please open an [issue](https://github.com/upbound/provider-upjet-nebius/issues)
for bug reports, feature requests, or questions.
