package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Window struct {
	ID      string
	Title   string
	X       int
	Y       int
	Width   int
	Height  int
	Content template.HTML
	Active  bool
}

var (
	windows       = make(map[string]*Window)
	windowsMutex  sync.RWMutex
	windowCounter = 0
)

func init() {
	tmpl := template.Must(template.ParseFiles(
		"./templates/pages/partials/storage-content.html",
	))
	// Execute the template to get the HTML content
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "storage", nil) // Pass your data as the second argument if needed
	if err != nil {
		// Handle error appropriately
		panic(err) // or log.Fatal(err)
	}
	// Create initial window
	windows["0"] = &Window{
		ID:      "0",
		Title:   "Welcome",
		X:       650,
		Y:       0,
		Width:   600,
		Height:  750,
		Content: template.HTML("Hello World! This is the first window."),
		Active:  false,
	}

	windows["1"] = &Window{
		ID:      "1",
		Title:   "Storage",
		X:       50,
		Y:       400,
		Width:   650,
		Height:  300,
		Content: template.HTML(buf.String()),
		Active:  true,
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"./templates/pages/home.html",
		"./templates/pages/partials/window.html",
	))

	windowsMutex.RLock()
	data := struct {
		Windows map[string]*Window
	}{
		Windows: windows,
	}
	windowsMutex.RUnlock()

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func CreateWindowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	windowsMutex.Lock()
	windowCounter++
	id := strconv.Itoa(windowCounter)

	// Create new window with random offset
	newWindow := &Window{
		ID:      id,
		Title:   fmt.Sprintf("Window %s", id),
		X:       100 + (windowCounter * 30), // Cascade effect
		Y:       100 + (windowCounter * 30),
		Width:   400,
		Height:  300,
		Content: template.HTML(fmt.Sprintf("This is window #%s created at %s", id, time.Now().Format("15:04:05"))),
		Active:  true,
	}

	windows[id] = newWindow
	windowsMutex.Unlock()

	// Parse and execute only the window template
	tmpl := template.Must(template.ParseFiles("./templates/pages/partials/window.html"))
	err := tmpl.ExecuteTemplate(w, "window", newWindow)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func DeleteWindowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract window ID from URL path
	id := r.URL.Path[len("/window/"):]

	windowsMutex.Lock()
	delete(windows, id)
	windowsMutex.Unlock()

	// Return empty response for HTMX to remove the element
	w.WriteHeader(http.StatusOK)
}

func main() {
	port := ":8080"
	fmt.Printf("Server Starting on http://localhost%s\n", port)

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/window/create", CreateWindowHandler)
	http.HandleFunc("/window/", DeleteWindowHandler)

	log.Fatal(http.ListenAndServe(port, nil))
}

