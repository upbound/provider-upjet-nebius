# provider-upjet-nebius — Architecture

A Crossplane provider for Nebius Cloud, built with [Upjet](https://github.com/crossplane/upjet)
on top of the [terraform-provider-nebius](https://github.com/nebius/terraform-provider-nebius)
**plugin-framework** provider. Upjet ingests the Terraform provider's schema +
docs and code-generates Crossplane CRDs and controllers.

This document describes only the project's structure and conventions so an agent
can quickly find the files and configuration involved in common tasks
(especially **adding a new resource**). It is not a usage guide.

## Key facts

- **Provider name:** `nebius` (Makefile `PROVIDER_NAME`). Module path:
  `github.com/upbound/provider-nebius`.
- **Upstream TF provider:** `nebius/nebius` v0.6.12, pinned in the `Makefile`
  (`TERRAFORM_PROVIDER_VERSION`). It is a Terraform **plugin-framework** provider,
  so resources are wired through `WithTerraformPluginFrameworkProvider` and use
  framework-style external names — not the legacy SDKv2 path.
- **Dual scope:** every resource is generated twice — once **cluster-scoped**
  (group suffix `.upbound.io`) and once **namespaced** (group suffix
  `.m.upbound.io`). Configuration is mirrored under `config/cluster/` and
  `config/namespaced/`.
- **TF CLI not used at runtime:** the framework provider is linked in-process;
  there is no terraform binary or provider mirror in the runtime image.

## Directory layout

```
config/                       Upjet provider configuration (hand-written) — the heart of resource setup
  external_name.go            ExternalNameConfigs map: TF resource -> external-name strategy
  groups.go                   GroupMap + ReplaceGroupWords: TF resource -> (ShortGroup, Kind)
  provider.go                 GetProvider() (cluster); embeds schema.json + provider-metadata.yaml
  provider_namespaced.go      GetProviderNamespaced() + registerTFConversions() (singleton-list -> embedded)
  generated.lst               JSON array of TF resource names with a generated CRD (kept in sync by codegen)
  schema.json                 Scraped TF provider schema (generated; input to codegen)
  provider-metadata.yaml      Scraped TF docs: per-resource examples + argument docs (input to codegen)
  cluster/
    provider.go               init(): registers each group's Configure into ProviderConfiguration
    configuration.go          Configurator registry type + global ProviderConfiguration var
    vpcv1/config.go           Per-group AddResourceConfigurator calls: cross-resource references live here
  namespaced/                 Mirror of cluster/ for the namespaced provider (vpcv1/config.go is duplicated)

apis/                         Generated Go API types (zz_*.go). Subtree: {cluster,namespaced}/<group>/<version>/
  cluster/vpcv1/v1beta1/      zz_<resource>_types.go, zz_<resource>_terraformed.go, zz_generated.*.go
  namespaced/vpcv1/v1beta1/   (mirror)

internal/
  controller/{cluster,namespaced}/<group>/<resource>/   Generated controllers (one dir per resource)
  controller/{cluster,namespaced}/zz_setup*.go           Generated controller registration
  clients/                    TerraformSetupBuilder + provider client wiring
  config/                     IdentityType auth-selection logic (hand-written)
  features/                   Feature flags

cmd/
  generator/main.go           Codegen entrypoint (run by `make generate`)
  provider/<subpkg>/zz_main.go  Provider binary entrypoints (monolith is the default SUBPACKAGE)

hack/                         main.go.tmpl (provider main template) + embed.go (MainTemplate)
package/crds/                 Generated CRD YAML (one file per resource per scope)
examples/                     Hand-curated examples (the source of truth shipped to users / used by uptest)
examples-generated/           Codegen output from provider-metadata.yaml examples (copy into examples/, then edit)
docs/resources/               Sparse-checkout target for pulled upstream TF docs (populated by `make generate`)
.work/                        Scratch: cloned TF provider docs, terraform workdir, schema generation
```

## Generation pipeline (`make generate`)

1. `generate.init` runs `$(TERRAFORM_PROVIDER_SCHEMA)` → installs terraform 1.5.7,
   `terraform providers schema -json` → writes `config/schema.json`.
   (Env `NEBIUS_TERRAFORM_PROVIDER_DISABLE_WRITE_ONLY` / `..._DISABLE_DYNAMIC` are
   set **only at generate time** to keep dynamic/write-only attrs out of the schema.)
2. `pull-docs` sparse-checks-out `docs/resources` from the upstream TF provider repo
   at the pinned tag into `.work/<source>/`.
3. `cmd/generator` loads `GetProvider` + `GetProviderNamespaced`, applies the
   config in `config/`, and writes `apis/`, `internal/controller/`, `package/crds/`,
   and `examples-generated/`.

`provider-metadata.yaml` is the doc-derived metadata (examples + argument
descriptions). A resource only gets an `examples-generated/` file if it has an
`examples:` block in that metadata; otherwise build the example from the CRD schema.

## Resource-configuration knobs (where to edit when adding a resource)

Adding a TF resource requires touching, in order:

1. **`config/external_name.go`** — add to `ExternalNameConfigs`. Nebius resources
   use a server-computed `id`, so:
   `config.FrameworkResourceWithComputedIdentifier("id", "<placeholder-NID>")`.
   The placeholder must be a syntactically valid Nebius ID (NID) of the form
   `<type>-<routingCode><weakID>`; `e0t` is the default SDK routing code, so the
   convention is `<type>-e0t000000000000000`. The `<type>` prefix is the SDK
   resource type (see **NID type prefixes** below). This map drives both the
   external-name behavior and `WithTerraformPluginFrameworkIncludeList` (only
   resources listed here are generated).

2. **`config/groups.go`** — add to `GroupMap`. Almost always
   `ReplaceGroupWords("<group>", <count>)`. `ReplaceGroupWords(group, count)`
   strips the `nebius_` prefix, drops the first `count` words to form the group,
   and camel-cases the remainder into the Kind. For `nebius_vpc_v1_*` resources
   use `ReplaceGroupWords("vpcv1", 2)` → group `vpcv1`, e.g.
   `nebius_vpc_v1_security_group` → Kind `SecurityGroup`.

3. **`config/cluster/vpcv1/config.go` AND `config/namespaced/vpcv1/config.go`**
   (keep them identical) — add an `AddResourceConfigurator` block **only if the
   resource has ID references to other resources**. Each reference:
   ```go
   r.References["<dotted_snake_case_tf_path>"] = config.Reference{
       TerraformName: "nebius_vpc_v1_<target>",
       Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
   }
   ```
   Reference keys are the **dotted snake_case Terraform attribute paths**
   (e.g. `network_id`, `ipv4_private_pools.pools.id`). Do **not** use `[*]` for
   list elements — Upjet silently drops such references.

After editing config, run `make generate`, then curate `examples/`.

## Naming & convention reference

- **NID type prefixes** (for external-name placeholders) come from the Nebius
  Go SDK (`github.com/nebius/gosdk`) protobuf field annotations. VPC examples:
  `vpcnetwork`, `vpcsubnet`, `vpcpool`, `vpcroutetable`, `vpcallocation`,
  `vpcroute`, `vpcsecuritygroup`, `vpcsecurityrule`.
- **`parent_id`** is present on most Nebius resources. It is *not* always the
  project — check `metadata.parent_id` in `provider-metadata.yaml`:
  e.g. for `nebius_vpc_v1_route` it "represents the RouteTable" and for
  `nebius_vpc_v1_security_rule` it "represents the SecurityGroup". When `parent_id`
  represents another managed resource, configure it as a reference; when it
  represents the project/IAM-container, leave it as a literal/`${data.nebius_project_id}`.
- **Group/Kind**: `<group>.nebius.upbound.io` (cluster) and
  `<group>.nebius.m.upbound.io` (namespaced); version is `v1beta1` (set by
  `ExternalNameConfigurations`).
- **Examples**: `examples/<scope>/<group>/<version>/<kind-lowercase>.yaml`. Project
  parent is templated as `parentId: ${data.nebius_project_id}` for uptest; cross-
  resource fields use `<field>Selector.matchLabels` against
  `testing.upbound.io/example-name`.

## E2E testing

`make uptest` / `make e2e` drive examples through a live control plane.
`uptest.json` holds the test credentials; `uptest-data.ini` / `UPTEST_DATASOURCE_PATH`
inject dynamic values such as `${data.nebius_project_id}`.
