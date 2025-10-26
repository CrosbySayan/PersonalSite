package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
)

// BlogPost represents a blog post file
type BlogPost struct {
	Filename string
	Title    string // Title without extension
}

func SpawnWindowHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to spawn a new window")

	// Read blog posts from posts/ directory
	posts, err := getBlogPosts("./posts")
	if err != nil {
		log.Printf("Error reading blog posts: %v", err)
		// Continue with empty posts slice if error occurs
		posts = []BlogPost{}
	}

	// Pull file explorer style
	storageTemplate := template.Must(template.ParseFiles("./templates/pages/partials/file-explorer.html"))
	var storageContent bytes.Buffer

	// Pass the posts data to the template
	data := struct {
		Posts []BlogPost
	}{
		Posts: posts,
	}

	if err := storageTemplate.ExecuteTemplate(&storageContent, "explorer", data); err != nil {
		log.Printf("Error executing storage template: %v", err)
	}
	// Generate unique ID
	windowID := fmt.Sprintf("window-%d", time.Now().UnixNano())

	// windows[windowID] = newWindow

	// Create window data for template only
	newWindow := &Window{
		Title: "Blog Posts",
		ID:    windowID,
		Position: Position{
			X:      Percent(7.5),
			Y:      Percent(10),
			Anchor: "top-left",
		},
		Size: Size{
			Width:  Percent(85),
			Height: Percent(80),
		},
		Constraints: Constraints{
			MinWidth:  550,
			MinHeight: 300,
			// MaxWidth:  1200,
			// MaxHeight: 900,
		},
		ZIndex:  1000, // Fixed z-index, client will manage
		Content: template.HTML(storageContent.String()),
	}

	// Parse and execute template
	tmpl := template.Must(template.ParseFiles("./templates/pages/partials/window.html"))

	w.Header().Set("Content-Type", "text/html")

	// Just render the HTML - no state saved
	err = tmpl.ExecuteTemplate(w, "window", newWindow)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Failed to render window", http.StatusInternalServerError)
		return
	}
}

func getBlogPosts(postsDir string) ([]BlogPost, error) {
	var posts []BlogPost
	files, err := os.ReadDir(postsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read posts directory: %w", err)
	}

	// Process each file
	for _, file := range files {
		if !file.IsDir() {
			filename := file.Name()
			// You can filter by extension if needed (e.g., only .md or .html files)
			// For now, we'll include all files

			// Remove extension for title
			title := strings.TrimSuffix(filename, filepath.Ext(filename))

			posts = append(posts, BlogPost{
				Filename: filename,
				Title:    title,
			})
		}
	}
	return posts, nil
}

func AddPostHandler(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/window/add/post/")

	// Prevent directory traversal
	filename = filepath.Base(filename)

	postPath := filepath.Join("./posts/", filename)

	content, err := os.ReadFile(postPath)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	htmlContent := string(markdown.ToHTML(content, nil, nil))

	windowID := fmt.Sprintf("window-%d", time.Now().UnixNano())

	blogPost := &Window{
		Title: "Blog Posts",
		ID:    windowID,
		Position: Position{
			X:      Percent(2.5),
			Y:      Percent(5),
			Anchor: "top-left",
		},
		Size: Size{
			Width:  Percent(95),
			Height: Percent(90),
		},
		Constraints: Constraints{
			MinWidth:  550,
			MinHeight: 300,
			// MaxWidth:  1200,
			// MaxHeight: 900,
		},
		ZIndex:  1000, // Fixed z-index, client will manage
		Content: template.HTML(htmlContent),
	}

	// Parse and execute template
	tmpl := template.Must(template.ParseFiles("./templates/pages/partials/window.html"))

	w.Header().Set("Content-Type", "text/html")

	// Just render the HTML - no state saved
	err = tmpl.ExecuteTemplate(w, "window", blogPost)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Failed to render window", http.StatusInternalServerError)
		return
	}
}
