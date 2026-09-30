package pipeline

import (
	"regexp"
	"strings"
)

type clicheReplacement struct {
	phrase      string
	replacement string
}

var cliches = []clicheReplacement{
	{phrase: "I am thrilled to announce", replacement: "Lanzamos"},
	{phrase: "thrilled to announce", replacement: "presentamos"},
	{phrase: "game-changer", replacement: "mejora clave"},
	{phrase: "game changer", replacement: "mejora clave"},
	{phrase: "synergy", replacement: "integración"},
	{phrase: "paradigm shift", replacement: "cambio de enfoque"},
	{phrase: "unleash the power", replacement: "aprovechar"},
	{phrase: "next-level", replacement: "avanzado"},
	{phrase: "next level", replacement: "avanzado"},
}

func ApplyAntiCringeFilter(raw string) string {
	clean := raw

	// PASO 1: Reemplazar o remover clichés corporativos (orden descendente por longitud)
	for _, item := range cliches {
		re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(item.phrase))
		clean = re.ReplaceAllString(clean, item.replacement)
	}

	// PASO 2: Normalizar densidad excesiva de emojis repetidos consecutivamente (máximo 2)
	reEmojiRepeat := regexp.MustCompile(`([\x{1F600}-\x{1F64F}\x{1F300}-\x{1F5FF}\x{1F680}-\x{1F6FF}\x{1F900}-\x{1F9FF}\x{2600}-\x{26FF}\x{2700}-\x{27BF}]){3,}`)
	clean = reEmojiRepeat.ReplaceAllString(clean, "$1$1")

	// PASO 3: Limpiar espacios extras
	clean = strings.TrimSpace(clean)

	return clean
}