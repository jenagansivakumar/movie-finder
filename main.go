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

var redisClient *redis.Client
var limiterMap = make(map[string]*rate.Limiter)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error loading .env: %s", err)
	}
}

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func getApiKey() string {
	return os.Getenv("API_KEY")
}

func createRateLimiter(ip string) *rate.Limiter {
	limiter := rate.NewLimiter(rate.Every(5*time.Second), 1)
	limiterMap[ip] = limiter
	return limiter
}

func getTmdbResults(w http.ResponseWriter, r *http.Request) {
	apiKey := getApiKey()
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%s", apiKey)
	ip := r.RemoteAddr
	key := fmt.Sprintf("rate_limit:%s", ip)
	ctx := context.Background()

	rateLimiter, exists := limiterMap[ip]
	if !exists {
		rateLimiter = createRateLimiter(ip)
	}

	if !rateLimiter.Allow() {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	count, err := redisClient.Incr(ctx, key).Result()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if count == 1 {
		err = redisClient.Expire(ctx, key, 5*time.Second).Err()
		if err != nil {
			return
		}
	}

	if count > 1 {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	cachedItem, err := redisClient.Get(ctx, url).Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedItem))
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Cannot retrieve URL", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var totalResponse TotalResponse
	err = json.NewDecoder(resp.Body).Decode(&totalResponse)
	if err != nil {
		http.Error(w, "Cannot decode JSON", http.StatusInternalServerError)
		return
	}

	encodedJson, err := json.Marshal(totalResponse)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	err = redisClient.Set(ctx, url, string(encodedJson), 10*time.Minute).Err()
	if err != nil {
		http.Error(w, "Error setting data in Redis", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(totalResponse)
}

func main() {
	initRedis()
	http.HandleFunc("/", getTmdbResults)
	http.ListenAndServe(":8080", nil)
}
