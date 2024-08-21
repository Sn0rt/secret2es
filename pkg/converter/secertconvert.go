package converter

import (
	"fmt"
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"sigs.k8s.io/yaml"
)

const (
	patternVerifyOfAVP = `^<.*>$`
	patternExtraction  = `^<(.*)>$`
)

// ConvertSecret converts a Kubernetes Secret to an ExternalSecret
func ConvertSecret(inputFile, storeType, storeName, namespace, secretName string, verbose bool) error {
	bytes, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	inputSecretList, err := parseUnstructuredSecret(bytes)
	if err != nil {
		return fmt.Errorf("error parsing input secret: %w", err)
	}

	for _, inputSecret := range inputSecretList {
		externalSecret, err := convertSecret2ExtSecret(inputSecret, storeType, storeName, namespace, secretName)
		if err != nil {
			return fmt.Errorf("error converting secret to external secret: %s", err.Error())
		}

		yamlData, err := yaml.Marshal(externalSecret)
		if err != nil {
			return fmt.Errorf("error encoding external secret: %w", err)
		}
		fmt.Printf("%s\n", yamlData)
	}

	return nil
}

func convertSecret2ExtSecret(inputSecret UnstructuredSecret, storeType, storeName, namespace, secretName string) (*esv1beta1.ExternalSecret, error) {
	if inputSecret.Annotations == nil {
		return nil, fmt.Errorf(NotEmptyAnnotations, inputSecret.Name)
	} else {
		if inputSecret.Annotations["avp.kubernetes.io/path"] == "" {
			return nil, fmt.Errorf(NotSetAVPAnnotations, inputSecret.Name)
		}
	}

	externalSecret := &esv1beta1.ExternalSecret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "external-secrets.io/v1beta1",
			Kind:       "ExternalSecret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      inputSecret.Name,
			Namespace: inputSecret.Namespace,
			Labels:    inputSecret.ObjectMeta.Labels,
		},
		Spec: esv1beta1.ExternalSecretSpec{
			SecretStoreRef: esv1beta1.SecretStoreRef{
				Name: storeName,
				Kind: storeType,
			},
			Target: esv1beta1.ExternalSecretTarget{
				Name:           secretName,
				CreationPolicy: esv1beta1.CreatePolicyMerge,
				DeletionPolicy: esv1beta1.DeletionPolicyRetain,
			},
			Data: []esv1beta1.ExternalSecretData{},
		},
	}

	switch inputSecret.Type {
	case corev1.SecretTypeOpaque:
		return generateOpaqueSecret(inputSecret, externalSecret)
	case corev1.SecretTypeDockerConfigJson:
	case corev1.SecretTypeBasicAuth:
	case corev1.SecretTypeTLS:
	default:
		return nil, fmt.Errorf(NotImplSecretType, inputSecret.Type, inputSecret.Name)
	}

	return externalSecret, nil
}
