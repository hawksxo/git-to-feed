package pipeline

import "fmt"

type PromptBuilder struct{}

func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{}
}

func (pb *PromptBuilder) BuildPrompt(ctx RichContext, archetype Archetype) (string, error) {
	notesText := "Ninguna"
	if ctx.HumanNotes != nil && ctx.HumanNotes.Notes != "" {
		notesText = ctx.HumanNotes.Notes
	}

	tmpl, err := GetTemplateForArchetype(archetype)
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`[INSTRUCCIÓN DE SISTEMA]
%s

REGLAS OBLIGATORIAS DE REDACCIÓN:
1. NO te limites a resumir o copiar literalmente los textos del Pull Request o Release.
2. Debes IMITAR EXACTAMENTE la estructura, tono y estilo visual del [EJEMPLO DE REFERENCIA FEW-SHOT] que te proporcionamos abajo.
3. Transforma los datos de entrada técnicos en una publicación atractiva para la comunidad de desarrolladores en LinkedIn.
4. Si el arquetipo es FEATURE, explica el problema que resuelve la nueva funcionalidad, cómo usarla y qué ventaja le otorga al proyecto.
5. Si el arquetipo es RELEASE, agrupa los cambios en viñetas limpias con viñetas de impacto o métricas claras.
6. Si el arquetipo es REFACTOR, resalta los beneficios de performance o arquitectura logrados.
7. Incluye únicamente enlaces reales que estén presentes en la información enviada (o usa el URL provisto).
8. Responde EXCLUSIVAMENTE con el contenido final listo para ser publicado en LinkedIn. Sin explicaciones adicionales ni etiquetas meta.

[EJEMPLO DE REFERENCIA FEW-SHOT]
Entrada de ejemplo: %s
Publicación de ejemplo (IMITA ESTA ESTRUCTURA EXACTAMENTE):
%s

[DATOS DEL EVENTO EN GITHUB (TARGET CONTEXT)]
Tipo de Evento: %s
Versión / Tag: %s
Título del PR / Release: %s
Repositorio: %s
Autor: %s
URL de GitHub: %s
Descripción detallada:
%s

[NOTAS ADICIONALES DEL DESARROLLADOR]
%s`,
		tmpl.SystemInstruction,
		tmpl.Example.InputContext,
		tmpl.Example.ExpectedOutput,
		ctx.GitData.EventType,
		ctx.GitData.Tag,
		ctx.GitData.Title,
		ctx.GitData.Repository,
		ctx.GitData.Author,
		ctx.GitData.URL,
		ctx.GitData.Description,
		notesText,
	)

	return prompt, nil
}
