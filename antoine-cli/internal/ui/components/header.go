// Package components provides reusable UI components for Antoine CLI
// This file implements the header component with ASCII art and branding
package components

import (
	"fmt"
	"strings"
	"time"

	"antoine-cli/internal/ui/styles"
	"antoine-cli/pkg/ascii"
	"github.com/charmbracelet/lipgloss"
)

// HeaderType defines different header styles
type HeaderType string

const (
	HeaderTypeFull    HeaderType = "full"    // Full logo with subtitle
	HeaderTypeCompact HeaderType = "compact" // Compact logo only
	HeaderTypeMinimal HeaderType = "minimal" // Text only
	HeaderTypeBanner  HeaderType = "banner"  // Banner with context
)

// HeaderConfig configures the header component
type HeaderConfig struct {
	Type        HeaderType
	Title       string
	Subtitle    string
	Context     string
	ShowTime    bool
	ShowVersion bool
	Animated    bool
	Width       int
	Theme       string
}

// Header represents the header component
type Header struct {
	config HeaderConfig
	style  lipgloss.Style
}

// NewHeader creates a new header component
func NewHeader(config HeaderConfig) *Header {
	// Set defaults
	if config.Width == 0 {
		config.Width = 120
	}
	if config.Theme == "" {
		config.Theme = "dark"
	}

	// Create base style
	baseStyle := lipgloss.NewStyle().
		Width(config.Width).
		Align(lipgloss.Center).
		Padding(1, 0)

	return &Header{
		config: config,
		style:  baseStyle,
	}
}

// Render renders the header component
func (h *Header) Render() string {
	switch h.config.Type {
	case HeaderTypeFull:
		return h.renderFullHeader()
	case HeaderTypeCompact:
		return h.renderCompactHeader()
	case HeaderTypeMinimal:
		return h.renderMinimalHeader()
	case HeaderTypeBanner:
		return h.renderBannerHeader()
	default:
		return h.renderFullHeader()
	}
}

// renderFullHeader renders the full header with logo and details
func (h *Header) renderFullHeader() string {
	var sections []string

	// ASCII Art Logo
	logo := ascii.GetLogo(80)
	if h.config.Animated {
		logo = ascii.GetAnimatedLogo(80, time.Now())
	}

	styledLogo := styles.BrandStyle.Render(logo)
	sections = append(sections, styledLogo)

	// Title and subtitle
	if h.config.Title != "" {
		title := styles.H1Style.Render(h.config.Title)
		sections = append(sections, title)
	}

	if h.config.Subtitle != "" {
		subtitle := styles.H2Style.Render(h.config.Subtitle)
		sections = append(sections, subtitle)
	}

	// Additional info line
	var infoItems []string

	if h.config.ShowVersion {
		version := styles.AccentStyle.Render("v1.0.0")
		infoItems = append(infoItems, version)
	}

	if h.config.ShowTime {
		currentTime := time.Now().Format("15:04:05")
		timeStr := styles.BodySecondaryStyle.Render(currentTime)
		infoItems = append(infoItems, timeStr)
	}

	if h.config.Context != "" {
		contextStr := styles.InfoStyle.Render(h.config.Context)
		infoItems = append(infoItems, contextStr)
	}

	if len(infoItems) > 0 {
		infoLine := strings.Join(infoItems, " • ")
		sections = append(sections, infoLine)
	}

	// Add separator
	separator := styles.GoldDivider(h.config.Width - 4)
	sections = append(sections, separator)

	content := strings.Join(sections, "\n")
	return h.style.Render(content)
}

// renderCompactHeader renders a compact version of the header
func (h *Header) renderCompactHeader() string {
	var sections []string

	// Compact ASCII logo
	logo := ascii.GetLogo(40)
	if h.config.Animated {
		logo = ascii.GetAnimatedLogo(40, time.Now())
	}

	// Horizontal layout: logo + title + info
	var headerLine []string

	// Logo (left)
	styledLogo := styles.BrandStyle.Render(logo)
	headerLine = append(headerLine, styledLogo)

	// Title and context (center)
	var centerContent []string
	if h.config.Title != "" {
		centerContent = append(centerContent, styles.H2Style.Render(h.config.Title))
	}
	if h.config.Context != "" {
		centerContent = append(centerContent, styles.BodySecondaryStyle.Render(h.config.Context))
	}

	if len(centerContent) > 0 {
		center := strings.Join(centerContent, " | ")
		headerLine = append(headerLine, center)
	}

	// Info (right)
	var infoItems []string
	if h.config.ShowVersion {
		infoItems = append(infoItems, styles.AccentStyle.Render("v1.0.0"))
	}
	if h.config.ShowTime {
		infoItems = append(infoItems, styles.BodySecondaryStyle.Render(time.Now().Format("15:04")))
	}

	if len(infoItems) > 0 {
		info := strings.Join(infoItems, " ")
		headerLine = append(headerLine, info)
	}

	// Join with spacing
	content := strings.Join(headerLine, "   ")
	sections = append(sections, content)

	// Add thin separator
	separator := styles.SubtleDivider(h.config.Width - 4)
	sections = append(sections, separator)

	finalContent := strings.Join(sections, "\n")
	return h.style.Render(finalContent)
}

