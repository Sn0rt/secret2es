package converter

import (
	"fmt"
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"os"
	"sigs.k8s.io/yaml"
)

// ConvertSecret converts a Kubernetes Secret to an ExternalSecret
func ConvertSecret(inputFile, storeType, storeName string, outputPath string) error {
	bytes, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading inputSecret file: %w", err)
	}

	inputSecretList, err := parseUnstructuredSecret(bytes)
	if err != nil {
		return fmt.Errorf("error parsing inputSecret secret: %w", err)
	}

	for _, inputSecret := range inputSecretList {
		externalSecret, err := convertSecret2ExtSecret(inputSecret, storeType, storeName)
		if err != nil {
			return fmt.Errorf("error converting secret to external secret: %s", err.Error())
		}

		yamlData, err := yaml.Marshal(externalSecret)
		if err != nil {
			return fmt.Errorf("error encoding external secret: %w", err)
		}

		if outputPath == "" {
			fmt.Printf("---\n")
			fmt.Printf("%s", yamlData)
		} else {
			err = os.WriteFile(outputPath, yamlData, 0644)
			if err != nil {
				return fmt.Errorf("error writing external secret to file: %w", err)
			}
		}
	}

	return nil
}

func convertSecret2ExtSecret(inputSecret UnstructuredSecret, storeType, storeName string) (*esv1beta1.ExternalSecret, error) {
	if inputSecret.Annotations == nil {
		return nil, fmt.Errorf(NotEmptyAnnotations, inputSecret.Name)
	} else {
		if inputSecret.Annotations["avp.kubernetes.io/path"] == "" {
			return nil, fmt.Errorf(NotSetAVPAnnotations, inputSecret.Name)
		}
	}

	if storeType != SecretStoreType && storeType != ClusterSecretStoreType {
		return nil, fmt.Errorf(NotSupportedStoreType, storeType)
	}

	switch inputSecret.Type {
	case corev1.SecretTypeOpaque:
		return generateOpaqueSecret(inputSecret, storeType, storeName)
	case corev1.SecretTypeBasicAuth:
	case corev1.SecretTypeDockerConfigJson:
	case corev1.SecretTypeTLS:
	}

	return nil, fmt.Errorf(NotImplSecretType, inputSecret.Type, inputSecret.Name)
}
