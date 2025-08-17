// Package ascii provides ASCII art banners for Antoine CLI
// This file implements contextual banners and decorative elements
package ascii

import (
	"fmt"
	"strings"
	"time"
)

// BannerType defines different banner styles
type BannerType string

const (
	BannerTypeWelcome  BannerType = "welcome"  // Welcome screen banner
	BannerTypeSuccess  BannerType = "success"  // Success message banner
	BannerTypeError    BannerType = "error"    // Error message banner
	BannerTypeWarning  BannerType = "warning"  // Warning message banner
	BannerTypeInfo     BannerType = "info"     // Information banner
	BannerTypeLoading  BannerType = "loading"  // Loading state banner
	BannerTypeSearch   BannerType = "search"   // Search operation banner
	BannerTypeAnalysis BannerType = "analysis" // Analysis operation banner
	BannerTypeMentor   BannerType = "mentor"   // AI Mentor banner
	BannerTypeTrends   BannerType = "trends"   // Trends analysis banner
	BannerTypeConfig   BannerType = "config"   // Configuration banner
	BannerTypeComplete BannerType = "complete" // Task completion banner
)

// BannerConfig configures banner generation
type BannerConfig struct {
	Type     BannerType
	Title    string
	Subtitle string
	Message  string
	Width    int
	Height   int
	Animated bool
	Border   bool
	Centered bool
	Icon     string
	Color    string
}

// Banner represents a contextual banner
type Banner struct {
	config BannerConfig
	art    string
}

// NewBanner creates a new banner
func NewBanner(config BannerConfig) *Banner {
	// Set defaults
	if config.Width == 0 {
		config.Width = 80
	}
	if config.Height == 0 {
		config.Height = 8
	}

	banner := &Banner{config: config}
	banner.generateArt()
	return banner
}

// generateArt generates the banner art based on type
func (b *Banner) generateArt() {
	switch b.config.Type {
	case BannerTypeWelcome:
		b.art = b.generateWelcomeBanner()
	case BannerTypeSuccess:
		b.art = b.generateSuccessBanner()
	case BannerTypeError:
		b.art = b.generateErrorBanner()
	case BannerTypeWarning:
		b.art = b.generateWarningBanner()
	case BannerTypeInfo:
		b.art = b.generateInfoBanner()
	case BannerTypeLoading:
		b.art = b.generateLoadingBanner()
	case BannerTypeSearch:
		b.art = b.generateSearchBanner()
	case BannerTypeAnalysis:
		b.art = b.generateAnalysisBanner()
	case BannerTypeMentor:
		b.art = b.generateMentorBanner()
	case BannerTypeTrends:
		b.art = b.generateTrendsBanner()
	case BannerTypeConfig:
		b.art = b.generateConfigBanner()
	case BannerTypeComplete:
		b.art = b.generateCompleteBanner()
	default:
		b.art = b.generateGenericBanner()
	}
}

// generateWelcomeBanner creates a welcome banner
func (b *Banner) generateWelcomeBanner() string {
	var lines []string

	// Welcome ASCII art
	welcomeArt := []string{
		"██╗    ██╗███████╗██╗      ██████╗ ██████╗ ███╗   ███╗███████╗",
		"██║    ██║██╔════╝██║     ██╔════╝██╔═══██╗████╗ ████║██╔════╝",
		"██║ █╗ ██║█████╗  ██║     ██║     ██║   ██║██╔████╔██║█████╗  ",
		"██║███╗██║██╔══╝  ██║     ██║     ██║   ██║██║╚██╔╝██║██╔══╝  ",
		"╚███╔███╔╝███████╗███████╗╚██████╗╚██████╔╝██║ ╚═╝ ██║███████╗",
		" ╚══╝╚══╝ ╚══════╝╚══════╝ ╚═════╝ ╚═════╝ ╚═╝     ╚═╝╚══════╝",
	}

	// Add welcome art
	for _, line := range welcomeArt {
		lines = append(lines, b.centerText(line))
	}

	// Add title and subtitle
	if b.config.Title != "" {
		lines = append(lines, "")
		lines = append(lines, b.centerText(b.config.Title))
	}

	if b.config.Subtitle != "" {
		lines = append(lines, b.centerText(b.config.Subtitle))
	}

	return b.wrapWithBorder(lines)
}

