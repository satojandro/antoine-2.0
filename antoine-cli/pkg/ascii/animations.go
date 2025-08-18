// Package ascii provides ASCII art and animations for Antoine CLI
// This file implements animated ASCII art and visual effects
package ascii

import (
	"math"
	"strings"
	"time"
)

// AnimationType defines different animation types
type AnimationType string

const (
	AnimationTypeWave       AnimationType = "wave"       // Wave motion effect
	AnimationTypePulse      AnimationType = "pulse"      // Pulsing effect
	AnimationTypeRotate     AnimationType = "rotate"     // Rotation effect
	AnimationTypeSparkle    AnimationType = "sparkle"    // Sparkling effect
	AnimationTypeTypewriter AnimationType = "typewriter" // Typewriter effect
	AnimationTypeMatrix     AnimationType = "matrix"     // Matrix rain effect
	AnimationTypeGlitch     AnimationType = "glitch"     // Glitch effect
	AnimationTypeFade       AnimationType = "fade"       // Fade in/out effect
)

// AnimationConfig configures animation parameters
type AnimationConfig struct {
	Type    AnimationType
	Speed   time.Duration
	Frames  int
	Width   int
	Height  int
	Reverse bool
	Loop    bool
	Colors  []string
}

// Animation represents an animated ASCII sequence
type Animation struct {
	config       AnimationConfig
	startTime    time.Time
	currentFrame int
	frames       []string
	isPlaying    bool
}

// NewAnimation creates a new animation
func NewAnimation(config AnimationConfig) *Animation {
	animation := &Animation{
		config:    config,
		startTime: time.Now(),
		isPlaying: true,
	}

	animation.generateFrames()
	return animation
}

// generateFrames generates the animation frames based on type
func (a *Animation) generateFrames() {
	switch a.config.Type {
	case AnimationTypeWave:
		a.frames = a.generateWaveFrames()
	case AnimationTypePulse:
		a.frames = a.generatePulseFrames()
	case AnimationTypeRotate:
		a.frames = a.generateRotateFrames()
	case AnimationTypeSparkle:
		a.frames = a.generateSparkleFrames()
	case AnimationTypeTypewriter:
		a.frames = a.generateTypewriterFrames()
	case AnimationTypeMatrix:
		a.frames = a.generateMatrixFrames()
	case AnimationTypeGlitch:
		a.frames = a.generateGlitchFrames()
	case AnimationTypeFade:
		a.frames = a.generateFadeFrames()
	default:
		a.frames = []string{""}
	}
}

// generateWaveFrames creates wave animation frames
func (a *Animation) generateWaveFrames() []string {
	frames := make([]string, a.config.Frames)
	waveChars := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

	for frame := 0; frame < a.config.Frames; frame++ {
		var lines []string

		for y := 0; y < a.config.Height; y++ {
			var line strings.Builder

			for x := 0; x < a.config.Width; x++ {
				// Calculate wave position
				wavePos := math.Sin(float64(x)*0.3 + float64(frame)*0.2 + float64(y)*0.1)
				intensity := (wavePos + 1) / 2 // Normalize to 0-1

				charIndex := int(intensity * float64(len(waveChars)-1))
				if charIndex >= len(waveChars) {
					charIndex = len(waveChars) - 1
				}

				line.WriteString(waveChars[charIndex])
			}

			lines = append(lines, line.String())
		}

		frames[frame] = strings.Join(lines, "\n")
	}

	return frames
}

// generatePulseFrames creates pulsing animation frames
func (a *Animation) generatePulseFrames() []string {
	frames := make([]string, a.config.Frames)

	for frame := 0; frame < a.config.Frames; frame++ {
		intensity := (math.Sin(float64(frame)*0.3) + 1) / 2

		var char string
		switch {
		case intensity < 0.2:
			char = "░"
		case intensity < 0.4:
			char = "▒"
		case intensity < 0.6:
			char = "▓"
		case intensity < 0.8:
			char = "█"
		default:
			char = "██"
		}

		var lines []string
		centerY := a.config.Height / 2
		centerX := a.config.Width / 2
		radius := int(intensity * float64(centerX))

		for y := 0; y < a.config.Height; y++ {
			var line strings.Builder

			for x := 0; x < a.config.Width; x++ {
				// Calculate distance from center
				dx := x - centerX
				dy := y - centerY
				distance := math.Sqrt(float64(dx*dx + dy*dy))

				if distance <= float64(radius) {
					line.WriteString(char)
				} else {
					line.WriteString(" ")
				}
			}

			lines = append(lines, line.String())
		}

		frames[frame] = strings.Join(lines, "\n")
	}

	return frames
}

