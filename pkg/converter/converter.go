package converter

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
)

func ConvertSecretContent(content, storeType, storeName string, creationPolicy esv1beta1.ExternalSecretCreationPolicy, resolve bool) (string, error) {
	// Parse the content as YAML
	var secret corev1.Secret
	err := yaml.Unmarshal([]byte(content), &secret)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal secret: %w", err)
	}

	// Convert the secret
	externalSecret, err := convertToExternalSecret(&secret, storeType, storeName, creationPolicy, resolve)
	if err != nil {
		return "", fmt.Errorf("failed to convert secret: %w", err)
	}

	// Marshal the external secret back to YAML
	result, err := yaml.Marshal(externalSecret)
	if err != nil {
		return "", fmt.Errorf("failed to marshal external secret: %w", err)
	}

	return string(result), nil
}

func convertToExternalSecret(secret *corev1.Secret, storeType, storeName string, creationPolicy esv1beta1.ExternalSecretCreationPolicy, resolve bool) (*esv1beta1.ExternalSecret, error) {
	externalSecret := &esv1beta1.ExternalSecret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: esv1beta1.SchemeGroupVersion.String(),
			Kind:       "ExternalSecret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secret.Name,
			Namespace: secret.Namespace,
		},
		Spec: esv1beta1.ExternalSecretSpec{
			SecretStoreRef: esv1beta1.SecretStoreRef{
				Name: storeName,
				Kind: storeType,
			},
			Target: esv1beta1.ExternalSecretTarget{
				Name:           secret.Name,
				CreationPolicy: creationPolicy,
			},
			Data: []esv1beta1.ExternalSecretData{},
		},
	}

	for key, value := range secret.Data {
		externalSecret.Spec.Data = append(externalSecret.Spec.Data, esv1beta1.ExternalSecretData{
			SecretKey: key,
			RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
				Key:      key,
				Property: string(value),
			},
		})
	}

	return externalSecret, nil
}

