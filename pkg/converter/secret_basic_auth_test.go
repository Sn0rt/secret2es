package converter

import (
	"fmt"
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"testing"
)

func TestGenerateBasicAuthSecret(t *testing.T) {
	tests := []struct {
		name                 string
		inputSecret          UnstructuredSecret
		expectExternalSecret esv1beta1.ExternalSecret
		store                esv1beta1.SecretStoreRef
		envs                 map[string]string // for render <% ENV %>
		err                  error
	}{
		{
			name: "empty opaque type secret",
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
			},
			store: esv1beta1.SecretStoreRef{
				Name: "test",
				Kind: "ClusterSecretStore",
			},
			err: fmt.Errorf(ErrCommonNotAcceptNeitherSecretDataAndData, "input1"),
		},
		{
			name: "a simple case",
			inputSecret: UnstructuredSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
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
					"username": `<USER_ACCESS_KEY>`,
					"password": `sn0rt_<USER_SECRET_KEY>`,
					"host":     `localhost.local`,
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
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/test-foo",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "input1",
						CreationPolicy: esv1beta1.CreatePolicyMerge,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeBasicAuth,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Annotations: map[string]string{
									"avp.kubernetes.io/path": "secret/data/test-foo",
								},
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyMerge,
							Data: map[string]string{
								"username": `"{{ .USER_ACCESS_KEY }}"`,
								"password": `"sn0rt_{{ .USER_SECRET_KEY }}"`,
								"host":     `localhost.local`,
							},
						},
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "tenant-b",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "USER_ACCESS_KEY",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:      "test-foo",
								Property: "USER_ACCESS_KEY",
							},
						},
						{
							SecretKey: "USER_SECRET_KEY",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:      "test-foo",
								Property: "USER_SECRET_KEY",
							},
						},
					},
				},
			},
			store: esv1beta1.SecretStoreRef{
				Name: "tenant-b",
				Kind: "ClusterSecretStore",
			},
			envs: map[string]string{
				"ENV": "test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envs {
				_ = os.Setenv(k, v)
			}
			externalSecret, err := convertSecret2ExtSecret(tt.inputSecret, tt.store.Kind, tt.store.Name)
			if err != nil {
				if err == tt.err {
					t.Errorf("generateEsByOpaqueSecret() returned an unexpected error: got: %v, want: %v", err, tt.err)
				}
			} else {
				externalSecret.Status = esv1beta1.ExternalSecretStatus{}
				if diff := cmp.Diff(externalSecret, &tt.expectExternalSecret, cmpopts.SortSlices(func(a, b esv1beta1.ExternalSecretData) bool {
					return a.SecretKey > b.SecretKey
				})); diff != "" {
					t.Errorf("Mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
