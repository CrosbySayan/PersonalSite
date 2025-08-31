package main

import (
	"fmt"
	"html/template"
	// "html/template"
	"log"
	"net/http"
)

func process(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./templates/pages/home.html")

	t.Execute(w, "Hello World")
}

func main() {
	// Parse template
	// Setup Handler
	// Setup Website

	port := ":8080"
	fmt.Printf("Server Starting on http://localhost%s", port)

	http.HandleFunc("/", process)
	log.Fatal(http.ListenAndServe(port, nil))
}
