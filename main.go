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
	"golang.org/x/time/rate"
)

type Movie struct {
	Title string `json:"title"`
}

type TotalResults struct {
	Results []Movie `json:"results"`
}

var redisClient *redis.Client
var limiterMap = make(map[string]*rate.Limiter)

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

	apiKey := getApi()
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%s", apiKey)
	ctx := context.Background()
	ip := r.RemoteAddr

	cacheData, err := redisClient.Get(ctx, url).Result()

	if err == redis.Nil {
		fmt.Println("Cache miss")
		if _, ok := limiterMap[ip]; !ok {
			limiterMap[ip] = rate.NewLimiter(rate.Every(5*time.Minute), 1)

		}
		if !limiterMap[ip].Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
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
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.Write(jsonData)

	} else if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
	} else {
		fmt.Println("using cache")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cacheData))
		return
	}

}

func main() {
	initRedis()
	http.HandleFunc("/recommendations", getResults)
	http.HandleFunc("/health", getHealth)
	http.ListenAndServe(":8080", nil)
}
