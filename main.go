package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

type Movie struct {
	Title string `json:"title"`
}

type Cache struct {
	Item  []byte
	Found bool
}

type Results struct {
	Page         int     `json:"page"`
	Results      []Movie `json:"results"`
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}
}

var redisClient *redis.Client

func initRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return client
}

func getApiKey() string {
	return os.Getenv("API_KEY")
}

func getResults(w http.ResponseWriter, r *http.Request) {

	apiKey := getApiKey()
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%s", apiKey)

	ctx := context.Background()

	cachedData, err := redisClient.Get(ctx, url).Result()
	if err == redis.Nil {
		fmt.Println("Cannot find data in cache")

	} else if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)

	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedData))
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Cannot retrieve URL", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var results Results

	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		http.Error(w, "Cannot decode movie", http.StatusInternalServerError)
		return
	}

	encodedJson, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "Error encoding data for cache", http.StatusInternalServerError)
		return
	}

	err = redisClient.Set(ctx, url, encodedJson, 10*time.Minute).Err()
	if err != nil {
		log.Printf("Failed to set data in Redis: %v", err)
	} else {
		log.Println("Data added to cache")
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

}

func main() {
	redisClient = initRedis()

	fmt.Println(redisClient)

	http.HandleFunc("/", getResults)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
