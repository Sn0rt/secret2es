package converter

import (
	"fmt"
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"os"
	"sigs.k8s.io/yaml"
	"strings"
)

// ConvertSecret converts a AVP Secret to an ExternalSecret for CLI
func ConvertSecret(inputFile, storeType, storeName string, creationPolicy esv1beta1.ExternalSecretCreationPolicy, resolve bool) error {
	bytes, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading inputSecret file: %w", err)
	}

	output, warn, err := ConvertSecretContent(bytes, storeType, storeName, creationPolicy, resolve, nil)
	if err != nil {
		return fmt.Errorf("error converting secret: %w", err)
	}
	if warn != "" {
		_, _ = fmt.Fprintf(os.Stderr, "warn: %s", warn)
	}
	if output != "" {
		fmt.Println(output)
	}

	return nil
}

func ConvertSecretContent(input []byte, storeType, storeName string,
	creationPolicy esv1beta1.ExternalSecretCreationPolicy,
	resolve bool,
	EnvVars map[string]string) (string, string, error) {
	output := ""
	warn := ""

	if resolve && len(EnvVars) > 0 {
		for key, value := range EnvVars {
			_ = os.Setenv(key, value)
		}
	}

	inputSecretList, err := parseUnstructuredSecret(input)
	if err != nil {
		return "", "", fmt.Errorf("error parsing inputSecret secret: %w", err)
	}

	for _, inputSecret := range inputSecretList {
		externalSecret, err := convertSecret2ExtSecret(inputSecret, storeType, storeName, creationPolicy, resolve)
		if err != nil {
			switch err.Error() {
			case fmt.Errorf(ErrCommonNotIncludeAngleBrackets, inputSecret.Name).Error():
				warn += fmt.Sprintf("Error: %v\n", err)
				continue
			case fmt.Errorf(ErrCommonEmptyAnnotations, inputSecret.Name).Error():
				warn += fmt.Sprintf("Error: %v\n", err)
				continue
			}
			return "", "", fmt.Errorf("error converting secret to external secret: %s", err.Error())
		}
		yamlData, err := yaml.Marshal(externalSecret)
		if err != nil {
			return "", "", fmt.Errorf("error encoding external secret: %w", err)
		}
		output += fmt.Sprintf("---\n%s", postProcessOutputES(yamlData))
	}

	return output, warn, nil
}

func postProcessOutputES(yamlData []byte) string {
	var externalSecret map[string]interface{}
	if err := yaml.Unmarshal(yamlData, &externalSecret); err != nil {
		return string(yamlData)
	}

	// delete status field
	delete(externalSecret, "status")

	// delete metadata.creationTimestamp of yamlData
	if metadata, ok := externalSecret["metadata"].(map[string]interface{}); ok {
		delete(metadata, "creationTimestamp")
		externalSecret["metadata"] = metadata
	}

	// delete .spec.target.template.metadata
	if spec, ok := externalSecret["spec"].(map[string]interface{}); ok {
		if target, ok := spec["target"].(map[string]interface{}); ok {
			if template, ok := target["template"].(map[string]interface{}); ok {
				delete(template, "metadata")
			}
		}
	}

	var needReplace = false
	if spec, ok := externalSecret["spec"].(map[string]interface{}); ok {
		if target, ok := spec["target"].(map[string]interface{}); ok {
			if template, ok := target["template"].(map[string]interface{}); ok {
				if data, ok := template["data"].(map[string]interface{}); ok {
					for _, value := range data {
						if valString, ok := value.(string); ok {
							// if start with "{{ and end with }}" set a needReplace value
							if strings.HasPrefix(valString, `"{{`) && strings.HasSuffix(valString, `}}"`) {
								needReplace = true
							}

							// "<% ENV %>-{{ .TEST_DIST_LINUX }}"
							if strings.HasPrefix(valString, `"<`) && strings.HasSuffix(valString, `}}"`) {
								needReplace = true
							}

							// "{{ .TEST_DIST_LINUX }}-<% ENV %>"
							if strings.HasPrefix(valString, `"{{`) && strings.HasSuffix(valString, `>"`) {
								needReplace = true
							}
						}
					}
				}
			}
		}
	}

	newYamlData, err := yaml.Marshal(externalSecret)
	if err != nil {
		panic(err.Error())
	}
	if needReplace {
		removedLeft := strings.Replace(string(newYamlData), `'"`, `"`, -1)
		return strings.Replace(removedLeft, `"'`, `"`, -1)
	}
	return string(newYamlData)
}

