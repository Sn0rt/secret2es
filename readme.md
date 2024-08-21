# sercert2es

A tool to convert Kubernetes core v1 Secrets to External Secrets.

## Usage

```
sercert2es extsecret-gen \
  -i, --input <corev1-secret-file> \
  -o, --output <external-secret-file> \
  -s, --store <store-type> \
  -n, --storename <store-name> \
  [--namespace <external-namespace>] \
  [--secret-name <secret-name>] \
  [--verbose]
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

- `sercert2extsecret help`: Display help information about the tool.
- `sercert2extsecret version`: Show the current version of the tool.

## Building

To build the tool with version information:

```shell
make build
```

## License

BSD-3
