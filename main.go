package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

type Movie struct {
	Title string `json:"title"`
}

type TotalResults struct {
	Results []Movie `json:"results"`
}

var redisClient *redis.Client

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func getApi() string {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error loading .env: %s", err)
	}

	apiKey := os.Getenv("API_KEY")
	return apiKey

}

func getResults(w http.ResponseWriter, r *http.Request) {
	initRedis()
	apiKey := getApi()
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%s", apiKey)
	ctx := context.Background()

	cacheData, err := redisClient.Get(ctx, url).Result()

	if err == redis.Nil {
		fmt.Println("Cache miss")
		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, "Cannot retrieve response from url", http.StatusInternalServerError)
			return
		}

		var results TotalResults
		err = json.NewDecoder(resp.Body).Decode(&results)
		if err != nil {
			http.Error(w, "Error decoding json", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		jsonData, err := json.Marshal(results)
		if err != nil {
			http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
			return
		}

		err = redisClient.Set(ctx, url, string(jsonData), 10*time.Minute).Err()
		if err != nil {
			http.Error(w, "error adding json to the cache", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Println("using cache")
		if err := json.NewEncoder(w).Encode(results); err != nil {
			http.Error(w, "Error encoding json", http.StatusInternalServerError)
		}

	} else if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cacheData))
		return
	}

}

func main() {

	http.HandleFunc("/recommendations", getResults)
	http.HandleFunc("/health", getHealth)
	http.ListenAndServe(":8080", nil)
}