func convertSecret2ExtSecret(inputSecret internalSecret, storeType, storeName string,
	createPolicy esv1beta1.ExternalSecretCreationPolicy, resolve bool) (*esv1beta1.ExternalSecret, error) {
	if err := secretCommonVerify(inputSecret); err != nil {
		return nil, err
	}

	if storeType != SecretStoreType &&
		storeType != ClusterSecretStoreType {
		return nil, fmt.Errorf(illegalStoreType, storeType)
	}

	if createPolicy != esv1beta1.CreatePolicyOwner &&
		createPolicy != esv1beta1.CreatePolicyOrphan &&
		createPolicy != esv1beta1.CreatePolicyMerge {
		return nil, fmt.Errorf(illegalCreatePolicy, createPolicy)
	}

	// get the secret of vault path
	if resolve {
		var resolvedSecretPath, err = resolved(inputSecret.Annotations["avp.kubernetes.io/path"])
		if err != nil {
			return nil, err
		}
		inputSecret.Annotations["avp.kubernetes.io/path"] = resolvedSecretPath
	}

	switch inputSecret.Type {
	case corev1.SecretTypeOpaque:
		return generateEsByOpaqueSecret(&inputSecret, storeType, storeName, createPolicy, resolve)
	case corev1.SecretTypeBasicAuth:
		return generateEsByBasicAuthSecret(&inputSecret, storeType, storeName, createPolicy, resolve)
	case corev1.SecretTypeDockerConfigJson:
		return generateEsByDockerConfigJSON(&inputSecret, storeType, storeName, createPolicy, resolve)
	case corev1.SecretTypeTLS:
		return generateEsByTLS(&inputSecret, storeType, storeName, createPolicy, resolve)
	}

	return nil, fmt.Errorf(NotImplSecretType, inputSecret.Type, inputSecret.Name)
}

func secretCommonVerify(inputSecret internalSecret) error {
	if inputSecret.Annotations == nil {
		return fmt.Errorf(ErrCommonEmptyAnnotations, inputSecret.Name)
	}
	if inputSecret.Annotations["avp.kubernetes.io/path"] == "" {
		return fmt.Errorf(ErrCommonNotFoundAVPPath, inputSecret.Name)
	}
	if len(inputSecret.Data) != 0 && len(inputSecret.StringData) != 0 {
		return fmt.Errorf(ErrCommonNotAcceptBothSecretDataAndData, inputSecret.Name)
	}
	if len(inputSecret.Data) == 0 && len(inputSecret.StringData) == 0 {
		return fmt.Errorf(ErrCommonNotAcceptNeitherSecretDataAndData, inputSecret.Name)
	}

	var foundAngleBracketsData = false
	for _, value := range inputSecret.Data {
		if captureFromFile.MatchString(value) {
			foundAngleBracketsData = true
			break
		}
	}

	var foundAngleBracketsStringData = false
	for _, value := range inputSecret.StringData {
		if captureFromFile.MatchString(value) {
			foundAngleBracketsStringData = true
			break
		}
	}

	if !foundAngleBracketsData && !foundAngleBracketsStringData {
		return fmt.Errorf(ErrCommonNotIncludeAngleBrackets, inputSecret.Name)
	}

	return nil
}
