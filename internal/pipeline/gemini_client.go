package pipeline

import (
	"context"
	"errors"

	"google.golang.org/genai"
)

var ErrEmptyGeminiApiKey = errors.New("gemini API key cannot be empty")

type GeminiLLMClient struct {
	apiKey    string
	modelName string
}

func NewGeminiLLMClient(apiKey string, modelName string) (*GeminiLLMClient, error) {
	if apiKey == "" {
		return nil, ErrEmptyGeminiApiKey
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
	resp, err := client.Models.GenerateContent(ctx, c.modelName, genai.Text(prompt), nil)
	if err != nil {
		return "", err
	}
	return resp.Text(), nil
}
