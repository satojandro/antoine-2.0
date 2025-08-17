// Package styles provides consistent styling definitions for Antoine CLI
// This file defines the color palette based on Antoine's golden/blue theme
package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette based on Antoine's logo and branding
// Primary colors derived from the golden Antoine logo with dark blue accents
var (
	// Primary brand colors
	Gold     = lipgloss.Color("#FFD700") // Primary gold from Antoine logo
	DarkBlue = lipgloss.Color("#1a1b26") // Main background color
	Navy     = lipgloss.Color("#24283b") // Secondary background

	// Accent colors for UI elements
	Cyan   = lipgloss.Color("#7dcfff") // Highlights and accents
	Purple = lipgloss.Color("#bb9af7") // Secondary accents
	Blue   = lipgloss.Color("#7aa2f7") // Links and interactive elements

	// Semantic colors for different states
	Green  = lipgloss.Color("#9ece6a") // Success states
	Orange = lipgloss.Color("#ff9e64") // Warning states
	Red    = lipgloss.Color("#f7768e") // Error states
	Yellow = lipgloss.Color("#e0af68") // Attention/info states

	// Text colors
	White    = lipgloss.Color("#ffffff") // Primary text
	Gray     = lipgloss.Color("#9aa5ce") // Secondary text
	DarkGray = lipgloss.Color("#565f89") // Disabled/muted text

	// Background variations
	BlackBG = lipgloss.Color("#000000") // Pure black background
	SoftBG  = lipgloss.Color("#1f2335") // Soft background for cards
)

// ColorTheme represents a complete color theme for the application
type ColorTheme struct {
	Primary       lipgloss.Color
	Secondary     lipgloss.Color
	Background    lipgloss.Color
	Surface       lipgloss.Color
	Success       lipgloss.Color
	Warning       lipgloss.Color
	Error         lipgloss.Color
	Info          lipgloss.Color
	TextPrimary   lipgloss.Color
	TextSecondary lipgloss.Color
	TextMuted     lipgloss.Color
}

// DefaultTheme returns the default Antoine color theme
func DefaultTheme() ColorTheme {
	return ColorTheme{
		Primary:       Gold,
		Secondary:     Cyan,
		Background:    DarkBlue,
		Surface:       Navy,
		Success:       Green,
		Warning:       Orange,
		Error:         Red,
		Info:          Yellow,
		TextPrimary:   White,
		TextSecondary: Gray,
		TextMuted:     DarkGray,
	}
}

// LightTheme returns a light variant of the Antoine theme
func LightTheme() ColorTheme {
	return ColorTheme{
		Primary:       lipgloss.Color("#B8860B"), // Darker gold for light backgrounds
		Secondary:     lipgloss.Color("#4682B4"), // Steel blue
		Background:    lipgloss.Color("#f8f9fa"), // Light gray background
		Surface:       lipgloss.Color("#ffffff"), // White surface
		Success:       lipgloss.Color("#28a745"), // Success green
		Warning:       lipgloss.Color("#fd7e14"), // Warning orange
		Error:         lipgloss.Color("#dc3545"), // Error red
		Info:          lipgloss.Color("#17a2b8"), // Info blue
		TextPrimary:   lipgloss.Color("#212529"), // Dark text
		TextSecondary: lipgloss.Color("#6c757d"), // Gray text
		TextMuted:     lipgloss.Color("#adb5bd"), // Muted gray
	}
}

// MinimalTheme returns a minimal monochrome theme
func MinimalTheme() ColorTheme {
	return ColorTheme{
		Primary:       lipgloss.Color("#ffffff"),
		Secondary:     lipgloss.Color("#cccccc"),
		Background:    lipgloss.Color("#000000"),
		Surface:       lipgloss.Color("#111111"),
		Success:       lipgloss.Color("#ffffff"),
		Warning:       lipgloss.Color("#cccccc"),
		Error:         lipgloss.Color("#ffffff"),
		Info:          lipgloss.Color("#cccccc"),
		TextPrimary:   lipgloss.Color("#ffffff"),
		TextSecondary: lipgloss.Color("#cccccc"),
		TextMuted:     lipgloss.Color("#666666"),
	}
}

// GetTheme returns the appropriate theme based on the theme name
func GetTheme(themeName string) ColorTheme {
	switch themeName {
	case "light":
		return LightTheme()
	case "minimal":
		return MinimalTheme()
	case "dark":
		fallthrough
	default:
		return DefaultTheme()
	}
}

// Gradient color combinations for visual effects
var (
	GoldGradient    = []string{"#FFD700", "#FFA500", "#FF8C00"} // Gold gradient
	BlueGradient    = []string{"#7dcfff", "#7aa2f7", "#6366f1"} // Blue gradient
	SuccessGradient = []string{"#9ece6a", "#68d391", "#48bb78"} // Success gradient
)

// Special purpose colors for specific UI elements
var (
	// Border colors
	BorderPrimary   = Gold
	BorderSecondary = Cyan
	BorderMuted     = DarkGray

	// Highlight colors for selections and focus
	HighlightPrimary   = Gold
	HighlightSecondary = Cyan

	// Status indicator colors
	StatusOnline  = Green
	StatusOffline = Red
	StatusBusy    = Orange
	StatusIdle    = Yellow
)

// Color intensity variations
func Lighten(color lipgloss.Color, amount float64) lipgloss.Color {
	// This is a simplified approach - in a real implementation,
	// you might want to use a proper color manipulation library
	return color
}

func Darken(color lipgloss.Color, amount float64) lipgloss.Color {
	// This is a simplified approach - in a real implementation,
	// you might want to use a proper color manipulation library
	return color
}

// Alpha transparency values
const (
	Alpha100 = 1.0  // Fully opaque
	Alpha90  = 0.9  // Slightly transparent
	Alpha75  = 0.75 // Semi-transparent
	Alpha50  = 0.5  // Half transparent
	Alpha25  = 0.25 // Mostly transparent
	Alpha10  = 0.1  // Nearly transparent
)
