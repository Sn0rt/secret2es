package converter

import (
	"fmt"
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

func generateEsByTLS(inputSecret *internalSecret, storeType, storeName string, creationPolicy esv1beta1.ExternalSecretCreationPolicy) (*esv1beta1.ExternalSecret, error) {
	if len(inputSecret.StringData) != 0 {
		return nil, fmt.Errorf(ErrTLSNotAllowDataField, inputSecret.Name)
	}

	// prepare the ref of sensitive data
	output, err := generateEsByOpaqueSecret(inputSecret, storeType, storeName, creationPolicy)
	if err != nil {
		return nil, err
	}
	output.Spec.Target.Template.Type = corev1.SecretTypeTLS

	return output, nil
}
