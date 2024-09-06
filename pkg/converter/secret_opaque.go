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

func generateEsByOpaqueSecret(inputSecret *internalSecret, storeType, storeName string) (*esv1beta1.ExternalSecret, error) {
	var currentSecretOpaqueSubType int
	if len(inputSecret.Data) != 0 {
		currentSecretOpaqueSubType = opaqueDataType
	} else {
		currentSecretOpaqueSubType = opaqueStringDataType
	}

	// get the vault secret key
	var vaultSecretKey, err = getVaultSecretKey(inputSecret.Annotations["avp.kubernetes.io/path"])
	if err != nil {
		return nil, fmt.Errorf(illegalVaultPath, resolvedValueFromEnv)
	}

	// for specific secret opaque sub-type
	var externalSecretData []esv1beta1.ExternalSecretData
	var templateData = make(map[string]string)

	switch currentSecretOpaqueSubType {
	case opaqueDataType:
		// static value add to template directly
		// dynamic value add to externalSecretData
		for key, value := range inputSecret.Data {
			if resolvedValueFromEnv.MatchString(value) {
				resolvedValue, err := resolved(value)
				if err != nil {
					return nil, err
				}
				templateData[key] = resolvedValue
			} else {
				propertyFromSecretData := captureFromFile.FindStringSubmatch(value)
				if len(propertyFromSecretData) == 0 {
					continue
				}
				externalSecretData = append(externalSecretData, esv1beta1.ExternalSecretData{
					SecretKey: propertyFromSecretData[1],
					RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
						ConversionStrategy: esv1beta1.ExternalSecretConversionDefault,
						DecodingStrategy:   esv1beta1.ExternalSecretDecodeAuto,
						MetadataPolicy:     esv1beta1.ExternalSecretMetadataPolicyNone,
						Key:                vaultSecretKey,
						Property:           propertyFromSecretData[1],
					},
				})
				newFileContentWithoutQuote, err := resolveAngleBrackets(value)
				if err != nil {
					return nil, err
				}
				var newFileContent = addQuotesCurlyBraces(newFileContentWithoutQuote)
				templateData[key] = newFileContent
			}
		}
	case opaqueStringDataType:
		for fileName, fileContent := range inputSecret.StringData {
			// should ignore parse <% %> first
			// because it has been resolved by argo-vault-cd
			resolvedFileContent, err := resolved(fileContent)
			if err != nil {
				return nil, err
			}

			// map <>
			propertyFromSecretData := captureFromFile.FindAllSubmatch([]byte(resolvedFileContent), -1)
			if len(propertyFromSecretData) == 0 {
				templateData[fileName] = resolvedFileContent
			} else {
				for _, s := range propertyFromSecretData {
					output := strings.TrimSpace(string(s[1]))
					// if secret key not found in externalSecretData then append to slice
					if !contains(externalSecretData, output) {
						externalSecretData = append(externalSecretData, esv1beta1.ExternalSecretData{
							SecretKey: output,
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								ConversionStrategy: esv1beta1.ExternalSecretConversionDefault,
								DecodingStrategy:   esv1beta1.ExternalSecretDecodeNone,
								MetadataPolicy:     esv1beta1.ExternalSecretMetadataPolicyNone,
								Key:                vaultSecretKey,
								Property:           output,
							},
						})
					}
				}
				newFileContentWithoutQuote, err := resolveAngleBrackets(resolvedFileContent)
				if err != nil {
					return nil, err
				}
				var newFileContent = addQuotesCurlyBraces(newFileContentWithoutQuote)
				templateData[fileName] = newFileContent
			}
		}
	}

	if len(externalSecretData) == 0 {
		return nil, fmt.Errorf(ErrCommonNotNeedRefData, inputSecret.Name)
	}

	return &esv1beta1.ExternalSecret{
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
			RefreshInterval: stopRefreshInterval,
			SecretStoreRef: esv1beta1.SecretStoreRef{
				Name: storeName,
				Kind: storeType,
			},
			Target: esv1beta1.ExternalSecretTarget{
				Name:           inputSecret.Name,
				CreationPolicy: esv1beta1.CreatePolicyOrphan,
				DeletionPolicy: esv1beta1.DeletionPolicyRetain,
				Template: &esv1beta1.ExternalSecretTemplate{
					Type: corev1.SecretTypeOpaque,
					Metadata: esv1beta1.ExternalSecretTemplateMetadata{
						Labels: inputSecret.ObjectMeta.Labels,
					},
					MergePolicy: esv1beta1.MergePolicyReplace,
					Data:        templateData,
				},
			},
			Data: externalSecretData,
		},
	}, nil
}

func contains(data []esv1beta1.ExternalSecretData, output string) bool {
	for _, d := range data {
		if d.SecretKey == output {
			return true
		}
	}
	return false
}
