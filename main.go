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
	http.HandleFunc("/test", handlers.SpawnWindowHandler)
	// http.HandleFunc("/window/create", CreateWindowHandler)
	// http.HandleFunc("/window/", DeleteWindowHandler)

	log.Fatal(http.ListenAndServe(port, nil))
}
