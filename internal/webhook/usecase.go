package webhook

import (
	"fmt"
)

type ProcessWebhookUseCase struct {

}

func (u *ProcessWebhookUseCase) ProcessEvent(payload GitHubPayload) error {

	// PASO A: Validar reglas para release
	if payload.EventType == "release" && payload.Action != "published" {
		return nil
	}

	// PASO B: Validar reglas para pull_request
	if payload.EventType == "pull_request" && (payload.Action != "closed" || !payload.PullRequest.Merged) {
		return nil
	}

	// PASO C: Impresión de datos capturados
	fmt.Printf("Se han recibido los datos del evento correctamente")

	return nil
}