package converter

import (
	"fmt"
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"os"
	"sigs.k8s.io/yaml"
)

// ConvertSecret converts a Kubernetes Secret to an ExternalSecret
func ConvertSecret(inputFile, storeType, storeName string) error {
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
		// handle error
		if err != nil {
			switch err.Error() {
			case fmt.Errorf(ErrCommonNotIncludeAngleBrackets, inputSecret.Name).Error():
				continue
			}
			return fmt.Errorf("error converting secret to external secret: %s", err.Error())
		}
		yamlData, err := yaml.Marshal(externalSecret)
		if err != nil {
			return fmt.Errorf("error encoding external secret: %w", err)
		}

		// remove the status field
		yamlData = removeStatusField(yamlData)
		fmt.Printf("---\n")
		fmt.Printf("%s", yamlData)
	}

	return nil
}

func removeStatusField(yamlData []byte) []byte {
	var externalSecret map[string]interface{}
	if err := yaml.Unmarshal(yamlData, &externalSecret); err != nil {
		return yamlData
	}

	delete(externalSecret, "status")
	newYamlData, err := yaml.Marshal(externalSecret)
	if err != nil {
		return yamlData
	}

	return newYamlData
}

func convertSecret2ExtSecret(inputSecret internalSecret, storeType, storeName string) (*esv1beta1.ExternalSecret, error) {
	if err := secretCommonVerify(inputSecret); err != nil {
		return nil, err
	}

	if storeType != SecretStoreType &&
		storeType != ClusterSecretStoreType {
		return nil, fmt.Errorf(illegalStoreType, storeType)
	}

	// get the secret of vault path
	var resolvedSecretPath, err = resolved(inputSecret.Annotations["avp.kubernetes.io/path"])
	if err != nil {
		return nil, err
	}
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

func secretCommonVerify(inputSecret internalSecret) error {
	if inputSecret.Annotations == nil {
		return fmt.Errorf(ErrCommonEmptyAnnotations, inputSecret.Name)
	}
	if inputSecret.Annotations["avp.kubernetes.io/path"] == "" {
		return fmt.Errorf(ErrCommonNotFoundAVPPath, inputSecret.Name)
	}
	if len(inputSecret.Data) != 0 && len(inputSecret.StringData) != 0 {
		return fmt.Errorf(ErrCommonNotAcceptBothSecretDataAndData, inputSecret.Name)
	}
	if len(inputSecret.Data) == 0 && len(inputSecret.StringData) == 0 {
		return fmt.Errorf(ErrCommonNotAcceptNeitherSecretDataAndData, inputSecret.Name)
	}

	var foundAngleBracketsData = false
	for _, value := range inputSecret.Data {
		if captureFromFile.MatchString(value) {
			foundAngleBracketsData = true
			break
		}
	}

	var foundAngleBracketsStringData = false
	for _, value := range inputSecret.StringData {
		if captureFromFile.MatchString(value) {
			foundAngleBracketsStringData = true
			break
		}
	}

	if !foundAngleBracketsData && !foundAngleBracketsStringData {
		return fmt.Errorf(ErrCommonNotIncludeAngleBrackets, inputSecret.Name)
	}

	return nil
}
