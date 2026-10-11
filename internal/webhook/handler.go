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
	// STEP D: Retrieve Secret and Signature Header
	secret := os.Getenv("WEBHOOK_SECRET")
	if secret == "" {
		secret = "default_dev_secret_git_to_feed"
	}
	signature := r.Header.Get("X-Hub-Signature-256")

	// STEP E: Calculate local signature
	// STEP E-1: Read request bytes
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	// STEP E-2: Sign bytes using HMAC secret
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(bodyBytes)
	// STEP E-3: Convert calculated signature to hex
	calculatedSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	// STEP F: Constant time signature comparison
	if !hmac.Equal([]byte(signature), []byte(calculatedSignature)) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// STEP G: Event evaluator
	eventType := r.Header.Get("X-GitHub-Event")

	switch eventType {
	case "ping":
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, `{"status": "OK", "message": "Pong"}`)
		return
	case "release", "pull_request":
		// Allowed event: proceed with execution
	default:
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, `{"status": "OK", "message": "Event ignored"}`)
		return
	}

	// STEP A: Decode incoming JSON
	var info GitHubPayload
	info.EventType = eventType
	err = json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&info)

	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// STEP B: Invoke Use Case
	err = h.ProcessWebhookUseCase.ProcessEvent(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// STEP C: Respond to HTTP client
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintln(w, `{"status": "OK", "message": "Listening Event"}`)
}