// generateRotateFrames creates rotation animation frames
func (a *Animation) generateRotateFrames() []string {
	frames := make([]string, a.config.Frames)
	rotateChars := []string{"│", "╱", "─", "╲"}

	for frame := 0; frame < a.config.Frames; frame++ {
		charIndex := frame % len(rotateChars)
		char := rotateChars[charIndex]

		var lines []string
		for y := 0; y < a.config.Height; y++ {
			var line strings.Builder

			for x := 0; x < a.config.Width; x++ {
				if x == a.config.Width/2 && y == a.config.Height/2 {
					line.WriteString(char)
				} else {
					line.WriteString(" ")
				}
			}

			lines = append(lines, line.String())
		}

		frames[frame] = strings.Join(lines, "\n")
	}

	return frames
}

// generateSparkleFrames creates sparkling animation frames
func (a *Animation) generateSparkleFrames() []string {
	frames := make([]string, a.config.Frames)
	sparkleChars := []string{"·", "✦", "✧", "✩", "✪", "✫", "✬", "✭", "✮", "✯", "✰"}

	for frame := 0; frame < a.config.Frames; frame++ {
		var lines []string

		for y := 0; y < a.config.Height; y++ {
			var line strings.Builder

			for x := 0; x < a.config.Width; x++ {
				// Pseudo-random sparkle based on position and frame
				seed := (x*7 + y*11 + frame*3) % 100

				if seed < 5 { // 5% chance of sparkle
					charIndex := seed % len(sparkleChars)
					line.WriteString(sparkleChars[charIndex])
				} else {
					line.WriteString(" ")
				}
			}

			lines = append(lines, line.String())
		}

		frames[frame] = strings.Join(lines, "\n")
	}

	return frames
}

// generateTypewriterFrames creates typewriter animation frames
func (a *Animation) generateTypewriterFrames() []string {
	text := "ANTOINE CLI - Your Ultimate Hackathon Mentor"
	frames := make([]string, len(text)+1)

	for frame := 0; frame <= len(text); frame++ {
		displayText := text[:frame]
		if frame < len(text) {
			displayText += "│" // Cursor
		}

		// Center the text
		padding := (a.config.Width - len(displayText)) / 2
		if padding > 0 {
			displayText = strings.Repeat(" ", padding) + displayText
		}

		frames[frame] = displayText
	}

	return frames
}

// generateMatrixFrames creates Matrix-style rain animation
func (a *Animation) generateMatrixFrames() []string {
	frames := make([]string, a.config.Frames)
	matrixChars := []string{"0", "1", "ア", "イ", "ウ", "エ", "オ", "カ", "キ", "ク"}

	// Create columns with different speeds
	columns := make([][]string, a.config.Width)
	for x := 0; x < a.config.Width; x++ {
		columns[x] = make([]string, a.config.Height*2) // Longer for smooth scrolling

		// Fill column with random chars
		for y := 0; y < len(columns[x]); y++ {
			if (x*7+y*11)%3 == 0 { // Pseudo-random distribution
				charIndex := (x + y) % len(matrixChars)
				columns[x][y] = matrixChars[charIndex]
			} else {
				columns[x][y] = " "
			}
		}
	}

	for frame := 0; frame < a.config.Frames; frame++ {
		var lines []string

		for y := 0; y < a.config.Height; y++ {
			var line strings.Builder

			for x := 0; x < a.config.Width; x++ {
				// Calculate scrolling position
				speed := (x % 3) + 1 // Different speeds per column
				offset := (frame * speed) % len(columns[x])
				charY := (y + offset) % len(columns[x])

				line.WriteString(columns[x][charY])
			}

			lines = append(lines, line.String())
		}

		frames[frame] = strings.Join(lines, "\n")
	}

	return frames
}

// generateGlitchFrames creates glitch effect animation
func (a *Animation) generateGlitchFrames() []string {
	frames := make([]string, a.config.Frames)
	baseText := "ANTOINE"
	glitchChars := []string{"@", "#", "$", "%", "&", "*", "+", "=", "~"}

	for frame := 0; frame < a.config.Frames; frame++ {
		var result strings.Builder

		for i, char := range baseText {
			// Pseudo-random glitch based on frame and position
			seed := (frame*7 + i*11) % 100

			if seed < 10 { // 10% chance of glitch
				glitchIndex := seed % len(glitchChars)
				result.WriteString(glitchChars[glitchIndex])
			} else {
				result.WriteRune(char)
			}
		}

		frames[frame] = result.String()
	}

	return frames
}

