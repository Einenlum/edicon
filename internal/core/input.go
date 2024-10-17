package core

import (
	"regexp"
	"strings"
)

func DecomposeKey(notationStyle NotationStyle, key string) []string {
	if notationStyle == DotNotation {
		return DecomposeKeyWithDotNotation(key)
	}

	return DecomposeKeyWithBracketNotation(key)
}

func DecomposeKeyWithBracketNotation(key string) []string {
	re := regexp.MustCompile(`([^\[\]]+)|\[([^\[\]]+)\]`)
	matches := re.FindAllStringSubmatch(key, -1)

	var result []string
	for _, match := range matches {
		// Add the first part (outside the brackets) or any part inside the brackets
		if match[1] != "" {
			result = append(result, match[1])
		} else if match[2] != "" {
			result = append(result, match[2])
		}
	}

	return result
}

func DecomposeKeyWithDotNotation(key string) []string {
	if !strings.Contains(key, ".") {
		return []string{key}
	}

	keyParts := strings.Split(key, ".")

	return keyParts
}
