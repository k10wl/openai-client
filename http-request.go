package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

const baseUrl = "https://api.openai.com/"
const chatCompletion = "v1/chat/completions"

func (c *OpenAIClient) post(url string, body io.Reader, writer io.Writer) error {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode < 200 || res.StatusCode >= 400 {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(data))
	}
	_, err = io.Copy(writer, res.Body)
	return err
}
