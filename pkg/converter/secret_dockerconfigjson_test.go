package converter

import (
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestGenEsByDockerConfigJSON(t *testing.T) {
	var tests = []struct {
		name                 string
		inputSecret          internalSecret
		store                esv1beta1.SecretStoreRef
		expectExternalSecret esv1beta1.ExternalSecret
	}{
		{
			name: "basic",
			store: esv1beta1.SecretStoreRef{
				Name: "tenant-b",
				Kind: "ClusterSecretStore",
			},
			inputSecret: internalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeDockerConfigJson,
				ObjectMeta: metav1.ObjectMeta{
					Name: "input1",
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/test-foo",
					},
					Labels: map[string]string{
						"app": "test",
					},
				},
				StringData: map[string]string{
					".dockerconfigjson": `{
  "auths": {
    "https://index.docker.io/v1": {
      "auth": "<PASSWD_FROM_VAULT>"
    },
    "https://index.docker.io:8443/v1": {
      "auth": "<PASSWD_FROM_VAULT>"
    }
  }
}`,
				},
			},
			expectExternalSecret: esv1beta1.ExternalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "external-secrets.io/v1beta1",
					Kind:       "ExternalSecret",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "input1",
					Namespace: "",
					Labels: map[string]string{
						"app": "test",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "input1",
						CreationPolicy: esv1beta1.CreatePolicyOrphan,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeDockerConfigJson,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyReplace,
							Data: map[string]string{
								".dockerconfigjson": `{
  "auths": {
    "https://index.docker.io/v1": {
      "auth": "{{ .PASSWD_FROM_VAULT }}"
    },
    "https://index.docker.io:8443/v1": {
      "auth": "{{ .PASSWD_FROM_VAULT }}"
    }
  }
}`,
							},
						},
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "tenant-b",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "PASSWD_FROM_VAULT",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "test-foo",
								MetadataPolicy:     "None",
								Property:           "PASSWD_FROM_VAULT",
								ConversionStrategy: "Default",
								DecodingStrategy:   "None",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := convertSecret2ExtSecret(tt.inputSecret, tt.store.Kind, tt.store.Name, esv1beta1.CreatePolicyOrphan, true)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			} else {
				if diff := cmp.Diff(out, &tt.expectExternalSecret); diff != "" {
					t.Errorf("Mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestSerializeDockerConfigJSON(t *testing.T) {
	var tests = []struct {
		name     string
		input    []byte
		expected Auths
	}{
		{
			name: "basic",
			input: []byte(`{
      "auths": {
        "https://index.docker.io/v1": {
          "auth": "<PASSWD_FROM_VAULT>"
        },
        "https://index.docker.io:8443/v1": {
          "auth": "<PASSWD_FROM_VAULT>"
        }      
      }
    }`),
			expected: Auths{
				Auths: map[string]Auth{
					"https://index.docker.io/v1": {
						"<PASSWD_FROM_VAULT>",
					},
					"https://index.docker.io:8443/v1": {
						"<PASSWD_FROM_VAULT>",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := serializeDockerConfigJSON(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			} else {
				if diff := cmp.Diff(out, &tt.expected); diff != "" {
					t.Errorf("Mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
