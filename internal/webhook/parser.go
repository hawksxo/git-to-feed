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
)

func ParseAndValidate(r *http.Request, secret string) (GitHubPayload, error) {
	if secret == "" {
		return GitHubPayload{}, fmt.Errorf("se requiere configurar la clave secreta del webhook de github")
	}

	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		return GitHubPayload{}, fmt.Errorf("no se ha podido localizar la firma del encabezado del webhook")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return GitHubPayload{}, fmt.Errorf("error al leer el cuerpo de la petición")
	}

	validate := validateSignature(secret, signature, body)
	if !validate {
		return GitHubPayload{}, fmt.Errorf("firma inválida")
	}

	eventType := r.Header.Get("X-GitHub-Event")

	switch eventType {
	case "ping":
		return GitHubPayload{EventType: "ping"}, nil
	case "release", "pull_request", "issue_comment":
	default:
		return GitHubPayload{}, fmt.Errorf("evento no soportado: %s", eventType)
	}

	var info GitHubPayload
	info.EventType = eventType
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&info)

	if err != nil {
		return GitHubPayload{}, fmt.Errorf("cuerpo de carga del JSON inválido")
	}

	payload, err := validateEvent(info)
	if err != nil {
		return GitHubPayload{}, fmt.Errorf("error al validar el evento: %s", err)
	}

	if payload.EventType == "pull_request" {
		shouldProcess, reason := ShouldProcessPullRequest(payload.PullRequest)
		if !shouldProcess {
			return GitHubPayload{}, fmt.Errorf("el PullRequest ha sido ignorado por: %s", reason)
		}
	}

	return payload, nil
}

func validateSignature(secret, signature string, body []byte) bool {
	if secret == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)

	expectedMAC := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}

func validateEvent(payload GitHubPayload) (GitHubPayload, error) {
	if payload.EventType == "release" && payload.Action != "published" {
		return GitHubPayload{}, fmt.Errorf("evento %s con acción %s ignorado", payload.EventType, payload.Action)
	}

	if payload.EventType == "pull_request" && (payload.Action != "closed" || !payload.PullRequest.Merged) {
		return GitHubPayload{}, fmt.Errorf("evento %s con acción %s ignorado", payload.EventType, payload.Action)
	}

	if payload.EventType == "issue_comment" && (payload.Action != "created" || payload.Issue.PullRequest == nil) {
		return GitHubPayload{}, fmt.Errorf("evento %s con acción %s ignorado", payload.EventType, payload.Action)
	}

	return payload, nil
}