// generateFadeFrames creates fade in/out animation
func (a *Animation) generateFadeFrames() []string {
	frames := make([]string, a.config.Frames)
	text := "ANTOINE CLI"
	fadeChars := []string{" ", "░", "▒", "▓", "█"}

	for frame := 0; frame < a.config.Frames; frame++ {
		// Calculate fade intensity (0 to 1 and back)
		progress := float64(frame) / float64(a.config.Frames-1)
		intensity := math.Sin(progress * math.Pi) // Fade in and out

		charIndex := int(intensity * float64(len(fadeChars)-1))
		if charIndex >= len(fadeChars) {
			charIndex = len(fadeChars) - 1
		}

		fadeChar := fadeChars[charIndex]

		var result strings.Builder
		for _, char := range text {
			if char == ' ' {
				result.WriteString(" ")
			} else {
				result.WriteString(fadeChar)
			}
		}

		frames[frame] = result.String()
	}

	return frames
}

// GetCurrentFrame returns the current animation frame
func (a *Animation) GetCurrentFrame() string {
	if !a.isPlaying || len(a.frames) == 0 {
		return ""
	}

	elapsed := time.Since(a.startTime)
	frameIndex := int(elapsed/a.config.Speed) % len(a.frames)

	if a.config.Reverse {
		frameIndex = len(a.frames) - 1 - frameIndex
	}

	a.currentFrame = frameIndex
	return a.frames[frameIndex]
}

// GetFrame returns a specific frame by index
func (a *Animation) GetFrame(index int) string {
	if index < 0 || index >= len(a.frames) {
		return ""
	}
	return a.frames[index]
}

// GetFrameCount returns the total number of frames
func (a *Animation) GetFrameCount() int {
	return len(a.frames)
}

// Play starts or resumes the animation
func (a *Animation) Play() {
	a.isPlaying = true
	a.startTime = time.Now()
}

// Pause pauses the animation
func (a *Animation) Pause() {
	a.isPlaying = false
}

// Stop stops the animation and resets to first frame
func (a *Animation) Stop() {
	a.isPlaying = false
	a.currentFrame = 0
	a.startTime = time.Now()
}

// IsPlaying returns true if the animation is currently playing
func (a *Animation) IsPlaying() bool {
	return a.isPlaying
}

// SetSpeed changes the animation speed
func (a *Animation) SetSpeed(speed time.Duration) {
	a.config.Speed = speed
}

// Animated logo functions for Antoine CLI

// GetAnimatedLogo returns an animated version of the Antoine logo
// GetAnimatedLogo returns an animated version of the Antoine logo
func GetAnimatedLogo(size int, currentTime time.Time) string {
	baseArt := GetLogo(size) // Cambiado LogoSize por int
	if baseArt == "" {
		return ""
	}

	// Simple animation - cycle through different "glow" effects
	elapsed := currentTime.UnixNano() / int64(time.Millisecond)
	phase := (elapsed / 200) % 4 // Change every 200ms, 4 phases

	switch phase {
	case 0:
		return addGlow(baseArt, "●")
	case 1:
		return addGlow(baseArt, "◉")
	case 2:
		return addGlow(baseArt, "⬢")
	case 3:
		return addGlow(baseArt, "⬡")
	default:
		return baseArt
	}
}

// addGlow adds glowing effects around the logo
func addGlow(logo, glowChar string) string {
	lines := strings.Split(logo, "\n")
	if len(lines) == 0 {
		return logo
	}

	// Add glow characters at the beginning and end of non-empty lines
	var result []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			line = glowChar + " " + line + " " + glowChar
		}
		result = append(result, line)
	}

	// Add glow line at top and bottom
	if len(result) > 0 {
		width := len(result[0])
		glowLine := strings.Repeat(glowChar+" ", width/2)
		result = append([]string{glowLine}, result...)
		result = append(result, glowLine)
	}

	return strings.Join(result, "\n")
}

// AnimatedBanner creates an animated banner with text
func AnimatedBanner(text string, animType AnimationType, width int) *Animation {
	config := AnimationConfig{
		Type:   animType,
		Speed:  200 * time.Millisecond,
		Frames: 20,
		Width:  width,
		Height: 3,
		Loop:   true,
	}

	return NewAnimation(config)
}

// LoadingAnimation creates a loading animation
func LoadingAnimation(message string) *Animation {
	config := AnimationConfig{
		Type:   AnimationTypeWave,
		Speed:  100 * time.Millisecond,
		Frames: 30,
		Width:  len(message) + 10,
		Height: 1,
		Loop:   true,
	}

	return NewAnimation(config)
}

