package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func SpawnWindowHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to spawn a new window")

	// Generate unique ID
	windowID := fmt.Sprintf("window-%d", time.Now().UnixNano())

	// Create new window
	newWindow := &Window{
		Title: "New Window",
		ID:    windowID,
		Position: Position{
			X:      Px(150), // Cascade windows
			Y:      Px(150),
			Anchor: "top-left",
		},
		Size: Size{
			Width:  Percent(40),
			Height: Percent(50),
		},
		Constraints: Constraints{
			MinWidth:  400,
			MinHeight: 300,
			MaxWidth:  1200,
			MaxHeight: 900,
		},
		ZIndex: 1000 + len(windows), // Increment z-index
		Content: template.HTML(fmt.Sprintf("<h1>Window %d</h1><p>Created at %s</p>",
			len(windows), time.Now().Format("15:04:05"))),
	}

	// Add to windows map
	windows[windowID] = newWindow

	// Parse and execute ONLY the window template
	tmpl := template.Must(template.ParseFiles("./templates/pages/partials/window.html"))

	// Set content type
	w.Header().Set("Content-Type", "text/html")

	// Execute template - this should return ONLY the window HTML
	err := tmpl.ExecuteTemplate(w, "window", newWindow)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		// Since we might have already written headers, just log the error
		return
	}
}
