package converter

import (
	"fmt"
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
	"strings"
)

var (
	captureFromFile = regexp.MustCompile(`<([^<>]+)>`)
)

const (
	opaqueDataType = iota
	opaqueStringDataType
)

func generateEsByOpaqueSecret(inputSecret *UnstructuredSecret, storeType, storeName string) (*esv1beta1.ExternalSecret, error) {
	var currentSecretOpaqueSubType int
	if len(inputSecret.Data) != 0 {
		currentSecretOpaqueSubType = opaqueDataType
	} else {
		currentSecretOpaqueSubType = opaqueStringDataType
	}

	// get the vault secret key
	var vaultSecretKey, err = getVaultSecretKey(inputSecret.Annotations["avp.kubernetes.io/path"])
	if err != nil {
		return nil, fmt.Errorf(illegalVaultPath, resolvedSecretPath)
	}

	// for specific secret opaque sub-type
	switch currentSecretOpaqueSubType {
	case opaqueDataType:
		var externalSecretData []esv1beta1.ExternalSecretData
		for key, value := range inputSecret.Data {
			propertyFromSecretData := captureFromFile.FindStringSubmatch(value)
			externalSecretData = append(externalSecretData, esv1beta1.ExternalSecretData{
				SecretKey: key,
				RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
					Key:      vaultSecretKey,
					Property: propertyFromSecretData[1],
				},
			})
		}

		return &esv1beta1.ExternalSecret{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "external-secrets.io/v1beta1",
				Kind:       "ExternalSecret",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:        inputSecret.Name,
				Namespace:   inputSecret.Namespace,
				Labels:      inputSecret.ObjectMeta.Labels,
				Annotations: inputSecret.Annotations,
			},
			Spec: esv1beta1.ExternalSecretSpec{
				SecretStoreRef: esv1beta1.SecretStoreRef{
					Name: storeName,
					Kind: storeType,
				},
				Target: esv1beta1.ExternalSecretTarget{
					Name:           inputSecret.Name,
					CreationPolicy: esv1beta1.CreatePolicyMerge,
					DeletionPolicy: esv1beta1.DeletionPolicyRetain,
				},
				Data: externalSecretData,
			},
		}, nil
	case opaqueStringDataType:
		var externalSecretData []esv1beta1.ExternalSecretData
		var templateData = map[string]string{}

		for fileName, fileContent := range inputSecret.StringData {
			propertyFromSecretData := captureFromFile.FindAllSubmatch([]byte(fileContent), -1)
			for _, s := range propertyFromSecretData {
				output := strings.TrimSpace(string(s[1]))
				// if secret key not found in externalSecretData then append to slice
				if !contains(externalSecretData, output) {
					externalSecretData = append(externalSecretData, esv1beta1.ExternalSecretData{
						SecretKey: output,
						RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
							Key:      vaultSecretKey,
							Property: output,
						},
					})
				}
			}
			newFileContentWithoutQuote, err := resolveAngleBrackets(fileContent)
			if err != nil {
				return nil, err
			}
			var newFileContent = addQuotesCurlyBraces(newFileContentWithoutQuote)
			templateData[fileName] = newFileContent
		}

		return &esv1beta1.ExternalSecret{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "external-secrets.io/v1beta1",
				Kind:       "ExternalSecret",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:        inputSecret.Name,
				Namespace:   inputSecret.Namespace,
				Labels:      inputSecret.ObjectMeta.Labels,
				Annotations: inputSecret.Annotations,
			},
			Spec: esv1beta1.ExternalSecretSpec{
				SecretStoreRef: esv1beta1.SecretStoreRef{
					Name: storeName,
					Kind: storeType,
				},
				Target: esv1beta1.ExternalSecretTarget{
					Name:           inputSecret.Name,
					CreationPolicy: esv1beta1.CreatePolicyMerge,
					DeletionPolicy: esv1beta1.DeletionPolicyRetain,
					Template: &esv1beta1.ExternalSecretTemplate{
						Type: corev1.SecretTypeOpaque,
						Metadata: esv1beta1.ExternalSecretTemplateMetadata{
							Annotations: inputSecret.Annotations,
							Labels:      inputSecret.ObjectMeta.Labels,
						},
						MergePolicy: esv1beta1.MergePolicyMerge,
						Data:        templateData,
					},
				},
				Data: externalSecretData,
			},
		}, nil
	}

	return nil, fmt.Errorf("error converting secret to external secret: %s", NotSupportedSecretData)
}

func contains(data []esv1beta1.ExternalSecretData, output string) bool {
	for _, d := range data {
		if d.SecretKey == output {
			return true
		}
	}
	return false
}
