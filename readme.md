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

2. Because the Rewriting Keys in DataFrom feature provided by External secret is too complex, and the current implementation is to simply use Data instead of DataFrom. This is simple to implement. But the disadvantage is that there will be an extra k-v pair in the secret.

## License

BSD-3
