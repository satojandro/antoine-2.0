// Package components provides reusable UI components for Antoine CLI
// This file implements animated spinners with Antoine-themed designs
package components

import (
	"fmt"
	"strings"
	"time"

	"antoine-cli/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// SpinnerType defines different spinner animations
type SpinnerType string

const (
	SpinnerTypeDots      SpinnerType = "dots"      // Classic dots spinner
	SpinnerTypeBar       SpinnerType = "bar"       // Progress bar style
	SpinnerTypePulse     SpinnerType = "pulse"     // Pulsing animation
	SpinnerTypeArrow     SpinnerType = "arrow"     // Rotating arrow
	SpinnerTypeBounce    SpinnerType = "bounce"    // Bouncing ball
	SpinnerTypeWave      SpinnerType = "wave"      // Wave animation
	SpinnerTypeAntoine   SpinnerType = "antoine"   // Antoine-branded spinner
	SpinnerTypeMatrix    SpinnerType = "matrix"    // Matrix-style rain
	SpinnerTypeGears     SpinnerType = "gears"     // Rotating gears
	SpinnerTypeHeartbeat SpinnerType = "heartbeat" // Heartbeat pulse
)

// SpinnerConfig configures the spinner component
type SpinnerConfig struct {
	Type     SpinnerType
	Message  string
	Speed    time.Duration
	Color    lipgloss.Color
	Width    int
	Centered bool
	ShowTime bool
	Prefix   string
	Suffix   string
}

// Spinner represents an animated spinner component
type Spinner struct {
	config    SpinnerConfig
	startTime time.Time
	frame     int
	style     lipgloss.Style
}

// NewSpinner creates a new spinner component
func NewSpinner(config SpinnerConfig) *Spinner {
	// Set defaults
	if config.Speed == 0 {
		config.Speed = 100 * time.Millisecond
	}
	if config.Color == "" {
		config.Color = styles.Gold
	}
	if config.Width == 0 {
		config.Width = 40
	}

	// Create base style
	baseStyle := lipgloss.NewStyle().
		Foreground(config.Color)

	if config.Centered {
		baseStyle = baseStyle.
			Width(config.Width).
			Align(lipgloss.Center)
	}

	return &Spinner{
		config:    config,
		startTime: time.Now(),
		frame:     0,
		style:     baseStyle,
	}
}

// Render renders the current frame of the spinner
func (s *Spinner) Render() string {
	animation := s.getAnimation()
	frameCount := len(animation)

	if frameCount == 0 {
		return s.config.Message
	}

	// Calculate current frame based on time
	elapsed := time.Since(s.startTime)
	frameIndex := int(elapsed/s.config.Speed) % frameCount
	currentFrame := animation[frameIndex]

	// Build the complete spinner display
	var parts []string

	if s.config.Prefix != "" {
		parts = append(parts, s.config.Prefix)
	}

	parts = append(parts, currentFrame)

	if s.config.Message != "" {
		parts = append(parts, s.config.Message)
	}

	if s.config.ShowTime {
		duration := time.Since(s.startTime)
		timeStr := fmt.Sprintf("(%s)", formatDuration(duration))
		timeStyled := styles.BodySecondaryStyle.Render(timeStr)
		parts = append(parts, timeStyled)
	}

	if s.config.Suffix != "" {
		parts = append(parts, s.config.Suffix)
	}

	content := strings.Join(parts, " ")
	return s.style.Render(content)
}

// getAnimation returns the animation frames for the current spinner type
func (s *Spinner) getAnimation() []string {
	switch s.config.Type {
	case SpinnerTypeDots:
		return []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	case SpinnerTypeBar:
		return []string{
			"▱▱▱▱▱▱▱",
			"▰▱▱▱▱▱▱",
			"▰▰▱▱▱▱▱",
			"▰▰▰▱▱▱▱",
			"▰▰▰▰▱▱▱",
			"▰▰▰▰▰▱▱",
			"▰▰▰▰▰▰▱",
			"▰▰▰▰▰▰▰",
			"▱▰▰▰▰▰▰",
			"▱▱▰▰▰▰▰",
			"▱▱▱▰▰▰▰",
			"▱▱▱▱▰▰▰",
			"▱▱▱▱▱▰▰",
			"▱▱▱▱▱▱▰",
		}

	case SpinnerTypePulse:
		return []string{"●", "○", "◐", "◑", "◒", "◓"}

	case SpinnerTypeArrow:
		return []string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"}

	case SpinnerTypeBounce:
		return []string{
			"⠁      ",
			" ⠂     ",
			"  ⠄    ",
			"   ⠂   ",
			"    ⠁  ",
			"     ⠂ ",
			"      ⠄",
			"     ⠂ ",
			"    ⠁  ",
			"   ⠂   ",
			"  ⠄    ",
			" ⠂     ",
		}

	case SpinnerTypeWave:
		return []string{
			"▁▁▁▁▁",
			"▂▁▁▁▁",
			"▃▂▁▁▁",
			"▄▃▂▁▁",
			"▅▄▃▂▁",
			"▆▅▄▃▂",
			"▇▆▅▄▃",
			"█▇▆▅▄",
			"▇█▇▆▅",
			"▆▇█▇▆",
			"▅▆▇█▇",
			"▄▅▆▇█",
			"▃▄▅▆▇",
			"▂▃▄▅▆",
			"▁▂▃▄▅",
			"▁▁▂▃▄",
			"▁▁▁▂▃",
			"▁▁▁▁▂",
		}

	case SpinnerTypeAntoine:
		return []string{
			"🤖   ",
			" 🤖  ",
			"  🤖 ",
			"   🤖",
			"  🤖 ",
			" 🤖  ",
		}

	case SpinnerTypeMatrix:
		return []string{
			"┌─┐┌─┐┌─┐",
			"│ ││ ││ │",
			"└─┘└─┘└─┘",
			"┏━┓┏━┓┏━┓",
			"┃ ┃┃ ┃┃ ┃",
			"┗━┛┗━┛┗━┛",
		}

	case SpinnerTypeGears:
		return []string{
			"⚙ ⚙",
			"⚙⚙ ",
			" ⚙⚙",
			"⚙ ⚙",
		}

	case SpinnerTypeHeartbeat:
		return []string{
			"♡     ",
			"♥     ",
			"♥♡    ",
			"♥♥    ",
			"♥♥♡   ",
			"♥♥♥   ",
			"♥♥♥♡  ",
			"♥♥♥♥  ",
			"♥♥♥♥♡ ",
			"♥♥♥♥♥ ",
			"♥♥♥♥♥♡",
			"♥♥♥♥♥♥",
			" ♥♥♥♥♥",
			"  ♥♥♥♥",
			"   ♥♥♥",
			"    ♥♥",
			"     ♥",
			"      ",
		}

	default:
		return []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	}
}

// Start starts the spinner (resets the timer)
func (s *Spinner) Start() {
	s.startTime = time.Now()
	s.frame = 0
}

// Stop stops the spinner and returns the final display
func (s *Spinner) Stop() string {
	if s.config.Message != "" {
		return styles.SuccessStyle.Render("✓ " + s.config.Message)
	}
	return styles.SuccessStyle.Render("✓ Complete")
}

// Error stops the spinner and shows an error state
func (s *Spinner) Error(message string) string {
	if message == "" {
		message = s.config.Message
	}
	if message == "" {
		message = "Failed"
	}
	return styles.ErrorStyle.Render("✗ " + message)
}

// Update updates the spinner configuration
func (s *Spinner) Update(config SpinnerConfig) {
	// Preserve start time
	startTime := s.startTime

	// Update config
	s.config = config

	// Set defaults for new config
	if config.Speed == 0 {
		s.config.Speed = 100 * time.Millisecond
	}
	if config.Color == "" {
		s.config.Color = styles.Gold
	}
	if config.Width == 0 {
		s.config.Width = 40
	}

	// Update style
	s.style = lipgloss.NewStyle().Foreground(s.config.Color)
	if s.config.Centered {
		s.style = s.style.Width(s.config.Width).Align(lipgloss.Center)
	}

	// Restore start time
	s.startTime = startTime
}

// SetMessage updates the spinner message
func (s *Spinner) SetMessage(message string) {
	s.config.Message = message
}

// SetColor updates the spinner color
func (s *Spinner) SetColor(color lipgloss.Color) {
	s.config.Color = color
	s.style = s.style.Foreground(color)
}

// GetElapsedTime returns the elapsed time since the spinner started
func (s *Spinner) GetElapsedTime() time.Duration {
	return time.Since(s.startTime)
}

// SpinnerBuilder provides a fluent interface for building spinners
type SpinnerBuilder struct {
	config SpinnerConfig
}

// NewSpinnerBuilder creates a new spinner builder
func NewSpinnerBuilder() *SpinnerBuilder {
	return &SpinnerBuilder{
		config: SpinnerConfig{
			Type:     SpinnerTypeDots,
			Speed:    100 * time.Millisecond,
			Color:    styles.Gold,
			Width:    40,
			Centered: false,
		},
	}
}

// Type sets the spinner type
func (sb *SpinnerBuilder) Type(spinnerType SpinnerType) *SpinnerBuilder {
	sb.config.Type = spinnerType
	return sb
}

// Message sets the spinner message
func (sb *SpinnerBuilder) Message(message string) *SpinnerBuilder {
	sb.config.Message = message
	return sb
}

// Speed sets the animation speed
func (sb *SpinnerBuilder) Speed(speed time.Duration) *SpinnerBuilder {
	sb.config.Speed = speed
	return sb
}

// Color sets the spinner color
func (sb *SpinnerBuilder) Color(color lipgloss.Color) *SpinnerBuilder {
	sb.config.Color = color
	return sb
}

// Width sets the spinner width
func (sb *SpinnerBuilder) Width(width int) *SpinnerBuilder {
	sb.config.Width = width
	return sb
}

// Centered centers the spinner
func (sb *SpinnerBuilder) Centered(centered bool) *SpinnerBuilder {
	sb.config.Centered = centered
	return sb
}

// ShowTime shows elapsed time
func (sb *SpinnerBuilder) ShowTime(show bool) *SpinnerBuilder {
	sb.config.ShowTime = show
	return sb
}

// Prefix sets a prefix for the spinner
func (sb *SpinnerBuilder) Prefix(prefix string) *SpinnerBuilder {
	sb.config.Prefix = prefix
	return sb
}

// Suffix sets a suffix for the spinner
func (sb *SpinnerBuilder) Suffix(suffix string) *SpinnerBuilder {
	sb.config.Suffix = suffix
	return sb
}

// Build creates the spinner component
func (sb *SpinnerBuilder) Build() *Spinner {
	return NewSpinner(sb.config)
}

// Predefined spinners for common operations

// SearchSpinner creates a spinner for search operations
func SearchSpinner(query string) *Spinner {
	return NewSpinnerBuilder().
		Type(SpinnerTypeDots).
		Message(fmt.Sprintf("Searching for '%s'", query)).
		Color(styles.Cyan).
		ShowTime(true).
		Build()
}

// AnalysisSpinner creates a spinner for analysis operations
func AnalysisSpinner(target string) *Spinner {
	return NewSpinnerBuilder().
		Type(SpinnerTypeWave).
		Message(fmt.Sprintf("Analyzing %s", target)).
		Color(styles.Gold).
		ShowTime(true).
		Speed(150 * time.Millisecond).
		Build()
}

// MentorSpinner creates a spinner for AI mentor operations
func MentorSpinner() *Spinner {
	return NewSpinnerBuilder().
		Type(SpinnerTypeAntoine).
		Message("AI Mentor is thinking").
		Color(styles.Purple).
		ShowTime(true).
		Speed(500 * time.Millisecond).
		Build()
}

// DownloadSpinner creates a spinner for download operations
func DownloadSpinner(item string) *Spinner {
	return NewSpinnerBuilder().
		Type(SpinnerTypeBar).
		Message(fmt.Sprintf("Downloading %s", item)).
		Color(styles.Green).
		ShowTime(true).
		Speed(200 * time.Millisecond).
		Build()
}

// LoadingSpinner creates a generic loading spinner
func LoadingSpinner(operation string) *Spinner {
	return NewSpinnerBuilder().
		Type(SpinnerTypePulse).
		Message(operation).
		Color(styles.Blue).
		Speed(300 * time.Millisecond).
		Build()
}

// ConnectingSpinner creates a spinner for connection operations
func ConnectingSpinner(service string) *Spinner {
	return NewSpinnerBuilder().
		Type(SpinnerTypeArrow).
		Message(fmt.Sprintf("Connecting to %s", service)).
		Color(styles.Orange).
		ShowTime(true).
		Speed(250 * time.Millisecond).
		Build()
}

// ProcessingSpinner creates a spinner for processing operations
func ProcessingSpinner(task string) *Spinner {
	return NewSpinnerBuilder().
		Type(SpinnerTypeGears).
		Message(fmt.Sprintf("Processing %s", task)).
		Color(styles.Gold).
		ShowTime(true).
		Speed(400 * time.Millisecond).
		Build()
}

// MCPSpinner creates a spinner for MCP operations
func MCPSpinner(server, operation string) *Spinner {
	return NewSpinnerBuilder().
		Type(SpinnerTypeMatrix).
		Message(fmt.Sprintf("MCP %s: %s", server, operation)).
		Color(styles.Cyan).
		ShowTime(true).
		Prefix("🔗").
		Speed(300 * time.Millisecond).
		Build()
}

// GitSpinner creates a spinner for Git operations
func GitSpinner(operation string) *Spinner {
	return NewSpinnerBuilder().
		Type(SpinnerTypeBounce).
		Message(fmt.Sprintf("Git: %s", operation)).
		Color(styles.Green).
		ShowTime(true).
		Prefix("🌿").
		Speed(200 * time.Millisecond).
		Build()
}

// HeartbeatSpinner creates a heartbeat spinner for system status
func HeartbeatSpinner() *Spinner {
	return NewSpinnerBuilder().
		Type(SpinnerTypeHeartbeat).
		Message("System status").
		Color(styles.Red).
		Speed(100 * time.Millisecond).
		Build()
}

// SpinnerManager manages multiple spinners
type SpinnerManager struct {
	spinners map[string]*Spinner
	active   string
}

// NewSpinnerManager creates a new spinner manager
func NewSpinnerManager() *SpinnerManager {
	return &SpinnerManager{
		spinners: make(map[string]*Spinner),
	}
}

// Add adds a spinner to the manager
func (sm *SpinnerManager) Add(id string, spinner *Spinner) {
	sm.spinners[id] = spinner
}

// Start starts a specific spinner
func (sm *SpinnerManager) Start(id string) {
	if spinner, exists := sm.spinners[id]; exists {
		spinner.Start()
		sm.active = id
	}
}

// Stop stops a specific spinner
func (sm *SpinnerManager) Stop(id string) string {
	if spinner, exists := sm.spinners[id]; exists {
		if sm.active == id {
			sm.active = ""
		}
		return spinner.Stop()
	}
	return ""
}

// Error stops a spinner with error state
func (sm *SpinnerManager) Error(id, message string) string {
	if spinner, exists := sm.spinners[id]; exists {
		if sm.active == id {
			sm.active = ""
		}
		return spinner.Error(message)
	}
	return ""
}

// RenderActive renders the currently active spinner
func (sm *SpinnerManager) RenderActive() string {
	if sm.active == "" {
		return ""
	}

	if spinner, exists := sm.spinners[sm.active]; exists {
		return spinner.Render()
	}

	return ""
}

// GetActive returns the ID of the active spinner
func (sm *SpinnerManager) GetActive() string {
	return sm.active
}

// Remove removes a spinner from the manager
func (sm *SpinnerManager) Remove(id string) {
	if sm.active == id {
		sm.active = ""
	}
	delete(sm.spinners, id)
}

// Clear removes all spinners
func (sm *SpinnerManager) Clear() {
	sm.spinners = make(map[string]*Spinner)
	sm.active = ""
}

// Helper function to format duration
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	return fmt.Sprintf("%.1fm", d.Minutes())
}
