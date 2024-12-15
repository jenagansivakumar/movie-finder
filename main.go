package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Movie struct {
	Title string `json:"title"`
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func getApi() string {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env: %s", err)
	}

	apiKey := os.Getenv("API_KEY")
	return apiKey

}

func getResults(w http.ResponseWriter, r *http.Request) {
	apiKey := getApi()
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%s", apiKey)

	resp, err := http.Get(url)

	if err != nil {
		http.Error(w, "Cannot retrieve response from url", http.StatusInternalServerError)
	}

	var movie Movie
	err = json.NewDecoder(resp.Body).Decode(&movie)
	if err != nil {
		http.Error(w, "Error decoding json", http.StatusInternalServerError)
	}

	fmt.Println(movie)
}

func main() {

	http.HandleFunc("/recommendations", getResults)
	http.HandleFunc("/health", getHealth)
	http.ListenAndServe(":8080", nil)
}
