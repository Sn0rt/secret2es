package converter

import (
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	"github.com/google/go-cmp/cmp"
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
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/test-foo",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "tls_secret_case1",
						CreationPolicy: esv1beta1.CreatePolicyMerge,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "tenant-b",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "tls.key",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:      "test-foo",
								Property: "TLS_KEY_VAULT",
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
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/test-foo",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "open-source-secret-with-github-action-test-sn0rt-dev",
						CreationPolicy: esv1beta1.CreatePolicyMerge,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "tenant-b",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "tls.key",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:      "test-foo",
								Property: "TLS_KEY_VAULT",
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
					Annotations: map[string]string{
						"avp.kubernetes.io/path": "secret/data/test-foo",
					},
				},
				Spec: esv1beta1.ExternalSecretSpec{
					RefreshInterval: stopRefreshInterval,
					Target: esv1beta1.ExternalSecretTarget{
						Name:           "open-source-secret-with-github-action-test-sn0rt-dev",
						CreationPolicy: esv1beta1.CreatePolicyMerge,
						DeletionPolicy: esv1beta1.DeletionPolicyRetain,
					},
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "tenant-b",
						Kind: "ClusterSecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "tls.key",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key:      "test-foo",
								Property: "TLS_KEY_VAULT",
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
			out, err := convertSecret2ExtSecret(inputSecretList[0], tt.store.Kind, tt.store.Name)
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
