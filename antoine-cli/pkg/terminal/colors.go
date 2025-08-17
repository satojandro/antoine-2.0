// Package terminal provides color support detection and utilities
// This file implements color support detection and ANSI color utilities
package terminal

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Color represents an ANSI color code
type Color int

// Basic ANSI colors (3/4-bit)
const (
	ColorBlack Color = iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
	ColorBrightBlack   // 8
	ColorBrightRed     // 9
	ColorBrightGreen   // 10
	ColorBrightYellow  // 11
	ColorBrightBlue    // 12
	ColorBrightMagenta // 13
	ColorBrightCyan    // 14
	ColorBrightWhite   // 15
)

// ANSI color codes
const (
	// Reset
	ANSIReset = "\033[0m"

	// Foreground colors (30-37)
	ANSIBlack   = "\033[30m"
	ANSIRed     = "\033[31m"
	ANSIGreen   = "\033[32m"
	ANSIYellow  = "\033[33m"
	ANSIBlue    = "\033[34m"
	ANSIMagenta = "\033[35m"
	ANSICyan    = "\033[36m"
	ANSIWhite   = "\033[37m"

	// Bright foreground colors (90-97)
	ANSIBrightBlack   = "\033[90m"
	ANSIBrightRed     = "\033[91m"
	ANSIBrightGreen   = "\033[92m"
	ANSIBrightYellow  = "\033[93m"
	ANSIBrightBlue    = "\033[94m"
	ANSIBrightMagenta = "\033[95m"
	ANSIBrightCyan    = "\033[96m"
	ANSIBrightWhite   = "\033[97m"

	// Background colors (40-47)
	ANSIBgBlack   = "\033[40m"
	ANSIBgRed     = "\033[41m"
	ANSIBgGreen   = "\033[42m"
	ANSIBgYellow  = "\033[43m"
	ANSIBgBlue    = "\033[44m"
	ANSIBgMagenta = "\033[45m"
	ANSIBgCyan    = "\033[46m"
	ANSIBgWhite   = "\033[47m"

	// Bright background colors (100-107)
	ANSIBgBrightBlack   = "\033[100m"
	ANSIBgBrightRed     = "\033[101m"
	ANSIBgBrightGreen   = "\033[102m"
	ANSIBgBrightYellow  = "\033[103m"
	ANSIBgBrightBlue    = "\033[104m"
	ANSIBgBrightMagenta = "\033[105m"
	ANSIBgBrightCyan    = "\033[106m"
	ANSIBgBrightWhite   = "\033[107m"

	// Text formatting
	ANSIBold      = "\033[1m"
	ANSIDim       = "\033[2m"
	ANSIItalic    = "\033[3m"
	ANSIUnderline = "\033[4m"
	ANSIBlink     = "\033[5m"
	ANSIReverse   = "\033[7m"
	ANSIStrike    = "\033[9m"
)

// RGB represents an RGB color
type RGB struct {
	R, G, B uint8
}

// NewRGB creates a new RGB color
func NewRGB(r, g, b uint8) RGB {
	return RGB{R: r, G: g, B: b}
}

// ToANSI returns the ANSI escape sequence for this RGB color as foreground
func (rgb RGB) ToANSI() string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", rgb.R, rgb.G, rgb.B)
}

// ToANSIBackground returns the ANSI escape sequence for this RGB color as background
func (rgb RGB) ToANSIBackground() string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", rgb.R, rgb.G, rgb.B)
}

// To256 converts RGB to closest 256-color palette index
func (rgb RGB) To256() uint8 {
	// Convert to 256-color palette
	// Colors 0-15: standard colors
	// Colors 16-231: 6×6×6 color cube
	// Colors 232-255: grayscale

	// Check if it's a grayscale color
	if rgb.R == rgb.G && rgb.G == rgb.B {
		// Grayscale (232-255)
		if rgb.R < 8 {
			return 16 // Black
		}
		if rgb.R > 248 {
			return 231 // White
		}
		return 232 + uint8((int(rgb.R)-8)*23/247)
	}

	// Color cube (16-231)
	// Each component is mapped to 0-5
	r := uint8(rgb.R * 5 / 255)
	g := uint8(rgb.G * 5 / 255)
	b := uint8(rgb.B * 5 / 255)

	return 16 + 36*r + 6*g + b
}

