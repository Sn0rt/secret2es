package converter

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	captureFromFileNew = regexp.MustCompile(`<([^<>]+)>|<(<[^>]*>.*?)>`)
	captureFromFile    = regexp.MustCompile(`<([^<>]+)>`)
)

const (
	opaqueDataType = iota
	opaqueStringDataType
)

func generateEsByOpaqueSecret(inputSecret *internalSecret, storeType, storeName string,
	creationPolicy esv1beta1.ExternalSecretCreationPolicy, isResolve bool) (*esv1beta1.ExternalSecret, error) {
	var currentSecretOpaqueSubType int
	if len(inputSecret.Data) != 0 {
		currentSecretOpaqueSubType = opaqueDataType
	} else {
		currentSecretOpaqueSubType = opaqueStringDataType
	}

	// get the vault secret key
	vaultSecretKey, err := getVaultSecretKey(inputSecret.Annotations["avp.kubernetes.io/path"])
	if err != nil {
		return nil, fmt.Errorf(illegalVaultPath, resolvedValueFromEnv)
	}

	// for specific secret opaque sub-type
	var externalSecretData []esv1beta1.ExternalSecretData
	var templateData = make(map[string]string)

	switch currentSecretOpaqueSubType {
	case opaqueDataType:
		// 1. resolve the <% KEY %> from ENV
		if isResolve {
			if err := resolveSecret(inputSecret); err != nil {
				return nil, err
			}
		}

		// static value add to template directly
		// dynamic value add to externalSecretData
		for key, value := range inputSecret.Data {
			propertyFromSecretData := captureFromFileNew.FindAllStringSubmatch(value, -1)
			if len(propertyFromSecretData) == 0 {
				if IsBase64(value) {
					templateData[key] = fmt.Sprintf(`{{ "%s" | b64dec }}`, value)
				} else {
					templateData[key] = value
				}
				continue
			}

			if len(propertyFromSecretData) != 1 {
				return nil, fmt.Errorf(ErrCommonNotSupportMultipleValue, inputSecret.Name)
			}

			if strings.HasPrefix(propertyFromSecretData[0][0], "<%") &&
				strings.HasSuffix(propertyFromSecretData[0][0], "%>") {
				templateData[key] = propertyFromSecretData[0][0]
				continue
			}

			var propertyName string
			if propertyFromSecretData[0][1] != "" {
				propertyName = propertyFromSecretData[0][1]
			} else {
				propertyName = propertyFromSecretData[0][2]
			}

			if !contains(externalSecretData, propertyName) {
				externalSecretData = append(externalSecretData, esv1beta1.ExternalSecretData{
					SecretKey: propertyName,
					RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
						ConversionStrategy: esv1beta1.ExternalSecretConversionDefault,
						DecodingStrategy:   esv1beta1.ExternalSecretDecodeBase64,
						MetadataPolicy:     esv1beta1.ExternalSecretMetadataPolicyNone,
						Key:                vaultSecretKey,
						Property:           propertyName,
					},
				})
			}

			newFileContentWithoutQuote, err := resolveAngleBrackets(value)
			if err != nil {
				return nil, err
			}
			var newFileContent = addQuotesCurlyBraces(newFileContentWithoutQuote)
			templateData[key] = newFileContent
		}
	case opaqueStringDataType:
		// 1. resolve the <% KEY %> from ENV
		if isResolve {
			if err := resolveSecret(inputSecret); err != nil {
				return nil, err
			}
		}

		// 2. process the secret key from file content
		for fileName, fileContent := range inputSecret.StringData {
			propertyFromSecretData := captureFromFileNew.FindAllStringSubmatch(fileContent, -1)
			// simple case, no need to resolve
			if len(propertyFromSecretData) == 0 {
				templateData[fileName] = fileContent
				continue
			}

			// prepare the resolved file content
			resolvedFileContent := fileContent

			// resolve the secret key from file content
			for idx, _ := range propertyFromSecretData {
				if strings.HasPrefix(propertyFromSecretData[idx][0], "<%") &&
					strings.HasSuffix(propertyFromSecretData[idx][0], "%>") {
					continue
				}

				var propertyName string
				if propertyFromSecretData[idx][1] != "" {
					propertyName = propertyFromSecretData[idx][1]
				} else {
					propertyName = propertyFromSecretData[idx][2]
				}

				// if secret key not found in externalSecretData then append to slice
				if !contains(externalSecretData, propertyName) {
					externalSecretData = append(externalSecretData, esv1beta1.ExternalSecretData{
						SecretKey: strings.TrimSpace(propertyName),
						RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
							ConversionStrategy: esv1beta1.ExternalSecretConversionDefault,
							DecodingStrategy:   esv1beta1.ExternalSecretDecodeNone,
							MetadataPolicy:     esv1beta1.ExternalSecretMetadataPolicyNone,
							Key:                vaultSecretKey,
							Property:           strings.TrimSpace(propertyName),
						},
					})
				}
			}

			newFileContentWithoutQuote, err := resolveAngleBrackets(resolvedFileContent)
			if err != nil {
				return nil, err
			}
			if !strings.Contains(newFileContentWithoutQuote, "\n") {
				var newFileContent = addQuotesCurlyBraces(newFileContentWithoutQuote)
				templateData[fileName] = newFileContent
			} else {
				templateData[fileName] = newFileContentWithoutQuote
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
				CreationPolicy: creationPolicy,
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

func resolveSecret(inputSecret *internalSecret) (err error) {
	for fileName, fileContent := range inputSecret.Data {
		propertyFromSecretData := captureFromFile.FindAllStringSubmatch(fileContent, -1)
		// simple case, no need to resolve

		// resolve the secret key from file content
		for idx, _ := range propertyFromSecretData {
			if strings.HasPrefix(propertyFromSecretData[idx][0], "<%") &&
				strings.HasSuffix(propertyFromSecretData[idx][0], "%>") {
				inputSecret.Data[fileName], err = resolved(fileContent)
				if err != nil {
					return err
				}
				continue
			}
		}
	}

	for fileName, fileContent := range inputSecret.StringData {
		propertyFromSecretData := captureFromFile.FindAllStringSubmatch(fileContent, -1)

		// resolve the secret key from file content
		for idx, _ := range propertyFromSecretData {
			// process if match <% ... %>
			if strings.HasPrefix(propertyFromSecretData[idx][0], "<%") &&
				strings.HasSuffix(propertyFromSecretData[idx][0], "%>") {
				inputSecret.StringData[fileName], err = resolved(fileContent)
				if err != nil {
					return err
				}
				continue
			}
		}
	}

	return nil
}

func IsBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

func contains(data []esv1beta1.ExternalSecretData, output string) bool {
	for _, d := range data {
		if d.SecretKey == output {
			return true
		}
	}
	return false
}
