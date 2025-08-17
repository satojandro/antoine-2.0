// Package components provides reusable UI components for Antoine CLI
// This file implements progress bars and progress indicators
package components

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"antoine-cli/internal/ui/styles"
	"antoine-cli/internal/utils"
)

// ProgressType defines different progress bar styles
type ProgressType string

const (
	ProgressTypeBar      ProgressType = "bar"      // Classic progress bar
	ProgressTypeCircle   ProgressType = "circle"   // Circular progress
	ProgressTypeDots     ProgressType = "dots"     // Dotted progress
	ProgressTypeBlocks   ProgressType = "blocks"   // Block-based progress
	ProgressTypeWave     ProgressType = "wave"     // Wave animation
	ProgressTypeGradient ProgressType = "gradient" // Gradient progress bar
	ProgressTypeSteps    ProgressType = "steps"    // Step-based progress
	ProgressTypeAntoine  ProgressType = "antoine"  // Antoine-themed progress
)

// ProgressConfig configures the progress component
type ProgressConfig struct {
	Type          ProgressType
	Width         int
	Height        int
	ShowPercent   bool
	ShowETA       bool
	ShowRate      bool
	ShowLabel     bool
	Label         string
	Color         lipgloss.Color
	BackgroundColor lipgloss.Color
	TextColor     lipgloss.Color
	Animated      bool
	Centered      bool
	BorderStyle   lipgloss.Border
}

// Progress represents a progress indicator component
type Progress struct {
	config    ProgressConfig
	value     float64        // Current progress (0.0 - 1.0)
	startTime time.Time
	lastUpdate time.Time
	totalWork  int64
	doneWork   int64
	rate       float64        // Items per second
	style     lipgloss.Style
}

// NewProgress creates a new progress component
func NewProgress(config ProgressConfig) *Progress {
	// Set defaults
	if config.Width == 0 {
		config.Width = 50
	}
	if config.Height == 0 {
		config.Height = 1
	}
	if config.Color == "" {
		config.Color = styles.Gold
	}
	if config.BackgroundColor == "" {
		config.BackgroundColor = styles.DarkGray
	}
	if config.TextColor == "" {
		config.TextColor = styles.White
	}

	// Create base style
	baseStyle := lipgloss.NewStyle()
	if config.Centered {
		baseStyle = baseStyle.Align(lipgloss.Center)
	}

	return &Progress{
		config:    config,
		startTime: time.Now(),
		lastUpdate: time.Now(),
		style:     baseStyle,
	}
}

// Render renders the progress component
func (p *Progress) Render() string {
	switch p.config.Type {
	case ProgressTypeBar:
		return p.renderBar()
	case ProgressTypeCircle:
		return p.renderCircle()
	case ProgressTypeDots:
		return p.renderDots()
	case ProgressTypeBlocks:
		return p.renderBlocks()
	case ProgressTypeWave:
		return p.renderWave()
	case ProgressTypeGradient:
		return p.renderGradient()
	case ProgressTypeSteps:
		return p.renderSteps()
	case ProgressTypeAntoine:
		return p.renderAntoine()
	default:
		return p.renderBar()
	}
}

// renderBar renders a classic progress bar
func (p *Progress) renderBar() string {
	filledWidth := int(float64(p.config.Width) * p.value)
	emptyWidth := p.config.Width - filledWidth

	// Create filled and empty sections
	filled := strings.Repeat("█", filledWidth)
	empty := strings.Repeat("░", emptyWidth)

	// Style the sections
	filledStyled := lipgloss.NewStyle().Foreground(p.config.Color).Render(filled)
	emptyStyled := lipgloss.NewStyle().Foreground(p.config.BackgroundColor).Render(empty)

	bar := filledStyled + emptyStyled

	// Add borders if specified
	if p.config.BorderStyle.Left != "" {
		borderStyle := lipgloss.NewStyle().
			Border(p.config.BorderStyle).
			BorderForeground(p.config.Color).
			Width(p.config.Width)
		bar = borderStyle.Render(bar)
	}

	return p.addMetadata(bar)
}