// ColorManager manages color support and provides color utilities
type ColorManager struct {
	detector     *TerminalDetector
	forceColor   bool
	disableColor bool
	colorProfile ColorSupport
}

// NewColorManager creates a new color manager
func NewColorManager() *ColorManager {
	return &ColorManager{
		detector: GetGlobalDetector(),
	}
}

// SetForceColor forces color output regardless of detection
func (cm *ColorManager) SetForceColor(force bool) {
	cm.forceColor = force
}

// SetDisableColor disables all color output
func (cm *ColorManager) SetDisableColor(disable bool) {
	cm.disableColor = disable
}

// SupportsColor returns true if colors should be used
func (cm *ColorManager) SupportsColor() bool {
	if cm.disableColor {
		return false
	}
	if cm.forceColor {
		return true
	}

	info := cm.detector.GetInfo()
	return info.SupportsColor
}

// Supports256Color returns true if 256 colors are supported
func (cm *ColorManager) Supports256Color() bool {
	if !cm.SupportsColor() {
		return false
	}

	info := cm.detector.GetInfo()
	return info.Supports256
}

// SupportsTrueColor returns true if true color (24-bit) is supported
func (cm *ColorManager) SupportsTrueColor() bool {
	if !cm.SupportsColor() {
		return false
	}

	info := cm.detector.GetInfo()
	return info.SupportsTrueColor
}

// GetColorProfile returns the color support profile
func (cm *ColorManager) GetColorProfile() ColorSupport {
	if cm.disableColor {
		return ColorSupportNone
	}
	if cm.forceColor {
		return ColorSupportTrueColor
	}

	return cm.detector.GetColorSupport()
}

// WrapColor wraps text with ANSI color codes if supported
func (cm *ColorManager) WrapColor(text, colorCode string) string {
	if !cm.SupportsColor() {
		return text
	}
	return colorCode + text + ANSIReset
}

// WrapRGB wraps text with RGB color if supported
func (cm *ColorManager) WrapRGB(text string, rgb RGB) string {
	if cm.SupportsTrueColor() {
		return rgb.ToANSI() + text + ANSIReset
	} else if cm.Supports256Color() {
		return cm.Wrap256(text, rgb.To256())
	} else if cm.SupportsColor() {
		// Fallback to basic color
		basicColor := cm.rgbToBasicColor(rgb)
		return cm.WrapBasicColor(text, basicColor)
	}
	return text
}

// Wrap256 wraps text with 256-color palette
func (cm *ColorManager) Wrap256(text string, colorIndex uint8) string {
	if !cm.Supports256Color() {
		return text
	}
	return fmt.Sprintf("\033[38;5;%dm%s%s", colorIndex, text, ANSIReset)
}

// WrapBasicColor wraps text with basic ANSI color
func (cm *ColorManager) WrapBasicColor(text string, color Color) string {
	if !cm.SupportsColor() {
		return text
	}

	var colorCode string
	switch color {
	case ColorBlack:
		colorCode = ANSIBlack
	case ColorRed:
		colorCode = ANSIRed
	case ColorGreen:
		colorCode = ANSIGreen
	case ColorYellow:
		colorCode = ANSIYellow
	case ColorBlue:
		colorCode = ANSIBlue
	case ColorMagenta:
		colorCode = ANSIMagenta
	case ColorCyan:
		colorCode = ANSICyan
	case ColorWhite:
		colorCode = ANSIWhite
	case ColorBrightBlack:
		colorCode = ANSIBrightBlack
	case ColorBrightRed:
		colorCode = ANSIBrightRed
	case ColorBrightGreen:
		colorCode = ANSIBrightGreen
	case ColorBrightYellow:
		colorCode = ANSIBrightYellow
	case ColorBrightBlue:
		colorCode = ANSIBrightBlue
	case ColorBrightMagenta:
		colorCode = ANSIBrightMagenta
	case ColorBrightCyan:
		colorCode = ANSIBrightCyan
	case ColorBrightWhite:
		colorCode = ANSIBrightWhite
	default:
		return text
	}

	return colorCode + text + ANSIReset
}

