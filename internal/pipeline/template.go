package pipeline

import "errors"

type FewShotExample struct {
	InputContext   string
	ExpectedOutput string
}

type PromptTemplate struct {
	Archetype         Archetype
	SystemInstruction string
	Example           FewShotExample
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
				InputContext:  "Release v3.8-Omni-Flash: Major release with multi-modal audio-video understanding, agentic tools, and 89% cost reduction.",
				ExpectedOutput: "🚀 ¡Presentamos la nueva versión v3.8-Omni-Flash de nuestro sistema!\n\nCapacidades nativas de procesamiento multimodal, razonamiento largo y orquestación de herramientas reunidos en un solo lanzamiento.\n\nAspectos destacados:\n- Inteligencia multimodal: Razonamiento sobre audio y video para workflows complejos y edición automatizada.\n- Reducción del 89% en costos de procesamiento en comparación con la versión anterior.\n- Ventana de contexto de 1M tokens con percepción precisa de eventos en tiempo real.\n\nRevisa la documentación y release oficial aquí: https://github.com/hawksxo/git-to-feed/releases",
			},
		}, nil
	case ArchetypeFeature:
		return PromptTemplate{
			Archetype:         ArchetypeFeature,
			SystemInstruction: "Eres un Desarrollador Senior explicando una nueva funcionalidad técnica en LinkedIn. Explica el problema resuelto mediante una analogía o concepto claro, cómo usarlo y qué cambia para el desarrollador.",
			Example: FewShotExample{
				InputContext:  "Feature: Added interactive Skills/Commands framework for automated repetitive workflows.",
				ExpectedOutput: "✨ Incorporamos un nuevo sistema de Skills interactivas para automatizar tus tareas más repetitivas.\n\nPiensa en una Skill como un atajo reutilizable para procesos recurrentes que funciona similar a un slash command. En lugar de configurar manualmente el contexto en cada petición, guardas las instrucciones una sola vez.\n\n¿Cómo funciona?\nSimplemente ejecuta la instrucción asignada y el sistema procesa el flujo completo por ti de manera transparente.\n\nRevisa el Pull Request y la implementación en GitHub: https://github.com/hawksxo/git-to-feed/pull/38",
			},
		}, nil
	case ArchetypeRefactor:
		return PromptTemplate{
			Archetype:         ArchetypeRefactor,
			SystemInstruction: "Eres un Arquitecto de Software explicando refactorizaciones internas, optimización de performance y deuda técnica en LinkedIn. Muestra métricas % claras de mejora, por qué se hizo el cambio y el impacto positivo.",
			Example: FewShotExample{
				InputContext:  "Refactor: Optimized execution engine, reduced token consumption by 30% and improved latency.",
				ExpectedOutput: "🛠️ Refactorización y optimización de rendimiento en el motor central.\n\nLogramos una mejora sustancial en la velocidad de respuesta, reduciendo en más de un 30% el tiempo de procesamiento y el consumo de recursos en ejecuciones pesadas.\n\nDetalles del cambio:\n- Reestructuración de llamadas concurrentes para evitar asignaciones innecesarias en memoria.\n- Menor uso de recursos logrando exactamente el mismo resultado y consistencia.\n\nRevisa los detalles técnicos de la refactorización: https://github.com/hawksxo/git-to-feed",
			},
		}, nil
	default:
		return PromptTemplate{}, errors.New("invalid or unsupported archetype")
	}
}