// renderCircle renders a circular progress indicator
func (p *Progress) renderCircle() string {
	// Simple circle representation using Unicode
	percentage := p.value * 100

	var circle string
	switch {
	case percentage < 12.5:
		circle = "○"
	case percentage < 25:
		circle = "◔"
	case percentage < 37.5:
		circle = "◑"
	case percentage < 50:
		circle = "◕"
	case percentage < 62.5:
		circle = "◕"
	case percentage < 75:
		circle = "◑"
	case percentage < 87.5:
		circle = "◔"
	default:
		circle = "●"
	}

	styledCircle := lipgloss.NewStyle().
		Foreground(p.config.Color).
		Render(circle)

	return p.addMetadata(styledCircle)
}

// renderDots renders a dotted progress indicator
func (p *Progress) renderDots() string {
	dots := int(math.Ceil(float64(p.config.Width) / 2))
	filledDots := int(float64(dots) * p.value)

	var result strings.Builder
	for i := 0; i < dots; i++ {
		if i < filledDots {
			result.WriteString("●")
		} else {
			result.WriteString("○")
		}
		if i < dots-1 {
			result.WriteString(" ")
		}
	}

	styled := lipgloss.NewStyle().
		Foreground(p.config.Color).
		Render(result.String())

	return p.addMetadata(styled)
}

// renderBlocks renders a block-based progress bar
func (p *Progress) renderBlocks() string {
	blocks := []string{"░", "▏", "▎", "▍", "▌", "▋", "▊", "▉", "█"}
	totalBlocks := p.config.Width
	progress := p.value * float64(totalBlocks)

	var result strings.Builder
	for i := 0; i < totalBlocks; i++ {
		remaining := progress - float64(i)
		if remaining <= 0 {
			result.WriteString(blocks[0]) // Empty
		} else if remaining >= 1 {
			result.WriteString(blocks[8]) // Full
		} else {
			// Partial block
			blockIndex := int(remaining * 8)
			if blockIndex >= len(blocks) {
				blockIndex = len(blocks) - 1
			}
			result.WriteString(blocks[blockIndex])
		}
	}

	styled := lipgloss.NewStyle().
		Foreground(p.config.Color).
		Render(result.String())

	return p.addMetadata(styled)
}

// renderWave renders a wave animation progress
func (p *Progress) renderWave() string {
	waveChars := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	time := time.Since(p.startTime).Seconds()

	var result strings.Builder
	for i := 0; i < p.config.Width; i++ {
		// Create wave effect
		wavePos := math.Sin(float64(i)*0.5 + time*3) * 0.5 + 0.5

		// Modulate by progress
		intensity := wavePos * p.value
		charIndex := int(intensity * float64(len(waveChars)-1))
		if charIndex >= len(waveChars) {
			charIndex = len(waveChars) - 1
		}

		result.WriteString(waveChars[charIndex])
	}

	styled := lipgloss.NewStyle().
		Foreground(p.config.Color).
		Render(result.String())

	return p.addMetadata(styled)
}

// renderGradient renders a gradient progress bar
func (p *Progress) renderGradient() string {
	// Simulate gradient using different Unicode characters
	gradientChars := []string{"░", "▒", "▓", "█"}
	filledWidth := int(float64(p.config.Width) * p.value)

	var result strings.Builder
	for i := 0; i < p.config.Width; i++ {
		if i < filledWidth {
			// Filled area with gradient effect
			gradientPos := float64(i) / float64(filledWidth)
			charIndex := int(gradientPos * float64(len(gradientChars)-1))
			if charIndex >= len(gradientChars) {
				charIndex = len(gradientChars) - 1
			}
			result.WriteString(gradientChars[charIndex])
		} else {
			result.WriteString("░")
		}
	}

	styled := lipgloss.NewStyle().
		Foreground(p.config.Color).
		Render(result.String())

	return p.addMetadata(styled)
}

// renderSteps renders a step-based progress indicator
func (p *Progress) renderSteps() string {
	steps := 5 // Default number of steps
	currentStep := int(p.value * float64(steps))

	var result strings.Builder
	for i := 0; i < steps; i++ {
		if i < currentStep {
			result.WriteString("✓")
		} else if i == currentStep {
			result.WriteString("◐")
		} else {
			result.WriteString("○")
		}

		if i < steps-1 {
			result.WriteString("━━")
		}
	}

	styled := lipgloss.NewStyle().
		Foreground(p.config.Color).
		Render(result.String())

	return p.addMetadata(styled)
}

