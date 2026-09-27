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
			SystemInstruction: "Eres un Tech Lead anunciando una nueva versión del software en LinkedIn. Resalta el tag de la versión, notas principales y enlace.",
			Example: FewShotExample{
				InputContext:  "Release v1.0.0: HMAC verification added.",
				ExpectedOutput: "¡Lanzamos la versión v1.0.0 de git-to-feed!\n\nEn esta entrega incorporamos validación de firma criptográfica HMAC-SHA256 para asegurar la ingesta de webhooks.\n\nNovedades principales:\n- Verificación de firmas con cabeceras de GitHub.\n- Middleware de autenticación constante.\n\nRevisa los detalles en GitHub: https://github.com/hawksxo/git-to-feed",
			},
		}, nil
	case ArchetypeFeature:
		return PromptTemplate{
			Archetype:         ArchetypeFeature,
			SystemInstruction: "Eres un Desarrollador Senior explicando una nueva funcionalidad técnica en LinkedIn. Destaca el problema resuelto y aprendizajes.",
			Example: FewShotExample{
				InputContext:  "Feature: Added webhooks event filtering.",
				ExpectedOutput: "¿Cómo filtramos eventos eficientemente en webhooks de GitHub?\n\nImplementamos un evaluador estricto para ignorar eventos irrelevantes y procesar únicamente releases y PRs fusionadas.\n\nAspectos clave:\n- Inspección rápida del header X-GitHub-Event.\n- Respuestas 200 OK tempranas para evitar procesamiento innecesario.\n\nRevisa los cambios: https://github.com/hawksxo/git-to-feed",
			},
		}, nil
	case ArchetypeRefactor:
		return PromptTemplate{
			Archetype:         ArchetypeRefactor,
			SystemInstruction: "Eres un Arquitecto de Software explicando refactorizaciones internas y Clean Architecture en LinkedIn.",
			Example: FewShotExample{
				InputContext:  "Refactor: Migrated package layout to Package-by-Feature",
				ExpectedOutput: "Refactorizando para escalar: Migración a Package-by-Feature\n\nReorganizamos el código del servicio para agrupar componentes por funcionalidad en lugar de por capas técnicas.\n\nBeneficios:\n- Mayor cohesión dentro de cada módulo.\n- Desacoplamiento estricto siguiendo Clean Architecture.\n\nDetalles del commit: https://github.com/hawksxo/git-to-feed",
			},
		}, nil
	default:
		return PromptTemplate{}, errors.New("invalid or unsupported archetype")
	}
}