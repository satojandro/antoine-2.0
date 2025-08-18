// Package styles provides typography and text styling for Antoine CLI
// This file defines consistent typography styles, fonts, and text formatting
package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Font weights and styles
//const (
//	FontWeightNormal = lipgloss.Normal
//	FontWeightBold   = lipgloss.Bold
//)

// Text alignment options
const (
	AlignLeft   = lipgloss.Left
	AlignCenter = lipgloss.Center
	AlignRight  = lipgloss.Right
)

// Typography styles for different text elements
var (
	// Heading styles - using Antoine's golden theme
	H1Style = lipgloss.NewStyle().
		Foreground(Gold).
		Bold(true).
		MarginBottom(1).
		MarginTop(1)

	H2Style = lipgloss.NewStyle().
		Foreground(Cyan).
		Bold(true).
		MarginBottom(1)

	H3Style = lipgloss.NewStyle().
		Foreground(White).
		Bold(true).
		MarginBottom(1)

	H4Style = lipgloss.NewStyle().
		Foreground(Gray).
		Bold(true)

	// Body text styles
	BodyStyle = lipgloss.NewStyle().
			Foreground(White).
			MarginBottom(1)

	BodySecondaryStyle = lipgloss.NewStyle().
				Foreground(Gray)

	BodyMutedStyle = lipgloss.NewStyle().
			Foreground(DarkGray)

	// Interactive text styles
	LinkStyle = lipgloss.NewStyle().
			Foreground(Blue).
			Underline(true)

	LinkHoverStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Underline(true).
			Bold(true)

	ButtonStyle = lipgloss.NewStyle().
			Foreground(DarkBlue).
			Background(Gold).
			Bold(true).
			Padding(0, 2).
			MarginRight(1)

	ButtonHoverStyle = lipgloss.NewStyle().
				Foreground(DarkBlue).
				Background(Cyan).
				Bold(true).
				Padding(0, 2).
				MarginRight(1)

	ButtonDisabledStyle = lipgloss.NewStyle().
				Foreground(DarkGray).
				Background(Navy).
				Padding(0, 2).
				MarginRight(1)

	// Code and technical text styles
	CodeStyle = lipgloss.NewStyle().
			Foreground(Green).
			Background(Navy).
			Padding(0, 1)

	CodeBlockStyle = lipgloss.NewStyle().
			Foreground(Green).
			Background(Navy).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	// Status and semantic text styles
	SuccessStyle = lipgloss.NewStyle().
			Foreground(Green).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Orange).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Red).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(Yellow).
			Bold(true)

	// Special purpose text styles
	BrandStyle = lipgloss.NewStyle().
			Foreground(Gold).
			Bold(true)

	AccentStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true)

	HighlightStyle = lipgloss.NewStyle().
			Foreground(DarkBlue).
			Background(Gold).
			Bold(true)

	// Label and form styles
	LabelStyle = lipgloss.NewStyle().
			Foreground(Gray).
			Bold(true)

	InputStyle = lipgloss.NewStyle().
			Foreground(White).
			Background(Navy).
			Padding(0, 1)

	InputFocusStyle = lipgloss.NewStyle().
			Foreground(White).
			Background(Navy).
			Border(lipgloss.NormalBorder()).
			BorderForeground(Gold).
			Padding(0, 1)

	// Table styles
	TableHeaderStyle = lipgloss.NewStyle().
				Foreground(Gold).
				Background(Navy).
				Bold(true).
				Padding(0, 1).
				Align(AlignCenter)

	TableCellStyle = lipgloss.NewStyle().
			Foreground(White).
			Padding(0, 1)

	TableSelectedStyle = lipgloss.NewStyle().
				Foreground(DarkBlue).
				Background(Gold).
				Bold(true).
				Padding(0, 1)

	// List and menu styles
	ListItemStyle = lipgloss.NewStyle().
			Foreground(White).
			MarginRight(2)

	ListSelectedStyle = lipgloss.NewStyle().
				Foreground(DarkBlue).
				Background(Gold).
				Bold(true).
				MarginRight(2)

	ListBulletStyle = lipgloss.NewStyle().
			Foreground(Gold).
			Bold(true)
)

// Typography helper functions

// Title creates a large, prominent title with Antoine branding
func Title(text string) string {
	return H1Style.Render(text)
}

// Subtitle creates a secondary heading
func Subtitle(text string) string {
	return H2Style.Render(text)
}

// SectionHeader creates a section header
func SectionHeader(text string) string {
	return H3Style.Render(text)
}

// Body creates standard body text
func Body(text string) string {
	return BodyStyle.Render(text)
}

// Muted creates muted/secondary text
func Muted(text string) string {
	return BodyMutedStyle.Render(text)
}

// Success creates green success text
func Success(text string) string {
	return SuccessStyle.Render(text)
}

// Warning creates orange warning text
func Warning(text string) string {
	return WarningStyle.Render(text)
}

// Error creates red error text
func Error(text string) string {
	return ErrorStyle.Render(text)
}

// Info creates yellow info text
func Info(text string) string {
	return InfoStyle.Render(text)
}

// Brand creates text with Antoine's brand styling
func Brand(text string) string {
	return BrandStyle.Render(text)
}

// Accent creates text with accent color
func Accent(text string) string {
	return AccentStyle.Render(text)
}

// Highlight creates highlighted text with background
func Highlight(text string) string {
	return HighlightStyle.Render(text)
}

// Code creates inline code styling
func Code(text string) string {
	return CodeStyle.Render(text)
}

// CodeBlock creates block code styling
func CodeBlock(text string) string {
	return CodeBlockStyle.Render(text)
}

// Link creates a clickable link style
func Link(text string) string {
	return LinkStyle.Render(text)
}

// Button creates button-like text
func Button(text string) string {
	return ButtonStyle.Render(text)
}

// Label creates form label text
func Label(text string) string {
	return LabelStyle.Render(text)
}

// Responsive typography based on terminal width
type ResponsiveText struct {
	Small  string
	Medium string
	Large  string
}

// GetResponsiveText returns appropriate text based on terminal width
func GetResponsiveText(width int, responsive ResponsiveText) string {
	if width < 60 {
		return responsive.Small
	} else if width < 100 {
		return responsive.Medium
	}
	return responsive.Large
}

// Text truncation helpers
func TruncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	if maxLength <= 3 {
		return "..."
	}
	return text[:maxLength-3] + "..."
}

// Word wrapping for long text
func WrapText(text string, width int) []string {
	if len(text) <= width {
		return []string{text}
	}

	var lines []string
	words := []rune(text)

	for len(words) > width {
		// Find the last space within the width
		breakPoint := width
		for i := width - 1; i >= 0; i-- {
			if words[i] == ' ' {
				breakPoint = i
				break
			}
		}

		lines = append(lines, string(words[:breakPoint]))
		words = words[breakPoint+1:] // Skip the space
	}

	if len(words) > 0 {
		lines = append(lines, string(words))
	}

	return lines
}

// Advanced typography combinations
func TitleWithSubtitle(title, subtitle string) string {
	return Title(title) + "\n" + Subtitle(subtitle)
}

func ErrorWithCode(message, code string) string {
	return Error(message) + "\n" + Code(code)
}

func SuccessWithHighlight(message, highlight string) string {
	return Success(message) + " " + Highlight(highlight)
}
