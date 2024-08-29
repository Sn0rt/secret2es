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

func TestGenerateOpaqueSecret(t *testing.T) {
	tests := []struct {
		name                 string
		inputSecret          internalSecret
		expectExternalSecret esv1beta1.ExternalSecret
		store                esv1beta1.SecretStoreRef
		envs                 map[string]string // for render <% ENV %>
		err                  error
	}{
		{
			name: "empty opaque type secret",
			inputSecret: internalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name: "empty_data",
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
			err: fmt.Errorf(ErrCommonNotAcceptNeitherSecretDataAndData, "empty_data"),
		},
		{
			name: "simple opaque type secret",
			inputSecret: internalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name: "simple_example",
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
					Name:      "simple_example",
					Namespace: "",
					Labels: map[string]string{
						"app": "test",
					},
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/foo",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "simple_example",
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
			inputSecret: internalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name: "multiple_env_with_path",
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
					Name:      "multiple_env_with_path",
					Namespace: "",
					Labels: map[string]string{
						"app": "test",
					},
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/ubuntu-22.04-foo",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "multiple_env_with_path",
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
			name: "opaque type secret with path <% ENV %> and multiple property",
			inputSecret: internalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name: "multiple_property",
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
					Name:      "multiple_property",
					Namespace: "",
					Labels: map[string]string{
						"app": "test",
					},
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/ubuntu-22.04-foo",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "multiple_property",
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
			inputSecret: internalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name: "string_data_example",
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/<% DIST %>-<% VER %>-foo",
					},
					Labels: map[string]string{
						"app": "test",
					},
				},
				StringData: map[string]string{
					"mylogin.conf": `[client]
host = example.com
user = < USER >
password = <MYSQL_PASSWD>
port = 4000`,
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
					Name:      "string_data_example",
					Namespace: "",
					Labels: map[string]string{
						"app": "test",
					},
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/ubuntu-22.04-foo",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
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
						Name:           "string_data_example",
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
								"mylogin.conf": `[client]
host = example.com
user = "{{ .USER }}"
password = "{{ .MYSQL_PASSWD }}"
port = 4000`,
							},
						},
					},
				},
			},
		},
		{
			name: "opaque type secret with path <% ENV %> and multiple stringData",
			inputSecret: internalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name: "string_data_multiple_example",
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/<% DIST %>-<% VER %>-foo",
					},
					Labels: map[string]string{
						"app": "test",
					},
				},
				StringData: map[string]string{
					"sn0rt.github.io.default.access_key": "< USER_ACCESS_KEY >",
					"sn0rt.github.io.default.secret_key": "<USER_SECRET_KEY>",
					"sn0rt.github.io.default.cmt":        "sn0rt-<USER_SECRET_KEY>",
					"sn0rt.github.io.default.key":        "key", // merge policy should ignore this
				},
			},
			store: esv1beta1.SecretStoreRef{
				Kind: "ClusterSecretStore",
				Name: "tenant-b",
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
					Name:      "string_data_multiple_example",
					Namespace: "",
					Labels: map[string]string{
						"app": "test",
					},
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/ubuntu-22.04-foo",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "tenant-b",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "USER_ACCESS_KEY",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:      "ubuntu-22.04-foo",
								Property: "USER_ACCESS_KEY",
							},
						},
						{
							SecretKey: "USER_SECRET_KEY",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:      "ubuntu-22.04-foo",
								Property: "USER_SECRET_KEY",
							},
						},
					},
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "string_data_multiple_example",
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
								"sn0rt.github.io.default.access_key": `"{{ .USER_ACCESS_KEY }}"`,
								"sn0rt.github.io.default.secret_key": `"{{ .USER_SECRET_KEY }}"`,
								"sn0rt.github.io.default.cmt":        `"sn0rt-{{ .USER_SECRET_KEY }}"`,
							},
						},
					},
				},
			},
		},
		{
			name: "resolve the value from env",
			inputSecret: internalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name: "set_env_with_body_no_gen_es",
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/foo",
					},
					Labels: map[string]string{
						"app": "test",
					},
				},
				Data: map[string]string{
					"data1": "data1",
				},
			},
			store: esv1beta1.SecretStoreRef{
				Name: "test",
				Kind: "ClusterSecretStore",
			},
			err: fmt.Errorf(ErrCommonNotIncludeAngleBrackets, "set_env_with_body_no_gen_es"),
		},
		{
			name: "resolve the value from env case 1",
			inputSecret: internalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name: "set_env_with_body_no_gen_es_1",
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/foo",
					},
					Labels: map[string]string{
						"app": "test",
					},
				},
				Data: map[string]string{
					"data1": "data1",
					"data2": "<% DIST %>",
				},
			},
			store: esv1beta1.SecretStoreRef{
				Name: "test",
				Kind: "ClusterSecretStore",
			},
			err: fmt.Errorf(ErrCommonNotNeedRefData, "set_env_with_body_no_gen_es_1"),
		},
		{
			name: "resolve the value from env case 2",
			inputSecret: internalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name: "set_env_with_body_no_gen_es_2",
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/foo",
					},
					Labels: map[string]string{
						"app": "test",
					},
				},
				Data: map[string]string{
					"data1": "data1",
					"data2": "<% DIST %>",
					"data3": "<FROM_VAULT_DATA3>",
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
					Name:      "set_env_with_body_no_gen_es_2",
					Namespace: "",
					Labels: map[string]string{
						"app": "test",
					},
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/foo",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "set_env_with_body_no_gen_es_2",
						CreationPolicy: esv1beta1.CreatePolicyMerge,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "test",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "data3",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:      "foo",
								Property: "FROM_VAULT_DATA3",
							},
						},
					},
				},
			},
		},
		{
			name: "opaque type secret with <% ENV %> with stringData and multiple stringData",
			inputSecret: internalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name: "multiple_example_env_with_stringData",
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/foo",
					},
					Labels: map[string]string{
						"app": "test",
					},
				},
				StringData: map[string]string{
					"sn0rt.github.io.default.access_key": "<USER_ACCESS_KEY>",
					"sn0rt.github.io.default.secret_key": "<% USER_SECRET_KEY %>",
					"sn0rt.github.io.default.key":        "key", // merge policy should ignore this
				},
			},
			store: esv1beta1.SecretStoreRef{
				Kind: "ClusterSecretStore",
				Name: "tenant-b",
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
					Name:      "multiple_example_env_with_stringData",
					Namespace: "",
					Labels: map[string]string{
						"app": "test",
					},
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/foo",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "tenant-b",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "USER_ACCESS_KEY",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:      "foo",
								Property: "USER_ACCESS_KEY",
							},
						},
					},
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "multiple_example_env_with_stringData",
						CreationPolicy: esv1beta1.CreatePolicyMerge,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeOpaque,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Annotations: map[string]string{
									"avp.kubernetes.io/path": "secret/data/foo",
								},
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyMerge,
							Data: map[string]string{
								"sn0rt.github.io.default.access_key": `"{{ .USER_ACCESS_KEY }}"`,
							},
						},
					},
				},
			},
		},
		{
			name: "resolve <% ENV %> with stringData and multiple stringData empty ref",
			inputSecret: internalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name: "multiple_stringData_should_empty_ref",
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/foo",
					},
					Labels: map[string]string{
						"app": "test",
					},
				},
				StringData: map[string]string{
					"sn0rt.github.io.default.secret_key": "<% USER_SECRET_KEY %>",
					"sn0rt.github.io.default.key":        "key", // merge policy should ignore this
				},
			},
			store: esv1beta1.SecretStoreRef{
				Kind: "ClusterSecretStore",
				Name: "tenant-b",
			},
			envs: map[string]string{
				"DIST": "ubuntu",
				"VER":  "22.04",
			},
			err: fmt.Errorf(ErrCommonNotNeedRefData, "multiple_stringData_should_empty_ref"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envs {
				_ = os.Setenv(k, v)
			}
			externalSecret, err := convertSecret2ExtSecret(tt.inputSecret, tt.store.Kind, tt.store.Name)
			if err != nil {
				if tt.err.Error() != err.Error() {
					t.Errorf("Err Mismatch (+goot: %s)\n", err)
					t.Errorf("Err Mismatch (+want: %s)\n", tt.err)
				}
			} else {
				diff := cmp.Diff(externalSecret, &tt.expectExternalSecret, cmpopts.SortSlices(
					func(a, b esv1beta1.ExternalSecretData) bool {
						return a.SecretKey > b.SecretKey
					}))
				if diff != "" {
					t.Errorf("%s case Mismatch (-want +got):\n%s", tt.name, diff)
				}
			}
		})
	}
}
