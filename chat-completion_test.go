package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
)

type MockHttpClient struct {
	Responses map[string]string
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	response, found := m.Responses[req.URL.String()]
	if !found {
		return nil, errors.New("model not found")
	}
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(response)),
	}, nil
}

func TestChatCompletion(t *testing.T) {
	client := NewOpenAIClient("secret-api-key")
	client.modelsLimits = map[string]limit{
		GPT3_5Turbo: {name: GPT3_5Turbo, tokens: 16385},
	}
	client.httpClient = &MockHttpClient{
		Responses: map[string]string{
			fmt.Sprintf("%s%s", baseUrl, chatCompletion): `{
				"id": "test-id-1",
				"model": "gpt-3.5-turbo",
				"choices": [{
					"message": {
						"role": "assistant",
						"content": "Hello, how can I assist you today?"
					}
				}]
			}`,
		},
	}

	tests := []struct {
		modelName   string
		messages    []Message
		expected    *ChatCompletionObject
		expectedErr bool
	}{
		{
			modelName: GPT3_5Turbo,
			messages: []Message{
				{Role: "user", Content: "Hello"},
			},
			expected: &ChatCompletionObject{
				ID:    "test-id-1",
				Model: GPT3_5Turbo,
				Choices: []Choice{
					{
						Message: Message{
							Role:    "assistant",
							Content: "Hello, how can I assist you today?",
						},
					},
				},
			},
			expectedErr: false,
		},
		{
			modelName: "random model",
			messages: []Message{
				{Role: "user", Content: "Tell me a joke"},
			},
			expected: &ChatCompletionObject{
				ID:    "test-id-1",
				Model: GPT3_5Turbo,
				Choices: []Choice{
					{
						Message: Message{
							Role:    "assistant",
							Content: "Why don't scientists trust atoms? Because they make up everything!",
						},
					},
				},
			},
			expectedErr: true,
		},
		{
			modelName: "nonexistent-model",
			messages: []Message{
				{Role: "user", Content: "This should fail"},
			},
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.modelName, func(t *testing.T) {
			model, err := NewChatCompletionModel(client, test.modelName)
			if (test.expectedErr && err == nil) || (!test.expectedErr && err != nil) {
				t.Fatalf("expected error: %v, got: %v", test.expectedErr, err)
			}

			if err == nil {
				actual, err := model.ChatCompletion(test.messages)
				if (test.expectedErr && err == nil) || (!test.expectedErr && err != nil) {
					t.Fatalf("expected error: %v, got: %v", test.expectedErr, err)
				}
				if !reflect.DeepEqual(actual, test.expected) {
					t.Errorf("expected: %v, got: %v", test.expected, actual)
				}
			}
		})
	}
}
