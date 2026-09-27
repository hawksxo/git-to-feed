package pipeline

import (
	"regexp"
	"strings"
)

var cliches = map[string]string{
	"I am thrilled to announce": "Lanzamos",
	"thrilled to announce":      "presentamos",
	"game-changer":              "mejora clave",
	"game changer":              "mejora clave",
	"synergy":                   "integración",
	"paradigm shift":            "cambio de enfoque",
	"unleash the power":         "aprovechar",
	"next-level":                "avanzado",
	"next level":                "avanzado",
}

func ApplyAntiCringeFilter(raw string) string {
	clean := raw

	// PASO 1: Reemplazar o remover clichés corporativos
	for cliche, replacement := range cliches {
		re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(cliche))
		clean = re.ReplaceAllString(clean, replacement)
	}

	// PASO 2: Normalizar densidad excesiva de emojis repetidos consecutivamente (máximo 2)
	reEmojiRepeat := regexp.MustCompile(`([\x{1F600}-\x{1F64F}\x{1F300}-\x{1F5FF}\x{1F680}-\x{1F6FF}\x{1F900}-\x{1F9FF}\x{2600}-\x{26FF}\x{2700}-\x{27BF}]){3,}`)
	clean = reEmojiRepeat.ReplaceAllString(clean, "$1$1")

	// PASO 3: Limpiar espacios extras
	clean = strings.TrimSpace(clean)

	return clean
}