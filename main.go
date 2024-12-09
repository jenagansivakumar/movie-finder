package main

import (
	"net/http"
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
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func handleRecommendations(w http.ResponseWriter, r *http.Request) {
	var movieName string = ""
	for _, movie := range Movies {
		movieName += movie.Title + "\n"
	}
	w.Write([]byte(movieName))
}

func main() {
	http.HandleFunc("/recommendations", handleRecommendations)
	http.HandleFunc("/health", handleHealth)

	http.ListenAndServe(":8080", nil)
}
