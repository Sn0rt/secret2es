# secret2es

A tool to convert Kubernetes core v1 Secrets to External Secrets.

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
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  annotations:
    avp.kubernetes.io/path: secret/data/foo
  creationTimestamp: null
  name: simple-2
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
    name: simple-2
status:
  binding: {}
  refreshTime: null
---
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  annotations:
    avp.kubernetes.io/path: secret/data/foo-test
  creationTimestamp: null
  name: include-temp-value
spec:
  data:
  - remoteRef:
      key: foo-test
      property: dist-name-of-linux
    secretKey: dist
  secretStoreRef:
    kind: ClusterSecretStore
    name: test
  target:
    creationPolicy: Merge
    deletionPolicy: Retain
    name: include-temp-value
status:
  binding: {}
  refreshTime: null
```

### Options

- `-i, --input <corev1-secret-file>`: Required. Path to the input core v1 Secret file. Must include special `argocd-vault-plugin` annotations.
- `-o, --output <external-secret-file>`: Required. Path to the output External Secret file.
- `-s, --store <store-type>`: Optional. Type of secret store. Default is "ClusterSecretStore".
- `-n, --storename <store-name>`: Required. Name of the secret store.
- `--namespace <external-namespace>`: Optional. Namespace for the External Secret. Default is "default" if not specified in the input Secret.
- `--secret-name <secret-name>`: Optional. Name for the External Secret. Default is the name of the input Secret.
- `--verbose`: Optional. Enable verbose output for more detailed process information.

### Additional Commands

- `secret2es help`: Display help information about the tool.
- `secret2es version`: Show the current version of the tool.

## Building

To build the tool with version information:

```shell
make build
```

## License

BSD-3
