package webhook

import (
	"fmt"
)

type ProcessWebhookUseCase struct {

}

func (u *ProcessWebhookUseCase) ProcessEvent(rf RepositoryInfo) error {

	// Paso A: Validar contenido del objeto
	if rf.Action == "" {
		return fmt.Errorf("la acción no puede estar vacía")
	}

	if rf.Name == "" {
		return fmt.Errorf("el nombre no puede estar vacío")
	}
	
	if rf.User == "" {
		return fmt.Errorf("el usuario no puede estar vacío")
	}

	if rf.Description == "" {
		return fmt.Errorf("la descripción no puede estar vacía")
	}

	if rf.URL == "" {
		return fmt.Errorf("la url no puede estar vacía")
	}

	return nil
}