// generateSuccessBanner creates a success banner
func (b *Banner) generateSuccessBanner() string {
	var lines []string

	successArt := []string{
		"  ██████╗██╗   ██╗ ██████╗ ██████╗███████╗███████╗███████╗",
		" ██╔════╝██║   ██║██╔════╝██╔════╝██╔════╝██╔════╝██╔════╝",
		" ╚█████╗ ██║   ██║██║     ██║     █████╗  ███████╗███████╗",
		"  ╚═══██╗██║   ██║██║     ██║     ██╔══╝  ╚════██║╚════██║",
		" ██████╔╝╚██████╔╝╚██████╗╚██████╗███████╗███████║███████║",
		" ╚═════╝  ╚═════╝  ╚═════╝ ╚═════╝╚══════╝╚══════╝╚══════╝",
	}

	// Add success icon
	lines = append(lines, b.centerText("✅ SUCCESS! ✅"))
	lines = append(lines, "")

	// Add success art
	for _, line := range successArt {
		lines = append(lines, b.centerText(line))
	}

	if b.config.Message != "" {
		lines = append(lines, "")
		lines = append(lines, b.centerText(b.config.Message))
	}

	return b.wrapWithBorder(lines)
}

// generateErrorBanner creates an error banner
func (b *Banner) generateErrorBanner() string {
	var lines []string

	errorArt := []string{
		" ███████╗██████╗ ██████╗  ██████╗ ██████╗ ",
		" ██╔════╝██╔══██╗██╔══██╗██╔═══██╗██╔══██╗",
		" █████╗  ██████╔╝██████╔╝██║   ██║██████╔╝",
		" ██╔══╝  ██╔══██╗██╔══██╗██║   ██║██╔══██╗",
		" ███████╗██║  ██║██║  ██║╚██████╔╝██║  ██║",
		" ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝",
	}

	// Add error icon
	lines = append(lines, b.centerText("❌ ERROR! ❌"))
	lines = append(lines, "")

	// Add error art
	for _, line := range errorArt {
		lines = append(lines, b.centerText(line))
	}

	if b.config.Message != "" {
		lines = append(lines, "")
		lines = append(lines, b.centerText(b.config.Message))
	}

	return b.wrapWithBorder(lines)
}

// generateWarningBanner creates a warning banner
func (b *Banner) generateWarningBanner() string {
	var lines []string

	warningIcon := []string{
		"        ⚠️  WARNING  ⚠️",
		"    ╔══════════════════╗",
		"    ║   ATTENTION!     ║",
		"    ╚══════════════════╝",
	}

	for _, line := range warningIcon {
		lines = append(lines, b.centerText(line))
	}

	if b.config.Message != "" {
		lines = append(lines, "")
		lines = append(lines, b.centerText(b.config.Message))
	}

	return b.wrapWithBorder(lines)
}

// generateInfoBanner creates an info banner
func (b *Banner) generateInfoBanner() string {
	var lines []string

	infoIcon := []string{
		"        ℹ️  INFORMATION  ℹ️",
		"    ╔══════════════════╗",
		"    ║      NOTICE      ║",
		"    ╚══════════════════╝",
	}

	for _, line := range infoIcon {
		lines = append(lines, b.centerText(line))
	}

	if b.config.Message != "" {
		lines = append(lines, "")
		lines = append(lines, b.centerText(b.config.Message))
	}

	return b.wrapWithBorder(lines)
}

// generateLoadingBanner creates a loading banner
func (b *Banner) generateLoadingBanner() string {
	var lines []string

	// Animated loading elements
	spinner := "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"
	elapsed := time.Now().UnixMilli() / 200
	spinnerChar := string([]rune(spinner)[elapsed%int64(len([]rune(spinner)))])

	loadingText := fmt.Sprintf("%s LOADING %s", spinnerChar, spinnerChar)
	lines = append(lines, b.centerText(loadingText))

	// Progress bar
	progressBar := "▓▓▓▓▓▓▓░░░░░░░░░░░░░"
	lines = append(lines, b.centerText(progressBar))

	if b.config.Message != "" {
		lines = append(lines, "")
		lines = append(lines, b.centerText(b.config.Message))
	}

	return b.wrapWithBorder(lines)
}

// generateSearchBanner creates a search banner
func (b *Banner) generateSearchBanner() string {
	var lines []string

	searchArt := []string{
		"  🔍 SEARCH MODE ACTIVATED 🔍",
		"     ╔═══════════════════╗",
		"     ║   SCANNING...     ║",
		"     ╚═══════════════════╝",
	}

	for _, line := range searchArt {
		lines = append(lines, b.centerText(line))
	}

	if b.config.Message != "" {
		lines = append(lines, "")
		lines = append(lines, b.centerText(b.config.Message))
	}

	return b.wrapWithBorder(lines)
}

