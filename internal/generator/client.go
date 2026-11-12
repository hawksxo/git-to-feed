package generator

import (
	"context"
	"errors"

	"google.golang.org/genai"
)

type LLMClient interface {
	GeneratePostWithLLM(ctx context.Context, prompt string) (string, error)
}

type GeminiLLMClient struct {
	apiKey    string
	modelName string
	config    *genai.GenerateContentConfig
}

type mockLLMClient struct {
	response string
}

func NewMockLLMClient(response string) LLMClient {
	return &mockLLMClient{response: response}
}

func (m *mockLLMClient) GeneratePostWithLLM(ctx context.Context, prompt string) (string, error) {
	return m.response, nil
}

func NewGeminiLLMClient(apiKey string, modelName string) (*GeminiLLMClient, error) {
	if apiKey == "" {
		return nil, errors.New("la API KEY de Gemini no puede estar vacía")
	}

	if modelName == "" {
		modelName = "gemini-3.5-flash-lite"
	}

	return &GeminiLLMClient{
		apiKey:    apiKey,
		modelName: modelName,
	}, nil
}

func (c *GeminiLLMClient) GeneratePostWithLLM(ctx context.Context, prompt string) (string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  c.apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", err
	}

	resp, err := client.Models.GenerateContent(ctx, c.modelName, genai.Text(prompt), c.config)
	if err != nil {
		return "", err
	}

	return resp.Text(), nil
}