// renderAntoine renders an Antoine-themed progress bar
func (p *Progress) renderAntoine() string {
	robotEmojis := []string{"🤖", "⚙️", "🔧", "💻", "🚀"}
	totalRobots := 5
	activeRobots := int(p.value * float64(totalRobots))

	var result strings.Builder
	for i := 0; i < totalRobots; i++ {
		if i < activeRobots {
			emojiIndex := i % len(robotEmojis)
			result.WriteString(robotEmojis[emojiIndex])
		} else {
			result.WriteString("⚪")
		}

		if i < totalRobots-1 {
			result.WriteString(" ")
		}
	}

	return p.addMetadata(result.String())
}

// addMetadata adds percentage, ETA, and other metadata to the progress display
func (p *Progress) addMetadata(progressBar string) string {
	var parts []string

	// Add label if specified
	if p.config.ShowLabel && p.config.Label != "" {
		label := styles.BodyStyle.Render(p.config.Label)
		parts = append(parts, label)
	}

	// Add the progress bar itself
	parts = append(parts, progressBar)

	// Add percentage if enabled
	if p.config.ShowPercent {
		percentage := fmt.Sprintf("%.1f%%", p.value*100)
		percentStyled := lipgloss.NewStyle().
			Foreground(p.config.TextColor).
			Render(percentage)
		parts = append(parts, percentStyled)
	}

	// Add ETA if enabled and we have enough data
	if p.config.ShowETA && p.value > 0 && p.value < 1 {
		eta := p.calculateETA()
		if eta > 0 {
			etaStr := fmt.Sprintf("ETA: %s", utils.FormatDuration(eta))
			etaStyled := styles.BodySecondaryStyle.Render(etaStr)
			parts = append(parts, etaStyled)
		}
	}

	// Add rate if enabled
	if p.config.ShowRate && p.rate > 0 {
		rateStr := fmt.Sprintf("%.1f/s", p.rate)
		rateStyled := styles.BodySecondaryStyle.Render(rateStr)
		parts = append(parts, rateStyled)
	}

	result := strings.Join(parts, " ")

	if p.config.Centered {
		return p.style.Render(result)
	}

	return result
}

// calculateETA calculates estimated time of arrival
func (p *Progress) calculateETA() time.Duration {
	if p.value <= 0 {
		return 0
	}

	elapsed := time.Since(p.startTime)
	remaining := (1.0 - p.value) / p.value
	return time.Duration(float64(elapsed) * remaining)
}

// Progress control methods

// SetProgress sets the progress value (0.0 - 1.0)
func (p *Progress) SetProgress(value float64) {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}

	p.value = value
	p.lastUpdate = time.Now()
}

// IncrementProgress increments the progress by a delta
func (p *Progress) IncrementProgress(delta float64) {
	p.SetProgress(p.value + delta)
}

// SetWorkProgress sets progress based on work done vs total work
func (p *Progress) SetWorkProgress(done, total int64) {
	p.doneWork = done
	p.totalWork = total

	if total > 0 {
		p.SetProgress(float64(done) / float64(total))
	}

	// Calculate rate
	elapsed := time.Since(p.startTime).Seconds()
	if elapsed > 0 {
		p.rate = float64(done) / elapsed
	}
}

// IncrementWork increments the work done
func (p *Progress) IncrementWork(delta int64) {
	p.SetWorkProgress(p.doneWork+delta, p.totalWork)
}

// GetProgress returns the current progress value
func (p *Progress) GetProgress() float64 {
	return p.value
}

// IsComplete returns true if progress is 100%
func (p *Progress) IsComplete() bool {
	return p.value >= 1.0
}

// Reset resets the progress to 0
func (p *Progress) Reset() {
	p.value = 0
	p.startTime = time.Now()
	p.lastUpdate = time.Now()
	p.doneWork = 0
	p.rate = 0
}

// Configuration methods

// SetLabel sets the progress label
func (p *Progress) SetLabel(label string) {
	p.config.Label = label
}

// SetColor sets the progress color
func (p *Progress) SetColor(color lipgloss.Color) {
	p.config.Color = color
}

// SetWidth sets the progress bar width
func (p *Progress) SetWidth(width int) {
	p.config.Width = width
}

// TogglePercent toggles percentage display
func (p *Progress) TogglePercent() {
	p.config.ShowPercent = !p.config.ShowPercent
}

// ToggleETA toggles ETA display
func (p *Progress) ToggleETA() {
	p.config.ShowETA = !p.config.ShowETA
}

