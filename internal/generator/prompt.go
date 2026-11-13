package generator

import (
	"errors"
	"fmt"
)

func (a Archetype) IsValid() bool {
	switch a {
	case ArchetypeRelease, ArchetypeFeature, ArchetypeRefactor:
		return true
	default:
		return false
	}
}

func GetTemplateForArchetype(a Archetype) (PromptTemplate, error) {
	if !a.IsValid() {
		return PromptTemplate{}, errors.New("invalid or unsupported archetype")
	}

	switch a {
	case ArchetypeRelease:
		return PromptTemplate{
			Archetype:         ArchetypeRelease,
			SystemInstruction: "Eres un Tech Lead anunciando un lanzamiento o versión importante de software en LinkedIn. Usa un gancho directo, resalta métricas e hitos clave con viñetas estructuradas y un Call-To-Action final limpio.",
			Example: FewShotExample{
				InputContext:   "Release v3.8-Omni-Flash: Major release with multi-modal audio-video understanding, agentic tools, and 89% cost reduction.",
				ExpectedOutput: "🚀 ¡Presentamos la nueva versión v3.8-Omni-Flash de nuestro sistema!\n\nCapacidades nativas de procesamiento multimodal, razonamiento largo y orquestación de herramientas reunidos en un solo lanzamiento.\n\nAspectos destacados:\n- Inteligencia multimodal: Razonamiento sobre audio y video para workflows complejos y edición automatizada.\n- Reducción del 89% en costos de procesamiento en comparación con la versión anterior.\n- Ventana de contexto de 1M tokens con percepción precisa de eventos en tiempo real.\n\nRevisa la documentación y release oficial aquí: https://github.com/hawksxo/git-to-feed/releases",
			},
		}, nil
	case ArchetypeFeature:
		return PromptTemplate{
			Archetype:         ArchetypeFeature,
			SystemInstruction: "Eres un Desarrollador Senior explicando una nueva funcionalidad técnica en LinkedIn. Explica el problema resuelto mediante una analogía o concepto claro, cómo usarlo y qué cambia para el desarrollador.",
			Example: FewShotExample{
				InputContext:   "Feature: Added interactive Skills/Commands framework for automated repetitive workflows.",
				ExpectedOutput: "✨ Incorporamos un nuevo sistema de Skills interactivas para automatizar tus tareas más repetitivas.\n\nPiensa en una Skill como un atajo reutilizable para procesos recurrentes que funciona similar a un slash command. En lugar de configurar manualmente el contexto en cada petición, guardas las instrucciones una sola vez.\n\n¿Cómo funciona?\nSimplemente ejecuta la instrucción asignada y el sistema procesa el flujo completo por ti de manera transparente.\n\nRevisa el Pull Request y la implementación en GitHub: https://github.com/hawksxo/git-to-feed/pull/38",
			},
		}, nil
	case ArchetypeRefactor:
		return PromptTemplate{
			Archetype:         ArchetypeRefactor,
			SystemInstruction: "Eres un Arquitecto de Software explicando refactorizaciones internas, optimización de performance y deuda técnica en LinkedIn. Muestra métricas % claras de mejora, por qué se hizo el cambio y el impacto positivo.",
			Example: FewShotExample{
				InputContext:   "Refactor: Optimized execution engine, reduced token consumption by 30% and improved latency.",
				ExpectedOutput: "🛠️ Refactorización y optimización de rendimiento en el motor central.\n\nLogramos una mejora sustancial en la velocidad de respuesta, reduciendo en más de un 30% el tiempo de procesamiento y el consumo de recursos en ejecuciones pesadas.\n\nDetalles del cambio:\n- Reestructuración de llamadas concurrentes para evitar asignaciones innecesarias en memoria.\n- Menor uso de recursos logrando exactamente el mismo resultado y consistencia.\n\nRevisa los detalles técnicos de la refactorización: https://github.com/hawksxo/git-to-feed",
			},
		}, nil
	default:
		return PromptTemplate{}, errors.New("invalid or unsupported archetype")
	}
}

func BuildPrompt(richCtx RichContext, tmpl PromptTemplate) (string, error) {
	humanNotes := "Ninguna"
	if richCtx.HumanNotes != nil && richCtx.HumanNotes.Notes != "" {
		humanNotes = richCtx.HumanNotes.Notes
	}

	diffStats := "Ninguno"
	if richCtx.DiffStats != nil {
		diffStats = fmt.Sprintf("Archivos modificados: %d (+%d / -%d)\n", richCtx.DiffStats.TotalFiles, richCtx.DiffStats.Additions, richCtx.DiffStats.Deletions)
		for _, f := range richCtx.DiffStats.Files {
			diffStats += fmt.Sprintf("- %s (+%d / -%d)\n", f.Path, f.Additions, f.Deletions)
		}
	}

	identityProject := "No provista (opcional)"
	if richCtx.ProjectIdentity != "" {
		identityProject = richCtx.ProjectIdentity
	}

	userPrompt := fmt.Sprintf(
		`REGLAS OBLIGATORIAS DE REDACCIÓN:
1. NO te limites a resumir o copiar literalmente los textos del Pull Request o Release.
2. Debes IMITAR EXACTAMENTE la estructura, tono y estilo visual del [EJEMPLO DE REFERENCIA FEW-SHOT] que te proporcionamos abajo.
3. Transforma los datos de entrada técnicos en una publicación atractiva para la comunidad de desarrolladores en LinkedIn.
4. Si el arquetipo es FEATURE, explica el problema que resuelve la nueva funcionalidad, cómo usarla y qué ventaja le otorga al proyecto.
5. Si el arquetipo es RELEASE, agrupa los cambios en viñetas limpias con viñetas de impacto o métricas claras.
6. Si el arquetipo es REFACTOR, resalta los beneficios de performance o arquitectura logrados.
7. Incluye únicamente enlaces reales que estén presentes en la información enviada (o usa el URL provisto).
8. Responde EXCLUSIVAMENTE con el contenido final listo para ser publicado en LinkedIn. Sin explicaciones adicionales ni etiquetas meta.

[EJEMPLO DE REFERENCIA FEW-SHOT]
Entrada de ejemplo: 
%s
Publicación de ejemplo:
%s

[DATOS DEL EVENTO EN GITHUB]
Tipo de evento: %s
Versión: %s
Título del PR / Release: %s
Repositorio: %s
Autor: %s
URL de GitHub: %s
Descripción detallada: %s

[NOTAS ADICIONALES DEL DESARROLLADOR]
%s

[MÉTRICAS DE DIFF Y ARCHIVOS TOCADOS]
%s

[IDENTIDAD Y ARQUITECTURA DEL PROYECTO]
%s
`,
		tmpl.Example.InputContext,
		tmpl.Example.ExpectedOutput,
		richCtx.GitData.EventType,
		richCtx.GitData.Tag,
		richCtx.GitData.Title,
		richCtx.GitData.Repository,
		richCtx.GitData.Author,
		richCtx.GitData.URL,
		richCtx.GitData.Description,
		humanNotes,
		diffStats,
		identityProject,
	)

	return userPrompt, nil
}
