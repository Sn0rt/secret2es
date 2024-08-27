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
  VAULT_POD_NAME=$(kubectl get pods -l app.kubernetes.io/name=vault -o jsonpath='{.items[0].metadata.name}')
  kubectl exec -it "$VAULT_POD_NAME" -- vault kv put secret/test-foo my-value=from_vault_sn0rt
  kubectl exec -it "$VAULT_POD_NAME" -- vault kv patch secret/test-foo TEST_USERNAME=from_vault_true_username
  kubectl exec -it "$VAULT_POD_NAME" -- vault kv patch secret/test-foo TEST_PASSWORD=from_vault_true_password
  kubectl exec -it "$VAULT_POD_NAME" -- vault kv patch secret/test-foo TEST_DIST_LINUX=from_vault_Gentoo
  kubectl exec -it "$VAULT_POD_NAME" -- vault kv patch secret/test-foo MYSQL_PASSWD=from_vault_mysql_password
  kubectl exec -it "$VAULT_POD_NAME" -- vault kv patch secret/test-foo MYSQL_USER=from_vault_mysql_user
  kubectl exec -it "$VAULT_POD_NAME" -- vault kv patch secret/test-foo ACCESS_KEY=from_vault_ac_123123123
  kubectl exec -it "$VAULT_POD_NAME" -- vault kv patch secret/test-foo SECRET_KEY=from_vault_sc_321321321
  kubectl exec -it "$VAULT_POD_NAME" -- vault kv patch secret/test-foo TLS_CRT=LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNyakNDQVpZQ0NRQ1N0dllQRmZKUmF6QU5CZ2txaGtpRzl3MEJBUXNGQURBWk1SY3dGUVlEVlFRRERBNTUKYjNWeVpHOXRZV2x1TG1OdmJUQWVGdzB5TkRBNE1qY3dNakV6TVRaYUZ3MHlOVEE0TWpjd01qRXpNVFphTUJreApGekFWQmdOVkJBTU1Ebmx2ZFhKa2IyMWhhVzR1WTI5dE1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBCk1JSUJDZ0tDQVFFQW85amtKaHY5SXFoWDN0YXVSd3VGR28vdzhzZUlrUnZicmxCSEpxQUFrRTZQdGFiQlJyR2oKQnRyaDhSZlN0UHRpR2poYmdpM3dJUllEdU83NjZNKzFZaVljdC9YU1ZYd3hMRUpWWXBJWHhxQmlJd0hKL2xSTwoyNUVGN0pDeStmWlFuMFNvTnJiYlJSNUs3WnJHcUNUdk9aeFR5aTY1NUZGY1dEM2IwQTVBZkNLL3dxcDg2cW5zCktzc2YyaHRXR2c1U0dnR3lkcVpuZzc2ZVRCQzlRbGxGZGM2ZWR5NXpnZzRtaEkvN08wWEdVVys0VVoydmhxRTIKRnpkM0ZIWkR4MTVOL2IvVmZhZnlvclZYL3ZDcVVKNm5RUlZkcHJEQWlUdjJxOEJBN0t5SzdTNENpTnkwYkZSdgpaYzVxQUFoQzA5VkptanJyM3ZaSGpEQVBIakpSVGk0YXl3SURBUUFCTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCCkFRQUVYMExMNWtRU25nYmh1b2NHdlFwanNUVm9Cb05CelppYkorcHg0Mk9WRTFRd3M1ajN0US9udzBWVFZLZmUKMGlUTEw0QkxSaGNweW54RDJPQmZmWVEycTFydDBNeUMxRWdVZHg4VXl2ckhKNVBENDltcXdqVjdtR292cis4YQpCZWd3VTcybnFIQ3d4clpkWkl0ZVlHYjBTRm1rMklTYm5JdzVkQk5hTDJlb2hBV0g2Qm4xdncyNlZ3WnQ4eEtiCmZJRXhqczlYVU9VNXhyYUtPWDRrbTdBZFRyenBET2VRTXBOV042Z21zdjR4WVlRUHl3Y3lYUk1ZOTJ4dm1XM2kKMVpzQiswbkFqdWVkdWNGb1B0R21GVTR0MlB5UjBpQlZSZCtzQUl0akJvbVpvVzJnbGNRZDhHdUZ4aGxXRnBTRApHV0pROFFwZTdWV0ZmL0hMcEZPYU1rbS8KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
  kubectl exec -it "$VAULT_POD_NAME" -- vault kv patch secret/test-foo TLS_KEY=LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBbzlqa0podjlJcWhYM3RhdVJ3dUZHby93OHNlSWtSdmJybEJISnFBQWtFNlB0YWJCClJyR2pCdHJoOFJmU3RQdGlHamhiZ2kzd0lSWUR1Tzc2Nk0rMVlpWWN0L1hTVlh3eExFSlZZcElYeHFCaUl3SEoKL2xSTzI1RUY3SkN5K2ZaUW4wU29OcmJiUlI1SzdackdxQ1R2T1p4VHlpNjU1RkZjV0QzYjBBNUFmQ0svd3FwOAo2cW5zS3NzZjJodFdHZzVTR2dHeWRxWm5nNzZlVEJDOVFsbEZkYzZlZHk1emdnNG1oSS83TzBYR1VXKzRVWjJ2CmhxRTJGemQzRkhaRHgxNU4vYi9WZmFmeW9yVlgvdkNxVUo2blFSVmRwckRBaVR2MnE4QkE3S3lLN1M0Q2lOeTAKYkZSdlpjNXFBQWhDMDlWSm1qcnIzdlpIakRBUEhqSlJUaTRheXdJREFRQUJBb0lCQUNrM1dEMFY4Vm1VaTNZcwovdTQwUWFscTZDdktjZG8rN2NZdHY1aEJ5NktCZ0xrclY1ZFcvREd2UWdNS0FTRXgwMzNSQzRQMTFtQWNUNWRuCjFvcFdKY1NvM2JTUkMvWWhKYVdDa2tRWGlBK1pMTmF0am9pQjRNeHU4TlNQbWRZelZoaWFoczRzdFgvdm5OMmsKZjdDd2lkVXVOQTI2TDF6MThvcm9GTEdEeEVqMWRadzJmUE5pMnNFOFNvdno5LzdJM2R2dzlZaDg2dTBwZTdQawpvVE1kdGdWV3c5SnNWM0NyQld4NlBKazhIKzV2R2Z6VGk5aG8yWE5tNEx5aE1UMXZVaTVoY3RrQllmZFlTQUhDCmhFOUdNL2RiSmVrZDY2RmFDM1R0UTlvTXUybWJJdy9JUGJoOFRkZWxNT1RyQXR2MFhSSFJHMnRjeHdiMXJXTm0KVC9GY1RQRUNnWUVBMlNjR3FOVEMvTjM5OXgwYlVpYk5EZUgxYlYxUWRGTmlUcmsxaFhIa1FKZzE0ZzEzYlh0aQpEWnhGMnMxeXJSa3oxT2RBUndKa3ZPL3d0NHREdCtnVHVmeVp2NXVFN0h0NGtiM0ZJQzl1WnQ0d25pVkhmWnFICkZwVmhDMmRCK0Mrbzg5b21tQzFNOTNDV0VDbWV4QkZidXp0UlpmQUlSaE1jd2pET3RqZkpmNGNDZ1lFQXdTaWkKdVRGRUdqLzlaTHhtNkN5Q2taR0tYanQ5dGlsTjZ3Mk9WZmxvSE9NY3pVb1F6M21USUZaejZZZHZZM2VxcVM5MQpKL3REL2JyNkNUMnpKMnplc1h0Y3ljbVpOZWc0T0J1TlplUElFWGR2alRoK0pxYXc4eHdyRlhzNGZMZEUwbU1oClcwK2M4RzFzcEZjOGowUk53QmJ6Z1JUSTR4S1QxcTJLajZIMU01MENnWUFyeWxMdGVQcFpRK3NUQ2l1WVJYclUKY2R5c1VVVUlNRWlDMTVhVGNvUTFBbnpiT1J2OFdBVk4rVldjNmhGV0Z0Nzg4Q1ZtTEhWa0pIN0doSzhEUnltegpOOTFKWm5OSHZSNXpSWEdiSy9WM2lSY0V6VCs5ZEl3SllkWlFGbUtYU2dVb0o3WGd1a0hySkNrZTJVWExCRFViCmJMcmRjNm8zZDJNMVJlSnBuSlpsd1FLQmdRQ0NHb3ZZYjQydW5MRmgyK0Q0dTVwSzBKeEJzcEtQVXl6dmlSYjUKWTkrenJXb21BS0JvRHp5QlNKb1VqeXdBOUlhWUpLWW1Bd0dkOHdZZG1WaUYwcmdCRmRXKytUSmdkQVVDRGRUawo3MU5BS0pHVHJweVNEaThiNFRwSDR0SitkcmM5ZXBYcU9pcThhd2dGZmRrRnF2MHZ5SVhGeVNreWdiM2dtTTIrCngxa3dwUUtCZ0hQeFBjMjhTaCs5TERLbUEzOHBMck0zZHljd1o1WDFwUUNlWFZ4TW0vZGY2dXBZZ25IQ1dKYWsKekNNZU0wcVR1UG9wM2VDOFdEVkFnamxhT1N3VmFFcThJTmhxRmE3M2tINVcxdjVueFk1aHVMN2Q2bWlBWmdYawp3WXZlaUE4L3V4b0N5NmpnU0lPR1ZwMVZpUnh4bEM3ZVZxNHN1ckJndnZGV0ROS25SY1dBCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==
  kubectl exec -it "$VAULT_POD_NAME" -- vault kv get secret/test-foo
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