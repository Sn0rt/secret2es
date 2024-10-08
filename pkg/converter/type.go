package converter

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"regexp"
	"runtime/debug"
	"sigs.k8s.io/yaml"
	"strings"
	"time"
)

const (
	SecretStoreType        = "SecretStore"
	ClusterSecretStoreType = "ClusterSecretStore"
)

var (
	stopRefreshInterval = &metav1.Duration{Duration: time.Second * 0}
)

type internalSecret struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Immutable, if set to true, ensures that data stored in the Secret cannot
	// be updated (only object metadata can be modified).
	// If not set to true, the field can be modified at any time.
	// Defaulted to nil.
	// +optional
	Immutable *bool `json:"immutable,omitempty" protobuf:"varint,5,opt,name=immutable"`

	// Data contains the secret data. Each key must consist of alphanumeric
	// characters, '-', '_' or '.'. The serialized form of the secret data is a
	// base64 encoded string, representing the arbitrary (possibly non-string)
	// data value here. Described in https://tools.ietf.org/html/rfc4648#section-4
	// +optional
	Data map[string]string `json:"data,omitempty" protobuf:"bytes,2,rep,name=data"`

	// stringData allows specifying non-binary secret data in string form.
	// It is provided as a write-only inputSecret field for convenience.
	// All keys and values are merged into the data field on write, overwriting any existing values.
	// The stringData field is never output when reading from the API.
	// +k8s:conversion-gen=false
	// +optional
	StringData map[string]string `json:"stringData,omitempty" protobuf:"bytes,4,rep,name=stringData"`

	// Used to facilitate programmatic handling of secret data.
	// More info: https://kubernetes.io/docs/concepts/configuration/secret/#secret-types
	// +optional
	Type corev1.SecretType `json:"type,omitempty" protobuf:"bytes,3,opt,name=type,casttype=SecretType"`
}

func splitYAMLDocuments(fileBody []byte) []string {
	re := regexp.MustCompile(`(?m)^---$`)
	fileBodyWithoutCommented := processCommented(fileBody)
	potentialDocs := re.Split(string(fileBodyWithoutCommented), -1)
	return potentialDocs
}

func parseUnstructuredSecret(body []byte) ([]internalSecret, error) {
	defer func() {
		if r := recover(); r != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Panic occurred: %v\n", r)
			_, _ = fmt.Fprintf(os.Stderr, "Body content:\n%s\n", string(body))
			_, _ = fmt.Fprintf(os.Stderr, "Stack trace:\n%s\n", debug.Stack())
		}
	}()

	var secrets []internalSecret
	for _, yamlContent := range splitYAMLDocuments(body) {
		if !strings.Contains(yamlContent, "kind: Secret") {
			continue
		}
		inputSecret := &internalSecret{}
		if err := yaml.Unmarshal([]byte(yamlContent), &inputSecret); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "yaml content: %s\n", yamlContent)
			return nil, fmt.Errorf("error unmarshalling inputSecret secret: %w", err)
		}
		secrets = append(secrets, *inputSecret)
	}
	return secrets, nil
}
