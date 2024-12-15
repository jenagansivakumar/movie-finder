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
var redisTimeOut = 10 * time.Second

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "redis-cache:6379",
		// Addr: "localhost:6379",
	})
	fmt.Println("redis initialised")
}

func getOrCreateLimiter(ip string) *rate.Limiter {
	if _, exists := limiterMap[ip]; !exists {
		limiterMap[ip] = rate.NewLimiter(rate.Every(5*time.Second), 1)
	}
	return limiterMap[ip]
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func getApi() string {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error loading .env: %s", err)
	}
	return os.Getenv("API_KEY")
}

func getResults(w http.ResponseWriter, r *http.Request, redisClient *redis.Client) {
	apiKey := getApi()
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%s", apiKey)
	ctx := context.Background()
	ip := r.RemoteAddr

	cacheData, err := redisClient.Get(ctx, url).Result()

	if err == redis.Nil {
		fmt.Println("cache misssed")
		getOrCreateLimiter(ip)
		limiter := getOrCreateLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		{
			fmt.Println("Go routine: fetching data from api")
			resp, err := http.Get(url)
			if err != nil {
				http.Error(w, "Cannot retrieve response from url", http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			var results TotalResults
			err = json.NewDecoder(resp.Body).Decode(&results)
			if err != nil {
				http.Error(w, "Error decoding json", http.StatusInternalServerError)
				return
			}

			jsonData, err := json.Marshal(results)
			if err != nil {
				http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
				return
			}

			err = redisClient.Set(ctx, url, string(jsonData), redisTimeOut).Err()
			if err != nil {
				http.Error(w, "error adding json to the cache", http.StatusInternalServerError)
				return
			}
			fmt.Println("setting cached data")
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
		}

	} else if err != nil {
		fmt.Printf("Error connecting to redis: %v\n", err)
		http.Error(w, "Internal error, redis error", http.StatusInternalServerError)
	} else {
		fmt.Println("using cache")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cacheData))
	}
}

func main() {
	initRedis()
	http.HandleFunc("/recommendations", func(w http.ResponseWriter, r *http.Request) {
		getResults(w, r, redisClient)
	})
	http.HandleFunc("/health", getHealth)
	http.ListenAndServe(":8080", nil)
}
