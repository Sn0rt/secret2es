package converter

import (
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	"path/filepath"
	"regexp"
)

func generateOpaqueSecret(inputSecret UnstructuredSecret, externalSecret *esv1beta1.ExternalSecret) (*esv1beta1.ExternalSecret, error) {
	verify := regexp.MustCompile(patternVerifyOfAVP)
	capture := regexp.MustCompile(patternExtraction)

	for k, v := range inputSecret.StringData {
		if verify.MatchString(v) {
			refString := capture.FindStringSubmatch(v)
			if len(refString) != 1 {
				continue
			}
			externalSecret.Spec.Data = append(externalSecret.Spec.Data, esv1beta1.ExternalSecretData{
				SecretKey: k,
				RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
					Key:      filepath.Join(inputSecret.Namespace, inputSecret.Name, k),
					Property: refString[0],
				},
			})
		}
	}

	for key := range inputSecret.Data {
		externalSecret.Spec.Data = append(externalSecret.Spec.Data, esv1beta1.ExternalSecretData{
			SecretKey: key,
			RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
				Key:      filepath.Join(inputSecret.Namespace, inputSecret.Name, key),
				Property: key,
			},
		})
	}
	return externalSecret, nil
}
