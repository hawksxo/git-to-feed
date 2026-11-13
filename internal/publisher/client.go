package publisher

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	ErrEmptyAccessToken = errors.New("el token de acceso a linkedin no puede estar vacio")
	ErrEmptyAuthorURN   = errors.New("el autor urn de linkedin no puede estar vacio")
	ErrEmptyPostText    = errors.New("el texto del post a publicar no puede estar vacio")
	ErrLinkedInAPI      = errors.New("error en la api de linkedin")
)

type LinkedInClient interface {
	SharePost(ctx context.Context, authorURN string, text string) (string, error)
}

type HTTPLinkedInClient struct {
	accessToken string
	httpClient  *http.Client
	baseURL     string
}

func NewHTTPLinkedInClient(accessToken string, httpClient *http.Client) (*HTTPLinkedInClient, error) {
	if accessToken == "" {
		return nil, ErrEmptyAccessToken
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &HTTPLinkedInClient{
		accessToken: accessToken,
		httpClient:  httpClient,
		baseURL:     "https://api.linkedin.com/v2/ugcPosts",
	}, nil
}

type ugcPostRequest struct {
	Author          string          `json:"author"`
	LifecycleState  string          `json:"lifecycleState"`
	SpecificContent specificContent `json:"specificContent"`
	Visibility      visibility      `json:"visibility"`
}

type specificContent struct {
	ShareContent shareContent `json:"com.linkedin.ugc.ShareContent"`
}

type shareContent struct {
	ShareCommentary    shareCommentary `json:"shareCommentary"`
	ShareMediaCategory string          `json:"shareMediaCategory"`
}

type shareCommentary struct {
	Text string `json:"text"`
}

type visibility struct {
	MemberNetworkVisibility string `json:"com.linkedin.ugc.MemberNetworkVisibility"`
}

type ugcPostResponse struct {
	ID string `json:"id"`
}

func (c *HTTPLinkedInClient) SharePost(ctx context.Context, authorURN string, text string) (string, error) {
	if authorURN == "" {
		return "", ErrEmptyAuthorURN
	}
	if text == "" {
		return "", ErrEmptyPostText
	}

	payload := ugcPostRequest{
		Author:         authorURN,
		LifecycleState: "PUBLISHED",
		SpecificContent: specificContent{
			ShareContent: shareContent{
				ShareCommentary: shareCommentary{
					Text: text,
				},
				ShareMediaCategory: "NONE",
			},
		},
		Visibility: visibility{
			MemberNetworkVisibility: "PUBLIC",
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshalling ugc request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("network error sharing to linkedin: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("%w: status %d, body: %s", ErrLinkedInAPI, resp.StatusCode, string(respBody))
	}

	var ugcResp ugcPostResponse
	if err := json.Unmarshal(respBody, &ugcResp); err != nil {
		return "", fmt.Errorf("error parsing linkedin response: %w", err)
	}

	if ugcResp.ID == "" {
		return "", errors.New("linkedin respondio con un identificador de post vacio")
	}

	return ugcResp.ID, nil
}

type MockLinkedInClient struct {
	ReturnURN   string
	ShouldFail  bool
	LastText    string
	LastAuthor  string
	SharedCount int
}

func NewMockLinkedInClient(returnURN string, shouldFail bool) *MockLinkedInClient {
	return &MockLinkedInClient{
		ReturnURN:  returnURN,
		ShouldFail: shouldFail,
	}
}

func (m *MockLinkedInClient) SharePost(ctx context.Context, authorURN string, text string) (string, error) {
	if m.ShouldFail {
		return "", ErrLinkedInAPI
	}
	if text == "" {
		return "", ErrEmptyPostText
	}
	if authorURN == "" {
		return "", ErrEmptyAuthorURN
	}
	m.LastAuthor = authorURN
	m.LastText = text
	m.SharedCount++
	return m.ReturnURN, nil
}
