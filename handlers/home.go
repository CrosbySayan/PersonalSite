package handlers

import (
	"html/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"../templates/pages/home.html",
		"../templates/pages/partials/window.html",
	))

	data := struct {
		ID      string
		Title   string
		X       int
		Y       int
		Width   int
		Height  int
		Content template.HTML
	}{
		ID:      "0",
		Title:   "Test",
		X:       0,
		Y:       0,
		Width:   500,
		Height:  400,
		Content: template.HTML("Hello World"),
	}
	tmpl.Execute(w, data)
}
