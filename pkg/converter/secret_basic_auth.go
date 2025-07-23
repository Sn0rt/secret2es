package converter

import (
	"fmt"

	esv1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generateEsByBasicAuthSecret(inputSecret *internalSecret, storeType, storeName string,
	creationPolicy esv1.ExternalSecretCreationPolicy, resolve bool,
	refreshPolicy esv1.ExternalSecretRefreshPolicy, refreshInterval *metav1.Duration) (*esv1.ExternalSecret, error) {
	if len(inputSecret.Data) != 0 {
		return nil, fmt.Errorf(ErrBasicAuthNotAllowDataField, inputSecret.Name)
	}
	if inputSecret.StringData[corev1.BasicAuthUsernameKey] == "" {
		return nil, fmt.Errorf(ErrBasicAuthWithEmptyUsername, inputSecret.Name)
	}
	if inputSecret.StringData[corev1.BasicAuthPasswordKey] == "" {
		return nil, fmt.Errorf(ErrBasicAuthWithEmptyPassword, inputSecret.Name)
	}

	output, err := generateEsByOpaqueSecret(inputSecret, storeType, storeName, creationPolicy, resolve, refreshPolicy, refreshInterval)
	if err != nil {
		return nil, err
	}
	output.Spec.Target.Template.Type = corev1.SecretTypeBasicAuth

	return output, nil
}
