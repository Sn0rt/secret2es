package converter

import (
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	"os"
	"sigs.k8s.io/yaml"
	"testing"
)

func TestUnstructuredSecret(t *testing.T) {
	tests := []struct {
		name string
		body []byte
	}{
		{
			name: "opaque type secret",
			body: []byte(`
kind: Secret
apiVersion: v1
metadata:
  name: input1
  annotations:
    avp.kubernetes.io/path: "secret/data/foo"
type: Opaque
data:
  dist: <dist-name-of-linux>
`),
		},
		{
			name: "opaque type secret list",
			body: []byte(`
---
---
apiVersion: v1
kind: Secret
metadata:
  name: simple
  annotations:
    avp.kubernetes.io/path: "secret/data/foo"
type: Opaque
data:
  dist: <dist-name-of-linux>
---
kind: Secret
apiVersion: v1
metadata:
  name: simple-2
  annotations:
    avp.kubernetes.io/path: "secret/data/foo"
type: Opaque
data:
  dist: <dist-name-of-linux>
---
apiVersion: v1
kind: Secret
metadata:
  name: include-temp-value
  annotations:
    avp.kubernetes.io/path: "secret/data/foo-<% ENV %>"
type: Opaque
data:
  dist: <dist-name-of-linux>
  user: <user-name-of-github>
  passwd: <token>
---
apiVersion: v1
kind: ConfigMap
metadata:
  annotations:
    kubernetes.io/description: Contains a CA bundle that can be used to verify the
      kube-apiserver when using internal endpoints such as the internal service IP
      or kubernetes.default.svc. No other usage is guaranteed across distributions
      of Kubernetes clusters.
  name: should-ignore-the-cm
  namespace: default
data:
  ca.crt: |
    -----BEGIN CERTIFICATE-----
    MIIDBjCCAe6gAwIBAgIBATANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDEwptaW5p
    a3ViZUNBMB4XDTI0MDcwMTAzMzEwN1oXDTM0MDYzMDAzMzEwN1owFTETMBEGA1UE
    AxMKbWluaWt1YmVDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAM8r
    DMCiFKKLpS9GfGGUZMSpTNOc3be5RPVYN28J4N6JUsx0qgCpyd9eSdn5kv8t+Kku
    r16ilVF6w06cWy1KVStEoKL4i2XH9ZjbIYEmqR3UuIxVJgl1rjw7I85tT/yt6oNR
    oup4GTwwqlXgxQb39TkOwu89Jcb/rsjONsfUaPaXdv4vLxRMswLQOtdi1c6U1ZQ7
    IL1yDkdT7mwuetBnQhP+bY+UFgsaRhzAJK4Rqx7pLdJq8oi8B/XkoiDwThpoMqnG
    yQqGwAtG5xCBnJpzxs0H3YyYLJ1yNnfVkDiuBQz/ByS1dBJsdNmt2jeTgz/Vabty
    yzRTeV0Gl5fh65h1LG8CAwEAAaNhMF8wDgYDVR0PAQH/BAQDAgKkMB0GA1UdJQQW
    MBQGCCsGAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQW
    BBTkA7/xXAGTPzgRg7DOUZAsmsm0XzANBgkqhkiG9w0BAQsFAAOCAQEADmDFzwJJ
    OEZOKH/JVtvoDZ/7bMJVwN/H2wqP5HEPDeYDwPL1xd3nxxisnoxE9sxyAFE3KWKy
    2qswmBsNfc4+JHBAESOrckCftbIVdANGUSEEkLGgXuZV8xxS9F76D8cKCnprb/po
    Nvf7/UCSeNDVOEfFlmRNA2o36mvQFfWUttP7ULCJE2RXoe3pOJFnwwc9N6BGaUxz
    z6HFfMFZj6QCTEVMdAfpoJqmm44LEftGW3t8BEskieirjW9AGGOLKHFAeDkXMcPc
    9/+N0zKtCHBFOG+0aZQyYwxb3vxVGkdPpkcfGv59a30/UvEOzQeAbvYzOmYjY4xa
    3gRBomAS54ZYvw==
    -----END CERTIFICATE-----
---
apiVersion: v1
kind: Secret
metadata:
  name: simple
  annotations:
    avp.kubernetes.io/path: "secret/data/<% ENV %>-foo"
type: Opaque
stringData:
  mylogin.conf: |
    [client]
    host = example.com
    user = < USER >
    password = <MYSQL_PASSWD>
    port = 4000
---
apiVersion: v1
kind: Secret
metadata:
  name: simple
  annotations:
    avp.kubernetes.io/path: "secret/data/<% ENV %>-foo"
type: Opaque
stringData:
  sn0rt.github.io.default.access_key: < USER_ACCESS_KEY >
  sn0rt.github.io.default.secret_key: <USER_SECRET_KEY>
`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("ENV", "test")
			out, err := parseUnstructuredSecret(tt.body)
			if err != nil {
				t.Errorf("parseUnstructuredSecret() returned an unexpected error: got: %v", err)
			}
			for _, v := range out {
				externalSecret, err := convertSecret2ExtSecret(v, ClusterSecretStoreType, "test", esv1beta1.CreatePolicyOrphan, true)
				if err != nil {
					t.Errorf("convertSecret2ExtSecret() returned an unexpected error: got: %v", err)
				}
				_, err = yaml.Marshal(externalSecret)
				if err != nil {
					t.Errorf("yaml.Marshal() returned an unexpected error: got: %v", err)
				}
			}
		})
	}
}