// generateAnalysisBanner creates an analysis banner
func (b *Banner) generateAnalysisBanner() string {
	var lines []string

	analysisArt := []string{
		" █████╗ ███╗   ██╗ █████╗ ██╗  ██╗   ██╗███████╗██╗███████╗",
		"██╔══██╗████╗  ██║██╔══██╗██║  ╚██╗ ██╔╝██╔════╝██║██╔════╝",
		"███████║██╔██╗ ██║███████║██║   ╚████╔╝ ███████╗██║███████╗",
		"██╔══██║██║╚██╗██║██╔══██║██║    ╚██╔╝  ╚════██║██║╚════██║",
		"██║  ██║██║ ╚████║██║  ██║███████╗██║   ███████║██║███████║",
		"╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝╚══════╝╚═╝   ╚══════╝╚═╝╚══════╝",
	}

	lines = append(lines, b.centerText("🔬 DEEP ANALYSIS MODE 🔬"))
	lines = append(lines, "")

	for _, line := range analysisArt {
		lines = append(lines, b.centerText(line))
	}

	if b.config.Message != "" {
		lines = append(lines, "")
		lines = append(lines, b.centerText(b.config.Message))
	}

	return b.wrapWithBorder(lines)
}

// generateMentorBanner creates a mentor banner
func (b *Banner) generateMentorBanner() string {
	var lines []string

	mentorArt := []string{
		"🤖 AI MENTOR ACTIVATED 🤖",
		"     ╔═══════════════════╗",
		"     ║  READY TO HELP!   ║",
		"     ╚═══════════════════╝",
		"       YOUR GUIDE TO",
		"      HACKATHON SUCCESS",
	}

	for _, line := range mentorArt {
		lines = append(lines, b.centerText(line))
	}

	if b.config.Message != "" {
		lines = append(lines, "")
		lines = append(lines, b.centerText(b.config.Message))
	}

	return b.wrapWithBorder(lines)
}

// generateTrendsBanner creates a trends banner
func (b *Banner) generateTrendsBanner() string {
	var lines []string

	trendsArt := []string{
		"📈 TECHNOLOGY TRENDS 📈",
		"     ╔═══════════════════╗",
		"     ║   MARKET PULSE    ║",
		"     ╚═══════════════════╝",
		"    ANALYZING PATTERNS...",
	}

	for _, line := range trendsArt {
		lines = append(lines, b.centerText(line))
	}

	if b.config.Message != "" {
		lines = append(lines, "")
		lines = append(lines, b.centerText(b.config.Message))
	}

	return b.wrapWithBorder(lines)
}

// generateConfigBanner creates a config banner
func (b *Banner) generateConfigBanner() string {
	var lines []string

	configArt := []string{
		"⚙️  CONFIGURATION ⚙️",
		"     ╔═══════════════════╗",
		"     ║   SETTINGS MENU   ║",
		"     ╚═══════════════════╝",
	}

	for _, line := range configArt {
		lines = append(lines, b.centerText(line))
	}

	if b.config.Message != "" {
		lines = append(lines, "")
		lines = append(lines, b.centerText(b.config.Message))
	}

	return b.wrapWithBorder(lines)
}

// generateCompleteBanner creates a completion banner
func (b *Banner) generateCompleteBanner() string {
	var lines []string

	completeArt := []string{
		" ██████╗ ██████╗ ███╗   ███╗██████╗ ██╗     ███████╗████████╗███████╗",
		"██╔════╝██╔═══██╗████╗ ████║██╔══██╗██║     ██╔════╝╚══██╔══╝██╔════╝",
		"██║     ██║   ██║██╔████╔██║██████╔╝██║     █████╗     ██║   █████╗  ",
		"██║     ██║   ██║██║╚██╔╝██║██╔═══╝ ██║     ██╔══╝     ██║   ██╔══╝  ",
		"╚██████╗╚██████╔╝██║ ╚═╝ ██║██║     ███████╗███████╗   ██║   ███████╗",
		" ╚═════╝ ╚═════╝ ╚═╝     ╚═╝╚═╝     ╚══════╝╚══════╝   ╚═╝   ╚══════╝",
	}

	lines = append(lines, b.centerText("🎉 TASK COMPLETED! 🎉"))
	lines = append(lines, "")

	for _, line := range completeArt {
		lines = append(lines, b.centerText(line))
	}

	if b.config.Message != "" {
		lines = append(lines, "")
		lines = append(lines, b.centerText(b.config.Message))
	}

	return b.wrapWithBorder(lines)
}

