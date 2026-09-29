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
	ErrEmptyAccessToken = errors.New("access token cannot be empty")
	ErrEmptyAuthorURN   = errors.New("author URN cannot be empty")
	ErrEmptyPostText    = errors.New("post text cannot be empty")
	ErrLinkedInAPI      = errors.New("linkedin api error")
)

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
	Author         string         `json:"author"`
	LifecycleState string         `json:"lifecycleState"`
	SpecificContent specificContent `json:"specificContent"`
	Visibility     visibility     `json:"visibility"`
}

type specificContent struct {
	ShareContent shareContent `json:"com.linkedin.ugc.ShareContent"`
}

type shareContent struct {
	ShareCommentary shareCommentary `json:"shareCommentary"`
	ShareMediaCategory string       `json:"shareMediaCategory"`
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
				ShareCommentary: shareCommentary{Text: text},
				ShareMediaCategory: "NONE",
			},
		},
		Visibility: visibility{
			MemberNetworkVisibility: "PUBLIC",
		},
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error serializando payload ugcPost: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(jsonBytes))
	if err != nil {
		return "", fmt.Errorf("error creando http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error ejecutando request a linkedin: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("%w: status %d - %s", ErrLinkedInAPI, resp.StatusCode, string(bodyBytes))
	}

	var resPayload ugcPostResponse
	if err := json.NewDecoder(resp.Body).Decode(&resPayload); err != nil {
		return "", fmt.Errorf("error decodificando respuesta de linkedin: %w", err)
	}

	if resPayload.ID == "" {
		return "", fmt.Errorf("%w: URN no retornado por linkedin", ErrLinkedInAPI)
	}

	return resPayload.ID, nil
}

type MockLinkedInClient struct {
	ShouldFail bool
	ShareURN   string
}

func NewMockLinkedInClient(shareURN string, shouldFail bool) *MockLinkedInClient {
	if shareURN == "" {
		shareURN = "urn:li:share:999999999"
	}
	return &MockLinkedInClient{
		ShareURN:   shareURN,
		ShouldFail: shouldFail,
	}
}

func (m *MockLinkedInClient) SharePost(ctx context.Context, authorURN string, text string) (string, error) {
	if authorURN == "" {
		return "", ErrEmptyAuthorURN
	}
	if text == "" {
		return "", ErrEmptyPostText
	}
	if m.ShouldFail {
		return "", fmt.Errorf("%w: mock failure configured", ErrLinkedInAPI)
	}
	return m.ShareURN, nil
}