// rgbToBasicColor converts RGB to closest basic ANSI color
func (cm *ColorManager) rgbToBasicColor(rgb RGB) Color {
	// Simple mapping to basic colors
	if rgb.R > 128 && rgb.G < 64 && rgb.B < 64 {
		return ColorRed
	}
	if rgb.R < 64 && rgb.G > 128 && rgb.B < 64 {
		return ColorGreen
	}
	if rgb.R > 128 && rgb.G > 128 && rgb.B < 64 {
		return ColorYellow
	}
	if rgb.R < 64 && rgb.G < 64 && rgb.B > 128 {
		return ColorBlue
	}
	if rgb.R > 128 && rgb.G < 64 && rgb.B > 128 {
		return ColorMagenta
	}
	if rgb.R < 64 && rgb.G > 128 && rgb.B > 128 {
		return ColorCyan
	}
	if rgb.R > 128 && rgb.G > 128 && rgb.B > 128 {
		return ColorWhite
	}
	return ColorBlack
}

// StripANSI removes ANSI escape sequences from text
func (cm *ColorManager) StripANSI(text string) string {
	// Simple ANSI stripping - in production, you might want a more robust regex
	ansiRegex := regexp.MustCompile(`\033\[[0-9;]*m`)
	return ansiRegex.ReplaceAllString(text, "")
}

// Colorize provides convenient color methods
type Colorize struct {
	manager *ColorManager
}

// NewColorize creates a new colorize instance
func NewColorize() *Colorize {
	return &Colorize{
		manager: NewColorManager(),
	}
}

// Red colors text red
func (c *Colorize) Red(text string) string {
	return c.manager.WrapColor(text, ANSIRed)
}

// Green colors text green
func (c *Colorize) Green(text string) string {
	return c.manager.WrapColor(text, ANSIGreen)
}

// Yellow colors text yellow
func (c *Colorize) Yellow(text string) string {
	return c.manager.WrapColor(text, ANSIYellow)
}

// Blue colors text blue
func (c *Colorize) Blue(text string) string {
	return c.manager.WrapColor(text, ANSIBlue)
}

// Magenta colors text magenta
func (c *Colorize) Magenta(text string) string {
	return c.manager.WrapColor(text, ANSIMagenta)
}

// Cyan colors text cyan
func (c *Colorize) Cyan(text string) string {
	return c.manager.WrapColor(text, ANSICyan)
}

// White colors text white
func (c *Colorize) White(text string) string {
	return c.manager.WrapColor(text, ANSIWhite)
}

// Black colors text black
func (c *Colorize) Black(text string) string {
	return c.manager.WrapColor(text, ANSIBlack)
}

// Bold makes text bold
func (c *Colorize) Bold(text string) string {
	return c.manager.WrapColor(text, ANSIBold)
}

// Underline underlines text
func (c *Colorize) Underline(text string) string {
	return c.manager.WrapColor(text, ANSIUnderline)
}

// Italic makes text italic
func (c *Colorize) Italic(text string) string {
	return c.manager.WrapColor(text, ANSIItalic)
}

// Strike strikes through text
func (c *Colorize) Strike(text string) string {
	return c.manager.WrapColor(text, ANSIStrike)
}

// RGB colors text with RGB values
func (c *Colorize) RGB(text string, r, g, b uint8) string {
	return c.manager.WrapRGB(text, NewRGB(r, g, b))
}

// Hex colors text with hex color code
func (c *Colorize) Hex(text, hexColor string) string {
	rgb, err := HexToRGB(hexColor)
	if err != nil {
		return text
	}
	return c.manager.WrapRGB(text, rgb)
}

// Color utility functions

