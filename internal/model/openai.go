package model

import "time"

type OpenAIRequest struct {
	Prompt   string
	Text     string
	MaxToken int64
	APIKey   string
	Timeout  time.Duration
}

type OpenAIResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
