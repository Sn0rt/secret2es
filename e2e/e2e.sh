#!/bin/bash
set -x

function install_secret_to_k8s() {
  kubectl apply -f ./e2e/installed.yaml
  return 0
}

function build_secret2es() {
  make build
  return 0
}

function generate_external_secret_template() {
  ./secret2es es-gen -i e2e/templated.yaml -n tenant-b > e2e/render.yaml
  return 0
}

function apply_external_secret_template() {
  kubectl apply -f e2e/render.yaml
  return 0
}

function wait_external_secret_template_ready() {
  for i in $(seq 1 7);
  do
    kubectl wait --for=condition=Ready=True es/input"$i" --timeout=60s || (kubectl describe es/input"$i" && return 1)
  done
  kubectl get es -o wide
  return 0
}

function init_vault_kv_pair() {
  export APPROLE_NAME="approle1"
  export SECRET_PATH="secret/my-secret-approle1"

  VAULT_POD_NAME=$(kubectl get pods -l app.kubernetes.io/name=vault -o jsonpath='{.items[0].metadata.name}')
  kubectl cp ./e2e/vault_init.sh "$VAULT_POD_NAME":/home/vault/vault_init.sh
  kubectl exec "$VAULT_POD_NAME" -- /bin/sh -c "chmod +x /home/vault/vault_init.sh && /home/vault/vault_init.sh"

  ROLE_ID=$(kubectl exec -it "$VAULT_POD_NAME" -- vault read -field=role_id auth/approle/role/$APPROLE_NAME/role-id)
  SECRET_ID=$(kubectl exec -it "$VAULT_POD_NAME" -- vault write -field=secret_id -f auth/approle/role/$APPROLE_NAME/secret-id)
  echo "RoleID: $ROLE_ID"
  echo "RoleID=$ROLE_ID" >> "$GITHUB_ENV"
  echo "SecretID: $SECRET_ID"
  echo "SECRET_ID=$SECRET_ID" >> "$GITHUB_ENV"

  return 0
}

function install_cluster_secret_store() {
  kubectl apply -f ./e2e/ClusterSecretStore.yaml
  kubectl wait --for=condition=Ready=True clustersecretstores.external-secrets.io tenant-b --timeout=60s
  return 0
}

function get_secrets_from_k8s() {
  kubectl get secrets "$1" -o json
  return 0
}

function wait_external_secret_synced() {
  kubectl wait --for=condition=SecretSynced clustersecretstores.external-secrets.io tenant-b --timeout=60s
  return 0
}

function get_secret_content() {
    for i in $(seq 1 7);
    do
      echo -e "process input$i"
      kubectl get secret input"$i" -o jsonpath='{.data}' | jq -r 'to_entries[] | .key + "=" + (.value | @base64d)'
      echo -e "---"
    done
    return 0
}

function keep_secret_lifecycle_independent() {
    for i in $(seq 1 7);
    do
      kubectl delete es/input"$i"
      sleep 1
      kubectl get secret input"$i" || return 1
    done
    return 0
}