package converter

import (
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
	"strings"
)

const (
	patternExtraction = `^<(.*)>$`
)

var (
	capture = regexp.MustCompile(patternExtraction) // capture the value
)

func generateOpaqueSecret(inputSecret UnstructuredSecret, storeType, storeName string) (*esv1beta1.ExternalSecret, error) {
	var externalSecretDatas []esv1beta1.ExternalSecretData

	// get the secret of vault path
	var secretPath = inputSecret.Annotations["avp.kubernetes.io/path"]
	var resolvedSecretPath = resolved(secretPath)

	var secretPathList = strings.Split(resolvedSecretPath, "/")
	var vaultSecretKey = secretPathList[len(secretPathList)-1]

	for key, value := range inputSecret.Data {
		propertyFromSecretData := capture.FindStringSubmatch(value)
		externalSecretDatas = append(externalSecretDatas, esv1beta1.ExternalSecretData{
			SecretKey: key,
			RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
				Key:      vaultSecretKey,
				Property: propertyFromSecretData[1],
			},
		})
	}

	// new annotations
	var annotations = make(map[string]string)
	for annK, annV := range inputSecret.Annotations {
		if annK != "avp.kubernetes.io/path" {
			annotations[annK] = annV
		}
	}
	annotations["avp.kubernetes.io/path"] = resolvedSecretPath

	return &esv1beta1.ExternalSecret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "external-secrets.io/v1beta1",
			Kind:       "ExternalSecret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        inputSecret.Name,
			Namespace:   inputSecret.Namespace,
			Labels:      inputSecret.ObjectMeta.Labels,
			Annotations: annotations,
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
			Data: externalSecretDatas,
		},
	}, nil
}
