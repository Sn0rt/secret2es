package converter

import (
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
kind: Secret
apiVersion: v1
metadata:
  name: input1
  annotations:
    avp.kubernetes.io/path: "secret/data/foo"
type: Opaque
data:
  dist: <dist-name-of-linux>
---
kind: Secret
apiVersion: v1
metadata:
  name: input2
  annotations:
    avp.kubernetes.io/path: "secret/data/foo"
type: Opaque
data:
  dist: <dist-name-of-linux>
`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := parseUnstructuredSecret(tt.body)
			if err != nil {
				t.Errorf("parseUnstructuredSecret() returned an unexpected error: got: %v", err)
			}
			for _, v := range out {
				t.Logf("process %s/%s\n", v.Namespace, v.Name)
				externalSecret, err := convertSecret2ExtSecret(v, "vault", "test", "default", "test")
				if err != nil {
					t.Errorf("convertSecret2ExtSecret() returned an unexpected error: got: %v", err)
				}

				yamlData, err := yaml.Marshal(externalSecret)
				if err != nil {
					t.Errorf("yaml.Marshal() returned an unexpected error: got: %v", err)
				}

				t.Logf("%s\n", yamlData)
			}
		})
	}
}
