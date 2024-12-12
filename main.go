package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var tmdbApiKey string

func getDotEnv(key string) string {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file", err)

	}

	value := os.Getenv(key)
	if value == "" {
		log.Fatal("Environment variable %s is not set or empty", key)
	}

	return value
}

type Movie struct {
	Title  string
	Genre  string
	Rating float64
}

var Movies = []Movie{
	{Title: "Parasite", Genre: "Drama", Rating: 8.6},
	{Title: "Oldboy", Genre: "Thriller", Rating: 8.4},
	{Title: "The Handmaiden", Genre: "Romance", Rating: 8.1},
	{Title: "I Saw the Devil", Genre: "Horror", Rating: 7.8},
	{Title: "Train to Busan", Genre: "Action", Rating: 7.6},
	{Title: "Memories of Murder", Genre: "Crime", Rating: 8.1},
	{Title: "Sympathy for Lady Vengeance", Genre: "Thriller", Rating: 7.6},
	{Title: "The Wailing", Genre: "Horror", Rating: 7.4},
	{Title: "A Tale of Two Sisters", Genre: "Horror", Rating: 7.1},
	{Title: "Burning", Genre: "Mystery", Rating: 7.5},
}

func getRecommendation(w http.ResponseWriter, r *http.Request) {
	genre := r.URL.Query().Get("genre")
	movieList := []string{}

	if genre == "" {
		w.Write([]byte("Genre has not been speicified"))
		return
	}

	for _, movie := range Movies {
		if genre == movie.Genre {
			movieList = append(movieList, movie.Title)
		}
	}
	if len(movieList) == 0 {
		w.Write([]byte("Movielist is empty!"))
		return
	}

	jsonData, err := json.Marshal(movieList)
	if err != nil {
		http.Error(w, "Error encoding the  data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func main() {
	tmdbApiKey = getDotEnv("TMDB_API_KEY")
	fmt.Printf(tmdbApiKey)
	http.HandleFunc("/recommendations", getRecommendation)
	http.HandleFunc("/health", healthCheck)
	http.ListenAndServe(":8080", nil)
}