// Multi-progress bar for handling multiple concurrent operations
type MultiProgress struct {
	bars     map[string]*Progress
	order    []string
	title    string
	width    int
	showSummary bool
}

// NewMultiProgress creates a new multi-progress manager
func NewMultiProgress(title string, width int) *MultiProgress {
	return &MultiProgress{
		bars:     make(map[string]*Progress),
		title:    title,
		width:    width,
		showSummary: true,
	}
}

// AddProgress adds a progress bar
func (mp *MultiProgress) AddProgress(id string, config ProgressConfig) {
	if config.Width == 0 {
		config.Width = mp.width - 20 // Reserve space for labels
	}

	progress := NewProgress(config)
	mp.bars[id] = progress
	mp.order = append(mp.order, id)
}

// SetProgress sets progress for a specific bar
func (mp *MultiProgress) SetProgress(id string, value float64) {
	if bar, exists := mp.bars[id]; exists {
		bar.SetProgress(value)
	}
}

// IncrementProgress increments progress for a specific bar
func (mp *MultiProgress) IncrementProgress(id string, delta float64) {
	if bar, exists := mp.bars[id]; exists {
		bar.IncrementProgress(delta)
	}
}

// RemoveProgress removes a progress bar
func (mp *MultiProgress) RemoveProgress(id string) {
	delete(mp.bars, id)

	// Remove from order
	for i, orderId := range mp.order {
		if orderId == id {
			mp.order = append(mp.order[:i], mp.order[i+1:]...)
			break
		}
	}
}

// Render renders all progress bars
func (mp *MultiProgress) Render() string {
	var sections []string

	// Title
	if mp.title != "" {
		title := styles.H3Style.Render(mp.title)
		sections = append(sections, title)
	}

	// Individual progress bars
	for _, id := range mp.order {
		if bar, exists := mp.bars[id]; exists {
			// Create label with ID
			label := fmt.Sprintf("%-15s", id)
			labelStyled := styles.LabelStyle.Render(label)

			// Render progress bar
			progressBar := bar.Render()

			// Combine label and progress
			line := labelStyled + " " + progressBar
			sections = append(sections, line)
		}
	}

	// Summary if enabled
	if mp.showSummary && len(mp.bars) > 1 {
		summary := mp.renderSummary()
		sections = append(sections, "")
		sections = append(sections, summary)
	}

	return strings.Join(sections, "\n")
}

// renderSummary renders a summary of all progress bars
func (mp *MultiProgress) renderSummary() string {
	totalBars := len(mp.bars)
	completeBars := 0
	totalProgress := 0.0

	for _, bar := range mp.bars {
		totalProgress += bar.GetProgress()
		if bar.IsComplete() {
			completeBars++
		}
	}

	avgProgress := totalProgress / float64(totalBars)

	summaryText := fmt.Sprintf("Overall: %d/%d complete (%.1f%%)",
		completeBars, totalBars, avgProgress*100)

	// Create summary progress bar
	summaryConfig := ProgressConfig{
		Type:        ProgressTypeBar,
		Width:       mp.width - 30,
		ShowPercent: false,
		Color:       styles.Gold,
	}

	summaryBar := NewProgress(summaryConfig)
	summaryBar.SetProgress(avgProgress)

	summaryLabel := styles.LabelStyle.Render("Summary:")
	summaryText = styles.BodySecondaryStyle.Render(summaryText)

	return summaryLabel + " " + summaryBar.Render() + " " + summaryText
}

// GetOverallProgress returns the overall progress (0.0 - 1.0)
func (mp *MultiProgress) GetOverallProgress() float64 {
	if len(mp.bars) == 0 {
		return 0
	}

	totalProgress := 0.0
	for _, bar := range mp.bars {
		totalProgress += bar.GetProgress()
	}

	return totalProgress / float64(len(mp.bars))
}

// IsAllComplete returns true if all progress bars are complete
func (mp *MultiProgress) IsAllComplete() bool {
	for _, bar := range mp.bars {
		if !bar.IsComplete() {
			return false
		}
	}
	return len(mp.bars) > 0
}

// Clear removes all progress bars
func (mp *MultiProgress) Clear() {
	mp.bars = make(map[string]*Progress)
	mp.order = []string{}
}

// ProgressBuilder provides a fluent interface for building progress bars
type ProgressBuilder struct {
	config ProgressConfig
}

