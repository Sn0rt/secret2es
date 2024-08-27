# secret2es

This tool allows administrators to migrate secrets originally managed by [argocd-vault-plugin](https://argocd-vault-plugin.readthedocs.io/en/stable/) to [external-secrets](https://github.com/external-secrets/external-secrets) for continued maintenance.

## Usage

```shell
secret2es es-gen \
  -i, --input <corev1-secret-file> \
  -n --storename <store-name> \
```

example 

```shell
➜  secret2es git:(main) ✗ ./secret2es es-gen -i test/opaque-secret.yaml -n test
---
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  annotations:
    avp.kubernetes.io/path: secret/data/foo
  creationTimestamp: null
  name: simple
spec:
  data:
  - remoteRef:
      key: foo
      property: dist-name-of-linux
    secretKey: dist
  secretStoreRef:
    kind: ClusterSecretStore
    name: test
  target:
    creationPolicy: Merge
    deletionPolicy: Retain
    name: simple
status:
  binding: {}
  refreshTime: null
---
...
```

### `es-gen` subCommand Options

- `-i, --input <corev1-secret-file>`: Required. Path to the input core v1 Secret file. Must include special `argocd-vault-plugin` annotations.
- `-s, --storetype <store-type>`: Optional. Type of secret store. Default is "ClusterSecretStore".
- `-n, --storename <store-name>`: Required. Name of the secret store.

### Additional Commands

- `secret2es help`: Display help information about the tool.
- `secret2es version`: Show the current version of the tool.

## Building

To build the tool with version information:

```shell
make build
```

## known issues

1. the `label` and `annotation` of the secret has not been created if it has been set with `ExternalSecret` CRD.

2. for `kubernetes.io/basic-auth` secret only

In fact, we do not expect additional key value pairs to appear in secret. Although it does not affect the work, it is not a good experience. The specific reason needs to be further investigated, which may be a problem with the controller implementation.

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  annotations:
    avp.kubernetes.io/path: secret/data/test-foo
  name: simple
  namespace: default
spec:
  data:
  - remoteRef:
      conversionStrategy: Default
      decodingStrategy: None
      key: test-foo
      metadataPolicy: None
      property: USER_SECRET_KEY
    secretKey: USER_SECRET_KEY
  - remoteRef:
      conversionStrategy: Default
      decodingStrategy: None
      key: test-foo
      metadataPolicy: None
      property: USER_ACCESS_KEY
    secretKey: USER_ACCESS_KEY
  refreshInterval: 1h
  secretStoreRef:
    kind: ClusterSecretStore
    name: tenant-b
  target:
    creationPolicy: Owner
    deletionPolicy: Retain
    name: simple
    template:
      data:
        host: localhost.local
        password: '"sn0rt_{{ .USER_SECRET_KEY }}"'
        username: '"{{ .USER_ACCESS_KEY }}"'
      engineVersion: v2
      mergePolicy: Merge
      metadata:
        annotations:
          avp.kubernetes.io/path: secret/data/test-foo
      type: kubernetes.io/basic-auth
```

```yaml
apiVersion: v1
data:
  USER_ACCESS_KEY: YWNjZXNza2V5X3NuMHJ0
  USER_SECRET_KEY: c2VjcmV0a2V5X3NuMHJ0
  host: bG9jYWxob3N0LmxvY2Fs
  password: InNuMHJ0X3NlY3JldGtleV9zbjBydCI=
  username: ImFjY2Vzc2tleV9zbjBydCI=
immutable: false
kind: Secret
metadata:
  annotations:
    avp.kubernetes.io/path: secret/data/test-foo
    reconcile.external-secrets.io/data-hash: 867bd7cd4b34f8b857174ece5e1f2186
  labels:
    reconcile.external-secrets.io/created-by: 178edf16c659c2998642565fbe634a1a
  name: simple
  namespace: default
  ownerReferences:
  - apiVersion: external-secrets.io/v1beta1
    blockOwnerDeletion: true
    controller: true
    kind: ExternalSecret
    name: simple
    uid: 2c6b65e7-3c13-4039-9dda-845f0a64b23f
type: kubernetes.io/basic-auth
```

## License

BSD-3
