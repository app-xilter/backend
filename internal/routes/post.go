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

		// validate json
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

		// validate tweets uri
		validatedTweets := durable.ValidateTweets(req.Tweets)
		if validatedTweets == nil {
			http.Error(w, "No valid tweets", http.StatusBadRequest)
			return
		}

		// create response
		var responseModel model.Response
		responseModel.Results = make([]model.Result, 0)

		// check already exist database
		for _, tweet := range validatedTweets {
			var tweetModel model.Tweets
			result := durable.Connection().Where(model.Tweets{Link: tweet.Link}).First(&tweetModel)
			if result.Error != nil {
				log.Printf("Error handling tweet: %v", result.Error)
				continue
			}

			if result.RowsAffected > 0 {
				responseModel.Results = append(responseModel.Results, model.Result{
					Link: tweetModel.Link,
					Tag:  tweetModel.TagId,
				})

				// remove from validatedTweets
				for i, v := range validatedTweets {
					if v.Link == tweetModel.Link {
						validatedTweets = append(validatedTweets[:i], validatedTweets[i+1:]...)
						break
					}
				}
			}
		}

		// return existing data
		if len(validatedTweets) == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			res, err := json.Marshal(responseModel)
			if _, err = w.Write(res); err != nil {
				log.Printf("Error writing response: %v", err)
			}
			return
		}

		// create prompts
		tweetsPrompt := durable.CreateTweetsPrompt(validatedTweets)
		categoriesPrompt := durable.CreateCategoriesPrompt(req.Tags)

		// send openai request
		apiResponse, _ := durable.OpenAIRequest(model.OpenAIRequest{
			Prompt:   fmt.Sprintf("%s %s", os.Getenv("OPENAI_PROMPT"), categoriesPrompt),
			Text:     tweetsPrompt,
			MaxToken: 1000,
			APIKey:   os.Getenv("OPENAI_API_KEY"),
			Timeout:  5 * time.Second,
		})
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusPartialContent)

			res, err := json.Marshal(responseModel)
			if _, err = w.Write(res); err != nil {
				log.Printf("Error writing response: %v", err)
			}
			return
		}

		// parse openai response
		var contentData map[string]interface{}
		err = json.Unmarshal([]byte(apiResponse.Choices[0].Message.Content), &contentData)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusPartialContent)

			res, err := json.Marshal(responseModel)
			if _, err = w.Write(res); err != nil {
				log.Printf("Error writing response: %v", err)
			}
			return
		}

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

			// add to database
			createTweet := model.Tweets{
				Link:  validatedTweets[index-1].Link,
				TagId: value.(int),
			}

			result := durable.Connection().Where(model.Tweets{Link: validatedTweets[index-1].Link}).FirstOrCreate(&createTweet)
			if result.Error != nil {
				log.Printf("Error handling tweet: %v", result.Error)
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

		// write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(res); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})
}
