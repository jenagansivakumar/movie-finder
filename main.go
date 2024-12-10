package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/time/rate"
)

type TMDbMovie struct {
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	VoteAverage float64 `json:"vote_average"`
}

type TMDbResponse struct {
	Results []TMDbMovie `json:"results"`
}

var limiter = make(map[string]*rate.Limiter)

func getLimiter(ip string) *rate.Limiter {
	ipLimiter, exists := limiter[ip]
	if exists {
		return ipLimiter
	}
	newLimiter := rate.NewLimiter(rate.Every(5*time.Second), 1)
	limiter[ip] = newLimiter
	return newLimiter
}

func isThereIP(w http.ResponseWriter, r *http.Request) {
	userIp := r.Header.Get("X-Forwarded-For")
	if userIp == "" {
		userIp, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	ipLimiter := getLimiter(userIp)
	if !ipLimiter.Allow() {
		fmt.Println("Request denied for IP:", userIp)
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	fmt.Println("Request allowed for IP:", userIp)
}

func fetchMoviesFromTMDb(query string) ([]TMDbMovie, error) {
	apiKey := "6b951cc71ae3f04fedcccb3585b50bab"
	escapedQuery := url.QueryEscape(query)
	apiURL := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s&query=%s", apiKey, escapedQuery)

	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("HTTP Request Error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Error: %s", resp.Status)
	}

	var tmdbResponse TMDbResponse
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResponse); err != nil {
		fmt.Println("JSON Decode Error:", err)
		return nil, err
	}

	return tmdbResponse.Results, nil
}

func handleRecommendations(w http.ResponseWriter, r *http.Request) {
	isThereIP(w, r)

	query := r.URL.Query().Get("genre")
	if query == "" {
		http.Error(w, "Missing 'genre' query parameter", http.StatusBadRequest)
		return
	}

	results, err := fetchMoviesFromTMDb(query)
	if err != nil {
		http.Error(w, "Failed to fetch movies", http.StatusInternalServerError)
		fmt.Println("Error fetching movies:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		fmt.Println("Error encoding JSON:", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/recommendations", handleRecommendations)
	http.HandleFunc("/health", handleHealth)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