// renderMinimalHeader renders a text-only minimal header
func (h *Header) renderMinimalHeader() string {
	var parts []string

	// Antoine branding
	brand := styles.BrandStyle.Render("ANTOINE")
	parts = append(parts, brand)

	if h.config.Title != "" {
		parts = append(parts, styles.H3Style.Render(h.config.Title))
	}

	if h.config.Context != "" {
		parts = append(parts, styles.BodySecondaryStyle.Render(fmt.Sprintf("(%s)", h.config.Context)))
	}

	content := strings.Join(parts, " ")

	// Add minimal separator
	separator := styles.SubtleDivider(len(content))
	finalContent := content + "\n" + separator

	return h.style.Render(finalContent)
}

// renderBannerHeader renders a banner-style header with contextual information
func (h *Header) renderBannerHeader() string {
	var sections []string

	// Banner top border
	topBorder := styles.GoldDivider(h.config.Width - 4)
	sections = append(sections, topBorder)

	// Main banner content
	//var bannerContent []string

	// Left side: Logo or brand
	leftContent := ascii.GetLogo(40)
	if leftContent == "" {
		leftContent = styles.BrandStyle.Render("🤖 ANTOINE")
	} else {
		leftContent = styles.BrandStyle.Render(leftContent)
	}

	// Center: Title and subtitle
	var centerLines []string
	if h.config.Title != "" {
		centerLines = append(centerLines, styles.H2Style.Render(h.config.Title))
	}
	if h.config.Subtitle != "" {
		centerLines = append(centerLines, styles.BodyStyle.Render(h.config.Subtitle))
	}
	centerContent := strings.Join(centerLines, "\n")

	// Right side: Context and status
	var rightLines []string
	if h.config.Context != "" {
		rightLines = append(rightLines, styles.InfoStyle.Render(h.config.Context))
	}
	if h.config.ShowTime {
		rightLines = append(rightLines, styles.BodySecondaryStyle.Render(time.Now().Format("Mon 15:04")))
	}
	rightContent := strings.Join(rightLines, "\n")

	// Create three-column layout
	leftCol := lipgloss.NewStyle().Width(30).Render(leftContent)
	centerCol := lipgloss.NewStyle().Width(h.config.Width - 64).Align(lipgloss.Center).Render(centerContent)
	rightCol := lipgloss.NewStyle().Width(30).Align(lipgloss.Right).Render(rightContent)

	bannerRow := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, centerCol, rightCol)
	sections = append(sections, bannerRow)

	// Banner bottom border
	bottomBorder := styles.GoldDivider(h.config.Width - 4)
	sections = append(sections, bottomBorder)

	content := strings.Join(sections, "\n")
	return h.style.Render(content)
}

// Update updates the header configuration
func (h *Header) Update(config HeaderConfig) {
	// Preserve width if not specified
	if config.Width == 0 {
		config.Width = h.config.Width
	}

	// Preserve theme if not specified
	if config.Theme == "" {
		config.Theme = h.config.Theme
	}

	h.config = config

	// Update style width
	h.style = h.style.Width(config.Width)
}

// SetWidth sets the header width
func (h *Header) SetWidth(width int) {
	h.config.Width = width
	h.style = h.style.Width(width)
}

// SetContext sets the context information
func (h *Header) SetContext(context string) {
	h.config.Context = context
}

// SetTitle sets the header title
func (h *Header) SetTitle(title string) {
	h.config.Title = title
}

// SetSubtitle sets the header subtitle
func (h *Header) SetSubtitle(subtitle string) {
	h.config.Subtitle = subtitle
}

// ToggleAnimation toggles animation on/off
func (h *Header) ToggleAnimation() {
	h.config.Animated = !h.config.Animated
}

// GetHeight returns the rendered height of the header
func (h *Header) GetHeight() int {
	rendered := h.Render()
	return strings.Count(rendered, "\n") + 1
}

// HeaderBuilder provides a fluent interface for building headers
type HeaderBuilder struct {
	config HeaderConfig
}

// NewHeaderBuilder creates a new header builder
func NewHeaderBuilder() *HeaderBuilder {
	return &HeaderBuilder{
		config: HeaderConfig{
			Type:     HeaderTypeFull,
			Width:    120,
			Theme:    "dark",
			Animated: true,
		},
	}
}

// Type sets the header type
func (hb *HeaderBuilder) Type(headerType HeaderType) *HeaderBuilder {
	hb.config.Type = headerType
	return hb
}

