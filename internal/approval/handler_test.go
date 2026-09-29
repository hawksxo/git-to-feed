package approval

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hawksxo/git-to-feed/internal/pipeline"
)

func TestHandleListPending(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	uc := NewApprovalUseCase(repo)
	handler := NewHandler(uc)

	_, _ = uc.SubmitForApproval(context.Background(), pipeline.GeneratedPost{Content: "Pending Post 1"})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/approvals/pending", nil)
	rec := httptest.NewRecorder()

	handler.HandleListPending(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("HandleListPending status = %d, want %d", rec.Code, http.StatusOK)
	}

	var posts []ApprovalPost
	_ = json.NewDecoder(rec.Body).Decode(&posts)

	if len(posts) != 1 {
		t.Errorf("len(posts) = %d, want %d", len(posts), 1)
	}
}

func TestHandleApprove(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	uc := NewApprovalUseCase(repo)
	handler := NewHandler(uc)

	appPost, _ := uc.SubmitForApproval(context.Background(), pipeline.GeneratedPost{Content: "Pending Post"})

	// Setup ServeMux to parse {uuid} path value in Go 1.22+
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/approvals/{uuid}/approve", handler.HandleApprove)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/approvals/"+appPost.UUID+"/approve", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("HandleApprove status = %d, want %d", rec.Code, http.StatusOK)
	}

	var approved ApprovalPost
	_ = json.NewDecoder(rec.Body).Decode(&approved)

	if approved.Status != StatusApproved {
		t.Errorf("approved.Status = %q, want %q", approved.Status, StatusApproved)
	}
}

func TestHandleReject(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	uc := NewApprovalUseCase(repo)
	handler := NewHandler(uc)

	appPost, _ := uc.SubmitForApproval(context.Background(), pipeline.GeneratedPost{Content: "Pending Post"})

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/approvals/{uuid}/reject", handler.HandleReject)

	bodyBytes, _ := json.Marshal(RejectRequest{Reason: "Inappropriate tone"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/approvals/"+appPost.UUID+"/reject", bytes.NewBuffer(bodyBytes))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("HandleReject status = %d, want %d", rec.Code, http.StatusOK)
	}

	var rejected ApprovalPost
	_ = json.NewDecoder(rec.Body).Decode(&rejected)

	if rejected.Status != StatusRejected {
		t.Errorf("rejected.Status = %q, want %q", rejected.Status, StatusRejected)
	}
}

func TestHandleEditAndApprove(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	uc := NewApprovalUseCase(repo)
	handler := NewHandler(uc)

	appPost, _ := uc.SubmitForApproval(context.Background(), pipeline.GeneratedPost{Content: "Pending Post"})

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v1/approvals/{uuid}/edit", handler.HandleEditAndApprove)

	bodyBytes, _ := json.Marshal(EditRequest{Content: "Polished post content for LinkedIn"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/approvals/"+appPost.UUID+"/edit", bytes.NewBuffer(bodyBytes))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("HandleEditAndApprove status = %d, want %d", rec.Code, http.StatusOK)
	}

	var edited ApprovalPost
	_ = json.NewDecoder(rec.Body).Decode(&edited)

	if edited.Status != StatusApproved {
		t.Errorf("edited.Status = %q, want %q", edited.Status, StatusApproved)
	}

	if edited.EditedContent != "Polished post content for LinkedIn" {
		t.Errorf("edited.EditedContent = %q, want %q", edited.EditedContent, "Polished post content for LinkedIn")
	}
}

func TestHandleNotFound(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	uc := NewApprovalUseCase(repo)
	handler := NewHandler(uc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/approvals/{uuid}/approve", handler.HandleApprove)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/approvals/non-existent-uuid/approve", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("HandleApprove non-existent status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
