// Package terminal provides terminal detection and capability utilities
// This file implements terminal capability detection and feature support
package terminal

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

// TerminalInfo contains information about the terminal capabilities
type TerminalInfo struct {
	IsTTY             bool   // Whether output is to a terminal
	SupportsColor     bool   // 8/16 color support
	Supports256       bool   // 256 color support
	SupportsTrueColor bool   // 24-bit RGB color support
	SupportsUnicode   bool   // Unicode character support
	Width             int    // Terminal width in columns
	Height            int    // Terminal height in rows
	Type              string // Terminal type (from $TERM)
	Program           string // Terminal program name
	Version           string // Terminal version if detectable
	Platform          string // Operating system
	Shell             string // Shell being used
}

// ColorSupport defines the level of color support
type ColorSupport int

const (
	ColorSupportNone      ColorSupport = iota
	ColorSupportBasic                  // 8/16 colors
	ColorSupport256                    // 256 colors
	ColorSupportTrueColor              // 24-bit RGB
)

// Feature represents a terminal feature
type Feature string

const (
	FeatureColor             Feature = "color"
	FeatureUnicode           Feature = "unicode"
	FeatureMouseSupport      Feature = "mouse"
	FeatureCursorPositioning Feature = "cursor"
	FeatureAlternateScreen   Feature = "altscreen"
	FeatureBracketedPaste    Feature = "bracketed_paste"
	FeatureWindowTitle       Feature = "window_title"
	FeatureResize            Feature = "resize"
	FeatureScrollback        Feature = "scrollback"
)

// TerminalDetector provides terminal detection capabilities
type TerminalDetector struct {
	info     *TerminalInfo
	features map[Feature]bool
	cached   bool
}

// NewTerminalDetector creates a new terminal detector
func NewTerminalDetector() *TerminalDetector {
	return &TerminalDetector{
		features: make(map[Feature]bool),
		cached:   false,
	}
}

// DetectCapabilities detects all terminal capabilities
func (td *TerminalDetector) DetectCapabilities() *TerminalInfo {
	if td.cached && td.info != nil {
		return td.info
	}

	info := &TerminalInfo{
		Platform: runtime.GOOS,
		Shell:    detectShell(),
	}

	// Detect TTY
	info.IsTTY = td.isTTY()

	// Detect terminal dimensions
	info.Width, info.Height = td.getTerminalSize()

	// Detect terminal type and program
	info.Type = os.Getenv("TERM")
	info.Program = td.detectTerminalProgram()
	info.Version = td.detectTerminalVersion()

	// Detect color support
	colorSupport := td.detectColorSupport()
	info.SupportsColor = colorSupport >= ColorSupportBasic
	info.Supports256 = colorSupport >= ColorSupport256
	info.SupportsTrueColor = colorSupport >= ColorSupportTrueColor

	// Detect Unicode support
	info.SupportsUnicode = td.detectUnicodeSupport()

	// Detect features
	td.detectFeatures()

	td.info = info
	td.cached = true

	return info
}

// isTTY checks if the output is connected to a terminal
func (td *TerminalDetector) isTTY() bool {
	// Check stdout
	if fileInfo, err := os.Stdout.Stat(); err == nil {
		return (fileInfo.Mode() & os.ModeCharDevice) != 0
	}
	return false
}

// getTerminalSize returns the terminal dimensions
func (td *TerminalDetector) getTerminalSize() (width, height int) {
	width, height = 80, 24 // Default fallback

	// Try environment variables first
	if w := os.Getenv("COLUMNS"); w != "" {
		if parsed, err := strconv.Atoi(w); err == nil && parsed > 0 {
			width = parsed
		}
	}
	if h := os.Getenv("LINES"); h != "" {
		if parsed, err := strconv.Atoi(h); err == nil && parsed > 0 {
			height = parsed
		}
	}

	// Try system call for more accurate size
	if runtime.GOOS != "windows" {
		if w, h := td.getUnixTerminalSize(); w > 0 && h > 0 {
			width, height = w, h
		}
	} else {
		if w, h := td.getWindowsTerminalSize(); w > 0 && h > 0 {
			width, height = w, h
		}
	}

	return width, height
}

