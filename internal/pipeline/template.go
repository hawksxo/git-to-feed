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
        Archetype: ArchetypeRelease,
        SystemInstruction: "Eres un Tech Lead anunciando una nueva versión del software en LinkedIn. Resalta el tag de la versión, notas principales y enlace.",
        Example: FewShotExample{
            InputContext:  "Release v1.0.0: HMAC verification added.",
            ExpectedOutput: "¡Nueva versión v1.0.0 disponible! ...",
        },
    }, nil
	case ArchetypeFeature:
		return PromptTemplate{
			Archetype: ArchetypeFeature,
			SystemInstruction: "Eres un Desarrollador Senior explicando una nueva funcionalidad técnica en LinkedIn. Destaca el problema resuelto y aprendizajes.",
			Example: FewShotExample{
				InputContext: "Feature: Added webhooks event filtering.",
				ExpectedOutput: "¿Cómo filtramos eventos en webhooks? ...",
			},
		}, nil
	case ArchetypeRefactor:
		return PromptTemplate{
			Archetype: ArchetypeRefactor,
			SystemInstruction: "Eres un Arquitecto de Software explicando refactorizaciones internas y Clean Architecture en LinkedIn.",
			Example: FewShotExample{
				InputContext: "Refactor: Megrated package layout to Package-by-Feature",
				ExpectedOutput: "Refactorizando para escalar: Clean Architecture ...",
			},
		}, nil
	default:
		return PromptTemplate{}, errors.New("invalid or unsupported archetype")
	}
}