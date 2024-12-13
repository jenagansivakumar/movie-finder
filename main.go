package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var apiKey string

func getApi() string {
	godotenv.Load()
	apiKey = os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatalf("Cannot find API Key")
	}
	return apiKey

}

type Movie struct {
	Title string `json:"title"`
}

type MoviePage struct {
	Page         string  `json:"page"`
	Results      []Movie `json:"results"`
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
}

func fetchResults(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%s&language=en-US&page=1", apiKey)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Cannot fetch url", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)

}

func main() {
	getApi()
	http.HandleFunc("/", fetchResults)
	http.ListenAndServe(":8080", nil)
}
