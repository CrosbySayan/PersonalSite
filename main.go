package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CrosbySayan/PersonalSite/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server Starting on port %s\n", port)

	http.Handle("/photos/", http.StripPrefix("/photos/", http.FileServer(http.Dir("photos"))))

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/window/add", handlers.SpawnWindowHandler)
	http.HandleFunc("/window/add/post/", handlers.AddPostHandler)
	http.HandleFunc("/window/delete/", handlers.DeleteWindowHandler)
	http.HandleFunc("/preview/", handlers.PreviewHandler)
	// http.HandleFunc("/window/create", CreateWindowHandler)
	// http.HandleFunc("/window/", DeleteWindowHandler)

	log.Fatal(http.ListenAndServe(port, nil))
}
