package converter

import (
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"testing"
)

func TestGenerateOpaqueSecret(t *testing.T) {
	tests := []struct {
		name                 string
		inputSecret          UnstructuredSecret
		expectExternalSecret esv1beta1.ExternalSecret
		store                esv1beta1.SecretStoreRef
		envs                 map[string]string // for render <% ENV %>
	}{
		{
			name: "simple opaque type secret",
			inputSecret: UnstructuredSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "input1",
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/foo",
					},
					Labels: map[string]string{
						"app": "test",
					},
				},
				Data: map[string]string{
					"dist": "<dist-name-of-linux>",
				},
			},
			store: esv1beta1.SecretStoreRef{
				Name: "test",
				Kind: "ClusterSecretStore",
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
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/foo",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "input1",
						CreationPolicy: esv1beta1.CreatePolicyMerge,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "test",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "dist",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:      "foo",
								Property: "dist-name-of-linux",
							},
						},
					},
				},
			},
		},
		{
			name: "opaque type secret with path <% ENV %>",
			inputSecret: UnstructuredSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "input1",
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/<% DIST %>-<% VER %>-foo",
					},
					Labels: map[string]string{
						"app": "test",
					},
				},
				Data: map[string]string{
					"dist": "<dist-name-of-linux>",
				},
			},
			store: esv1beta1.SecretStoreRef{
				Name: "test",
				Kind: "ClusterSecretStore",
			},
			envs: map[string]string{
				"DIST": "ubuntu",
				"VER":  "22.04",
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
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/ubuntu-22.04-foo",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "input1",
						CreationPolicy: esv1beta1.CreatePolicyMerge,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "test",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "dist",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:      "ubuntu-22.04-foo",
								Property: "dist-name-of-linux",
							},
						},
					},
				},
			},
		},
		{
			name: "opaque type secret with path <% ENV %> and mutilate property",
			inputSecret: UnstructuredSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "input1",
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/<% DIST %>-<% VER %>-foo",
					},
					Labels: map[string]string{
						"app": "test",
					},
				},
				Data: map[string]string{
					"dist":   "<dist-name-of-linux>",
					"user":   "<github-username>",
					"passwd": "<github-passwd>",
				},
			},
			store: esv1beta1.SecretStoreRef{
				Name: "test",
				Kind: "ClusterSecretStore",
			},
			envs: map[string]string{
				"DIST": "ubuntu",
				"VER":  "22.04",
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
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/ubuntu-22.04-foo",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "input1",
						CreationPolicy: esv1beta1.CreatePolicyMerge,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "test",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "dist",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:      "ubuntu-22.04-foo",
								Property: "dist-name-of-linux",
							},
						},
						{
							SecretKey: "user",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:      "ubuntu-22.04-foo",
								Property: "github-username",
							},
						},
						{
							SecretKey: "passwd",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:      "ubuntu-22.04-foo",
								Property: "github-passwd",
							},
						},
					},
				},
			},
		},
		{
			name: "opaque type secret with path <% ENV %> and stringData",
			inputSecret: UnstructuredSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "string-data-example",
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/<% DIST %>-<% VER %>-foo",
					},
					Labels: map[string]string{
						"app": "test",
					},
				},
				StringData: map[string]string{
					"mylogin.conf": "    [client]\n    host = example.com\n    user = < USER >\n    password = <MYSQL_PASSWD>\n    port = 4000",
				},
			},
			store: esv1beta1.SecretStoreRef{
				Kind: "ClusterSecretStore",
				Name: "tenant-a",
			},
			envs: map[string]string{
				"DIST": "ubuntu",
				"VER":  "22.04",
			},
			expectExternalSecret: esv1beta1.ExternalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "external-secrets.io/v1beta1",
					Kind:       "ExternalSecret",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "string-data-example",
					Namespace: "",
					Labels: map[string]string{
						"app": "test",
					},
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/ubuntu-22.04-foo",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "tenant-a",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "USER",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:      "ubuntu-22.04-foo",
								Property: "USER",
							},
						},
						{
							SecretKey: "MYSQL_PASSWD",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:      "ubuntu-22.04-foo",
								Property: "MYSQL_PASSWD",
							},
						},
					},
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "string-data-example",
						CreationPolicy: esv1beta1.CreatePolicyMerge,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeOpaque,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Annotations: map[string]string{
									"avp.kubernetes.io/path": "secret/data/ubuntu-22.04-foo",
								},
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyMerge,
							Data: map[string]string{
								"mylogin.conf": "    [client]\n    host = example.com\n    user = \"{{ .USER }}\"\n    password = \"{{ .MYSQL_PASSWD }}\"\n    port = 4000",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envs {
				_ = os.Setenv(k, v)
			}
			externalSecret, err := generateOpaqueSecret(tt.inputSecret, tt.store.Kind, tt.store.Name)
			if err != nil {
				t.Errorf("generateOpaqueSecret() returned an unexpected error: got: %v", err)
			}
			externalSecret.Status = esv1beta1.ExternalSecretStatus{}
			if diff := cmp.Diff(externalSecret, &tt.expectExternalSecret); diff != "" {
				t.Errorf("Mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