// getUnixTerminalSize gets terminal size on Unix-like systems
func (td *TerminalDetector) getUnixTerminalSize() (width, height int) {
	type winsize struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}

	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(0x5413), // TIOCGWINSZ
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		// If ioctl fails, try stderr and stdout
		for _, fd := range []uintptr{uintptr(syscall.Stderr), uintptr(syscall.Stdout)} {
			retCode, _, errno = syscall.Syscall(syscall.SYS_IOCTL,
				fd,
				uintptr(0x5413),
				uintptr(unsafe.Pointer(ws)))
			if int(retCode) != -1 {
				break
			}
		}
	}

	if int(retCode) == -1 || errno != 0 {
		return 0, 0
	}

	return int(ws.Col), int(ws.Row)
}

// getWindowsTerminalSize gets terminal size on Windows
func (td *TerminalDetector) getWindowsTerminalSize() (width, height int) {
	// Windows implementation would use GetConsoleScreenBufferInfo
	// For now, return defaults and rely on environment variables
	return 0, 0
}

// detectColorSupport detects the level of color support
func (td *TerminalDetector) detectColorSupport() ColorSupport {
	// Check environment variables
	term := strings.ToLower(os.Getenv("TERM"))
	colorterm := strings.ToLower(os.Getenv("COLORTERM"))

	// Force no color if explicitly disabled
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return ColorSupportNone
	}

	// Check for truecolor support
	if colorterm == "truecolor" || colorterm == "24bit" ||
		strings.Contains(term, "truecolor") || strings.Contains(term, "24bit") {
		return ColorSupportTrueColor
	}

	// Check specific terminal types for truecolor
	if td.supportsTrueColor(term) {
		return ColorSupportTrueColor
	}

	// Check for 256 color support
	if strings.Contains(term, "256") || strings.Contains(term, "256color") {
		return ColorSupport256
	}

	// Check specific terminal types for 256 colors
	if td.supports256Color(term) {
		return ColorSupport256
	}

	// Check for basic color support
	if term != "" && term != "dumb" {
		if strings.Contains(term, "color") || td.supportsBasicColor(term) {
			return ColorSupportBasic
		}
	}

	// Check Windows specific
	if runtime.GOOS == "windows" {
		return td.detectWindowsColorSupport()
	}

	return ColorSupportNone
}

// supportsTrueColor checks if terminal supports true color
func (td *TerminalDetector) supportsTrueColor(term string) bool {
	trueColorTerminals := []string{
		"iterm2", "iterm", "kitty", "alacritty", "wezterm",
		"gnome-terminal", "konsole", "terminator", "tilix",
		"hyper", "terminus", "windows-terminal", "mintty",
		"tmux-256color", "screen-256color",
	}

	for _, supportedTerm := range trueColorTerminals {
		if strings.Contains(term, supportedTerm) {
			return true
		}
	}

	return false
}

// supports256Color checks if terminal supports 256 colors
func (td *TerminalDetector) supports256Color(term string) bool {
	color256Terminals := []string{
		"xterm", "screen", "tmux", "rxvt", "putty",
		"cygwin", "linux", "ansi",
	}

	for _, supportedTerm := range color256Terminals {
		if strings.Contains(term, supportedTerm) {
			return true
		}
	}

	return false
}

// supportsBasicColor checks if terminal supports basic colors
func (td *TerminalDetector) supportsBasicColor(term string) bool {
	basicColorTerminals := []string{
		"vt100", "vt220", "vt320", "vt420", "vt520",
		"xterm", "screen", "tmux", "color",
	}

	for _, supportedTerm := range basicColorTerminals {
		if strings.Contains(term, supportedTerm) {
			return true
		}
	}

	return false
}

// detectWindowsColorSupport detects color support on Windows
func (td *TerminalDetector) detectWindowsColorSupport() ColorSupport {
	// Windows Terminal and modern terminals support truecolor
	if os.Getenv("WT_SESSION") != "" { // Windows Terminal
		return ColorSupportTrueColor
	}

	// Check Windows version for ANSI support
	if td.windowsSupportsANSI() {
		return ColorSupport256
	}

	return ColorSupportBasic
}

