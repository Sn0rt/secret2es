package converter

import (
	"fmt"
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"strings"
)

type Auth struct {
	Auth string `json:"auth"`
}

type Auths struct {
	Auths map[string]Auth `json:"auths"`
}

func generateEsByDockerConfigJSON(inputSecret *UnstructuredSecret, storeType, storeName string) (*esv1beta1.ExternalSecret, error) {
	if len(inputSecret.Data) != 0 {
		return nil, fmt.Errorf(ErrDockerConfigJsonAcceptOnlyDataFields, inputSecret.Name)
	}
	if len(inputSecret.StringData) != 1 {
		return nil, fmt.Errorf(ErrDockerConfigJsonAcceptOnlyOneValue, inputSecret.Name)
	}

	// get the vault secret key
	vaultSecretKey, err := getVaultSecretKey(inputSecret.Annotations["avp.kubernetes.io/path"])
	if err != nil {
		return nil, fmt.Errorf(illegalVaultPath, resolvedSecretPath)
	}

	authFileContent, err := serializeDockerConfigJSON([]byte(inputSecret.StringData[".dockerconfigjson"]))
	if err != nil {
		return nil, err
	}

	// prepare the ref of sensitive data
	var externalSecretData []esv1beta1.ExternalSecretData
	for _, loginInfo := range authFileContent.Auths {
		propertyFromSecretData := captureFromFile.FindAllSubmatch([]byte(loginInfo.Auth), -1)
		if len(propertyFromSecretData) == 0 {
			continue
		}
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
	}

	// render template
	templateData := make(map[string]string)
	var dockerloginfo = Auths{
		Auths: make(map[string]Auth),
	}
	for key, value := range authFileContent.Auths {
		var singleLoginfo = Auth{}
		var t, _ = resolveAngleBrackets(value.Auth)
		singleLoginfo.Auth = t
		dockerloginfo.Auths[key] = singleLoginfo
	}
	var out, _ = json.Marshal(&dockerloginfo)
	templateData[".dockerconfigjson"] = string(out)

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
			RefreshInterval: stopRefreshInterval,
			SecretStoreRef: esv1beta1.SecretStoreRef{
				Name: storeName,
				Kind: storeType,
			},
			Target: esv1beta1.ExternalSecretTarget{
				Name:           inputSecret.Name,
				CreationPolicy: esv1beta1.CreatePolicyMerge,
				DeletionPolicy: esv1beta1.DeletionPolicyRetain,
				Template: &esv1beta1.ExternalSecretTemplate{
					Type: corev1.SecretTypeDockerConfigJson,
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

func serializeDockerConfigJSON(dockerConfigJson []byte) (*Auths, error) {
	var dockerConfigJSON Auths
	if err := json.Unmarshal(dockerConfigJson, &dockerConfigJSON); err != nil {
		return nil, err
	}
	return &dockerConfigJSON, nil
}
