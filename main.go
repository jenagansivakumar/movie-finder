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

type TotalResponse struct {
	Page         int     `json:"page"`
	Results      []Movie `json:"results"`
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
}

type Movie struct {
	Title      string  `json:"title"`
	Overview   string  `json:"overview"`
	Popularity float64 `json:"popularity"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error loading .env: %s", err)
	}
}

var redisClient *redis.Client

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

var limiter *rate.Limiter

var limiterMap = make(map[string]*rate.Limiter)

func initLimiter() {

	limiter = rate.NewLimiter(rate.Every(5*time.Second), 1)

}

func getApiKey() string {
	apiKey := os.Getenv("API_KEY")
	fmt.Println(apiKey)
	return apiKey
}

func getTmdbResults(w http.ResponseWriter, r *http.Request) {

	apiKey := getApiKey()
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%s", apiKey)
	ip := r.RemoteAddr

	rateLimiter, exists := limiterMap[ip]

	if !exists {
		rateLimiter = rate.NewLimiter(rate.Every(5*time.Second), 1)
		limiterMap[ip] = rateLimiter
	}

	rateLimiter.Allow()

	if !rateLimiter.Allow() {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
	}

	ctx := context.Background()
	cachedItem, err := redisClient.Get(ctx, url).Result()

	if err == redis.Nil {
		fmt.Println("Cache miss")
	} else if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
	} else {
		fmt.Println("Using cached data")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedItem))
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Cannot retrieve url", http.StatusInternalServerError)
		return
	}

	var totalResponse TotalResponse

	err = json.NewDecoder(resp.Body).Decode(&totalResponse)
	if err != nil {
		http.Error(w, "Cannot decode json", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	encodedJson, err := json.Marshal(totalResponse)
	if err != nil {
		http.Error(w, "Error encoding json", http.StatusInternalServerError)
		return
	}

	err = redisClient.Set(ctx, url, string(encodedJson), 10*time.Minute).Err()
	fmt.Println("Setting context in redis client")
	if err != nil {
		http.Error(w, "Error setting context in redis", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(totalResponse)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

}

func main() {
	initRedis()
	http.HandleFunc("/", getTmdbResults)
	http.ListenAndServe(":8080", nil)
}