// generateGenericBanner creates a generic banner
func (b *Banner) generateGenericBanner() string {
	var lines []string

	if b.config.Icon != "" {
		lines = append(lines, b.centerText(b.config.Icon))
		lines = append(lines, "")
	}

	if b.config.Title != "" {
		lines = append(lines, b.centerText(b.config.Title))
	}

	if b.config.Subtitle != "" {
		lines = append(lines, b.centerText(b.config.Subtitle))
	}

	if b.config.Message != "" {
		lines = append(lines, "")
		lines = append(lines, b.centerText(b.config.Message))
	}

	return b.wrapWithBorder(lines)
}

// centerText centers text within the banner width
func (b *Banner) centerText(text string) string {
	if !b.config.Centered {
		return text
	}

	textLen := len(text)
	if textLen >= b.config.Width {
		return text
	}

	padding := (b.config.Width - textLen) / 2
	return strings.Repeat(" ", padding) + text
}

// wrapWithBorder wraps content with a border if enabled
func (b *Banner) wrapWithBorder(lines []string) string {
	if !b.config.Border {
		return strings.Join(lines, "\n")
	}

	// Find the maximum line length
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// Ensure minimum width
	if maxLen < b.config.Width-4 {
		maxLen = b.config.Width - 4
	}

	var result []string

	// Top border
	topBorder := "╔" + strings.Repeat("═", maxLen+2) + "╗"
	result = append(result, topBorder)

	// Content lines with side borders
	for _, line := range lines {
		padding := maxLen - len(line)
		contentLine := "║ " + line + strings.Repeat(" ", padding) + " ║"
		result = append(result, contentLine)
	}

	// Bottom border
	bottomBorder := "╚" + strings.Repeat("═", maxLen+2) + "╝"
	result = append(result, bottomBorder)

	return strings.Join(result, "\n")
}

// Render returns the banner art
func (b *Banner) Render() string {
	return b.art
}

// Update updates banner content and regenerates art
func (b *Banner) Update(config BannerConfig) {
	// Preserve dimensions if not specified
	if config.Width == 0 {
		config.Width = b.config.Width
	}
	if config.Height == 0 {
		config.Height = b.config.Height
	}

	b.config = config
	b.generateArt()
}

// Predefined banner functions for common use cases

// WelcomeBanner creates a welcome banner for Antoine CLI
func WelcomeBanner(width int) string {
	banner := NewBanner(BannerConfig{
		Type:     BannerTypeWelcome,
		Title:    "ANTOINE CLI",
		Subtitle: "Your Ultimate Hackathon Mentor",
		Width:    width,
		Border:   true,
		Centered: true,
	})
	return banner.Render()
}

// SuccessBanner creates a success banner with message
func SuccessBanner(message string, width int) string {
	banner := NewBanner(BannerConfig{
		Type:     BannerTypeSuccess,
		Message:  message,
		Width:    width,
		Border:   true,
		Centered: true,
	})
	return banner.Render()
}

// ErrorBanner creates an error banner with message
func ErrorBanner(message string, width int) string {
	banner := NewBanner(BannerConfig{
		Type:     BannerTypeError,
		Message:  message,
		Width:    width,
		Border:   true,
		Centered: true,
	})
	return banner.Render()
}

// WarningBanner creates a warning banner with message
func WarningBanner(message string, width int) string {
	banner := NewBanner(BannerConfig{
		Type:     BannerTypeWarning,
		Message:  message,
		Width:    width,
		Border:   true,
		Centered: true,
	})
	return banner.Render()
}

// InfoBanner creates an info banner with message
func InfoBanner(message string, width int) string {
	banner := NewBanner(BannerConfig{
		Type:     BannerTypeInfo,
		Message:  message,
		Width:    width,
		Border:   true,
		Centered: true,
	})
	return banner.Render()
}

// LoadingBanner creates a loading banner with message
func LoadingBanner(message string, width int) string {
	banner := NewBanner(BannerConfig{
		Type:     BannerTypeLoading,
		Message:  message,
		Width:    width,
		Border:   true,
		Centered: true,
		Animated: true,
	})
	return banner.Render()
}

// SearchBanner creates a search operation banner
func SearchBanner(query string, width int) string {
	message := fmt.Sprintf("Searching for: %s", query)
	banner := NewBanner(BannerConfig{
		Type:     BannerTypeSearch,
		Message:  message,
		Width:    width,
		Border:   true,
		Centered: true,
	})
	return banner.Render()
}

