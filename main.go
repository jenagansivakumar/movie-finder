package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type Movie struct {
	Title  string
	Genre  string
	Rating float64
}

var Movies = []Movie{
	{Title: "test", Genre: "test", Rating: 0},
	{Title: "test1", Genre: "test1", Rating: 0},
	{Title: "test2", Genre: "test2", Rating: 0},
	{Title: "jdka", Genre: "kdjfka", Rating: 939},
	{Title: "jdka", Genre: "Sci-Fi", Rating: 939},
	{Title: "sci fi ", Genre: "Sci-Fi", Rating: 939},
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

	fmt.Println("Request allowed for IP:", userIp, ipLimiter)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func handleRecommendations(w http.ResponseWriter, r *http.Request) {
	isThereIP(w, r)
	genre := r.URL.Query().Get("genre")
	fmt.Println(genre)
	var movieName []Movie
	for _, movie := range Movies {
		if genre == "" || genre == movie.Genre {
			movieName = append(movieName, movie)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(movieName); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}

func main() {
	limiter = make(map[string]*rate.Limiter)
	http.HandleFunc("/recommendations", handleRecommendations)
	http.HandleFunc("/health", handleHealth)
	http.ListenAndServe(":8080", nil)
}
