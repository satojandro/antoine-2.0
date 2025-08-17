// Package terminal provides terminal size management utilities
// This file implements terminal size detection, monitoring, and responsive design helpers
package terminal

import (
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// SizeChangeCallback is called when terminal size changes
type SizeChangeCallback func(width, height int)

// SizeManager manages terminal size detection and monitoring
type SizeManager struct {
	currentWidth  int
	currentHeight int
	callbacks     []SizeChangeCallback
	monitoring    bool
	stopChan      chan bool
	mu            sync.RWMutex
}

// NewSizeManager creates a new size manager
func NewSizeManager() *SizeManager {
	sm := &SizeManager{
		stopChan: make(chan bool),
	}

	// Get initial size
	sm.currentWidth, sm.currentHeight = sm.detectSize()

	return sm
}

// detectSize detects the current terminal size
func (sm *SizeManager) detectSize() (width, height int) {
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

	// If we got both from environment, return them
	if width > 0 && height > 0 {
		return width, height
	}

	// Use detector for more accurate size
	detector := GetGlobalDetector()
	detectedWidth, detectedHeight := detector.getTerminalSize()

	// Use detected values or fallback to defaults
	if width == 0 {
		if detectedWidth > 0 {
			width = detectedWidth
		} else {
			width = 80 // Default width
		}
	}

	if height == 0 {
		if detectedHeight > 0 {
			height = detectedHeight
		} else {
			height = 24 // Default height
		}
	}

	return width, height
}

// GetSize returns the current terminal size
func (sm *SizeManager) GetSize() (width, height int) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.currentWidth, sm.currentHeight
}

// GetWidth returns the current terminal width
func (sm *SizeManager) GetWidth() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.currentWidth
}

// GetHeight returns the current terminal height
func (sm *SizeManager) GetHeight() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.currentHeight
}

// UpdateSize manually updates the terminal size
func (sm *SizeManager) UpdateSize(width, height int) {
	sm.mu.Lock()
	oldWidth, oldHeight := sm.currentWidth, sm.currentHeight
	sm.currentWidth = width
	sm.currentHeight = height
	sm.mu.Unlock()

	// Notify callbacks if size changed
	if width != oldWidth || height != oldHeight {
		sm.notifyCallbacks(width, height)
	}
}

// RefreshSize refreshes the terminal size from system
func (sm *SizeManager) RefreshSize() (width, height int) {
	width, height = sm.detectSize()
	sm.UpdateSize(width, height)
	return width, height
}

// AddCallback adds a callback for size changes
func (sm *SizeManager) AddCallback(callback SizeChangeCallback) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.callbacks = append(sm.callbacks, callback)
}

// RemoveCallback removes a callback (note: removes all instances)
func (sm *SizeManager) RemoveCallback(targetCallback SizeChangeCallback) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	var newCallbacks []SizeChangeCallback
	for _, callback := range sm.callbacks {
		// Note: This is a simplified comparison
		// In practice, you might want to use a more sophisticated method
		if &callback != &targetCallback {
			newCallbacks = append(newCallbacks, callback)
		}
	}
	sm.callbacks = newCallbacks
}

// notifyCallbacks notifies all callbacks of size change
func (sm *SizeManager) notifyCallbacks(width, height int) {
	sm.mu.RLock()
	callbacks := make([]SizeChangeCallback, len(sm.callbacks))
	copy(callbacks, sm.callbacks)
	sm.mu.RUnlock()

	for _, callback := range callbacks {
		go callback(width, height)
	}
}

// StartMonitoring starts monitoring for terminal size changes
func (sm *SizeManager) StartMonitoring() {
	sm.mu.Lock()
	if sm.monitoring {
		sm.mu.Unlock()
		return
	}
	sm.monitoring = true
	sm.mu.Unlock()

	// Monitor SIGWINCH signals for size changes
	go sm.monitorSignals()

	// Also poll periodically as a fallback
	go sm.monitorPolling()
}

// StopMonitoring stops monitoring for terminal size changes
func (sm *SizeManager) StopMonitoring() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if !sm.monitoring {
		return
	}

	sm.monitoring = false
	close(sm.stopChan)
	sm.stopChan = make(chan bool)
}

