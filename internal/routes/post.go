package routes

import (
	"backend/internal/durable"
	"backend/internal/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Post(mux *http.ServeMux) {
	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Panic: %v", err)
			}
		}()

		var req model.Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = durable.ValidateStruct(req)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		validatedTweets := durable.ValidateTweets(req.Tweets)
		if validatedTweets == nil {
			http.Error(w, "not valid tweets", http.StatusBadRequest)
			return
		}

		tweetsPrompt := durable.CreateTweetsPrompt(validatedTweets)

		categoriesPrompt, err := durable.CreateCategoriesPrompt(req.Tags)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		apiResponse, err := durable.OpenAIRequest(model.OpenAIRequest{
			Prompt:   fmt.Sprintf("%s %s", os.Getenv("OPENAI_PROMPT"), categoriesPrompt),
			Text:     tweetsPrompt,
			MaxToken: 1000,
			APIKey:   os.Getenv("OPENAI_API_KEY"),
			Timeout:  3 * time.Second,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var contentData map[string]interface{}
		err = json.Unmarshal([]byte(apiResponse.Choices[0].Message.Content), &contentData)
		if err != nil {
			log.Printf("Error parsing JSON: %v", err)
			return
		}

		var responseModel model.Response
		for key, value := range contentData {
			value = int(value.(float64))
			if value.(int) == 0 {
				continue
			}

			index, err := strconv.Atoi(key)
			if err != nil {
				log.Printf("Error converting key to integer: %v", err)
				continue
			}

			responseModel.Results = append(responseModel.Results, model.Result{
				Link: validatedTweets[index-1].Link,
				Tag:  value.(int),
			})
		}

		res, err := json.Marshal(responseModel)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err = w.Write(res); err != nil {
			log.Printf("Error writing response: %v", err)
		}

		/*var responseModel model.Response
		for _, tweet := range req.Tweets {
			err := durable.ValidateUrl(tweet.Link)
			if err != nil {
				continue
			}

			createTweet := model.Tweets{
				Link:      tweet.Link,
				TagId:     1,
				CreatedAt: time.Now(),
			}

			result := durable.Connection().Where(model.Tweets{Link: tweet.Link}).FirstOrCreate(&createTweet)

			if result.Error != nil {
				log.Printf("Error handling tweet: %v", result.Error)
				continue
			}

			responseModel.Results = append(responseModel.Results, model.Result{
				Link:    createTweet.Link,
				Tag:     createTweet.TagId,
				TagName: "football",
			})
		}

		res, err := json.Marshal(responseModel)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err = w.Write(res); err != nil {
			log.Printf("Error writing response: %v", err)
		}*/
	})
}
