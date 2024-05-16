package client

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func TestPost(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		body    io.Reader
		mockDo  func(req *http.Request) (*http.Response, error)
		wantErr bool
	}{
		{
			name: "valid request",
			url:  "https://api.openai.com/v1/chat/completions",
			body: bytes.NewReader([]byte(`{"test": "data"}`)),
			mockDo: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{}`))),
				}, nil
			},
			wantErr: false,
		},
		{
			name: "invalid request",
			url:  "https://api.openai.com/v1/chat/completions",
			body: bytes.NewReader([]byte(`{"test": "data"}`)),
			mockDo: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body: io.NopCloser(
						bytes.NewReader([]byte(`{"error": "something went wrong"}`)),
					),
				}, nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &OpenAIClient{
				apiKey: "test",
				httpClient: &http.Client{
					Transport: roundTripperFunc(tt.mockDo),
				},
			}
			var buff bytes.Buffer
			err := client.post(tt.url, tt.body, &buff)
			if (err != nil) != tt.wantErr {
				t.Errorf("post() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
