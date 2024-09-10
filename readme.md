# secret2es

This tool allows administrators to migrate secrets originally managed by [argocd-vault-plugin](https://argocd-vault-plugin.readthedocs.io/en/stable/) to [external-secrets](https://github.com/external-secrets/external-secrets) for continued maintenance.

## Usage

```shell
./secret2es --help
A tool to convert AVP secrets to ExternalSecrets

Usage:
  secret2es [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  es-gen      Generate external secrets from corev1 secrets
  help        Help about any command
  version     Print the version number of secret2es

Flags:
  -h, --help   help for secret2es

Use "secret2es [command] --help" for more information about a command.
```

```shell
/secret2es es-gen --help
Generate external secrets from corev1 secrets

Usage:
  secret2es es-gen [flags]

Flags:
  -c, --creation-policy string   Create policy (default: Orphan), only Owner, Orphan (default "Orphan")
  -h, --help                     help for es-gen
  -i, --input string             Input path of corev1 secret file (required)
  -r, --resolve                  Resolve the <% ENV %> from env
  -n, --storename string         Store name (required)
  -s, --storetype string         Store type (optional) (default "SecretStore")
```

example 

```shell
./secret2es es-gen -i e2e/templated.yaml -s ClusterSecretStore -n tenant-b -r true
---
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: input1
...
```

## Building

To build the tool with version information:

```shell
make build
```

## known issues

1. the `label` and `annotation` of the secret has not been created if it has been set with `ExternalSecret` CRD.

## License

BSD-3
