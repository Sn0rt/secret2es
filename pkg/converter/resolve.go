package converter

import (
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