// windowsSupportsANSI checks if Windows supports ANSI escape sequences
func (td *TerminalDetector) windowsSupportsANSI() bool {
	// Windows 10 version 1511 and later support ANSI
	// This is a simplified check
	return true // Assume modern Windows for now
}

// detectUnicodeSupport detects Unicode support
func (td *TerminalDetector) detectUnicodeSupport() bool {
	// Check locale settings
	for _, env := range []string{"LC_ALL", "LC_CTYPE", "LANG"} {
		if locale := os.Getenv(env); locale != "" {
			if strings.Contains(strings.ToLower(locale), "utf") ||
				strings.Contains(strings.ToLower(locale), "unicode") {
				return true
			}
		}
	}

	// Check terminal type
	term := strings.ToLower(os.Getenv("TERM"))
	unicodeTerminals := []string{
		"xterm-256color", "screen-256color", "tmux-256color",
		"iterm", "kitty", "alacritty", "gnome-terminal",
	}

	for _, supportedTerm := range unicodeTerminals {
		if strings.Contains(term, supportedTerm) {
			return true
		}
	}

	// Windows Terminal supports Unicode
	if runtime.GOOS == "windows" && os.Getenv("WT_SESSION") != "" {
		return true
	}

	return false
}

// detectTerminalProgram detects the terminal program being used
func (td *TerminalDetector) detectTerminalProgram() string {
	// Check environment variables
	if program := os.Getenv("TERM_PROGRAM"); program != "" {
		return program
	}

	// Check specific environment variables
	if os.Getenv("ITERM_SESSION_ID") != "" {
		return "iTerm2"
	}
	if os.Getenv("KITTY_WINDOW_ID") != "" {
		return "Kitty"
	}
	if os.Getenv("ALACRITTY_SOCKET") != "" {
		return "Alacritty"
	}
	if os.Getenv("WEZTERM_EXECUTABLE") != "" {
		return "WezTerm"
	}
	if os.Getenv("WT_SESSION") != "" {
		return "Windows Terminal"
	}
	if os.Getenv("GNOME_TERMINAL_SERVICE") != "" {
		return "GNOME Terminal"
	}
	if os.Getenv("KONSOLE_VERSION") != "" {
		return "Konsole"
	}

	// Check TERM variable for hints
	term := os.Getenv("TERM")
	if strings.Contains(term, "screen") {
		return "GNU Screen"
	}
	if strings.Contains(term, "tmux") {
		return "tmux"
	}
	if strings.Contains(term, "xterm") {
		return "xterm"
	}

	return "Unknown"
}

// detectTerminalVersion detects the terminal version if possible
func (td *TerminalDetector) detectTerminalVersion() string {
	// Check environment variables for versions
	if version := os.Getenv("TERM_PROGRAM_VERSION"); version != "" {
		return version
	}
	if version := os.Getenv("KITTY_VERSION"); version != "" {
		return version
	}
	if version := os.Getenv("KONSOLE_VERSION"); version != "" {
		return version
	}

	return ""
}

// detectShell detects the shell being used
func detectShell() string {
	// Check SHELL environment variable
	if shell := os.Getenv("SHELL"); shell != "" {
		// Extract just the shell name
		parts := strings.Split(shell, "/")
		return parts[len(parts)-1]
	}

	// Windows specific
	if runtime.GOOS == "windows" {
		if os.Getenv("PSModulePath") != "" {
			return "PowerShell"
		}
		return "cmd"
	}

	return "Unknown"
}

