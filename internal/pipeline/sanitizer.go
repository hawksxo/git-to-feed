package pipeline

import (
	"regexp"
	"strings"
)

func SanitizeText(raw string) string {
	// PASO A: Eliminar comentarios HTML
	reHTML := regexp.MustCompile(`(?s)<!--.*?-->`)
	clean := reHTML.ReplaceAllString(raw, "")

	// PASO B: Eliminar checkboxes sin marcar de GitHub PRs
	reUnchecked := regexp.MustCompile(`(?m)^\s*-\s*\[\s*\]\s*.*$\n?`)
	clean = reUnchecked.ReplaceAllString(clean, "")

	// PASO C: Colapsar espacios múltiples horizontales a uno solo
	reMultispace := regexp.MustCompile(`[ \t]+`)
	clean = reMultispace.ReplaceAllString(clean, " ")

	// PASO D: Normalizar múltiples saltos de línea consecutivos (máximo 2)
	reNewlines := regexp.MustCompile(`\n{3,}`)
	clean = reNewlines.ReplaceAllString(clean, "\n\n")

	// PASO E: Limpiar espacios en blanco al inicio y final
	clean = strings.TrimSpace(clean)
	return clean
}