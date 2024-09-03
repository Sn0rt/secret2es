package converter

import (
	"fmt"
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generateEsByTLS(inputSecret *internalSecret, storeType, storeName string) (*esv1beta1.ExternalSecret, error) {
	if len(inputSecret.StringData) != 0 {
		return nil, fmt.Errorf(ErrTLSNotAllowDataField, inputSecret.Name)
	}

	// get the vault secret key
	var vaultSecretKey, err = getVaultSecretKey(inputSecret.Annotations["avp.kubernetes.io/path"])
	if err != nil {
		return nil, fmt.Errorf(illegalVaultPath, resolvedValueFromEnv)
	}

	// for specific secret opaque sub-type
	var externalSecretData []esv1beta1.ExternalSecretData
	for fileName, pemContent := range inputSecret.Data {
		propertyFromSecretData := captureFromFile.FindStringSubmatch(pemContent)
		if len(propertyFromSecretData) == 0 {
			continue
		}
		externalSecretData = append(externalSecretData, esv1beta1.ExternalSecretData{
			SecretKey: fileName,
			RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
				ConversionStrategy: esv1beta1.ExternalSecretConversionDefault,
				DecodingStrategy:   esv1beta1.ExternalSecretDecodeNone,
				MetadataPolicy:     esv1beta1.ExternalSecretMetadataPolicyNone,
				Key:                vaultSecretKey,
				Property:           propertyFromSecretData[1],
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
			Annotations: inputSecret.Annotations,
		},
		Spec: esv1beta1.ExternalSecretSpec{
			RefreshInterval: stopRefreshInterval,
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
