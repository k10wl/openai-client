package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

type CreateChatCompletion struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

type Choice struct {
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
	Message      Message `json:"message"`
}

type ChatCompletionObject struct {
	ID      string   `json:"id,omitempty"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Created int      `json:"created,omitempty"`
}

type ChatCompletionModel struct {
	Name         string
	TokenLimit   int
	openAIClient *OpenAIClient
}

func NewChatCompletionModel(
	client *OpenAIClient,
	name string,
) (*ChatCompletionModel, error) {
	modelLimits, ok := modelsLimits[name]
	if !ok {
		available := []string{}
		for k := range modelsLimits {
			available = append(available, k)
		}
		return nil, errors.New(
			fmt.Sprintf(
				"model %q not found, available: %v",
				name,
				strings.Join(available, ", "),
			),
		)
	}
	return &ChatCompletionModel{
		Name:         modelLimits.name,
		TokenLimit:   modelLimits.tokens,
		openAIClient: client,
	}, nil
}

func (m *ChatCompletionModel) ChatCompletion(
	messages []Message,
) (*ChatCompletionObject, error) {
	return m.openAIClient.ChatCompletion(
		&CreateChatCompletion{
			Messages: messages,
			Model:    m.Name,
		})
}
func (c *OpenAIClient) ChatCompletion(
	completion *CreateChatCompletion,
) (*ChatCompletionObject, error) {
	d, err := json.Marshal(completion)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(d)
	var buff bytes.Buffer
	err = c.post(fmt.Sprintf("%s%s", baseUrl, chatCompletion), r, &buff)
	if err != nil {
		return nil, err
	}
	var obj ChatCompletionObject
	json.Unmarshal(buff.Bytes(), &obj)
	return &obj, nil
}