// detectFeatures detects specific terminal features
func (td *TerminalDetector) detectFeatures() {
	term := strings.ToLower(os.Getenv("TERM"))
	program := strings.ToLower(td.info.Program)

	// Color support (already detected)
	td.features[FeatureColor] = td.info.SupportsColor

	// Unicode support (already detected)
	td.features[FeatureUnicode] = td.info.SupportsUnicode

	// Mouse support
	td.features[FeatureMouseSupport] = td.detectMouseSupport(term, program)

	// Cursor positioning
	td.features[FeatureCursorPositioning] = td.detectCursorSupport(term, program)

	// Alternate screen
	td.features[FeatureAlternateScreen] = td.detectAlternateScreen(term, program)

	// Bracketed paste
	td.features[FeatureBracketedPaste] = td.detectBracketedPaste(term, program)

	// Window title support
	td.features[FeatureWindowTitle] = td.detectWindowTitle(term, program)

	// Resize support
	td.features[FeatureResize] = td.info.IsTTY

	// Scrollback support
	td.features[FeatureScrollback] = td.detectScrollback(term, program)
}

// detectMouseSupport detects mouse support
func (td *TerminalDetector) detectMouseSupport(term, program string) bool {
	// Most modern terminals support mouse
	modernTerminals := []string{
		"xterm", "gnome-terminal", "konsole", "iterm", "kitty",
		"alacritty", "wezterm", "windows terminal", "hyper",
	}

	for _, modernTerm := range modernTerminals {
		if strings.Contains(term, modernTerm) || strings.Contains(program, modernTerm) {
			return true
		}
	}

	return false
}

// detectCursorSupport detects cursor positioning support
func (td *TerminalDetector) detectCursorSupport(term, program string) bool {
	// Almost all terminals support cursor positioning
	return term != "dumb" && td.info.IsTTY
}

// detectAlternateScreen detects alternate screen support
func (td *TerminalDetector) detectAlternateScreen(term, program string) bool {
	// Most modern terminals support alternate screen
	return strings.Contains(term, "xterm") ||
		strings.Contains(term, "screen") ||
		strings.Contains(term, "tmux") ||
		strings.Contains(program, "iterm") ||
		strings.Contains(program, "kitty") ||
		strings.Contains(program, "alacritty")
}

// detectBracketedPaste detects bracketed paste support
func (td *TerminalDetector) detectBracketedPaste(term, program string) bool {
	// Modern terminals support bracketed paste
	modernTerminals := []string{
		"iterm", "kitty", "alacritty", "gnome-terminal",
		"konsole", "wezterm", "windows terminal",
	}

	for _, modernTerm := range modernTerminals {
		if strings.Contains(program, modernTerm) {
			return true
		}
	}

	return false
}

// detectWindowTitle detects window title support
func (td *TerminalDetector) detectWindowTitle(term, program string) bool {
	// Most GUI terminals support window titles
	guiTerminals := []string{
		"iterm", "kitty", "alacritty", "gnome-terminal",
		"konsole", "wezterm", "windows terminal", "hyper",
	}

	for _, guiTerm := range guiTerminals {
		if strings.Contains(program, guiTerm) {
			return true
		}
	}

	return false
}

// detectScrollback detects scrollback support
func (td *TerminalDetector) detectScrollback(term, program string) bool {
	// Most terminals support scrollback
	return td.info.IsTTY && term != "dumb"
}

// Public methods for accessing capabilities

// GetInfo returns the terminal information
func (td *TerminalDetector) GetInfo() *TerminalInfo {
	if !td.cached {
		return td.DetectCapabilities()
	}
	return td.info
}

// SupportsFeature checks if a specific feature is supported
func (td *TerminalDetector) SupportsFeature(feature Feature) bool {
	if !td.cached {
		td.DetectCapabilities()
	}
	return td.features[feature]
}

// GetColorSupport returns the level of color support
func (td *TerminalDetector) GetColorSupport() ColorSupport {
	info := td.GetInfo()
	if info.SupportsTrueColor {
		return ColorSupportTrueColor
	} else if info.Supports256 {
		return ColorSupport256
	} else if info.SupportsColor {
		return ColorSupportBasic
	}
	return ColorSupportNone
}

// GetTerminalSize returns the current terminal size
func (td *TerminalDetector) GetTerminalSize() (width, height int) {
	// Always get fresh size as it can change
	return td.getTerminalSize()
}