// NewProgressBuilder creates a new progress builder
func NewProgressBuilder() *ProgressBuilder {
	return &ProgressBuilder{
		config: ProgressConfig{
			Type:        ProgressTypeBar,
			Width:       50,
			Height:      1,
			ShowPercent: true,
			ShowETA:     false,
			ShowRate:    false,
			ShowLabel:   false,
			Color:       styles.Gold,
			Animated:    false,
			Centered:    false,
		},
	}
}

// Type sets the progress type
func (pb *ProgressBuilder) Type(progressType ProgressType) *ProgressBuilder {
	pb.config.Type = progressType
	return pb
}

// Width sets the progress width
func (pb *ProgressBuilder) Width(width int) *ProgressBuilder {
	pb.config.Width = width
	return pb
}

// Label sets the progress label
func (pb *ProgressBuilder) Label(label string) *ProgressBuilder {
	pb.config.Label = label
	pb.config.ShowLabel = true
	return pb
}

// Color sets the progress color
func (pb *ProgressBuilder) Color(color lipgloss.Color) *ProgressBuilder {
	pb.config.Color = color
	return pb
}

// ShowPercent enables/disables percentage display
func (pb *ProgressBuilder) ShowPercent(show bool) *ProgressBuilder {
	pb.config.ShowPercent = show
	return pb
}

// ShowETA enables/disables ETA display
func (pb *ProgressBuilder) ShowETA(show bool) *ProgressBuilder {
	pb.config.ShowETA = show
	return pb
}

// ShowRate enables/disables rate display
func (pb *ProgressBuilder) ShowRate(show bool) *ProgressBuilder {
	pb.config.ShowRate = show
	return pb
}

// Animated enables/disables animations
func (pb *ProgressBuilder) Animated(animated bool) *ProgressBuilder {
	pb.config.Animated = animated
	return pb
}

// Centered centers the progress bar
func (pb *ProgressBuilder) Centered(centered bool) *ProgressBuilder {
	pb.config.Centered = centered
	return pb
}

// Build creates the progress component
func (pb *ProgressBuilder) Build() *Progress {
	return NewProgress(pb.config)
}

// Predefined progress bars for common operations

// DownloadProgress creates a progress bar for downloads
func DownloadProgress(filename string) *Progress {
	return NewProgressBuilder().
		Type(ProgressTypeBar).
		Label(fmt.Sprintf("Downloading %s", filename)).
		ShowPercent(true).
		ShowETA(true).
		ShowRate(true).
		Color(styles.Green).
		Build()
}

// AnalysisProgress creates a progress bar for analysis operations
func AnalysisProgress(target string) *Progress {
	return NewProgressBuilder().
		Type(ProgressTypeWave).
		Label(fmt.Sprintf("Analyzing %s", target)).
		ShowPercent(true).
		ShowETA(true).
		Color(styles.Gold).
		Animated(true).
		Build()
}

// SearchProgress creates a progress bar for search operations
func SearchProgress() *Progress {
	return NewProgressBuilder().
		Type(ProgressTypeDots).
		Label("Searching").
		ShowPercent(false).
		Color(styles.Cyan).
		Animated(true).
		Build()
}

// UploadProgress creates a progress bar for upload operations
func UploadProgress(filename string) *Progress {
	return NewProgressBuilder().
		Type(ProgressTypeBlocks).
		Label(fmt.Sprintf("Uploading %s", filename)).
		ShowPercent(true).
		ShowETA(true).
		ShowRate(true).
		Color(styles.Blue).
		Build()
}

// InstallProgress creates a progress bar for installation operations
func InstallProgress(package string) *Progress {
return NewProgressBuilder().
Type(ProgressTypeSteps).
Label(fmt.Sprintf("Installing %s", package)).
ShowPercent(true).
Color(styles.Green).
Build()
}

// MentorProgress creates a progress bar for AI mentor operations
func MentorProgress() *Progress {
	return NewProgressBuilder().
		Type(ProgressTypeAntoine).
		Label("AI Mentor Processing").
		ShowPercent(false).
		Color(styles.Purple).
		Animated(true).
		Build()
}

// GitProgress creates a progress bar for Git operations
func GitProgress(operation string) *Progress {
	return NewProgressBuilder().
		Type(ProgressTypeBar).
		Label(fmt.Sprintf("Git: %s", operation)).
		ShowPercent(true).
		ShowETA(true).
		Color(styles.Orange).
		Build()
}
