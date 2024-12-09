package main

import "net/http"

func healthChecker(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
func main() {
	http.HandleFunc("/health", healthChecker)
	http.ListenAndServe(":8080", nil)
}
