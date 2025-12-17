package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type LLMClient interface {
	Generate(prompt string) (string, error)
}

func NewLLMClient(cfg LLMConfig) LLMClient {
	switch cfg.Backend {
	case "openai":
		return &OpenAIClient{model: cfg.Model}
	case "anthropic":
		return &AnthropicClient{model: cfg.Model}
	default:
		return &OllamaClient{model: cfg.Model}
	}
}

func BuildPrompt(commitType, scope, diff string) string {
	scopeInfo := ""
	if scope != "" {
		scopeInfo = fmt.Sprintf("Scope: %s\n", scope)
	}

	if len(diff) > 4000 {
		diff = diff[:4000] + "\n... (truncated)"
	}

	return fmt.Sprintf(`Generate a concise commit message for this diff.
Commit type: %s
%sOnly output the message text, no quotes, no prefix.

Diff:
%s`, commitType, scopeInfo, diff)
}

// Ollama
type OllamaClient struct {
	model string
}

func (c *OllamaClient) Generate(prompt string) (string, error) {
	body := map[string]any{
		"model":  c.model,
		"prompt": prompt,
		"stream": false,
	}
	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("ollama unavailable: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return strings.TrimSpace(result.Response), nil
}

// OpenAI
type OpenAIClient struct {
	model string
}

func (c *OpenAIClient) Generate(prompt string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY not set")
	}

	model := c.model
	if model == "" {
		model = "gpt-4o-mini"
	}

	body := map[string]any{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens": 100,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("openai error: %s", string(b))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from openai")
	}
	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

// Anthropic
type AnthropicClient struct {
	model string
}

func (c *AnthropicClient) Generate(prompt string) (string, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("ANTHROPIC_API_KEY not set")
	}

	model := c.model
	if model == "" {
		model = "claude-3-haiku-20240307"
	}

	body := map[string]any{
		"model":      model,
		"max_tokens": 100,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("anthropic error: %s", string(b))
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Content) == 0 {
		return "", fmt.Errorf("no response from anthropic")
	}
	return strings.TrimSpace(result.Content[0].Text), nil
}