// monitorSignals monitors SIGWINCH signals for terminal resize
func (sm *SizeManager) monitorSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGWINCH)

	for {
		select {
		case <-sigChan:
			sm.RefreshSize()
		case <-sm.stopChan:
			signal.Stop(sigChan)
			return
		}
	}
}

// monitorPolling polls for size changes as a fallback
func (sm *SizeManager) monitorPolling() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Check if size has changed
			newWidth, newHeight := sm.detectSize()
			currentWidth, currentHeight := sm.GetSize()

			if newWidth != currentWidth || newHeight != currentHeight {
				sm.UpdateSize(newWidth, newHeight)
			}
		case <-sm.stopChan:
			return
		}
	}
}

// IsMonitoring returns true if currently monitoring size changes
func (sm *SizeManager) IsMonitoring() bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.monitoring
}

// Responsive design helpers

// BreakpointSize defines screen size breakpoints
type BreakpointSize int

const (
	BreakpointXS BreakpointSize = iota // < 40 columns
	BreakpointSM                       // 40-79 columns
	BreakpointMD                       // 80-119 columns
	BreakpointLG                       // 120-159 columns
	BreakpointXL                       // >= 160 columns
)

// Breakpoint thresholds
const (
	BreakpointXSMax = 39
	BreakpointSMMax = 79
	BreakpointMDMax = 119
	BreakpointLGMax = 159
)

// GetBreakpoint returns the current breakpoint based on terminal width
func (sm *SizeManager) GetBreakpoint() BreakpointSize {
	width := sm.GetWidth()

	switch {
	case width <= BreakpointXSMax:
		return BreakpointXS
	case width <= BreakpointSMMax:
		return BreakpointSM
	case width <= BreakpointMDMax:
		return BreakpointMD
	case width <= BreakpointLGMax:
		return BreakpointLG
	default:
		return BreakpointXL
	}
}

// IsSmallScreen returns true if terminal is considered small
func (sm *SizeManager) IsSmallScreen() bool {
	return sm.GetBreakpoint() <= BreakpointSM
}

// IsMediumScreen returns true if terminal is medium sized
func (sm *SizeManager) IsMediumScreen() bool {
	bp := sm.GetBreakpoint()
	return bp == BreakpointMD
}

// IsLargeScreen returns true if terminal is large
func (sm *SizeManager) IsLargeScreen() bool {
	return sm.GetBreakpoint() >= BreakpointLG
}

// GetOptimalWidth returns optimal width for content based on terminal size
func (sm *SizeManager) GetOptimalWidth() int {
	width := sm.GetWidth()

	// Leave some padding on the sides
	switch sm.GetBreakpoint() {
	case BreakpointXS:
		return width - 2 // Very little padding for tiny screens
	case BreakpointSM:
		return width - 4
	case BreakpointMD:
		return width - 8
	case BreakpointLG:
		return width - 12
	case BreakpointXL:
		return width - 16
	default:
		return width - 4
	}
}

// GetOptimalHeight returns optimal height for content
func (sm *SizeManager) GetOptimalHeight() int {
	height := sm.GetHeight()

	// Reserve space for header and footer
	if height <= 10 {
		return height - 2 // Minimal UI for very short terminals
	} else if height <= 20 {
		return height - 4
	} else {
		return height - 6
	}
}

// GetSafeArea returns safe area dimensions (width, height) for content
func (sm *SizeManager) GetSafeArea() (width, height int) {
	return sm.GetOptimalWidth(), sm.GetOptimalHeight()
}

// CalculateColumns calculates optimal number of columns for given item width
func (sm *SizeManager) CalculateColumns(itemWidth int) int {
	if itemWidth <= 0 {
		return 1
	}

	availableWidth := sm.GetOptimalWidth()
	columns := availableWidth / itemWidth

	if columns < 1 {
		columns = 1
	}

	return columns
}