// AnalysisBanner creates an analysis operation banner
func AnalysisBanner(target string, width int) string {
	message := fmt.Sprintf("Analyzing: %s", target)
	banner := NewBanner(BannerConfig{
		Type:     BannerTypeAnalysis,
		Message:  message,
		Width:    width,
		Border:   true,
		Centered: true,
	})
	return banner.Render()
}

// MentorBanner creates an AI mentor banner
func MentorBanner(sessionType string, width int) string {
	message := fmt.Sprintf("Session Type: %s", sessionType)
	banner := NewBanner(BannerConfig{
		Type:     BannerTypeMentor,
		Message:  message,
		Width:    width,
		Border:   true,
		Centered: true,
	})
	return banner.Render()
}

// TrendsBanner creates a trends analysis banner
func TrendsBanner(timeframe string, width int) string {
	message := fmt.Sprintf("Timeframe: %s", timeframe)
	banner := NewBanner(BannerConfig{
		Type:     BannerTypeTrends,
		Message:  message,
		Width:    width,
		Border:   true,
		Centered: true,
	})
	return banner.Render()
}

// ConfigBanner creates a configuration banner
func ConfigBanner(section string, width int) string {
	message := fmt.Sprintf("Configuring: %s", section)
	banner := NewBanner(BannerConfig{
		Type:     BannerTypeConfig,
		Message:  message,
		Width:    width,
		Border:   true,
		Centered: true,
	})
	return banner.Render()
}

// CompleteBanner creates a task completion banner
func CompleteBanner(task string, width int) string {
	message := fmt.Sprintf("Task completed: %s", task)
	banner := NewBanner(BannerConfig{
		Type:     BannerTypeComplete,
		Message:  message,
		Width:    width,
		Border:   true,
		Centered: true,
	})
	return banner.Render()
}

// CustomBanner creates a custom banner with specified parameters
func CustomBanner(title, subtitle, message, icon string, width int, border bool) string {
	banner := NewBanner(BannerConfig{
		Type:     BannerType("custom"),
		Title:    title,
		Subtitle: subtitle,
		Message:  message,
		Icon:     icon,
		Width:    width,
		Border:   border,
		Centered: true,
	})
	return banner.Render()
}

// StatusBanner creates a status banner with icon and color coding
func StatusBanner(status, message string, width int) string {
	var bannerType BannerType
	var icon string

	switch strings.ToLower(status) {
	case "success", "complete", "done":
		bannerType = BannerTypeSuccess
		icon = "✅"
	case "error", "failed", "failure":
		bannerType = BannerTypeError
		icon = "❌"
	case "warning", "warn":
		bannerType = BannerTypeWarning
		icon = "⚠️"
	case "info", "information":
		bannerType = BannerTypeInfo
		icon = "ℹ️"
	case "loading", "processing":
		bannerType = BannerTypeLoading
		icon = "⏳"
	default:
		bannerType = BannerTypeInfo
		icon = "📢"
	}

	banner := NewBanner(BannerConfig{
		Type:     bannerType,
		Message:  message,
		Icon:     icon,
		Width:    width,
		Border:   true,
		Centered: true,
	})
	return banner.Render()
}

// ProgressBanner creates a progress banner with percentage
func ProgressBanner(operation string, progress float64, width int) string {
	// Create progress bar
	barWidth := width - 20
	filled := int(progress * float64(barWidth))
	empty := barWidth - filled

	progressBar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	percentage := fmt.Sprintf("%.1f%%", progress*100)

	message := fmt.Sprintf("%s\n%s %s", operation, progressBar, percentage)

	banner := NewBanner(BannerConfig{
		Type:     BannerTypeLoading,
		Message:  message,
		Width:    width,
		Border:   true,
		Centered: true,
	})
	return banner.Render()
}

// MultiBanner creates a banner with multiple messages
func MultiBanner(messages []string, bannerType BannerType, width int) string {
	message := strings.Join(messages, "\n")

	banner := NewBanner(BannerConfig{
		Type:     bannerType,
		Message:  message,
		Width:    width,
		Border:   true,
		Centered: true,
	})
	return banner.Render()
}

// TimedBanner creates a banner with timestamp
func TimedBanner(title, message string, bannerType BannerType, width int) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	timedMessage := fmt.Sprintf("%s\n[%s]", message, timestamp)

	banner := NewBanner(BannerConfig{
		Type:     bannerType,
		Title:    title,
		Message:  timedMessage,
		Width:    width,
		Border:   true,
		Centered: true,
	})
	return banner.Render()
}