// Title sets the header title
func (hb *HeaderBuilder) Title(title string) *HeaderBuilder {
	hb.config.Title = title
	return hb
}

// Subtitle sets the header subtitle
func (hb *HeaderBuilder) Subtitle(subtitle string) *HeaderBuilder {
	hb.config.Subtitle = subtitle
	return hb
}

// Context sets the context information
func (hb *HeaderBuilder) Context(context string) *HeaderBuilder {
	hb.config.Context = context
	return hb
}

// Width sets the header width
func (hb *HeaderBuilder) Width(width int) *HeaderBuilder {
	hb.config.Width = width
	return hb
}

// ShowTime enables/disables time display
func (hb *HeaderBuilder) ShowTime(show bool) *HeaderBuilder {
	hb.config.ShowTime = show
	return hb
}

// ShowVersion enables/disables version display
func (hb *HeaderBuilder) ShowVersion(show bool) *HeaderBuilder {
	hb.config.ShowVersion = show
	return hb
}

// Animated enables/disables animations
func (hb *HeaderBuilder) Animated(animated bool) *HeaderBuilder {
	hb.config.Animated = animated
	return hb
}

// Theme sets the color theme
func (hb *HeaderBuilder) Theme(theme string) *HeaderBuilder {
	hb.config.Theme = theme
	return hb
}

// Build creates the header component
func (hb *HeaderBuilder) Build() *Header {
	return NewHeader(hb.config)
}

// Predefined header configurations for common use cases

// WelcomeHeader creates a welcome header for the main dashboard
func WelcomeHeader() *Header {
	return NewHeaderBuilder().
		Type(HeaderTypeFull).
		Title("Welcome to Antoine CLI").
		Subtitle("Your Ultimate Hackathon Mentor").
		ShowTime(true).
		ShowVersion(true).
		Animated(true).
		Build()
}

// SearchHeader creates a header for search operations
func SearchHeader(searchType string) *Header {
	return NewHeaderBuilder().
		Type(HeaderTypeCompact).
		Title("Search").
		Context(fmt.Sprintf("Searching %s", searchType)).
		ShowTime(true).
		Build()
}

// AnalysisHeader creates a header for analysis operations
func AnalysisHeader(target string) *Header {
	return NewHeaderBuilder().
		Type(HeaderTypeBanner).
		Title("Repository Analysis").
		Subtitle(fmt.Sprintf("Analyzing: %s", target)).
		Context("Deep Analysis Mode").
		ShowTime(true).
		Build()
}

// MentorHeader creates a header for mentor sessions
func MentorHeader(sessionType string) *Header {
	return NewHeaderBuilder().
		Type(HeaderTypeCompact).
		Title("AI Mentor").
		Context(fmt.Sprintf("%s Session", sessionType)).
		ShowTime(true).
		Animated(true).
		Build()
}

// TrendsHeader creates a header for trends analysis
func TrendsHeader(timeframe string) *Header {
	return NewHeaderBuilder().
		Type(HeaderTypeBanner).
		Title("Technology Trends").
		Subtitle("Market Analysis & Insights").
		Context(fmt.Sprintf("Timeframe: %s", timeframe)).
		ShowTime(true).
		Build()
}

// ConfigHeader creates a header for configuration management
func ConfigHeader() *Header {
	return NewHeaderBuilder().
		Type(HeaderTypeMinimal).
		Title("Configuration").
		Context("Settings Management").
		Build()
}

// ErrorHeader creates a header for error states
func ErrorHeader(errorType string) *Header {
	return NewHeaderBuilder().
		Type(HeaderTypeCompact).
		Title("Error").
		Context(errorType).
		ShowTime(true).
		Animated(false).
		Build()
}

// LoadingHeader creates a header for loading states
func LoadingHeader(operation string) *Header {
	return NewHeaderBuilder().
		Type(HeaderTypeCompact).
		Title("Loading").
		Context(operation).
		ShowTime(true).
		Animated(true).
		Build()
}

// ResponsiveHeader creates a header that adapts to terminal width
func ResponsiveHeader(width int, title string) *Header {
	var headerType HeaderType

	switch {
	case width < 60:
		headerType = HeaderTypeMinimal
	case width < 100:
		headerType = HeaderTypeCompact
	default:
		headerType = HeaderTypeFull
	}

	return NewHeaderBuilder().
		Type(headerType).
		Title(title).
		Width(width).
		ShowTime(width >= 80).
		ShowVersion(width >= 100).
		Animated(width >= 80).
		Build()
}

// ContextualHeader creates a header based on current context
func ContextualHeader(context, operation string, width int) *Header {
	baseHeader := ResponsiveHeader(width, operation)
	baseHeader.SetContext(context)
	return baseHeader
}