// ProgressAnimation creates a progress bar animation
func ProgressAnimation(width int, progress float64) string {
	filled := int(float64(width) * progress)
	empty := width - filled

	// Animated progress characters
	elapsed := time.Now().UnixNano() / int64(time.Millisecond*100)
	phase := elapsed % 4

	var progressChar string
	switch phase {
	case 0:
		progressChar = "█"
	case 1:
		progressChar = "▉"
	case 2:
		progressChar = "▊"
	case 3:
		progressChar = "▋"
	default:
		progressChar = "█"
	}

	return strings.Repeat(progressChar, filled) + strings.Repeat("░", empty)
}

// StatusAnimation creates animated status indicators
func StatusAnimation(status string, currentTime time.Time) string {
	elapsed := currentTime.UnixNano() / int64(time.Millisecond)
	phase := (elapsed / 300) % 4 // Change every 300ms

	var indicator string
	switch status {
	case "loading":
		indicators := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		indicator = indicators[int(phase)%len(indicators)] // Conversión a int
	case "success":
		indicators := []string{"✓", "✓", "✓", "✓"}
		indicator = indicators[int(phase)]
	case "error":
		indicators := []string{"✗", "✗", "✗", "✗"}
		indicator = indicators[int(phase)]
	case "warning":
		indicators := []string{"⚠", "⚠", " ", " "}
		indicator = indicators[int(phase)]
	default:
		indicator = "●"
	}

	return indicator
}

// RainbowText creates rainbow-colored text effect (using Unicode)
func RainbowText(text string, currentTime time.Time) string {
	// Simulate rainbow effect with different Unicode characters
	rainbowChars := [][]string{
		{"▓", "▓", "▓"}, // Red
		{"▒", "▒", "▒"}, // Orange
		{"░", "░", "░"}, // Yellow
		{"▓", "▒", "░"}, // Green
		{"▒", "░", "▓"}, // Blue
		{"░", "▓", "▒"}, // Indigo
		{"▓", "░", "▒"}, // Violet
	}

	elapsed := currentTime.UnixNano() / int64(time.Millisecond)
	offset := (elapsed / 200) % int64(len(rainbowChars))

	var result strings.Builder
	for i, char := range text {
		if char == ' ' {
			result.WriteString(" ")
		} else {
			colorIndex := (int(offset) + i) % len(rainbowChars)
			charIndex := i % len(rainbowChars[colorIndex])
			result.WriteString(rainbowChars[colorIndex][charIndex])
		}
	}

	return result.String()
}

// FireEffect creates a fire-like animation effect
func FireEffect(width, height int, intensity float64) *Animation {
	config := AnimationConfig{
		Type:   AnimationTypeWave,
		Speed:  50 * time.Millisecond,
		Frames: 40,
		Width:  width,
		Height: height,
		Loop:   true,
	}

	animation := NewAnimation(config)

	// Override with custom fire frames
	frames := make([]string, config.Frames)
	fireChars := []string{" ", ".", ":", "^", "*", "x", "s", "S", "#", "$"}

	for frame := 0; frame < config.Frames; frame++ {
		var lines []string

		for y := 0; y < height; y++ {
			var line strings.Builder

			for x := 0; x < width; x++ {
				// Fire effect: higher intensity at bottom, random flickering
				baseIntensity := float64(height-y) / float64(height) * intensity
				flicker := math.Sin(float64(frame)*0.3+float64(x)*0.2) * 0.3
				finalIntensity := baseIntensity + flicker

				if finalIntensity < 0 {
					finalIntensity = 0
				}
				if finalIntensity > 1 {
					finalIntensity = 1
				}

				charIndex := int(finalIntensity * float64(len(fireChars)-1))
				line.WriteString(fireChars[charIndex])
			}

			lines = append(lines, line.String())
		}

		frames[frame] = strings.Join(lines, "\n")
	}

	animation.frames = frames
	return animation
}

