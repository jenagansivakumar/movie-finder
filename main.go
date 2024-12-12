package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type MovieStruct struct {
	Title string
	Genre string
}

type DataBaseResult struct {
	Result []MovieStruct
}

type ApiResponse struct {
}

func getApiKey() string {
	godotenv.Load()
	apiKey := os.Getenv("API_KEY")
	return apiKey

}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func getRecommendations(w http.ResponseWriter, r *http.Request) {
	// genre := r.URL.Query().Get("genre")
	apiKey := getApiKey()
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%s", apiKey)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Error retrieving response from url", http.StatusInternalServerError)
		return
	}
	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		w.Write([]byte("json data not found"))
		return
	}
	resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
	fmt.Println(string(jsonData))

	var movieResponse MovieStruct

	err = json.Unmarshal(jsonData, &movieResponse)
	if err != nil {
		http.Error(w, "Error unmarshalling json", http.StatusInternalServerError)
	}

	fmt.Println(movieResponse)

}

func main() {
	http.HandleFunc("/recommendations", getRecommendations)
	http.HandleFunc("/health", healthCheck)
	http.ListenAndServe(":80", nil)
}
