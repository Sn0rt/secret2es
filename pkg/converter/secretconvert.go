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
	if err := secretCommonVerify(inputSecret); err != nil {
		return nil, err
	}

	if storeType != SecretStoreType &&
		storeType != ClusterSecretStoreType {
		return nil, fmt.Errorf(illegalStoreType, storeType)
	}

	// get the secret of vault path
	var resolvedSecretPath = resolved(inputSecret.Annotations["avp.kubernetes.io/path"])
	inputSecret.Annotations["avp.kubernetes.io/path"] = resolvedSecretPath

	switch inputSecret.Type {
	case corev1.SecretTypeOpaque:
		return generateEsByOpaqueSecret(&inputSecret, storeType, storeName)
	case corev1.SecretTypeBasicAuth:
		return generateEsByBasicAuthSecret(&inputSecret, storeType, storeName)
	case corev1.SecretTypeDockerConfigJson:
		return generateEsByDockerConfigJSON(&inputSecret, storeType, storeName)
	case corev1.SecretTypeTLS:
		return generateEsByTLS(&inputSecret, storeType, storeName)
	}

	return nil, fmt.Errorf(NotImplSecretType, inputSecret.Type, inputSecret.Name)
}

func secretCommonVerify(inputSecret UnstructuredSecret) error {
	if inputSecret.Annotations == nil {
		return fmt.Errorf(ErrCommonNotEmptyAnnotations, inputSecret.Name)
	}
	if inputSecret.Annotations["avp.kubernetes.io/path"] == "" {
		return fmt.Errorf(ErrCommonNotFoundAVPPath, inputSecret.Name)
	}
	if len(inputSecret.Data) != 0 && len(inputSecret.StringData) != 0 {
		return fmt.Errorf(ErrCommonNotAcceptBothSecretDataAndData, inputSecret.Name)
	}
	if len(inputSecret.Data) == 0 && len(inputSecret.StringData) == 0 {
		fmt.Println("ErrCommonNotAcceptNeitherSecretDataAndData")
		return fmt.Errorf(ErrCommonNotAcceptNeitherSecretDataAndData, inputSecret.Name)
	}
	return nil
}
