package converter

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
)

var (
	patternResolveFromEnv = `<%\s*(\w+)\s*%>`
	resolvedValueFromEnv  = regexp.MustCompile(patternResolveFromEnv)
)

func resolved(originalString string, enable bool) (string, error) {
	if !enable {
		return originalString, nil
	}
	needResolvedStrings := resolvedValueFromEnv.FindAllStringSubmatch(originalString, -1)
	for _, match := range needResolvedStrings {
		if len(match) > 1 {
			env := match[1]
			envResolved := os.Getenv(env)
			if envResolved == "" {
				return "", fmt.Errorf(ErrCommonNotSetEnv, env)
			}
			originalString = strings.Replace(originalString, match[0], envResolved, -1)
		}
	}
	return originalString, nil
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
	inPercentBracket := false

	for i := 0; i < len(s); i++ {
		char := rune(s[i])

		// Check for <% ... %> pattern and leave it unmodified
		if char == '<' && i+1 < len(s) && s[i+1] == '%' {
			inPercentBracket = true
			result.WriteRune(char)
			result.WriteRune('%')
			i++ // Skip the next '%'
			continue
		}
		if inPercentBracket {
			// Keep writing until we find %>
			if char == '%' && i+1 < len(s) && s[i+1] == '>' {
				result.WriteRune('%')
				result.WriteRune('>')
				i++ // Skip the next '>'
				inPercentBracket = false
				continue
			}
			result.WriteRune(char)
			continue
		}

		switch char {
		case '<':
			if inBracket {
				return s, fmt.Errorf(FileContentAngleBracketsParseSyntaxError, `nested or unclosed '<'`)
			}
			inBracket = true
			result.WriteString("{{ .")
		case '>':
			if !inBracket {
				return s, fmt.Errorf(FileContentAngleBracketsParseSyntaxError, `unpaired '>'`)
			}
			inBracket = false
			trimmedVariable := strings.TrimSpace(temp.String())
			result.WriteString(trimmedVariable)
			result.WriteString(" }}")
			temp.Reset()
		case ' ', '\n', '\r': // Handle spaces and newlines
			if inBracket {
				temp.WriteRune(char)
			} else {
				result.WriteRune(char)
			}
		default:
			if inBracket {
				temp.WriteRune(char)
			} else {
				result.WriteRune(char)
			}
		}
	}

	if inBracket {
		return s, fmt.Errorf(FileContentAngleBracketsParseSyntaxError, `unclosed '<'`)
	}

	return result.String(), nil
}

func addQuotesCurlyBraces(input string) string {
	var result []string
	inCurlyBraces := false

	lines := strings.Split(input, "\n")

	for lineIndex, line := range lines {
		processedLine := processLine(line, &inCurlyBraces)

		if lineIndex > 0 {
			result = append(result, "\n")
		}

		result = append(result, processedLine)
	}

	return strings.Join(result, "")
}

func processLine(line string, inCurlyBraces *bool) string {
	var result []string
	var currentWord strings.Builder
	var leadingSpaces strings.Builder
	processingLeadingSpaces := true
	inTemplateBraces := false // 新增：用于跟踪 <% %> 模板标记

	for i, char := range line {
		if processingLeadingSpaces && unicode.IsSpace(char) {
			leadingSpaces.WriteRune(char)
			continue
		} else if processingLeadingSpaces {
			processingLeadingSpaces = false
			result = append(result, leadingSpaces.String())
		}

		// Check if it is in curly brackets
		if char == '{' && i > 0 && line[i-1] == '{' {
			*inCurlyBraces = true
		} else if char == '}' && i > 0 && line[i-1] == '}' && *inCurlyBraces {
			*inCurlyBraces = false
		}

		// Determines whether it is in the <% %> template tag
		if char == '<' && i < len(line)-1 && line[i+1] == '%' {
			inTemplateBraces = true
		} else if char == '%' && i < len(line)-1 && line[i+1] == '>' {
			inTemplateBraces = false
		}

		if unicode.IsSpace(char) && !*inCurlyBraces && !inTemplateBraces {
			if currentWord.Len() > 0 {
				word := currentWord.String()
				// Check if the whole word is within curly braces or template tags
				if strings.Contains(word, "{{") && strings.Contains(word, "}}") {
					result = append(result, `"`+word+`"`)
				} else {
					result = append(result, word)
				}
				currentWord.Reset()
			}
			result = append(result, string(char))
		} else {
			currentWord.WriteRune(char)
		}
	}

	// Process the last word
	if currentWord.Len() > 0 {
		word := currentWord.String()
		if strings.Contains(word, "{{") && strings.Contains(word, "}}") {
			result = append(result, `"`+word+`"`)
		} else {
			result = append(result, word)
		}
	}

	return strings.Join(result, "")
}

func processCommented(input []byte) []byte {
	var output []byte
	lines := bytes.Split(input, []byte("\n"))

	for _, line := range lines {
		trimmedLine := bytes.TrimLeft(line, " \t")
		if len(trimmedLine) == 0 || trimmedLine[0] != '#' {
			if idx := bytes.IndexByte(line, '#'); idx != -1 {
				line = line[:idx]
			}
			output = append(output, line...)
			output = append(output, '\n')
		}
	}

	if len(output) > 0 && output[len(output)-1] == '\n' {
		output = output[:len(output)-1]
	}

	return output
}
