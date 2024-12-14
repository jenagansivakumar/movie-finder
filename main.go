package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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

func getApiKey() string {
	return os.Getenv("API_KEY")
}

func getResults(w http.ResponseWriter, r *http.Request) {

	apiKey := getApiKey()
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%s", apiKey)

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

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

}

func main() {

	http.HandleFunc("/", getResults)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
