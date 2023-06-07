package preprocessor

import (
	"regexp"
	"strconv"
	"strings"
)

func CleanTemplateParam(paramName string) string {
	if IsTemplateReplaceParam(paramName) {
		// Strip prefix and suffix
		paramName = paramName[2 : len(paramName)-2]
		return paramName
	}

	return paramName
}

func IsTemplateReplaceParam(paramName string) bool {
	return strings.HasPrefix(paramName, "{{") && strings.HasSuffix(paramName, "}}")
}

func ParseTemplateTokens(template string, replaceIdx map[string]string) string {
	// Regular expression to match {{var:length}} format
	re := regexp.MustCompile(`{{(\w+):(\d+)}}`)

	// Replace function to replace matched string with variable value
	replaceFunc := func(s string) string {
		// Parse matched string
		matches := re.FindStringSubmatch(s)
		varName := matches[1]
		length := matches[2]

		// Get variable value
		varValue := replaceIdx[varName]

		// Truncate variable value if necessary
		lengthNum, _ := strconv.Atoi(length)
		if len(varValue) > lengthNum {
			varValue = varValue[:lengthNum]
		}

		return varValue
	}

	// Replace all matches in pattern
	result := re.ReplaceAllStringFunc(template, replaceFunc)

	return result
}
