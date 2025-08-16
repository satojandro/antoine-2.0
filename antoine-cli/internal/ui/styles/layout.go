// Package styles provides layout and spacing utilities for Antoine CLI
// This file defines consistent layouts, spacing, borders, and structural styles
package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Spacing constants for consistent layout
const (
	SpacingXS = 1 // Extra small spacing
	SpacingSM = 2 // Small spacing
	SpacingMD = 4 // Medium spacing (default)
	SpacingLG = 6 // Large spacing
	SpacingXL = 8 // Extra large spacing
)

// Border styles using Antoine's color palette
var (
	// Primary borders with gold accent
	BorderPrimaryStyle = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}

	// Secondary borders with cyan accent
	BorderSecondaryStyle = lipgloss.Border{
		Top:         "━",
		Bottom:      "━",
		Left:        "┃",
		Right:       "┃",
		TopLeft:     "┏",
		TopRight:    "┓",
		BottomLeft:  "┗",
		BottomRight: "┛",
	}

	// Subtle borders for content separation
	BorderSubtleStyle = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "└",
		BottomRight: "┘",
	}

	// Double border for emphasis
	BorderDoubleStyle = lipgloss.Border{
		Top:         "═",
		Bottom:      "═",
		Left:        "║",
		Right:       "║",
		TopLeft:     "╔",
		TopRight:    "╗",
		BottomLeft:  "╚",
		BottomRight: "╝",
	}
)

// Container styles for different layout purposes
var (
	// Main container for the entire application
	AppContainerStyle = lipgloss.NewStyle().
		Width(100).
		Height(30).
		Background(DarkBlue).
		Foreground(White)

	// Primary content container with border
	ContentContainerStyle = lipgloss.NewStyle().
		Border(BorderPrimaryStyle).
		BorderForeground(Gold).
		Background(DarkBlue).
		Foreground(White).
		Padding(1, 2).
		Margin(1)

	// Secondary container for nested content
	SecondaryContainerStyle = lipgloss.NewStyle().
		Border(BorderSubtleStyle).
		BorderForeground(Gray).
		Background(Navy).
		Foreground(White).
		Padding(1).
		Margin(0, 1)

	// Card-like container for grouped content
	CardStyle = lipgloss.NewStyle().
		Border(BorderPrimaryStyle).
		BorderForeground(Gold).
		Background(SoftBG).
		Foreground(White).
		Padding(1, 2).
		Margin(1, 0)

	// Modal container for dialogs and popups
	ModalStyle = lipgloss.NewStyle().
		Border(BorderDoubleStyle).
		BorderForeground(Cyan).
		Background(Navy).
		Foreground(White).
		Padding(2, 4).
		Align(AlignCenter)

	// Header container for titles and navigation
	HeaderStyle = lipgloss.NewStyle().
		Border(BorderPrimaryStyle, false, false, true, false).
		BorderForeground(Gold).
		Background(DarkBlue).
		Foreground(Gold).
		Bold(true).
		Padding(1, 2).
		Width(100)

	// Footer container for status and help
	FooterStyle = lipgloss.NewStyle().
		Border(BorderPrimaryStyle, true, false, false, false).
		BorderForeground(Gold).
		Background(Navy).
		Foreground(Gray).
		Padding(1, 2).
		Width(100)

	// Sidebar container for navigation
	SidebarStyle = lipgloss.NewStyle().
		Border(BorderSubtleStyle, false, true, false, false).
		BorderForeground(Gray).
		Background(Navy).
		Foreground(White).
		Padding(1).
		Width(20).
		Height(25)
)

// Flex layout utilities
var (
	// Horizontal flex container
	FlexRowStyle = lipgloss.NewStyle().
		Display(lipgloss.Block)

	// Vertical flex container
	FlexColumnStyle = lipgloss.NewStyle().
		Display(lipgloss.Block)

	// Centered flex container
	FlexCenterStyle = lipgloss.NewStyle().
		Display(lipgloss.Block).
		Align(AlignCenter)
)

// Spacing utilities
var (
	// Margin utilities
	MarginTopStyle    = lipgloss.NewStyle().MarginTop(SpacingMD)
	MarginBottomStyle = lipgloss.NewStyle().MarginBottom(SpacingMD)
	MarginLeftStyle   = lipgloss.NewStyle().MarginLeft(SpacingMD)
	MarginRightStyle  = lipgloss.NewStyle().MarginRight(SpacingMD)

	// Padding utilities
	PaddingTopStyle    = lipgloss.NewStyle().PaddingTop(SpacingMD)
	PaddingBottomStyle = lipgloss.NewStyle().PaddingBottom(SpacingMD)
	PaddingLeftStyle   = lipgloss.NewStyle().PaddingLeft(SpacingMD)
	PaddingRightStyle  = lipgloss.NewStyle().PaddingRight(SpacingMD)

	// Combined spacing
	SpacingSmallStyle  = lipgloss.NewStyle().Margin(SpacingSM).Padding(SpacingSM)
	SpacingMediumStyle = lipgloss.NewStyle().Margin(SpacingMD).Padding(SpacingMD)
	SpacingLargeStyle  = lipgloss.NewStyle().Margin(SpacingLG).Padding(SpacingLG)
)

// Grid system for responsive layouts
type GridColumn struct {
	Width   int
	Content string
	Style   lipgloss.Style
}