// CalculateRows calculates optimal number of rows for given item height
func (sm *SizeManager) CalculateRows(itemHeight int) int {
	if itemHeight <= 0 {
		return 1
	}

	availableHeight := sm.GetOptimalHeight()
	rows := availableHeight / itemHeight

	if rows < 1 {
		rows = 1
	}

	return rows
}

// Layout calculation helpers

// LayoutInfo contains layout calculation results
type LayoutInfo struct {
	ContentWidth   int
	ContentHeight  int
	Columns        int
	Rows           int
	ItemWidth      int
	ItemHeight     int
	HorizontalGaps int
	VerticalGaps   int
	Breakpoint     BreakpointSize
}

// CalculateLayout calculates optimal layout for given constraints
func (sm *SizeManager) CalculateLayout(minItemWidth, minItemHeight int, gapSize int) LayoutInfo {
	termWidth, termHeight := sm.GetSize()
	contentWidth := sm.GetOptimalWidth()
	contentHeight := sm.GetOptimalHeight()

	// Calculate columns
	availableForItems := contentWidth
	columns := 1

	if minItemWidth > 0 {
		// Account for gaps between columns
		for cols := 1; cols*minItemWidth+(cols-1)*gapSize <= availableForItems; cols++ {
			columns = cols
		}
	}

	// Calculate actual item width
	itemWidth := minItemWidth
	if columns > 1 {
		totalGapSpace := (columns - 1) * gapSize
		itemWidth = (availableForItems - totalGapSpace) / columns
	} else {
		itemWidth = availableForItems
	}

	// Calculate rows
	availableForRows := contentHeight
	rows := 1

	if minItemHeight > 0 {
		// Account for gaps between rows
		for r := 1; r*minItemHeight+(r-1)*gapSize <= availableForRows; r++ {
			rows = r
		}
	}

	// Calculate actual item height
	itemHeight := minItemHeight
	if rows > 1 {
		totalGapSpace := (rows - 1) * gapSize
		itemHeight = (availableForRows - totalGapSpace) / rows
	} else {
		itemHeight = availableForRows
	}

	return LayoutInfo{
		ContentWidth:   contentWidth,
		ContentHeight:  contentHeight,
		Columns:        columns,
		Rows:           rows,
		ItemWidth:      itemWidth,
		ItemHeight:     itemHeight,
		HorizontalGaps: columns - 1,
		VerticalGaps:   rows - 1,
		Breakpoint:     sm.GetBreakpoint(),
	}
}

// ResponsiveValue returns different values based on breakpoint
func (sm *SizeManager) ResponsiveValue(xs, sm, md, lg, xl interface{}) interface{} {
	switch sm.GetBreakpoint() {
	case BreakpointXS:
		return xs
	case BreakpointSM:
		if sm != nil {
			return sm
		}
		return xs
	case BreakpointMD:
		if md != nil {
			return md
		}
		if sm != nil {
			return sm
		}
		return xs
	case BreakpointLG:
		if lg != nil {
			return lg
		}
		if md != nil {
			return md
		}
		if sm != nil {
			return sm
		}
		return xs
	case BreakpointXL:
		if xl != nil {
			return xl
		}
		if lg != nil {
			return lg
		}
		if md != nil {
			return md
		}
		if sm != nil {
			return sm
		}
		return xs
	default:
		return xs
	}
}

// ResponsiveInt returns different int values based on breakpoint
func (sm *SizeManager) ResponsiveInt(xs, sm, md, lg, xl int) int {
	result := sm.ResponsiveValue(xs, sm, md, lg, xl)
	if intResult, ok := result.(int); ok {
		return intResult
	}
	return xs
}

// ResponsiveString returns different string values based on breakpoint
func (sm *SizeManager) ResponsiveString(xs, sm, md, lg, xl string) string {
	result := sm.ResponsiveValue(xs, sm, md, lg, xl)
	if stringResult, ok := result.(string); ok {
		return stringResult
	}
	return xs
}

// Global size manager instance
var globalSizeManager *SizeManager
var sizeManagerOnce sync.Once

// GetGlobalSizeManager returns the global size manager
func GetGlobalSizeManager() *SizeManager {
	sizeManagerOnce.Do(func() {
		globalSizeManager = NewSizeManager()
		globalSizeManager.StartMonitoring()
	})
	return globalSizeManager
}

