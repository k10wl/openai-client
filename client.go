package client

import "net/http"

const GPT3_5Turbo = "gpt-3.5-turbo"
const GPT4Turbo = "gpt-4-turbo"

var modelsLimits = map[string]limit{
	GPT3_5Turbo: {name: GPT3_5Turbo, tokens: 16_385},
	GPT4Turbo:   {name: GPT4Turbo, tokens: 128_000},
}

type limit struct {
	name   string
	tokens int
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type OpenAIClient struct {
	apiKey       string
	modelsLimits map[string]limit
	httpClient   HttpClient
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		apiKey:       apiKey,
		modelsLimits: modelsLimits,
		httpClient:   &http.Client{},
	}
}

func (c *OpenAIClient) GetModels() []string {
	models := []string{}
	for k := range c.modelsLimits {
		models = append(models, k)
	}
	return models
}
