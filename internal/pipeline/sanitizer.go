package pipeline

import (
	"regexp"
	"strings"
)

func SanitizeText(raw string) string {
	// STEP A: Remove HTML comments
	reHTML := regexp.MustCompile(`(?s)<!--.*?-->`)
	clean := reHTML.ReplaceAllString(raw, "")

	// STEP B: Remove unchecked checkboxes from GitHub PRs
	reUnchecked := regexp.MustCompile(`(?m)^\s*-\s*\[\s*\]\s*.*$\n?`)
	clean = reUnchecked.ReplaceAllString(clean, "")

	// STEP C: Collapse multiple horizontal spaces to a single space
	reMultispace := regexp.MustCompile(`[ \t]+`)
	clean = reMultispace.ReplaceAllString(clean, " ")

	// STEP D: Normalize multiple consecutive newlines (maximum 2)
	reNewlines := regexp.MustCompile(`\n{3,}`)
	clean = reNewlines.ReplaceAllString(clean, "\n\n")

	// STEP E: Trim leading and trailing whitespace
	clean = strings.TrimSpace(clean)
	return clean
}