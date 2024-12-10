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
	{Title: "Oldboy", Genre: "Thriller", Rating: 9.2},
	{Title: "I Saw the Devil", Genre: "Horror", Rating: 8.5},
	{Title: "Sympathy for Lady Vengeance", Genre: "Thriller", Rating: 8.4},
	{Title: "Sympathy for Mr. Vengeance", Genre: "Thriller", Rating: 8.0},
	{Title: "Memories of Murder", Genre: "Crime", Rating: 8.9},
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
