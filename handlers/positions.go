package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
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

func ToCSS(pos Position, size Size, constraint Constraints, zIndex int) map[string]string {
	styles := make(map[string]string)

	switch pos.Anchor {
	case "center":
		styles["left"] = "50%"
		styles["top"] = "50%"
		styles["transform"] = "translate(-50%, -50%)"

	case "top-right":
		styles["right"] = pos.X.String()
		styles["top"] = pos.Y.String()

	case "bottom-left":
		styles["left"] = pos.X.String()
		styles["bottom"] = pos.Y.String()

	case "bottom-right":
		styles["right"] = pos.X.String()
		styles["bottom"] = pos.Y.String()

	case "bottom-center":
		styles["left"] = "50%"
		styles["bottom"] = pos.Y.String()
		styles["transform"] = "translateX(-50%)"

	case "top-center":
		styles["left"] = "50%"
		styles["top"] = pos.Y.String()
		styles["transform"] = "translateX(-50%)"

	case "middle-left":
		styles["top"] = "50%"
		styles["left"] = pos.X.String()
		styles["transform"] = "translateY(-50%)"

	case "middle-right":
		styles["top"] = "50%"
		styles["right"] = pos.X.String()
		styles["transform"] = "translateY(-50%)"

	default: // top-left or unspecified
		if pos.X.Unit != "center" {
			styles["left"] = pos.X.String()
		}
		if pos.Y.Unit != "center" {
			styles["top"] = pos.Y.String()
		}
	}

	// Size
	styles["width"] = size.Width.String()
	styles["height"] = size.Height.String()

	// Constraints
	if constraint.MinWidth > 0 {
		styles["min-width"] = fmt.Sprintf("%dpx", constraint.MinWidth)
	}
	if constraint.MinHeight > 0 {
		styles["min-height"] = fmt.Sprintf("%dpx", constraint.MinHeight)
	}
	if constraint.MaxWidth > 0 {
		styles["max-width"] = fmt.Sprintf("%dpx", constraint.MaxWidth)
	}
	if constraint.MaxHeight > 0 {
		styles["max-height"] = fmt.Sprintf("%dpx", constraint.MaxHeight)
	}

	styles["z-index"] = fmt.Sprintf("%d", zIndex)

	return styles
}

func (w *Window) StyleString() template.CSS {
	css := ToCSS(w.Position, w.Size, w.Constraints, w.ZIndex)
	var sb strings.Builder
	for k, v := range css {
		sb.WriteString(fmt.Sprintf("%s: %s; ", k, v))
	}
	return template.CSS(sb.String())
}

func (f *Folder) StyleString() template.CSS {
	css := ToCSS(f.Position, f.Size, f.Constraints, f.ZIndex)
	var sb strings.Builder
	for k, v := range css {
		sb.WriteString(fmt.Sprintf("%s: %s; ", k, v))
	}
	return template.CSS(sb.String())
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

type Folder struct {
	Title       string
	Position    Position
	Size        Size
	Constraints Constraints
	ZIndex      int
	// Hold some data about what is inside you
}

func getHighestZIndex(r *http.Request) int {
	highestZIndex := 1000 // default fallback
	if zIndexHeader := r.Header.Get("X-Highest-Z-Index"); zIndexHeader != "" {
		if parsed, err := strconv.Atoi(zIndexHeader); err == nil {
			// The client sends the current highest, so we add 1 for the new window
			highestZIndex = parsed + 1
			log.Printf("Received z-index from client: %d, using: %d", parsed, highestZIndex)
		} else {
			log.Printf("Failed to parse X-Highest-Z-Index header: %v", err)
		}
	} else {
		log.Printf("No X-Highest-Z-Index header found, using default: %d", highestZIndex)
	}
	return highestZIndex
}