// WaterEffect creates a water-like ripple animation
func WaterEffect(width, height int, dropX, dropY int) *Animation {
	config := AnimationConfig{
		Type:   AnimationTypeWave,
		Speed:  100 * time.Millisecond,
		Frames: 30,
		Width:  width,
		Height: height,
		Loop:   false,
	}

	animation := NewAnimation(config)

	// Override with custom ripple frames
	frames := make([]string, config.Frames)
	rippleChars := []string{" ", "·", ":", "o", "O", "0", "O", "o", ":", "·"}

	for frame := 0; frame < config.Frames; frame++ {
		var lines []string
		rippleRadius := float64(frame) * 0.8

		for y := 0; y < height; y++ {
			var line strings.Builder

			for x := 0; x < width; x++ {
				// Calculate distance from drop point
				dx := float64(x - dropX)
				dy := float64(y - dropY)
				distance := math.Sqrt(dx*dx + dy*dy)

				// Create ripple effect
				if distance <= rippleRadius && distance >= rippleRadius-2 {
					intensity := 1.0 - math.Abs(distance-rippleRadius)/2.0
					if intensity < 0 {
						intensity = 0
					}

					charIndex := int(intensity * float64(len(rippleChars)-1))
					line.WriteString(rippleChars[charIndex])
				} else {
					line.WriteString(" ")
				}
			}

			lines = append(lines, line.String())
		}

		frames[frame] = strings.Join(lines, "\n")
	}

	animation.frames = frames
	return animation
}

// ExplosionEffect creates an explosion animation
func ExplosionEffect(centerX, centerY, maxRadius int) *Animation {
	config := AnimationConfig{
		Type:   AnimationTypePulse,
		Speed:  80 * time.Millisecond,
		Frames: 25,
		Width:  maxRadius * 2,
		Height: maxRadius * 2,
		Loop:   false,
	}

	animation := NewAnimation(config)

	// Override with custom explosion frames
	frames := make([]string, config.Frames)
	explosionChars := []string{"·", ":", "*", "x", "X", "#", "@", "█"}

	for frame := 0; frame < config.Frames; frame++ {
		var lines []string
		radius := float64(frame) * float64(maxRadius) / float64(config.Frames)

		for y := 0; y < config.Height; y++ {
			var line strings.Builder

			for x := 0; x < config.Width; x++ {
				dx := float64(x - centerX)
				dy := float64(y - centerY)
				distance := math.Sqrt(dx*dx + dy*dy)

				if distance <= radius {
					// Intensity based on distance and time
					intensity := 1.0 - distance/radius
					intensity *= 1.0 - float64(frame)/float64(config.Frames) // Fade over time

					if intensity > 0 {
						charIndex := int(intensity * float64(len(explosionChars)-1))
						if charIndex >= len(explosionChars) {
							charIndex = len(explosionChars) - 1
						}
						line.WriteString(explosionChars[charIndex])
					} else {
						line.WriteString(" ")
					}
				} else {
					line.WriteString(" ")
				}
			}

			lines = append(lines, line.String())
		}

		frames[frame] = strings.Join(lines, "\n")
	}

	animation.frames = frames
	return animation
}

// AnimationManager manages multiple animations
type AnimationManager struct {
	animations map[string]*Animation
	active     []string
}

// NewAnimationManager creates a new animation manager
func NewAnimationManager() *AnimationManager {
	return &AnimationManager{
		animations: make(map[string]*Animation),
		active:     make([]string, 0),
	}
}

// AddAnimation adds an animation to the manager
func (am *AnimationManager) AddAnimation(id string, animation *Animation) {
	am.animations[id] = animation
}

// PlayAnimation starts playing an animation
func (am *AnimationManager) PlayAnimation(id string) {
	if animation, exists := am.animations[id]; exists {
		animation.Play()

		// Add to active list if not already there
		for _, activeId := range am.active {
			if activeId == id {
				return
			}
		}
		am.active = append(am.active, id)
	}
}

// StopAnimation stops an animation
func (am *AnimationManager) StopAnimation(id string) {
	if animation, exists := am.animations[id]; exists {
		animation.Stop()

		// Remove from active list
		for i, activeId := range am.active {
			if activeId == id {
				am.active = append(am.active[:i], am.active[i+1:]...)
				break
			}
		}
	}
}

// GetCurrentFrame returns the current frame of an animation
func (am *AnimationManager) GetCurrentFrame(id string) string {
	if animation, exists := am.animations[id]; exists {
		return animation.GetCurrentFrame()
	}
	return ""
}

// UpdateAll updates all active animations
func (am *AnimationManager) UpdateAll() map[string]string {
	result := make(map[string]string)

	for _, id := range am.active {
		if animation, exists := am.animations[id]; exists && animation.IsPlaying() {
			result[id] = animation.GetCurrentFrame()
		}
	}

	return result
}

// StopAll stops all animations
func (am *AnimationManager) StopAll() {
	for _, animation := range am.animations {
		animation.Stop()
	}
	am.active = make([]string, 0)
}

// RemoveAnimation removes an animation from the manager
func (am *AnimationManager) RemoveAnimation(id string) {
	am.StopAnimation(id)
	delete(am.animations, id)
}
