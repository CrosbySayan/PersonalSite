package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Value struct {
	Amount float64
	Unit   string
}

type Position struct {
	X      Value  // Can be "50px", "10%", "center", "left", "right"
	Y      Value  // Can be "100px", "20%", "center", "top", "bottom"
	Anchor string // "top-left" (default), "center", "top-right", etc.
}

// Size defines window dimensions
type Size struct {
	Width  Value // Can be "600px", "50%", "auto"
	Height Value // Can be "400px", "80%", "auto"
}

// Constraints define min/max bounds
type Constraints struct {
	MinWidth  int // Minimum width in pixels
	MinHeight int // Minimum height in pixels
	MaxWidth  int // Maximum width in pixels (0 = no limit)
	MaxHeight int // Maximum height in pixels (0 = no limit)
}

type Window struct {
	Title       string
	ID          string
	Position    Position
	Size        Size
	Constraints Constraints
	ZIndex      int
	Content     template.HTML
}

func (v Value) String() string {
	if v.Unit == "px" || v.Unit == "%" {
		return fmt.Sprintf("%.0f%s", v.Amount, v.Unit)
	}
	return v.Unit
}

// Helper constructors
func Px(amount float64) Value {
	return Value{Amount: amount, Unit: "px"}
}

func Percent(amount float64) Value {
	return Value{Amount: amount, Unit: "%"}
}

func Auto() Value {
	return Value{Unit: "auto"}
}

func Center() Value {
	return Value{Unit: "center"}
}

func (w *Window) ToCSS() map[string]string {
	styles := make(map[string]string)

	switch w.Position.Anchor {
	case "center":
		styles["left"] = "50%"
		styles["top"] = "50%"
		styles["transform"] = "translate(-50%, -50%)"

	case "top-right":
		styles["right"] = w.Position.X.String()
		styles["top"] = w.Position.Y.String()

	case "bottom-left":
		styles["left"] = w.Position.X.String()
		styles["bottom"] = w.Position.Y.String()

	case "bottom-right":
		styles["right"] = w.Position.X.String()
		styles["bottom"] = w.Position.Y.String()

	case "bottom-center":
		styles["left"] = "50%"
		styles["bottom"] = w.Position.Y.String()
		styles["transform"] = "translateX(-50%)"

	case "top-center":
		styles["left"] = "50%"
		styles["top"] = w.Position.Y.String()
		styles["transform"] = "translateX(-50%)"

	case "middle-left":
		styles["top"] = "50%"
		styles["left"] = w.Position.X.String()
		styles["transform"] = "translateY(-50%)"

	case "middle-right":
		styles["top"] = "50%"
		styles["right"] = w.Position.X.String()
		styles["transform"] = "translateY(-50%)"

	default: // top-left or unspecified
		if w.Position.X.Unit != "center" {
			styles["left"] = w.Position.X.String()
		}
		if w.Position.Y.Unit != "center" {
			styles["top"] = w.Position.Y.String()
		}
	}

	// Size
	styles["width"] = w.Size.Width.String()
	styles["height"] = w.Size.Height.String()

	// Constraints
	if w.Constraints.MinWidth > 0 {
		styles["min-width"] = fmt.Sprintf("%dpx", w.Constraints.MinWidth)
	}
	if w.Constraints.MinHeight > 0 {
		styles["min-height"] = fmt.Sprintf("%dpx", w.Constraints.MinHeight)
	}
	if w.Constraints.MaxWidth > 0 {
		styles["max-width"] = fmt.Sprintf("%dpx", w.Constraints.MaxWidth)
	}
	if w.Constraints.MaxHeight > 0 {
		styles["max-height"] = fmt.Sprintf("%dpx", w.Constraints.MaxHeight)
	}

	styles["z-index"] = fmt.Sprintf("%d", w.ZIndex)

	return styles
}

func (w *Window) StyleString() template.CSS {
	css := w.ToCSS()
	var sb strings.Builder
	for k, v := range css {
		sb.WriteString(fmt.Sprintf("%s: %s; ", k, v))
	}
	return template.CSS(sb.String())
}

// Makes an array
var windows = make([]*Window, 2)

func makeWindows() {
	storageTemplate := template.Must(template.ParseFiles("./templates/pages/partials/storage-content.html"))

	var storageContent bytes.Buffer
	if err := storageTemplate.ExecuteTemplate(&storageContent, "storage", nil); err != nil {
		log.Printf("Error executing storage template: %v", err)
	}

	windows[0] = &Window{
		Title: "Main Window",
		ID:    "0",
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
	windows[1] = &Window{
		Title: "Storage",
		ID:    "1",
		Position: Position{
			X:      Px(50),
			Y:      Px(50),
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
}

// Responsible for loading the home page
// - Two windows and two folder icons
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Template imports
	tmpl := template.Must(template.ParseFiles(
		"./templates/pages/home.html",
		"./templates/pages/partials/window.html",
	))
	makeWindows()
	data := struct {
		Windows []*Window
	}{
		Windows: windows,
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