// Convenience functions using global size manager

// GetTerminalSize returns the current terminal size
func GetTerminalSize() (width, height int) {
	return GetGlobalSizeManager().GetSize()
}

// GetTerminalWidth returns the current terminal width
func GetTerminalWidth() int {
	return GetGlobalSizeManager().GetWidth()
}

// GetTerminalHeight returns the current terminal height
func GetTerminalHeight() int {
	return GetGlobalSizeManager().GetHeight()
}

// RefreshTerminalSize refreshes the terminal size
func RefreshTerminalSize() (width, height int) {
	return GetGlobalSizeManager().RefreshSize()
}

// AddSizeChangeCallback adds a callback for terminal size changes
func AddSizeChangeCallback(callback SizeChangeCallback) {
	GetGlobalSizeManager().AddCallback(callback)
}

// GetCurrentBreakpoint returns the current responsive breakpoint
func GetCurrentBreakpoint() BreakpointSize {
	return GetGlobalSizeManager().GetBreakpoint()
}

// IsSmallTerminal returns true if terminal is small
func IsSmallTerminal() bool {
	return GetGlobalSizeManager().IsSmallScreen()
}

// IsLargeTerminal returns true if terminal is large
func IsLargeTerminal() bool {
	return GetGlobalSizeManager().IsLargeScreen()
}

// GetOptimalContentWidth returns optimal width for content
func GetOptimalContentWidth() int {
	return GetGlobalSizeManager().GetOptimalWidth()
}

// GetOptimalContentHeight returns optimal height for content
func GetOptimalContentHeight() int {
	return GetGlobalSizeManager().GetOptimalHeight()
}

// GetSafeContentArea returns safe area for content
func GetSafeContentArea() (width, height int) {
	return GetGlobalSizeManager().GetSafeArea()
}

// ResponsiveWidth returns width based on breakpoint
func ResponsiveWidth(xs, sm, md, lg, xl int) int {
	return GetGlobalSizeManager().ResponsiveInt(xs, sm, md, lg, xl)
}

// ResponsiveHeight returns height based on breakpoint
func ResponsiveHeight(xs, sm, md, lg, xl int) int {
	return GetGlobalSizeManager().ResponsiveInt(xs, sm, md, lg, xl)
}

// ResponsiveColumns returns number of columns based on breakpoint
func ResponsiveColumns(xs, sm, md, lg, xl int) int {
	return GetGlobalSizeManager().ResponsiveInt(xs, sm, md, lg, xl)
}

// Helper functions for common responsive patterns

// FitTextToWidth calculates how to fit text within given width
func FitTextToWidth(text string, maxWidth int) (fittedText string, needsTruncation bool) {
	if len(text) <= maxWidth {
		return text, false
	}

	if maxWidth <= 3 {
		return "...", true
	}

	return text[:maxWidth-3] + "...", true
}

// CalculateGridLayout calculates grid layout for given item count
func CalculateGridLayout(itemCount, itemWidth, itemHeight, gapSize int) (columns, rows int) {
	sm := GetGlobalSizeManager()
	layout := sm.CalculateLayout(itemWidth, itemHeight, gapSize)

	columns = layout.Columns
	rows = (itemCount + columns - 1) / columns // Ceiling division

	return columns, rows
}

// WrapTextToWidth wraps text to fit within specified width
func WrapTextToWidth(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	var lines []string
	var currentLine []string
	currentLength := 0

	for _, word := range words {
		wordLength := len(word)

		// If adding this word would exceed width, start new line
		if currentLength > 0 && currentLength+1+wordLength > width {
			lines = append(lines, strings.Join(currentLine, " "))
			currentLine = []string{word}
			currentLength = wordLength
		} else {
			currentLine = append(currentLine, word)
			if currentLength > 0 {
				currentLength += 1 // Space
			}
			currentLength += wordLength
		}
	}

	// Add the last line
	if len(currentLine) > 0 {
		lines = append(lines, strings.Join(currentLine, " "))
	}

	return lines
}
