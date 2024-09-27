package converter

import (
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestGenerateEsByTLS(t *testing.T) {
	var tests = []struct {
		name                 string
		input                []byte
		store                esv1beta1.SecretStoreRef
		expectExternalSecret esv1beta1.ExternalSecret
	}{
		{
			name: "simple use case",
			store: esv1beta1.SecretStoreRef{
				Name: "tenant-b",
				Kind: "ClusterSecretStore",
			},
			input: []byte(`
---
apiVersion: v1
kind: Secret
metadata:
  name: tls_secret_case1
  annotations:
    avp.kubernetes.io/path: "secret/data/test-foo"
  labels:
    "app": "test"
type: kubernetes.io/tls
data:
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNyakNDQVpZQ0NRQ1N4TjdEbUl3OVRqQU5CZ2txaGtpRzl3MEJBUXNGQURBWk1SY3dGUVlEVlFRRERBNTUKYjNWeVpHOXRZV2x1TG1OdmJUQWVGdzB5TkRBNE1qWXdOakV4TlRKYUZ3MHlOVEE0TWpZd05qRXhOVEphTUJreApGekFWQmdOVkJBTU1Ebmx2ZFhKa2IyMWhhVzR1WTI5dE1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBCk1JSUJDZ0tDQVFFQXpJZDZDMU12ZkN3V0xDanNnejEwa29Ga3M2RklIbHlVNElwUDVtcitERVRGTnFKT1p6dnoKZStreGFFNjBsYkNhVDV6U2YxZDllQWM0M0t2b0w1eXBieUxWVGJjdCtlNnNYMm9rbWlzdGtxUmRxcjNtMm9hSAoyY3pKeUhEVVpyT3Z6SkRHTDJoNGdUdE03QXpsb3VaN3ViOGZNQUJDR3B5bUppNjlzMEZRQ21DakltWUdxcm02CnlpOU83VXp4bTlabmgzUWhXZ2xzbFJuS05oVUhzdHIxbnQ0K1NsMWU2TEhBbHJtTzF5eVJHUmphdHh1d1NKYTMKTUZKeFJnTHRWbnlMNzJmTWY3c1R3RzcrbDVXMmhsM2x5QW1yeGpORnIvMGJ6WHBVZHFnc0dObW84Ny80NmdSego1UFMrZVc5UzNwVDZPN2NkUlQzcTB3NVk2VUhidGdIQ3d3SURBUUFCTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCCkFRQU1HS3paS2ZsTllwRkpDczNMMEt6TFgrWmEzdG9jQUlBODFjQXU0NzNEem9uc1B3cEZaUnRPeVAzV0Foc0EKalpNcitnaVhkY3lvWjVEQTdEUkkxN0UxSDduZTFiaDR6RmtYRE1HdGQxdnZXM0xQNVlhb2NxUjlzdGMyL3A0dgpxVE03bjZ0alRqY2RYNEQ2eG5KSHRzbmF1dVBwTUdiTzUwK04yK3JobU1NbjZPVmpFRkgrRWlQYmYzNWtSbkhXCi83ZnowWnVtYkxwNUlqdWFjSFM2YXJwR25KNGZON1I2NVNHa0FpNEtvMFZ6VTNNM1laclFneFdpK29aTHpTUHUKUUZveWpYRlgvQlhBRG9vaEFuTlpkN2FmVmFaMlU3MjJqaEpKaEkxM0tobHRXb2RUT2hQVytabWxYeHZmRy9acwprdU1SVmZraHowaGlQWGtMWUVvQTZlN3MKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
  tls.key: <TLS_KEY_VAULT>`),
			expectExternalSecret: esv1beta1.ExternalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "external-secrets.io/v1beta1",
					Kind:       "ExternalSecret",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tls_secret_case1",
					Namespace: "",
					Labels: map[string]string{
						"app": "test",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "tls_secret_case1",
						CreationPolicy: esv1beta1.CreatePolicyOrphan,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeTLS,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyReplace,
							Data: map[string]string{
								"tls.crt": `{{ "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNyakNDQVpZQ0NRQ1N4TjdEbUl3OVRqQU5CZ2txaGtpRzl3MEJBUXNGQURBWk1SY3dGUVlEVlFRRERBNTUKYjNWeVpHOXRZV2x1TG1OdmJUQWVGdzB5TkRBNE1qWXdOakV4TlRKYUZ3MHlOVEE0TWpZd05qRXhOVEphTUJreApGekFWQmdOVkJBTU1Ebmx2ZFhKa2IyMWhhVzR1WTI5dE1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBCk1JSUJDZ0tDQVFFQXpJZDZDMU12ZkN3V0xDanNnejEwa29Ga3M2RklIbHlVNElwUDVtcitERVRGTnFKT1p6dnoKZStreGFFNjBsYkNhVDV6U2YxZDllQWM0M0t2b0w1eXBieUxWVGJjdCtlNnNYMm9rbWlzdGtxUmRxcjNtMm9hSAoyY3pKeUhEVVpyT3Z6SkRHTDJoNGdUdE03QXpsb3VaN3ViOGZNQUJDR3B5bUppNjlzMEZRQ21DakltWUdxcm02CnlpOU83VXp4bTlabmgzUWhXZ2xzbFJuS05oVUhzdHIxbnQ0K1NsMWU2TEhBbHJtTzF5eVJHUmphdHh1d1NKYTMKTUZKeFJnTHRWbnlMNzJmTWY3c1R3RzcrbDVXMmhsM2x5QW1yeGpORnIvMGJ6WHBVZHFnc0dObW84Ny80NmdSego1UFMrZVc5UzNwVDZPN2NkUlQzcTB3NVk2VUhidGdIQ3d3SURBUUFCTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCCkFRQU1HS3paS2ZsTllwRkpDczNMMEt6TFgrWmEzdG9jQUlBODFjQXU0NzNEem9uc1B3cEZaUnRPeVAzV0Foc0EKalpNcitnaVhkY3lvWjVEQTdEUkkxN0UxSDduZTFiaDR6RmtYRE1HdGQxdnZXM0xQNVlhb2NxUjlzdGMyL3A0dgpxVE03bjZ0alRqY2RYNEQ2eG5KSHRzbmF1dVBwTUdiTzUwK04yK3JobU1NbjZPVmpFRkgrRWlQYmYzNWtSbkhXCi83ZnowWnVtYkxwNUlqdWFjSFM2YXJwR25KNGZON1I2NVNHa0FpNEtvMFZ6VTNNM1laclFneFdpK29aTHpTUHUKUUZveWpYRlgvQlhBRG9vaEFuTlpkN2FmVmFaMlU3MjJqaEpKaEkxM0tobHRXb2RUT2hQVytabWxYeHZmRy9acwprdU1SVmZraHowaGlQWGtMWUVvQTZlN3MKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=" | b64dec }}`,
								"tls.key": `"{{ .TLS_KEY_VAULT }}"`,
							},
						},
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "tenant-b",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "TLS_KEY_VAULT",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "test-foo",
								MetadataPolicy:     "None",
								Property:           "TLS_KEY_VAULT",
								ConversionStrategy: "Default",
								DecodingStrategy:   "Base64",
							},
						},
					},
				},
			},
		},
		{
			name: "long_name_of_secret with -",
			store: esv1beta1.SecretStoreRef{
				Name: "tenant-b",
				Kind: "ClusterSecretStore",
			},
			input: []byte(`
---
apiVersion: v1
kind: Secret
metadata:
  name: open-source-secret-with-github-action-test-sn0rt-dev
  annotations:
    avp.kubernetes.io/path: "secret/data/test-foo"
  labels:
    "app": "test"
type: kubernetes.io/tls
data:
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNyakNDQVpZQ0NRQ1N4TjdEbUl3OVRqQU5CZ2txaGtpRzl3MEJBUXNGQURBWk1SY3dGUVlEVlFRRERBNTUKYjNWeVpHOXRZV2x1TG1OdmJUQWVGdzB5TkRBNE1qWXdOakV4TlRKYUZ3MHlOVEE0TWpZd05qRXhOVEphTUJreApGekFWQmdOVkJBTU1Ebmx2ZFhKa2IyMWhhVzR1WTI5dE1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBCk1JSUJDZ0tDQVFFQXpJZDZDMU12ZkN3V0xDanNnejEwa29Ga3M2RklIbHlVNElwUDVtcitERVRGTnFKT1p6dnoKZStreGFFNjBsYkNhVDV6U2YxZDllQWM0M0t2b0w1eXBieUxWVGJjdCtlNnNYMm9rbWlzdGtxUmRxcjNtMm9hSAoyY3pKeUhEVVpyT3Z6SkRHTDJoNGdUdE03QXpsb3VaN3ViOGZNQUJDR3B5bUppNjlzMEZRQ21DakltWUdxcm02CnlpOU83VXp4bTlabmgzUWhXZ2xzbFJuS05oVUhzdHIxbnQ0K1NsMWU2TEhBbHJtTzF5eVJHUmphdHh1d1NKYTMKTUZKeFJnTHRWbnlMNzJmTWY3c1R3RzcrbDVXMmhsM2x5QW1yeGpORnIvMGJ6WHBVZHFnc0dObW84Ny80NmdSego1UFMrZVc5UzNwVDZPN2NkUlQzcTB3NVk2VUhidGdIQ3d3SURBUUFCTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCCkFRQU1HS3paS2ZsTllwRkpDczNMMEt6TFgrWmEzdG9jQUlBODFjQXU0NzNEem9uc1B3cEZaUnRPeVAzV0Foc0EKalpNcitnaVhkY3lvWjVEQTdEUkkxN0UxSDduZTFiaDR6RmtYRE1HdGQxdnZXM0xQNVlhb2NxUjlzdGMyL3A0dgpxVE03bjZ0alRqY2RYNEQ2eG5KSHRzbmF1dVBwTUdiTzUwK04yK3JobU1NbjZPVmpFRkgrRWlQYmYzNWtSbkhXCi83ZnowWnVtYkxwNUlqdWFjSFM2YXJwR25KNGZON1I2NVNHa0FpNEtvMFZ6VTNNM1laclFneFdpK29aTHpTUHUKUUZveWpYRlgvQlhBRG9vaEFuTlpkN2FmVmFaMlU3MjJqaEpKaEkxM0tobHRXb2RUT2hQVytabWxYeHZmRy9acwprdU1SVmZraHowaGlQWGtMWUVvQTZlN3MKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
  tls.key: <TLS_KEY_VAULT>`),
			expectExternalSecret: esv1beta1.ExternalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "external-secrets.io/v1beta1",
					Kind:       "ExternalSecret",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "open-source-secret-with-github-action-test-sn0rt-dev",
					Namespace: "",
					Labels: map[string]string{
						"app": "test",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "open-source-secret-with-github-action-test-sn0rt-dev",
						CreationPolicy: esv1beta1.CreatePolicyOrphan,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeTLS,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyReplace,
							Data: map[string]string{
								"tls.crt": `{{ "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNyakNDQVpZQ0NRQ1N4TjdEbUl3OVRqQU5CZ2txaGtpRzl3MEJBUXNGQURBWk1SY3dGUVlEVlFRRERBNTUKYjNWeVpHOXRZV2x1TG1OdmJUQWVGdzB5TkRBNE1qWXdOakV4TlRKYUZ3MHlOVEE0TWpZd05qRXhOVEphTUJreApGekFWQmdOVkJBTU1Ebmx2ZFhKa2IyMWhhVzR1WTI5dE1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBCk1JSUJDZ0tDQVFFQXpJZDZDMU12ZkN3V0xDanNnejEwa29Ga3M2RklIbHlVNElwUDVtcitERVRGTnFKT1p6dnoKZStreGFFNjBsYkNhVDV6U2YxZDllQWM0M0t2b0w1eXBieUxWVGJjdCtlNnNYMm9rbWlzdGtxUmRxcjNtMm9hSAoyY3pKeUhEVVpyT3Z6SkRHTDJoNGdUdE03QXpsb3VaN3ViOGZNQUJDR3B5bUppNjlzMEZRQ21DakltWUdxcm02CnlpOU83VXp4bTlabmgzUWhXZ2xzbFJuS05oVUhzdHIxbnQ0K1NsMWU2TEhBbHJtTzF5eVJHUmphdHh1d1NKYTMKTUZKeFJnTHRWbnlMNzJmTWY3c1R3RzcrbDVXMmhsM2x5QW1yeGpORnIvMGJ6WHBVZHFnc0dObW84Ny80NmdSego1UFMrZVc5UzNwVDZPN2NkUlQzcTB3NVk2VUhidGdIQ3d3SURBUUFCTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCCkFRQU1HS3paS2ZsTllwRkpDczNMMEt6TFgrWmEzdG9jQUlBODFjQXU0NzNEem9uc1B3cEZaUnRPeVAzV0Foc0EKalpNcitnaVhkY3lvWjVEQTdEUkkxN0UxSDduZTFiaDR6RmtYRE1HdGQxdnZXM0xQNVlhb2NxUjlzdGMyL3A0dgpxVE03bjZ0alRqY2RYNEQ2eG5KSHRzbmF1dVBwTUdiTzUwK04yK3JobU1NbjZPVmpFRkgrRWlQYmYzNWtSbkhXCi83ZnowWnVtYkxwNUlqdWFjSFM2YXJwR25KNGZON1I2NVNHa0FpNEtvMFZ6VTNNM1laclFneFdpK29aTHpTUHUKUUZveWpYRlgvQlhBRG9vaEFuTlpkN2FmVmFaMlU3MjJqaEpKaEkxM0tobHRXb2RUT2hQVytabWxYeHZmRy9acwprdU1SVmZraHowaGlQWGtMWUVvQTZlN3MKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=" | b64dec }}`,
								"tls.key": `"{{ .TLS_KEY_VAULT }}"`,
							},
						},
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "tenant-b",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "TLS_KEY_VAULT",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "test-foo",
								MetadataPolicy:     "None",
								Property:           "TLS_KEY_VAULT",
								ConversionStrategy: "Default",
								DecodingStrategy:   "Base64",
							},
						},
					},
				},
			},
		},
		{
			name: "commented secret",
			store: esv1beta1.SecretStoreRef{
				Name: "tenant-b",
				Kind: "ClusterSecretStore",
			},
			input: []byte(`
---
apiVersion: v1
kind: Secret
metadata:
  name: open-source-secret-with-github-action-test-sn0rt-dev
  annotations:
    avp.kubernetes.io/path: "secret/data/test-foo"
  labels:
    "app": "test"
type: kubernetes.io/tls
data:
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNyakNDQVpZQ0NRQ1N4TjdEbUl3OVRqQU5CZ2txaGtpRzl3MEJBUXNGQURBWk1SY3dGUVlEVlFRRERBNTUKYjNWeVpHOXRZV2x1TG1OdmJUQWVGdzB5TkRBNE1qWXdOakV4TlRKYUZ3MHlOVEE0TWpZd05qRXhOVEphTUJreApGekFWQmdOVkJBTU1Ebmx2ZFhKa2IyMWhhVzR1WTI5dE1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBCk1JSUJDZ0tDQVFFQXpJZDZDMU12ZkN3V0xDanNnejEwa29Ga3M2RklIbHlVNElwUDVtcitERVRGTnFKT1p6dnoKZStreGFFNjBsYkNhVDV6U2YxZDllQWM0M0t2b0w1eXBieUxWVGJjdCtlNnNYMm9rbWlzdGtxUmRxcjNtMm9hSAoyY3pKeUhEVVpyT3Z6SkRHTDJoNGdUdE03QXpsb3VaN3ViOGZNQUJDR3B5bUppNjlzMEZRQ21DakltWUdxcm02CnlpOU83VXp4bTlabmgzUWhXZ2xzbFJuS05oVUhzdHIxbnQ0K1NsMWU2TEhBbHJtTzF5eVJHUmphdHh1d1NKYTMKTUZKeFJnTHRWbnlMNzJmTWY3c1R3RzcrbDVXMmhsM2x5QW1yeGpORnIvMGJ6WHBVZHFnc0dObW84Ny80NmdSego1UFMrZVc5UzNwVDZPN2NkUlQzcTB3NVk2VUhidGdIQ3d3SURBUUFCTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCCkFRQU1HS3paS2ZsTllwRkpDczNMMEt6TFgrWmEzdG9jQUlBODFjQXU0NzNEem9uc1B3cEZaUnRPeVAzV0Foc0EKalpNcitnaVhkY3lvWjVEQTdEUkkxN0UxSDduZTFiaDR6RmtYRE1HdGQxdnZXM0xQNVlhb2NxUjlzdGMyL3A0dgpxVE03bjZ0alRqY2RYNEQ2eG5KSHRzbmF1dVBwTUdiTzUwK04yK3JobU1NbjZPVmpFRkgrRWlQYmYzNWtSbkhXCi83ZnowWnVtYkxwNUlqdWFjSFM2YXJwR25KNGZON1I2NVNHa0FpNEtvMFZ6VTNNM1laclFneFdpK29aTHpTUHUKUUZveWpYRlgvQlhBRG9vaEFuTlpkN2FmVmFaMlU3MjJqaEpKaEkxM0tobHRXb2RUT2hQVytabWxYeHZmRy9acwprdU1SVmZraHowaGlQWGtMWUVvQTZlN3MKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
  tls.key: <TLS_KEY_VAULT>
---
#apiVersion: v1
#kind: Secret
#metadata:
#  name: open-source-secret-with-github-action-test-sn0rt-dev
#  annotations:
#    avp.kubernetes.io/path: "secret/data/test-foo"
#  labels:
#    "app": "test"
#type: kubernetes.io/tls
#data:
#  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNyakNDQVpZQ0NRQ1N4TjdEbUl3OVRqQU5CZ2txaGtpRzl3MEJBUXNGQURBWk1SY3dGUVlEVlFRRERBNTUKYjNWeVpHOXRZV2x1TG1OdmJUQWVGdzB5TkRBNE1qWXdOakV4TlRKYUZ3MHlOVEE0TWpZd05qRXhOVEphTUJreApGekFWQmdOVkJBTU1Ebmx2ZFhKa2IyMWhhVzR1WTI5dE1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBCk1JSUJDZ0tDQVFFQXpJZDZDMU12ZkN3V0xDanNnejEwa29Ga3M2RklIbHlVNElwUDVtcitERVRGTnFKT1p6dnoKZStreGFFNjBsYkNhVDV6U2YxZDllQWM0M0t2b0w1eXBieUxWVGJjdCtlNnNYMm9rbWlzdGtxUmRxcjNtMm9hSAoyY3pKeUhEVVpyT3Z6SkRHTDJoNGdUdE03QXpsb3VaN3ViOGZNQUJDR3B5bUppNjlzMEZRQ21DakltWUdxcm02CnlpOU83VXp4bTlabmgzUWhXZ2xzbFJuS05oVUhzdHIxbnQ0K1NsMWU2TEhBbHJtTzF5eVJHUmphdHh1d1NKYTMKTUZKeFJnTHRWbnlMNzJmTWY3c1R3RzcrbDVXMmhsM2x5QW1yeGpORnIvMGJ6WHBVZHFnc0dObW84Ny80NmdSego1UFMrZVc5UzNwVDZPN2NkUlQzcTB3NVk2VUhidGdIQ3d3SURBUUFCTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCCkFRQU1HS3paS2ZsTllwRkpDczNMMEt6TFgrWmEzdG9jQUlBODFjQXU0NzNEem9uc1B3cEZaUnRPeVAzV0Foc0EKalpNcitnaVhkY3lvWjVEQTdEUkkxN0UxSDduZTFiaDR6RmtYRE1HdGQxdnZXM0xQNVlhb2NxUjlzdGMyL3A0dgpxVE03bjZ0alRqY2RYNEQ2eG5KSHRzbmF1dVBwTUdiTzUwK04yK3JobU1NbjZPVmpFRkgrRWlQYmYzNWtSbkhXCi83ZnowWnVtYkxwNUlqdWFjSFM2YXJwR25KNGZON1I2NVNHa0FpNEtvMFZ6VTNNM1laclFneFdpK29aTHpTUHUKUUZveWpYRlgvQlhBRG9vaEFuTlpkN2FmVmFaMlU3MjJqaEpKaEkxM0tobHRXb2RUT2hQVytabWxYeHZmRy9acwprdU1SVmZraHowaGlQWGtMWUVvQTZlN3MKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
#  tls.key: <TLS_KEY_VAULT>`),
			expectExternalSecret: esv1beta1.ExternalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "external-secrets.io/v1beta1",
					Kind:       "ExternalSecret",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "open-source-secret-with-github-action-test-sn0rt-dev",
					Namespace: "",
					Labels: map[string]string{
						"app": "test",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "open-source-secret-with-github-action-test-sn0rt-dev",
						CreationPolicy: esv1beta1.CreatePolicyOrphan,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeTLS,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyReplace,
							Data: map[string]string{
								"tls.crt": `{{ "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNyakNDQVpZQ0NRQ1N4TjdEbUl3OVRqQU5CZ2txaGtpRzl3MEJBUXNGQURBWk1SY3dGUVlEVlFRRERBNTUKYjNWeVpHOXRZV2x1TG1OdmJUQWVGdzB5TkRBNE1qWXdOakV4TlRKYUZ3MHlOVEE0TWpZd05qRXhOVEphTUJreApGekFWQmdOVkJBTU1Ebmx2ZFhKa2IyMWhhVzR1WTI5dE1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBCk1JSUJDZ0tDQVFFQXpJZDZDMU12ZkN3V0xDanNnejEwa29Ga3M2RklIbHlVNElwUDVtcitERVRGTnFKT1p6dnoKZStreGFFNjBsYkNhVDV6U2YxZDllQWM0M0t2b0w1eXBieUxWVGJjdCtlNnNYMm9rbWlzdGtxUmRxcjNtMm9hSAoyY3pKeUhEVVpyT3Z6SkRHTDJoNGdUdE03QXpsb3VaN3ViOGZNQUJDR3B5bUppNjlzMEZRQ21DakltWUdxcm02CnlpOU83VXp4bTlabmgzUWhXZ2xzbFJuS05oVUhzdHIxbnQ0K1NsMWU2TEhBbHJtTzF5eVJHUmphdHh1d1NKYTMKTUZKeFJnTHRWbnlMNzJmTWY3c1R3RzcrbDVXMmhsM2x5QW1yeGpORnIvMGJ6WHBVZHFnc0dObW84Ny80NmdSego1UFMrZVc5UzNwVDZPN2NkUlQzcTB3NVk2VUhidGdIQ3d3SURBUUFCTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCCkFRQU1HS3paS2ZsTllwRkpDczNMMEt6TFgrWmEzdG9jQUlBODFjQXU0NzNEem9uc1B3cEZaUnRPeVAzV0Foc0EKalpNcitnaVhkY3lvWjVEQTdEUkkxN0UxSDduZTFiaDR6RmtYRE1HdGQxdnZXM0xQNVlhb2NxUjlzdGMyL3A0dgpxVE03bjZ0alRqY2RYNEQ2eG5KSHRzbmF1dVBwTUdiTzUwK04yK3JobU1NbjZPVmpFRkgrRWlQYmYzNWtSbkhXCi83ZnowWnVtYkxwNUlqdWFjSFM2YXJwR25KNGZON1I2NVNHa0FpNEtvMFZ6VTNNM1laclFneFdpK29aTHpTUHUKUUZveWpYRlgvQlhBRG9vaEFuTlpkN2FmVmFaMlU3MjJqaEpKaEkxM0tobHRXb2RUT2hQVytabWxYeHZmRy9acwprdU1SVmZraHowaGlQWGtMWUVvQTZlN3MKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=" | b64dec }}`,
								"tls.key": `"{{ .TLS_KEY_VAULT }}"`,
							},
						},
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "tenant-b",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "TLS_KEY_VAULT",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "test-foo",
								MetadataPolicy:     "None",
								Property:           "TLS_KEY_VAULT",
								ConversionStrategy: "Default",
								DecodingStrategy:   "Base64",
							},
						},
					},
				},
			},
		},
		{
			name: "set cert with pain text",
			store: esv1beta1.SecretStoreRef{
				Name: "tenant-b",
				Kind: "ClusterSecretStore",
			},
			input: []byte(`
---
apiVersion: v1
kind: Secret
metadata:
  name: set-cert-with-pain-text
  annotations:
    avp.kubernetes.io/path: "secret/data/test-foo"
  labels:
    "app": "test"
type: kubernetes.io/tls
StringData:
  tls.crt: |
    -----BEGIN CERTIFICATE-----
    MIICrjCCAZYCCQCSxN7DmIw9TjANBgkqhkiG9w0BAQsFADAZMRcwFQYDVQQDDA55
    b3VyZG9tYWluLmNvbTAeFw0yNDA4MjYwNjExNTJaFw0yNTA4MjYwNjExNTJaMBkx
    FzAVBgNVBAMMDnlvdXJkb21haW4uY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
    MIIBCgKCAQEAzId6C1MvfCwWLCjsgz10koFks6FIHlyU4IpP5mr+DETFNqJOZzvz
    e+kxaE60lbCaT5zSf1d9eAc43KvoL5ypbyLVTbct+e6sX2okmistkqRdqr3m2oaH
    2czJyHDUZrOvzJDGL2h4gTtM7AzlouZ7ub8fMABCGpymJi69s0FQCmCjImYGqrm6
    yi9O7Uzxm9Znh3QhWglslRnKNhUHstr1nt4+Sl1e6LHAlrmO1yyRGRjatxuwSJa3
    MFJxRgLtVnyL72fMf7sTwG7+l5W2hl3lyAmrxjNFr/0bzXpUdqgsGNmo87/46gRz
    5PS+eW9S3pT6O7cdRT3q0w5Y6UHbtgHCwwIDAQABMA0GCSqGSIb3DQEBCwUAA4IB
    AQAMGKzZKflNYpFJCs3L0KzLX+Za3tocAIA81cAu473DzonsPwpFZRtOyP3WAhsA
    jZMr+giXdcyoZ5DA7DRI17E1H7ne1bh4zFkXDMGtd1vvW3LP5YaocqR9stc2/p4v
    qTM7n6tjTjcdX4D6xnJHtsnauuPpMGbO50+N2+rhmMMn6OVjEFH+EiPbf35kRnHW
    /7fz0ZumbLp5IjuacHS6arpGnJ4fN7R65SGkAi4Ko0VzU3M3YZrQgxWi+oZLzSPu
    QFoyjXFX/BXADoohAnNZd7afVaZ2U722jhJJhI13KhltWodTOhPW+ZmlXxvfG/Zs
    kuMRVfkhz0hiPXkLYEoA6e7s
    -----END CERTIFICATE----
  tls.key: <TLS_KEY_VAULT>`),
			expectExternalSecret: esv1beta1.ExternalSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "external-secrets.io/v1beta1",
					Kind:       "ExternalSecret",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "set-cert-with-pain-text",
					Namespace: "",
					Labels: map[string]string{
						"app": "test",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "set-cert-with-pain-text",
						CreationPolicy: esv1beta1.CreatePolicyOrphan,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
						Template: &esv1beta1.ExternalSecretTemplate{
							Type: corev1.SecretTypeTLS,
							Metadata: esv1beta1.ExternalSecretTemplateMetadata{
								Labels: map[string]string{
									"app": "test",
								},
							},
							MergePolicy: esv1beta1.MergePolicyReplace,
							Data: map[string]string{
								"tls.crt": `-----BEGIN CERTIFICATE-----
MIICrjCCAZYCCQCSxN7DmIw9TjANBgkqhkiG9w0BAQsFADAZMRcwFQYDVQQDDA55
b3VyZG9tYWluLmNvbTAeFw0yNDA4MjYwNjExNTJaFw0yNTA4MjYwNjExNTJaMBkx
FzAVBgNVBAMMDnlvdXJkb21haW4uY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEAzId6C1MvfCwWLCjsgz10koFks6FIHlyU4IpP5mr+DETFNqJOZzvz
e+kxaE60lbCaT5zSf1d9eAc43KvoL5ypbyLVTbct+e6sX2okmistkqRdqr3m2oaH
2czJyHDUZrOvzJDGL2h4gTtM7AzlouZ7ub8fMABCGpymJi69s0FQCmCjImYGqrm6
yi9O7Uzxm9Znh3QhWglslRnKNhUHstr1nt4+Sl1e6LHAlrmO1yyRGRjatxuwSJa3
MFJxRgLtVnyL72fMf7sTwG7+l5W2hl3lyAmrxjNFr/0bzXpUdqgsGNmo87/46gRz
5PS+eW9S3pT6O7cdRT3q0w5Y6UHbtgHCwwIDAQABMA0GCSqGSIb3DQEBCwUAA4IB
AQAMGKzZKflNYpFJCs3L0KzLX+Za3tocAIA81cAu473DzonsPwpFZRtOyP3WAhsA
jZMr+giXdcyoZ5DA7DRI17E1H7ne1bh4zFkXDMGtd1vvW3LP5YaocqR9stc2/p4v
qTM7n6tjTjcdX4D6xnJHtsnauuPpMGbO50+N2+rhmMMn6OVjEFH+EiPbf35kRnHW
/7fz0ZumbLp5IjuacHS6arpGnJ4fN7R65SGkAi4Ko0VzU3M3YZrQgxWi+oZLzSPu
QFoyjXFX/BXADoohAnNZd7afVaZ2U722jhJJhI13KhltWodTOhPW+ZmlXxvfG/Zs
kuMRVfkhz0hiPXkLYEoA6e7s
-----END CERTIFICATE----
`,
								"tls.key": `"{{ .TLS_KEY_VAULT }}"`,
							},
						},
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "tenant-b",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "TLS_KEY_VAULT",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:                "test-foo",
								MetadataPolicy:     "None",
								Property:           "TLS_KEY_VAULT",
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
			inputSecretList, _ := parseUnstructuredSecret(tt.input)
			out, err := convertSecret2ExtSecret(inputSecretList[0], tt.store.Kind, tt.store.Name, esv1beta1.CreatePolicyOrphan, true)
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
