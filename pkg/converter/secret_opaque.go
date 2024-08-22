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
	opaqueSubType = iota
	opaqueDataType
	opaqueStringDataType
)

func generateOpaqueSecret(inputSecret UnstructuredSecret, storeType, storeName string) (*esv1beta1.ExternalSecret, error) {
	var currentSecretOpaqueSubType int
	if len(inputSecret.Data) != 0 && len(inputSecret.StringData) != 0 {
		return nil, fmt.Errorf(NotSupportedSecretDataBothStringData, inputSecret.Name)
	} else if len(inputSecret.Data) != 0 {
		currentSecretOpaqueSubType = opaqueDataType
	} else {
		currentSecretOpaqueSubType = opaqueStringDataType
	}

	// get the secret of vault path
	var secretPath = inputSecret.Annotations["avp.kubernetes.io/path"]
	var resolvedSecretPath = resolved(secretPath)

	// bugfix: should split with 'data'
	var secretPathList = strings.Split(resolvedSecretPath, "/")
	var vaultSecretKey = secretPathList[len(secretPathList)-1]

	// new resolvedAnnotations
	var resolvedAnnotations = make(map[string]string)
	for annK, annV := range inputSecret.Annotations {
		if annK != "avp.kubernetes.io/path" {
			resolvedAnnotations[annK] = annV
		}
	}
	resolvedAnnotations["avp.kubernetes.io/path"] = resolvedSecretPath

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
				Annotations: resolvedAnnotations,
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
				output := strings.TrimSpace(fmt.Sprintf("%s", s[1]))
				fmt.Println(output)
				externalSecretData = append(externalSecretData, esv1beta1.ExternalSecretData{
					SecretKey: fmt.Sprintf("%s", output),
					RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
						Key:      vaultSecretKey,
						Property: output,
					},
				})
			}

			var replaceFunc = func(fileContent string) string {
				return captureFromFile.ReplaceAllStringFunc(fileContent, func(s string) string {
					return fmt.Sprintf("\"{{ .%s }}\"", strings.TrimSpace(s[1:len(s)-1]))
				})
			}
			var newFileContent = captureFromFile.ReplaceAllStringFunc(fileContent, replaceFunc)
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
				Annotations: resolvedAnnotations,
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
							Annotations: resolvedAnnotations,
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
