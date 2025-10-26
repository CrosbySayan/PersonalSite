package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
)

// Makes an array
var (
	windows = make(map[string]*Window)
	folders = make([]*Folder, 1)
)

func init() {
}

func initializeDesktop() {
	// Initialize windows
	storageTemplate := template.Must(template.ParseFiles("./templates/pages/partials/storage-content.html"))
	var storageContent bytes.Buffer
	if err := storageTemplate.ExecuteTemplate(&storageContent, "storage", nil); err != nil {
		log.Printf("Error executing storage template: %v", err)
	}

	windows["0"] = &Window{
		Title: "Main Window",
		ID:    "main-window",
		Position: Position{
			X:      Px(50),
			Y:      Px(50),
			Anchor: "bottom-right",
		},
		Size: Size{
			Width:  Percent(50),
			Height: Percent(85),
		},
		Constraints: Constraints{
			MinWidth:  400,
			MinHeight: 300,
			MaxWidth:  1200,
			MaxHeight: 900,
		},
		ZIndex:  1000,
		Content: template.HTML("<h1>Welcome!</h1>"),
	}

	windows["1"] = &Window{
		Title: "Storage",
		ID:    "storage-window",
		Position: Position{
			X:      Px(100),
			Y:      Px(100),
			Anchor: "bottom-left",
		},
		Size: Size{
			Width:  Percent(50),
			Height: Percent(30),
		},
		Constraints: Constraints{
			MinWidth:  600,
			MinHeight: 200,
		},
		ZIndex:  1001,
		Content: template.HTML(storageContent.String()),
	}

	// Initialize folders
	folders[0] = &Folder{
		Title: "Test",
		Position: Position{
			X:      Px(50),
			Y:      Px(50),
			Anchor: "top-left",
		},
		Size: Size{
			Width:  Percent(12),
			Height: Percent(12),
		},
		Constraints: Constraints{},
		ZIndex:      999,
	}
}

// Responsible for loading the home page
// - Two windows and two folder icons
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"./templates/pages/home.html",
		"./templates/pages/partials/folder.html",
		"./templates/pages/partials/window.html",
	))

	initializeDesktop()
	data := struct {
		Windows map[string]*Window
		Folders []*Folder
	}{
		Windows: windows,
		Folders: folders,
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