// HexToRGB converts hex color string to RGB
func HexToRGB(hex string) (RGB, error) {
	// Remove # if present
	hex = strings.TrimPrefix(hex, "#")

	if len(hex) == 3 {
		// Expand shorthand hex (e.g., "f0a" -> "ff00aa")
		hex = string(hex[0]) + string(hex[0]) + string(hex[1]) + string(hex[1]) + string(hex[2]) + string(hex[2])
	}

	if len(hex) != 6 {
		return RGB{}, fmt.Errorf("invalid hex color: %s", hex)
	}

	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return RGB{}, err
	}

	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return RGB{}, err
	}

	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return RGB{}, err
	}

	return NewRGB(uint8(r), uint8(g), uint8(b)), nil
}

// RGBToHex converts RGB to hex string
func RGBToHex(rgb RGB) string {
	return fmt.Sprintf("#%02x%02x%02x", rgb.R, rgb.G, rgb.B)
}

// HSVToRGB converts HSV to RGB
func HSVToRGB(h, s, v float64) RGB {
	c := v * s
	x := c * (1 - abs(mod(h/60, 2)-1))
	m := v - c

	var r, g, b float64

	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return NewRGB(
		uint8((r+m)*255),
		uint8((g+m)*255),
		uint8((b+m)*255),
	)
}

// Helper functions
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func mod(x, y float64) float64 {
	return x - y*float64(int(x/y))
}

// Gradient creates a color gradient between two RGB colors
func Gradient(start, end RGB, steps int) []RGB {
	if steps <= 0 {
		return []RGB{}
	}
	if steps == 1 {
		return []RGB{start}
	}

	gradient := make([]RGB, steps)
	gradient[0] = start
	gradient[steps-1] = end

	for i := 1; i < steps-1; i++ {
		ratio := float64(i) / float64(steps-1)

		r := uint8(float64(start.R)*(1-ratio) + float64(end.R)*ratio)
		g := uint8(float64(start.G)*(1-ratio) + float64(end.G)*ratio)
		b := uint8(float64(start.B)*(1-ratio) + float64(end.B)*ratio)

		gradient[i] = NewRGB(r, g, b)
	}

	return gradient
}

// Rainbow generates rainbow colors
func Rainbow(steps int) []RGB {
	colors := make([]RGB, steps)

	for i := 0; i < steps; i++ {
		hue := float64(i) * 360 / float64(steps)
		colors[i] = HSVToRGB(hue, 1.0, 1.0)
	}

	return colors
}

// Global color manager and colorize instances
var globalColorManager *ColorManager
var globalColorize *Colorize

// GetGlobalColorManager returns the global color manager
func GetGlobalColorManager() *ColorManager {
	if globalColorManager == nil {
		globalColorManager = NewColorManager()

		// Check environment variables for color preferences
		if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
			globalColorManager.SetDisableColor(true)
		}
		if os.Getenv("FORCE_COLOR") != "" {
			globalColorManager.SetForceColor(true)
		}
	}
	return globalColorManager
}

// GetGlobalColorize returns the global colorize instance
func GetGlobalColorize() *Colorize {
	if globalColorize == nil {
		globalColorize = NewColorize()
	}
	return globalColorize
}

// Convenience functions using global instances

// SupportsTerminalColor returns true if terminal supports color
func SupportsTerminalColor() bool {
	return GetGlobalColorManager().SupportsColor()
}

// Supports256TerminalColor returns true if terminal supports 256 colors
func Supports256TerminalColor() bool {
	return GetGlobalColorManager().Supports256Color()
}

// SupportsTrueTerminalColor returns true if terminal supports true color
func SupportsTrueTerminalColor() bool {
	return GetGlobalColorManager().SupportsTrueColor()
}

// WrapTerminalColor wraps text with color if supported
func WrapTerminalColor(text, colorCode string) string {
	return GetGlobalColorManager().WrapColor(text, colorCode)
}

// WrapTerminalRGB wraps text with RGB color if supported
func WrapTerminalRGB(text string, r, g, b uint8) string {
	return GetGlobalColorManager().WrapRGB(text, NewRGB(r, g, b))
}

// St
