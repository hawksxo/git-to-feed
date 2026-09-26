package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Handler struct {
	ProcessWebhookUseCase ProcessWebhookUseCase
}

func (h *Handler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// PASO D: Obtener el Secret y Firma del Header
	secret := os.Getenv("WEBHOOK_SECRET")
	if secret == "" {
		secret = "default_dev_secret_git_to_feed"
	}
	signature := r.Header.Get("X-Hub-Signature-256")

	// PASO E: Calcular firma local
	// PASO E-1: Leer los bytes
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer cuerpo de petición", http.StatusBadRequest)
		return
	}
	// PASO E-2: Firmar los bytes junto a la secret
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(bodyBytes)
	// PASO E-3: Conversión firma calculada a hex
	calculatedSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	// PASO F: Comparar en Tiempo Constante
	if !hmac.Equal([]byte(signature), []byte(calculatedSignature)) {
		http.Error(w, "Firma inválida", http.StatusUnauthorized)
		return
	}

	// PASO G: Evaluador de eventos
	eventType := r.Header.Get("X-GitHub-Event")

	switch eventType {
	case "ping":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status": "OK", "message": "Pong"}`)
		return
	case "release", "pull_request":
		// Evento permitido: dejamos continuar la ejecución
	default:
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status": "OK", "message": "Event ignored"}`)
		return
	}

	// Paso A: Decodificar el JSON entrante
	var info GitHubPayload
	info.EventType = eventType
	err = json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&info)

	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Paso B: Invocar el Caso de Uso
	err = h.ProcessWebhookUseCase.ProcessEvent(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Paso C: Responder el Cliente HTTP
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status": "OK", "message": "Listening Event"}`)
}
