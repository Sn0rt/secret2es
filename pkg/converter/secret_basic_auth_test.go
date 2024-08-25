package converter

import (
	"fmt"
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
			err: fmt.Errorf(NotSupportedSecretDataEmpty, "input1"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envs {
				_ = os.Setenv(k, v)
			}
			externalSecret, err := generateBasicAuthSecret(tt.inputSecret, tt.store.Kind, tt.store.Name)
			if err != nil {
				if err == tt.err {
					t.Errorf("generateOpaqueSecret() returned an unexpected error: got: %v, want: %v", err, tt.err)
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
