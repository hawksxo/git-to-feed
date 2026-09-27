package pipeline

type Archetype string

const (
	ArchetypeRelease Archetype = "release"
	ArchetypeFeature Archetype = "feature"
	ArchetypeRefactor Archetype = "refactor"
)

func (a Archetype) IsValid() bool  {
	switch a {
	case ArchetypeRelease, ArchetypeFeature, ArchetypeRefactor:
		return true
	default:
		return false
	}
}