package converter

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	patternResolveFromEnv = `<%\s*(\w+)\s*%>`
	resolvedSecretPath    = regexp.MustCompile(patternResolveFromEnv)
)

func resolved(originalString string) string {
	needResolvedStrings := resolvedSecretPath.FindAllStringSubmatch(originalString, -1)
	for _, match := range needResolvedStrings {
		if len(match) > 1 {
			env := match[1]
			originalString = strings.Replace(originalString, match[0], os.Getenv(env), -1)
		}
	}
	return originalString
}

func getVaultSecretKey(secretPath string) (string, error) {
	parts := strings.Split(secretPath, "/")

	index := -1
	for i, part := range parts {
		if part == "data" {
			index = i
			break
		}
	}

	if index == -1 || index+1 >= len(parts) {
		return "", fmt.Errorf(illegalVaultPath, secretPath)
	}

	result := strings.Join(parts[index+1:], "/")
	return result, nil
}
