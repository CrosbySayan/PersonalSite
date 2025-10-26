package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
)

func PreviewHandler(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/preview/")

	// Prevent directory traversal
	filename = filepath.Base(filename)

	postPath := filepath.Join("./posts", filename)

	content, err := os.ReadFile(postPath)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	htmlContent := string(markdown.ToHTML(content, nil, nil))

	// Wrap in a preview container
	preview := fmt.Sprintf(`
        <div class="preview-content">
            %s
        </div>
    `, htmlContent)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(preview))
}
