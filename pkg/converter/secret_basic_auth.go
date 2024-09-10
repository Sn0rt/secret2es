package converter

import (
	"fmt"
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

func generateEsByBasicAuthSecret(inputSecret *internalSecret, storeType, storeName string,
	creationPolicy esv1beta1.ExternalSecretCreationPolicy, resolve bool) (*esv1beta1.ExternalSecret, error) {
	if len(inputSecret.Data) != 0 {
		return nil, fmt.Errorf(ErrBasicAuthNotAllowDataField, inputSecret.Name)
	}
	if inputSecret.StringData[corev1.BasicAuthUsernameKey] == "" {
		return nil, fmt.Errorf(ErrBasicAuthWithEmptyUsername, inputSecret.Name)
	}
	if inputSecret.StringData[corev1.BasicAuthPasswordKey] == "" {
		return nil, fmt.Errorf(ErrBasicAuthWithEmptyPassword, inputSecret.Name)
	}

	output, err := generateEsByOpaqueSecret(inputSecret, storeType, storeName, creationPolicy, resolve)
	if err != nil {
		return nil, err
	}
	output.Spec.Target.Template.Type = corev1.SecretTypeBasicAuth

	return output, nil
}
