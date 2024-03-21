package durable

import (
	"backend/internal/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func OpenAIRequest(request model.OpenAIRequest) (model.OpenAIResponse, error) {
	validatePrompt := []map[string]interface{}{
		{
			"role":    "system",
			"content": request.Prompt,
		},
		{
			"role":    "assistant",
			"content": "ok",
		},
		{
			"role":    "user",
			"content": request.Text,
		},
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":             "gpt-3.5-turbo-1106",
		"messages":          validatePrompt,
		"temperature":       0,
		"max_tokens":        request.MaxToken,
		"top_p":             0,
		"frequency_penalty": 0,
		"presence_penalty":  0,
	})
	if err != nil {
		return model.OpenAIResponse{}, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return model.OpenAIResponse{}, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", request.APIKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: request.Timeout,
	}
	response, err := client.Do(req)
	if err != nil {
		return model.OpenAIResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return model.OpenAIResponse{}, err
	}

	var apiResponse model.OpenAIResponse
	err = json.Unmarshal(responseBody, &apiResponse)
	if err != nil {
		return model.OpenAIResponse{}, err
	}

	return apiResponse, nil
}

func ValidateTweets(t []model.Tweet) []model.Tweet {
	var validTweets []model.Tweet
	for _, tweet := range t {
		err := ValidateUrl(tweet.Link)
		if err != nil {
			continue
		}

		validTweets = append(validTweets, tweet)
	}

	return validTweets

}

func CreateTweetsPrompt(t []model.Tweet) string {
	var tweetsPrompt string
	var length int
	for _, tweet := range t {
		length++
		tweetsPrompt += fmt.Sprintf("%d:`%s`\n", length, tweet.Text)
	}

	return tweetsPrompt
}

func CreateCategoriesPrompt(c []model.Tag) (string, error) {
	var categoriesPrompt = "categories:"
	var length int

	for _, tag := range c {
		length++
		categoriesPrompt += fmt.Sprintf("%s(%d),", tag.Text, tag.Id)
	}

	categoriesPrompt = categoriesPrompt[:len(categoriesPrompt)-1]

	if length == 0 {
		return "", fmt.Errorf("no valid categories")

	}

	return categoriesPrompt, nil
}
