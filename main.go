package main

import (
	"fmt"
	//"html/template"
	"log"
	"net/http"

	"github.com/CrosbySayan/PersonalSite/handlers"
)

func main() {
	port := ":8080"
	fmt.Printf("Server Starting on http://localhost%s\n", port)

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/window/add", handlers.SpawnWindowHandler)
	http.HandleFunc("/window/add/post/", handlers.AddPostHandler)
	http.HandleFunc("/window/delete/", handlers.DeleteWindowHandler)
	http.HandleFunc("/preview/", handlers.PreviewHandler)
	// http.HandleFunc("/window/create", CreateWindowHandler)
	// http.HandleFunc("/window/", DeleteWindowHandler)

	log.Fatal(http.ListenAndServe(port, nil))
}
