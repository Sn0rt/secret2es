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

func resolveAngleBrackets(s string) (string, error) {
	var result strings.Builder
	var temp strings.Builder
	inBracket := false
	lastCharWasSpace := false

	for _, char := range s {
		switch char {
		case '<':
			if inBracket {
				return s, fmt.Errorf(FileContentAngleBracketsParseSyntaxError, `nested or unclosed '<'`)
			}
			inBracket = true
			result.WriteString("{{ .")
			lastCharWasSpace = false
		case '>':
			if !inBracket {
				return s, fmt.Errorf(FileContentAngleBracketsParseSyntaxError, `unpaired '>'`)
			}
			inBracket = false
			trimmedVariable := strings.TrimSpace(temp.String())
			result.WriteString(trimmedVariable)
			result.WriteString(" }}")
			temp.Reset()
			lastCharWasSpace = false
		case ' ':
			if inBracket {
				temp.WriteRune(char)
			} else if !lastCharWasSpace {
				result.WriteRune(char)
				lastCharWasSpace = true
			}
		default:
			if inBracket {
				temp.WriteRune(char)
			} else {
				result.WriteRune(char)
			}
			lastCharWasSpace = false
		}
	}

	if inBracket {
		return s, fmt.Errorf(FileContentAngleBracketsParseSyntaxError, `syntax error: unclosed '<'`)
	}

	return result.String(), nil
}

func addQuotesCurlyBraces(s string) string {
	re := regexp.MustCompile(`(?m)^.*\{\{.*\}\}.*$`)

	result := re.ReplaceAllStringFunc(s, func(line string) string {
		if strings.HasPrefix(line, `"`) && strings.HasSuffix(line, `"`) {
			return line
		}
		return `"` + line + `"`
	})

	return result
}
