package converter

import (
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

func generateEsByTLS(inputSecret *internalSecret, storeType, storeName string,
	creationPolicy esv1beta1.ExternalSecretCreationPolicy, resolve bool) (*esv1beta1.ExternalSecret, error) {

	// prepare the ref of sensitive data
	output, err := generateEsByOpaqueSecret(inputSecret, storeType, storeName, creationPolicy, resolve)
	if err != nil {
		return nil, err
	}
	output.Spec.Target.Template.Type = corev1.SecretTypeTLS

	return output, nil
}
