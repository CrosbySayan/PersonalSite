package handlers

import (
	"log"
	"net/http"
	"strings"
)

func DeleteWindowHandler(w http.ResponseWriter, r *http.Request) {
	windowID := strings.TrimPrefix(r.URL.Path, "/api/delete-window/")
	log.Printf("Deleting window: %s", windowID)
	w.WriteHeader(http.StatusOK)
}
