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

func TestGenerateStringDataOpaqueSecret(t *testing.T) {
	tests := []struct {
		name                 string
		inputSecret          internalSecret
		expectExternalSecret esv1beta1.ExternalSecret
		store                esv1beta1.SecretStoreRef
		envs                 map[string]string // for render <% ENV %>
		err                  error
		enableResolve        bool
	}{
		{
			name: "test single value is <% ENV %>",
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
				StringData: map[string]string{
					"dist": "<dist-name-of-linux>",
					"env1": "<% ENV %>",
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
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "simple_example",
						CreationPolicy: esv1beta1.CreatePolicyOrphan,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeOpaque,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyReplace,
							Data: map[string]string{
								"dist": `"{{ .dist-name-of-linux }}"`,
								"env1": "<% ENV %>",
							},
						},
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "test",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "dist-name-of-linux",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "foo",
								MetadataPolicy:     "None",
								Property:           "dist-name-of-linux",
								ConversionStrategy: "Default",
								DecodingStrategy:   "None",
							},
						},
					},
				},
			},
		},
		{
			name: "test <% ENV %>-<VAULT>",
			inputSecret: internalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name: "mix_env_value_with_vault",
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/foo",
					},
					Labels: map[string]string{
						"app": "test",
					},
				},
				StringData: map[string]string{
					"dist": "<dist-name-of-linux>",
					"env1": "<% ENV %>-<dist-name-of-linux>",
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
					Name:      "mix_env_value_with_vault",
					Namespace: "",
					Labels: map[string]string{
						"app": "test",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "mix_env_value_with_vault",
						CreationPolicy: esv1beta1.CreatePolicyOrphan,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeOpaque,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyReplace,
							Data: map[string]string{
								"dist": `"{{ .dist-name-of-linux }}"`,
								"env1": `"<% ENV %>-{{ .dist-name-of-linux }}"`,
							},
						},
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "test",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "dist-name-of-linux",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "foo",
								MetadataPolicy:     "None",
								Property:           "dist-name-of-linux",
								ConversionStrategy: "Default",
								DecodingStrategy:   "None",
							},
						},
					},
				},
			},
		},
		{
			name: "test <<% ENV %>_VAULT>",
			inputSecret: internalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name: "mix_two_style",
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/foo",
					},
					Labels: map[string]string{
						"app": "test",
					},
				},
				StringData: map[string]string{
					"env0": "<VAULT0>",
					"env1": "<<% ENV1 %>_VAULT1>",
					"env2": "<<% ENV2 %>_VAULT2>",
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
					Name:      "mix_two_style",
					Namespace: "",
					Labels: map[string]string{
						"app": "test",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "mix_two_style",
						CreationPolicy: esv1beta1.CreatePolicyOrphan,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeOpaque,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyReplace,
							Data: map[string]string{
								"env0": `"{{ .VAULT0 }}"`,
								"env1": `"{{ .<% ENV1 %>_VAULT1 }}"`,
								"env2": `"{{ .<% ENV2 %>_VAULT2 }}"`,
							},
						},
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "test",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "VAULT0",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "foo",
								MetadataPolicy:     "None",
								Property:           "VAULT0",
								ConversionStrategy: "Default",
								DecodingStrategy:   "None",
							},
						},
						{
							SecretKey: "<% ENV1 %>_VAULT1",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "foo",
								MetadataPolicy:     "None",
								Property:           "<% ENV1 %>_VAULT1",
								ConversionStrategy: "Default",
								DecodingStrategy:   "None",
							},
						},
						{
							SecretKey: "<% ENV2 %>_VAULT2",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "foo",
								MetadataPolicy:     "None",
								Property:           "<% ENV2 %>_VAULT2",
								ConversionStrategy: "Default",
								DecodingStrategy:   "None",
							},
						},
					},
				},
			},
		},
		{
			name:          "opaque type secret with path <% ENV %> and stringData",
			enableResolve: true,
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
								Key:                "ubuntu-22.04-foo",
								MetadataPolicy:     "None",
								Property:           "USER",
								ConversionStrategy: "Default",
								DecodingStrategy:   "None",
							},
						},
						{
							SecretKey: "MYSQL_PASSWD",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "ubuntu-22.04-foo",
								MetadataPolicy:     "None",
								Property:           "MYSQL_PASSWD",
								ConversionStrategy: "Default",
								DecodingStrategy:   "None",
							},
						},
					},
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "string_data_example",
						CreationPolicy: esv1beta1.CreatePolicyOrphan,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeOpaque,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyReplace,
							Data: map[string]string{
								"mylogin.conf": `[client]
host = example.com
user = {{ .USER }}
password = {{ .MYSQL_PASSWD }}
port = 4000`,
							},
						},
					},
				},
			},
		},
		{
			name:          "opaque type secret with path <% ENV %> and multiple stringData",
			enableResolve: true,
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
					"sn0rt.github.io.default.key":        "key",
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
								Key:                "ubuntu-22.04-foo",
								MetadataPolicy:     "None",
								Property:           "USER_ACCESS_KEY",
								ConversionStrategy: "Default",
								DecodingStrategy:   "None",
							},
						},
						{
							SecretKey: "USER_SECRET_KEY",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "ubuntu-22.04-foo",
								MetadataPolicy:     "None",
								Property:           "USER_SECRET_KEY",
								ConversionStrategy: "Default",
								DecodingStrategy:   "None",
							},
						},
					},
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "string_data_multiple_example",
						CreationPolicy: esv1beta1.CreatePolicyOrphan,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeOpaque,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyReplace,
							Data: map[string]string{
								"sn0rt.github.io.default.access_key": `"{{ .USER_ACCESS_KEY }}"`,
								"sn0rt.github.io.default.secret_key": `"{{ .USER_SECRET_KEY }}"`,
								"sn0rt.github.io.default.cmt":        `"sn0rt-{{ .USER_SECRET_KEY }}"`,
								"sn0rt.github.io.default.key":        "key",
							},
						},
					},
				},
			},
		},
		{
			name:          "stringData type secret with <% ENV %> and stringData and multiple",
			enableResolve: true,
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
					"sn0rt.github.io.default.key":        "key",
				},
			},
			store: esv1beta1.SecretStoreRef{
				Kind: "ClusterSecretStore",
				Name: "tenant-b",
			},
			envs: map[string]string{
				"DIST":            "ubuntu",
				"VER":             "22.04",
				"USER_SECRET_KEY": "secret_key",
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
								Key:                "foo",
								MetadataPolicy:     "None",
								Property:           "USER_ACCESS_KEY",
								ConversionStrategy: "Default",
								DecodingStrategy:   "None",
							},
						},
					},
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "multiple_example_env_with_stringData",
						CreationPolicy: esv1beta1.CreatePolicyOrphan,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeOpaque,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyReplace,
							Data: map[string]string{
								"sn0rt.github.io.default.access_key": `"{{ .USER_ACCESS_KEY }}"`,
								"sn0rt.github.io.default.key":        "key",
								"sn0rt.github.io.default.secret_key": "secret_key",
							},
						},
					},
				},
			},
		},
		{
			name:          "resolve <% ENV %> with stringData and multiple stringData empty ref",
			enableResolve: true,
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
		{
			name:          "both <% ENV %> and < KEY > from stringData",
			enableResolve: true,
			inputSecret: internalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name: "both_env_and_key_stringData",
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/foo",
					},
					Labels: map[string]string{
						"app": "test",
					},
				},
				StringData: map[string]string{
					"config.yaml": `type: S3
prefix: "test/<% DIST %>"
config:
  endpoint: "https://s3.amazonaws.com"
  access_key: <S3_ACCESS_KEY>
  secret_key: <S3_SECRET_KEY>`,
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
					Name:      "both_env_and_key_stringData",
					Namespace: "",
					Labels: map[string]string{
						"app": "test",
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
							SecretKey: "S3_ACCESS_KEY",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "foo",
								MetadataPolicy:     "None",
								Property:           "S3_ACCESS_KEY",
								ConversionStrategy: "Default",
								DecodingStrategy:   "None",
							},
						},
						{
							SecretKey: "S3_SECRET_KEY",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "foo",
								MetadataPolicy:     "None",
								Property:           "S3_SECRET_KEY",
								ConversionStrategy: "Default",
								DecodingStrategy:   "None",
							},
						},
					},
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "both_env_and_key_stringData",
						CreationPolicy: esv1beta1.CreatePolicyOrphan,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeOpaque,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyReplace,
							Data: map[string]string{
								"config.yaml": `type: S3
prefix: "test/ubuntu"
config:
  endpoint: "https://s3.amazonaws.com"
  access_key: {{ .S3_ACCESS_KEY }}
  secret_key: {{ .S3_SECRET_KEY }}`,
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
			externalSecret, err := convertSecret2ExtSecret(tt.inputSecret, tt.store.Kind, tt.store.Name, esv1beta1.CreatePolicyOrphan, tt.enableResolve)
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

func TestGenerateDataOpaqueSecret(t *testing.T) {
	tests := []struct {
		name                 string
		inputSecret          internalSecret
		expectExternalSecret esv1beta1.ExternalSecret
		store                esv1beta1.SecretStoreRef
		envs                 map[string]string // for render <% ENV %>
		err                  error
		enableResolve        bool
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
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "simple_example",
						CreationPolicy: esv1beta1.CreatePolicyOrphan,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeOpaque,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyReplace,
							Data: map[string]string{
								"dist": `"{{ .dist-name-of-linux }}"`,
							},
						},
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "test",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "dist-name-of-linux",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "foo",
								MetadataPolicy:     "None",
								Property:           "dist-name-of-linux",
								ConversionStrategy: "Default",
								DecodingStrategy:   "Base64",
							},
						},
					},
				},
			},
		},
		{
			name:          "opaque type secret with path <% ENV %>",
			enableResolve: true,
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
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "multiple_env_with_path",
						CreationPolicy: esv1beta1.CreatePolicyOrphan,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type:        corev1.SecretTypeOpaque,
							Metadata:    esv1beta1.ExternalSecretTemplateMetadata{Labels: map[string]string{"app": "test"}},
							MergePolicy: esv1beta1.MergePolicyReplace,
							Data:        map[string]string{"dist": `"{{ .dist-name-of-linux }}"`},
						},
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "test",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "dist-name-of-linux",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "ubuntu-22.04-foo",
								MetadataPolicy:     "None",
								Property:           "dist-name-of-linux",
								ConversionStrategy: "Default",
								DecodingStrategy:   "Base64",
							},
						},
					},
				},
			},
		},
		{
			name:          "opaque type secret with path <% ENV %> and multiple property",
			enableResolve: true,
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
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "multiple_property",
						CreationPolicy: esv1beta1.CreatePolicyOrphan,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeOpaque,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyReplace,
							Data: map[string]string{
								"dist":   `"{{ .dist-name-of-linux }}"`,
								"passwd": `"{{ .github-passwd }}"`,
								"user":   `"{{ .github-username }}"`,
							},
						},
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "test",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "dist-name-of-linux",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "ubuntu-22.04-foo",
								MetadataPolicy:     "None",
								Property:           "dist-name-of-linux",
								ConversionStrategy: "Default",
								DecodingStrategy:   "Base64",
							},
						},
						{
							SecretKey: "github-username",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "ubuntu-22.04-foo",
								MetadataPolicy:     "None",
								Property:           "github-username",
								ConversionStrategy: "Default",
								DecodingStrategy:   "Base64",
							},
						},
						{
							SecretKey: "github-passwd",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "ubuntu-22.04-foo",
								MetadataPolicy:     "None",
								Property:           "github-passwd",
								ConversionStrategy: "Default",
								DecodingStrategy:   "Base64",
							},
						},
					},
				},
			},
		},
		{
			name: "test <<% ENV %>_VAULT> with data",
			inputSecret: internalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name: "mix_two_style",
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/foo",
					},
					Labels: map[string]string{
						"app": "test",
					},
				},
				Data: map[string]string{
					"env0": "<VAULT0>",
					"env1": "<<% ENV1 %>_VAULT1>",
					"env2": "<<% ENV2 %>_VAULT2>",
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
					Name:      "mix_two_style",
					Namespace: "",
					Labels: map[string]string{
						"app": "test",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "mix_two_style",
						CreationPolicy: esv1beta1.CreatePolicyOrphan,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeOpaque,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyReplace,
							Data: map[string]string{
								"env0": `"{{ .VAULT0 }}"`,
								"env1": `"{{ .<% ENV1 %>_VAULT1 }}"`,
								"env2": `"{{ .<% ENV2 %>_VAULT2 }}"`,
							},
						},
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "test",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "VAULT0",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "foo",
								MetadataPolicy:     "None",
								Property:           "VAULT0",
								ConversionStrategy: "Default",
								DecodingStrategy:   "Base64",
							},
						},
						{
							SecretKey: "<% ENV1 %>_VAULT1",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "foo",
								MetadataPolicy:     "None",
								Property:           "<% ENV1 %>_VAULT1",
								ConversionStrategy: "Default",
								DecodingStrategy:   "Base64",
							},
						},
						{
							SecretKey: "<% ENV2 %>_VAULT2",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "foo",
								MetadataPolicy:     "None",
								Property:           "<% ENV2 %>_VAULT2",
								ConversionStrategy: "Default",
								DecodingStrategy:   "Base64",
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
			name:          "resolve the value from env case 2",
			enableResolve: true,
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
			envs: map[string]string{
				"DIST": "ubuntu",
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
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "set_env_with_body_no_gen_es_2",
						CreationPolicy: esv1beta1.CreatePolicyOrphan,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeOpaque,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyReplace,
							Data: map[string]string{
								"data1": "data1",
								"data2": "ubuntu",
								"data3": `"{{ .FROM_VAULT_DATA3 }}"`,
							},
						},
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "test",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "FROM_VAULT_DATA3",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "foo",
								MetadataPolicy:     "None",
								Property:           "FROM_VAULT_DATA3",
								ConversionStrategy: "Default",
								DecodingStrategy:   "Base64",
							},
						},
					},
				},
			},
		},
		{
			name: "not support <VAULT>_<VAULT1> with data field",
			inputSecret: internalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name: "mix_example1",
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/foo-test",
					},
					Labels: map[string]string{
						"app": "test",
					},
				},
				Data: map[string]string{
					"key1": "<admin>-<dist-name-of-linux>",
				},
			},
			store: esv1beta1.SecretStoreRef{
				Name: "test",
				Kind: "ClusterSecretStore",
			},
			err: fmt.Errorf(ErrCommonNotSupportMultipleValue, "mix_example1"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envs {
				_ = os.Setenv(k, v)
			}
			externalSecret, err := convertSecret2ExtSecret(tt.inputSecret, tt.store.Kind, tt.store.Name, esv1beta1.CreatePolicyOrphan, tt.enableResolve)
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

func TestGenerateCreatePolicy(t *testing.T) {
	tests := []struct {
		name                 string
		inputSecret          internalSecret
		expectExternalSecret esv1beta1.ExternalSecret
		store                esv1beta1.SecretStoreRef
		envs                 map[string]string // for render <% ENV %>
		err                  error
	}{
		{
			name: "simple case for CreatePolicyOwner",
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
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "simple_example",
						CreationPolicy: esv1beta1.CreatePolicyOwner,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeOpaque,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyReplace,
							Data: map[string]string{
								"dist": `"{{ .dist-name-of-linux }}"`,
							},
						},
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "test",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "dist-name-of-linux",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "foo",
								MetadataPolicy:     "None",
								Property:           "dist-name-of-linux",
								ConversionStrategy: "Default",
								DecodingStrategy:   "Base64",
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
			externalSecret, err := convertSecret2ExtSecret(tt.inputSecret, tt.store.Kind, tt.store.Name, esv1beta1.CreatePolicyOwner, true)
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

func TestSkipResolveValue(t *testing.T) {
	tests := []struct {
		name                 string
		inputSecret          internalSecret
		expectExternalSecret esv1beta1.ExternalSecret
		store                esv1beta1.SecretStoreRef
		envs                 map[string]string // for render <% ENV %>
		err                  error
	}{
		{
			name: "simple case for skip resolve value",
			inputSecret: internalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name: "simple_example",
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/<% ENV %>-test",
					},
					Labels: map[string]string{
						"app": "test",
					},
				},
				Data: map[string]string{
					"dist": "<dist-name-of-linux>",
					"env1": "<% ENV1 %>",
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
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "simple_example",
						CreationPolicy: esv1beta1.CreatePolicyOwner,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeOpaque,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyReplace,
							Data: map[string]string{
								"dist": `"{{ .dist-name-of-linux }}"`,
								"env1": "<% ENV1 %>",
							},
						},
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "test",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "dist-name-of-linux",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "<% ENV %>-test",
								MetadataPolicy:     "None",
								Property:           "dist-name-of-linux",
								ConversionStrategy: "Default",
								DecodingStrategy:   "Base64",
							},
						},
					},
				},
			},
		},
		{
			name: "mix the vault of <% ENV %>",
			inputSecret: internalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name: "mix_example1",
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/<% ENV %>-test",
					},
					Labels: map[string]string{
						"app": "test",
					},
				},
				StringData: map[string]string{
					"dist": "<dist-name-of-linux>",
					"env1": "<% ENV1 %>",
					"env2": "<% ENV1 %>-<dist-name-of-linux>",
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
					Name:      "mix_example1",
					Namespace: "",
					Labels: map[string]string{
						"app": "test",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "mix_example1",
						CreationPolicy: esv1beta1.CreatePolicyOwner,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeOpaque,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyReplace,
							Data: map[string]string{
								"dist": `"{{ .dist-name-of-linux }}"`,
								"env1": "<% ENV1 %>",
								"env2": `"<% ENV1 %>-{{ .dist-name-of-linux }}"`,
							},
						},
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "test",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "dist-name-of-linux",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "<% ENV %>-test",
								MetadataPolicy:     "None",
								Property:           "dist-name-of-linux",
								ConversionStrategy: "Default",
								DecodingStrategy:   "None",
							},
						},
					},
				},
			},
		},
		{
			name: "mix the env value and vault value",
			inputSecret: internalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name: "mix_example2",
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/foo-test",
					},
					Labels: map[string]string{
						"app": "test",
					},
				},
				Data: map[string]string{
					"key": "<% ENV1 %>-<dist-name-of-linux>",
				},
			},
			store: esv1beta1.SecretStoreRef{
				Name: "test",
				Kind: "ClusterSecretStore",
			},
			err: fmt.Errorf(ErrCommonNotSupportMultipleValue, "mix_example2"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envs {
				_ = os.Setenv(k, v)
			}
			externalSecret, err := convertSecret2ExtSecret(tt.inputSecret, tt.store.Kind, tt.store.Name, esv1beta1.CreatePolicyOwner, false)
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
