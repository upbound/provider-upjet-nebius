#!/usr/bin/env bash
set -aeuo pipefail

echo "Running setup.sh"
echo "Creating cloud credential secret..."
${KUBECTL} -n upbound-system  create secret generic provider-secret --from-literal=credentials="${UPTEST_CLOUD_CREDENTIALS}" --dry-run=client -o yaml | ${KUBECTL} apply -f -

echo "Waiting until provider is healthy..."
${KUBECTL} wait provider.pkg --all --for condition=Healthy --timeout 5m

echo "Waiting for all pods to come online..."
${KUBECTL} -n upbound-system  wait --for=condition=Available deployment --all --timeout=5m

# TODO: update identityType to ServiceAccount
echo "Creating a default provider config..."
cat <<EOF | ${KUBECTL} apply -f -
apiVersion: nebius.upbound.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  identity:
    type: ServiceAccount
  credentials:
    source: Secret
    secretRef:
      name: provider-secret
      namespace: upbound-system 
      key: credentials
EOF

# TODO: update identityType to ServiceAccount
echo "Creating a default cluster provider config (v2-style)..."
cat <<EOF | ${KUBECTL} apply -f -
apiVersion: nebius.m.upbound.io/v1beta1
kind: ClusterProviderConfig
metadata:
  name: default
spec:
  identity:
    type: ServiceAccount
  credentials:
    source: Secret
    secretRef:
      name: provider-secret
      namespace: upbound-system
      key: credentials
EOF

echo "Creating the impersonation provider configs..."
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
datasource="${UPTEST_DATASOURCE_PATH:-}"
if [[ -n "${datasource}" && "${datasource}" != /* ]]; then
  datasource="${repo_root}/${datasource}"
fi
delegated_service_account_id=""
if [[ -f "${datasource}" ]]; then
  delegated_service_account_id="$(sed -n 's/^nebius_delegated_service_account_id:[[:space:]]*//p' "${datasource}" | tr -d "\"'\r")"
fi
if [[ -z "${delegated_service_account_id}" ]]; then
  echo "WARNING: no nebius_delegated_service_account_id in the uptest datasource, skipping the impersonation provider configs"
else
  for scope in cluster namespaced; do
    sed "s|\${data.nebius_delegated_service_account_id}|${delegated_service_account_id}|g" "${repo_root}/examples/${scope}/providerconfig/impersonation.yaml" | ${KUBECTL} apply -f -
  done
fi

${KUBECTL} wait provider.pkg --all --for condition=Healthy --timeout 5m
${KUBECTL} -n upbound-system  wait --for=condition=Available deployment --all --timeout=5m