// Layout helper functions

// CreateContainer creates a styled container with specified dimensions
func CreateContainer(width, height int, withBorder bool) lipgloss.Style {
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Background(DarkBlue).
		Foreground(White).
		Padding(1)

	if withBorder {
		style = style.
			Border(BorderPrimaryStyle).
			BorderForeground(Gold)
	}

	return style
}

// CreateCard creates a card-like container with optional title
func CreateCard(width int, title string) lipgloss.Style {
	style := CardStyle.Copy().Width(width)

	if title != "" {
		// Add title styling - this would be combined with content
		style = style.BorderTop(true)
	}

	return style
}

// CreateGrid creates a responsive grid layout
func CreateGrid(columns []GridColumn, totalWidth int) string {
	if len(columns) == 0 {
		return ""
	}

	// Calculate available width
	availableWidth := totalWidth - (len(columns) - 1) // Account for spacing between columns

	// Auto-size columns if width not specified
	for i := range columns {
		if columns[i].Width == 0 {
			columns[i].Width = availableWidth / len(columns)
		}
	}

	// Render columns side by side
	var renderedColumns []string
	for _, col := range columns {
		columnStyle := col.Style.Copy().Width(col.Width)
		renderedColumns = append(renderedColumns, columnStyle.Render(col.Content))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedColumns...)
}

// CreateTwoColumnLayout creates a two-column layout
func CreateTwoColumnLayout(leftContent, rightContent string, totalWidth int, leftRatio float64) string {
	leftWidth := int(float64(totalWidth) * leftRatio)
	rightWidth := totalWidth - leftWidth - 2 // Account for spacing

	leftColumn := GridColumn{
		Width:   leftWidth,
		Content: leftContent,
		Style:   SecondaryContainerStyle.Copy(),
	}

	rightColumn := GridColumn{
		Width:   rightWidth,
		Content: rightContent,
		Style:   SecondaryContainerStyle.Copy(),
	}

	return CreateGrid([]GridColumn{leftColumn, rightColumn}, totalWidth)
}

// CreateThreeColumnLayout creates a three-column layout
func CreateThreeColumnLayout(leftContent, centerContent, rightContent string, totalWidth int) string {
	columnWidth := (totalWidth - 4) / 3 // Account for spacing

	columns := []GridColumn{
		{Width: columnWidth, Content: leftContent, Style: SecondaryContainerStyle.Copy()},
		{Width: columnWidth, Content: centerContent, Style: SecondaryContainerStyle.Copy()},
		{Width: columnWidth, Content: rightContent, Style: SecondaryContainerStyle.Copy()},
	}

	return CreateGrid(columns, totalWidth)
}

// Layout composition helpers

// CreatePage creates a full page layout with header, content, and footer
func CreatePage(header, content, footer string, width, height int) string {
	headerHeight := 3
	footerHeight := 3
	contentHeight := height - headerHeight - footerHeight

	headerStyled := HeaderStyle.Copy().
		Width(width).
		Height(headerHeight).
		Render(header)

	contentStyled := ContentContainerStyle.Copy().
		Width(width - 4). // Account for container padding/borders
		Height(contentHeight).
		Render(content)

	footerStyled := FooterStyle.Copy().
		Width(width).
		Height(footerHeight).
		Render(footer)

	return lipgloss.JoinVertical(lipgloss.Left,
		headerStyled,
		contentStyled,
		footerStyled,
	)
}

// CreateSidebarLayout creates a layout with sidebar and main content
func CreateSidebarLayout(sidebar, main string, totalWidth, height int) string {
	sidebarWidth := 25
	mainWidth := totalWidth - sidebarWidth - 2

	sidebarStyled := SidebarStyle.Copy().
		Width(sidebarWidth).
		Height(height).
		Render(sidebar)

	mainStyled := ContentContainerStyle.Copy().
		Width(mainWidth).
		Height(height).
		Render(main)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		sidebarStyled,
		mainStyled,
	)
}

// Responsive breakpoints
const (
	BreakpointSM = 60  // Small devices
	BreakpointMD = 80  // Medium devices
	BreakpointLG = 120 // Large devices
	BreakpointXL = 160 // Extra large devices
)

// GetResponsiveWidth returns appropriate width based on terminal size
func GetResponsiveWidth(terminalWidth int) int {
	switch {
	case terminalWidth >= BreakpointXL:
		return BreakpointXL
	case terminalWidth >= BreakpointLG:
		return BreakpointLG
	case terminalWidth >= BreakpointMD:
		return BreakpointMD
	case terminalWidth >= BreakpointSM:
		return BreakpointSM
	default:
		return terminalWidth
	}
}

// Utility for centering content
func CenterContent(content string, width int) string {
	return lipgloss.NewStyle().
		Width(width).
		Align(AlignCenter).
		Render(content)
}

// Utility for creating dividers
func CreateDivider(width int, char string, color lipgloss.Color) string {
	divider := ""
	for i := 0; i < width; i++ {
		divider += char
	}

	return lipgloss.NewStyle().
		Foreground(color).
		Render(divider)
}

// Common dividers
func GoldDivider(width int) string {
	return CreateDivider(width, "─", Gold)
}

func CyanDivider(width int) string {
	return CreateDivider(width, "━", Cyan)
}

func SubtleDivider(width int) string {
	return CreateDivider(width, "·", Gray)
}
