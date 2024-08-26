package converter

import (
	"fmt"
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generateEsByTLS(inputSecret UnstructuredSecret, storeType, storeName string) (*esv1beta1.ExternalSecret, error) {
	if len(inputSecret.StringData) != 0 {
		return nil, fmt.Errorf(ErrorTLSNotAllowDataField, inputSecret.Name)
	}

	// get the secret of vault path
	var secretPath = inputSecret.Annotations["avp.kubernetes.io/path"]
	var resolvedSecretPath = resolved(secretPath)

	// get the vault secret key
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

	// for specific secret opaque sub-type
	var externalSecretData []esv1beta1.ExternalSecretData
	for _, pemContent := range inputSecret.Data {
		propertyFromSecretData := captureFromFile.FindStringSubmatch(pemContent)
		if len(propertyFromSecretData) == 0 {
			continue
		}
		externalSecretData = append(externalSecretData, esv1beta1.ExternalSecretData{
			SecretKey: propertyFromSecretData[1],
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
			Labels:      inputSecret.Labels,
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
}
