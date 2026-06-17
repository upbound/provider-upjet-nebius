# Provider Nebius

`provider-nebius` is a [Crossplane](https://crossplane.io/) provider that is built using [Upjet](https://github.com/crossplane/upjet) code generation tools and exposes XRM-conformant managed resources for [Nebius AI Cloud](https://nebius.com/), based on the [Nebius Terraform provider](https://github.com/nebius/terraform-provider-nebius).

## Developing

Run code-generation pipeline:
```console
make generate
```

Run against a Kubernetes cluster:

```console
make run
```

Build, push, and install:

```console
make all
```

Build binary:

```console
make build
```

## Report a Bug

For filing bugs, suggesting improvements, or requesting new features, please
open an [issue](https://github.com/upbound/provider-nebius/issues).