// IsTerminal checks if we're running in a terminal
func (td *TerminalDetector) IsTerminal() bool {
	return td.GetInfo().IsTTY
}

// GetPlatformInfo returns platform-specific information
func (td *TerminalDetector) GetPlatformInfo() map[string]string {
	info := td.GetInfo()
	return map[string]string{
		"platform":         info.Platform,
		"terminal_program": info.Program,
		"terminal_version": info.Version,
		"terminal_type":    info.Type,
		"shell":            info.Shell,
	}
}

// GetCapabilitySummary returns a summary of all capabilities
func (td *TerminalDetector) GetCapabilitySummary() map[string]interface{} {
	info := td.GetInfo()
	summary := map[string]interface{}{
		"is_tty":              info.IsTTY,
		"supports_color":      info.SupportsColor,
		"supports_256_color":  info.Supports256,
		"supports_true_color": info.SupportsTrueColor,
		"supports_unicode":    info.SupportsUnicode,
		"width":               info.Width,
		"height":              info.Height,
		"terminal_program":    info.Program,
		"terminal_type":       info.Type,
		"platform":            info.Platform,
		"shell":               info.Shell,
	}

	// Add feature support
	features := make(map[string]bool)
	for feature, supported := range td.features {
		features[string(feature)] = supported
	}
	summary["features"] = features

	return summary
}

// RefreshCapabilities forces a refresh of terminal capabilities
func (td *TerminalDetector) RefreshCapabilities() *TerminalInfo {
	td.cached = false
	td.info = nil
	td.features = make(map[Feature]bool)
	return td.DetectCapabilities()
}

// Global detector instance
var globalDetector *TerminalDetector

// GetGlobalDetector returns the global terminal detector
func GetGlobalDetector() *TerminalDetector {
	if globalDetector == nil {
		globalDetector = NewTerminalDetector()
	}
	return globalDetector
}

// Convenience functions using global detector

// DetectTerminal detects terminal capabilities using global detector
func DetectTerminal() *TerminalInfo {
	return GetGlobalDetector().DetectCapabilities()
}

// IsColorSupported checks if color is supported
func IsColorSupported() bool {
	return GetGlobalDetector().GetInfo().SupportsColor
}

// IsTrueColorSupported checks if true color is supported
func IsTrueColorSupported() bool {
	return GetGlobalDetector().GetInfo().SupportsTrueColor
}

// IsUnicodeSupported checks if Unicode is supported
func IsUnicodeSupported() bool {
	return GetGlobalDetector().GetInfo().SupportsUnicode
}

// GetSize returns the terminal size
func GetSize() (width, height int) {
	return GetGlobalDetector().GetTerminalSize()
}

// IsRunningInTerminal checks if running in a terminal
func IsRunningInTerminal() bool {
	return GetGlobalDetector().IsTerminal()
}

// GetTerminalProgram returns the terminal program name
func GetTerminalProgram() string {
	return GetGlobalDetector().GetInfo().Program
}

// PrintCapabilities prints all terminal capabilities (for debugging)
func PrintCapabilities() {
	detector := GetGlobalDetector()
	info := detector.GetInfo()

	fmt.Printf("Terminal Capabilities:\n")
	fmt.Printf("  TTY: %v\n", info.IsTTY)
	fmt.Printf("  Size: %dx%d\n", info.Width, info.Height)
	fmt.Printf("  Program: %s\n", info.Program)
	fmt.Printf("  Type: %s\n", info.Type)
	fmt.Printf("  Version: %s\n", info.Version)
	fmt.Printf("  Platform: %s\n", info.Platform)
	fmt.Printf("  Shell: %s\n", info.Shell)
	fmt.Printf("  Color Support: %v\n", info.SupportsColor)
	fmt.Printf("  256 Color: %v\n", info.Supports256)
	fmt.Printf("  True Color: %v\n", info.SupportsTrueColor)
	fmt.Printf("  Unicode: %v\n", info.SupportsUnicode)

	fmt.Printf("  Features:\n")
	for feature, supported := range detector.features {
		fmt.Printf("    %s: %v\n", feature, supported)
	}
}
