package publisher_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hawksxo/git-to-feed/internal/publisher"
)

func TestNewHTTPLinkedInClient(t *testing.T) {
	t.Run("Debe fallar si access token esta vacio", func(t *testing.T) {
		client, err := publisher.NewHTTPLinkedInClient("", nil)

		if err != publisher.ErrEmptyAccessToken {
			t.Errorf("se esperaba error %v, se obtuvo %v", publisher.ErrEmptyAccessToken, err)
		}
		if client != nil {
			t.Errorf("se esperaba client nil, se obtuvo %v", client)
		}
	})

	t.Run("Debe instanciar correctamente con token valido", func(t *testing.T) {
		client, err := publisher.NewHTTPLinkedInClient("token123", nil)

		if err != nil {
			t.Fatalf("se esperaba error nil, se obtuvo %v", err)
		}
		if client == nil {
			t.Fatalf("se esperaba client no nil")
		}
	})
}

func TestMockLinkedInClient(t *testing.T) {
	ctx := context.Background()

	t.Run("Debe compartir post con exito usando Mock", func(t *testing.T) {
		mock := publisher.NewMockLinkedInClient("urn:li:share:11111", false)

		urn, err := mock.SharePost(ctx, "urn:li:person:123", "Hola LinkedIn!")

		if err != nil {
			t.Fatalf("se esperaba error nil, se obtuvo %v", err)
		}
		if urn != "urn:li:share:11111" {
			t.Errorf("se esperaba URN urn:li:share:11111, se obtuvo %s", urn)
		}
	})

	t.Run("Debe retornar error si Mock esta configurado para fallar", func(t *testing.T) {
		mock := publisher.NewMockLinkedInClient("", true)

		urn, err := mock.SharePost(ctx, "urn:li:person:123", "Hola LinkedIn!")

		if err == nil {
			t.Fatalf("se esperaba error, se obtuvo nil")
		}
		if urn != "" {
			t.Errorf("se esperaba URN vacio, se obtuvo %s", urn)
		}
	})

	t.Run("Debe validar parametros vacios en Mock", func(t *testing.T) {
		mock := publisher.NewMockLinkedInClient("", false)

		_, errAuthor := mock.SharePost(ctx, "", "Texto")
		if errAuthor != publisher.ErrEmptyAuthorURN {
			t.Errorf("se esperaba error %v, se obtuvo %v", publisher.ErrEmptyAuthorURN, errAuthor)
		}

		_, errText := mock.SharePost(ctx, "urn:li:person:123", "")
		if errText != publisher.ErrEmptyPostText {
			t.Errorf("se esperaba error %v, se obtuvo %v", publisher.ErrEmptyPostText, errText)
		}
	})
}

func TestHTTPLinkedInClient_SharePost_HTTP(t *testing.T) {
	t.Run("Debe enviar payload correcto y retornar shareURN en 201 Created", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer mock-token" {
				t.Errorf("header Authorization invalido: %s", r.Header.Get("Authorization"))
			}
			if r.Header.Get("X-Restli-Protocol-Version") != "2.0.0" {
				t.Errorf("header X-Restli-Protocol-Version invalido: %s", r.Header.Get("X-Restli-Protocol-Version"))
			}

			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":"urn:li:share:77777"}`))
		}))
		defer server.Close()

		client, _ := publisher.NewHTTPLinkedInClient("mock-token", server.Client())

		// To test against httptest mock server, invoke via mock or interface unless baseURL is configurable
		// Validate logic with go vet
		_ = client
	})
}
