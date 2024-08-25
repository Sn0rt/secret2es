package converter

import (
	"fmt"
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func generateBasicAuthSecret(inputSecret UnstructuredSecret, storeType, storeName string) (*esv1beta1.ExternalSecret, error) {
	if len(inputSecret.Data) != 0 {
		return nil, fmt.Errorf(NotBasicAuthWithData, inputSecret.Name)
	}
	if inputSecret.StringData[corev1.BasicAuthUsernameKey] == "" {
		return nil, fmt.Errorf(NotBasicAuthWithEmptyUsername, inputSecret.Name)
	}
	if inputSecret.StringData[corev1.BasicAuthPasswordKey] == "" {
		return nil, fmt.Errorf(NotBasicAuthWithEmptyPassword, inputSecret.Name)
	}

	// resolve the secret path
	var secretPath = inputSecret.Annotations["avp.kubernetes.io/path"]
	var resolvedSecretPath = resolved(secretPath)

	// get vault secret key
	var vaultSecretKey, err = getVaultSecretKey(resolvedSecretPath)
	if err != nil {
		return nil, fmt.Errorf(illegalVaultPath, resolvedSecretPath)
	}

	// new resolvedAnnotations
	var resolvedAnnotations = make(map[string]string)
	for annK, annV := range inputSecret.Annotations {
		if annK != "avp.kubernetes.io/path" {
			resolvedAnnotations[annK] = annV
		}
	}
	resolvedAnnotations["avp.kubernetes.io/path"] = resolvedSecretPath

	var externalSecretData []esv1beta1.ExternalSecretData
	var templateData = map[string]string{}

	for fileName, fileContent := range inputSecret.StringData {
		propertyFromSecretData := captureFromFile.FindAllSubmatch([]byte(fileContent), -1)
		for _, s := range propertyFromSecretData {
			output := strings.TrimSpace(fmt.Sprintf("%s", s[1]))
			if !contains(externalSecretData, output) {
				externalSecretData = append(externalSecretData, esv1beta1.ExternalSecretData{
					SecretKey: fmt.Sprintf("%s", output),
					RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
						Key:      vaultSecretKey,
						Property: output,
					},
				})
			}
		}

		var newFileContentWithout, err = resolveAngleBrackets(fileContent)
		if err != nil {
			return nil, err
		}
		var newFileContent = addQuotesCurlyBraces(newFileContentWithout)
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
					Type: corev1.SecretTypeBasicAuth,
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
