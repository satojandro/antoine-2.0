// Package utils provides helper utilities for Antoine CLI
// This file contains common helper functions used throughout the application
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// String utilities

// TruncateString truncates a string to the specified length with ellipsis
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	if maxLength <= 3 {
		return "..."
	}
	return s[:maxLength-3] + "..."
}

// PadString pads a string to the specified length
func PadString(s string, length int, padChar rune, leftPad bool) string {
	if len(s) >= length {
		return s
	}

	padding := strings.Repeat(string(padChar), length-len(s))
	if leftPad {
		return padding + s
	}
	return s + padding
}

// CenterString centers a string within the specified width
func CenterString(s string, width int) string {
	if len(s) >= width {
		return s
	}

	leftPad := (width - len(s)) / 2
	rightPad := width - len(s) - leftPad

	return strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad)
}

// WrapText wraps text to fit within the specified width
func WrapText(text string, width int) []string {
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

		// If adding this word would exceed the width, start a new line
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

// SlugifyString converts a string to a URL-friendly slug
func SlugifyString(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")

	// Remove leading and trailing hyphens
	s = strings.Trim(s, "-")

	return s
}

// CamelCase converts a string to camelCase
func CamelCase(s string) string {
	if s == "" {
		return s
	}

	// Split by spaces, hyphens, underscores
	reg := regexp.MustCompile(`[-_\s]+`)
	words := reg.Split(s, -1)

	var result strings.Builder
	for i, word := range words {
		if word == "" {
			continue
		}

		if i == 0 {
			result.WriteString(strings.ToLower(word))
		} else {
			result.WriteString(strings.Title(strings.ToLower(word)))
		}
	}

	return result.String()
}

// PascalCase converts a string to PascalCase
func PascalCase(s string) string {
	camel := CamelCase(s)
	if camel == "" {
		return camel
	}

	// Capitalize first letter
	runes := []rune(camel)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// SnakeCase converts a string to snake_case
func SnakeCase(s string) string {
	// Insert underscores before uppercase letters
	reg := regexp.MustCompile(`([a-z])([A-Z])`)
	s = reg.ReplaceAllString(s, `${1}_${2}`)

	// Replace spaces and hyphens with underscores
	reg = regexp.MustCompile(`[-\s]+`)
	s = reg.ReplaceAllString(s, "_")

	// Convert to lowercase
	return strings.ToLower(s)
}

// KebabCase converts a string to kebab-case
func KebabCase(s string) string {
	return strings.ReplaceAll(SnakeCase(s), "_", "-")
}

// Slice utilities

// ContainsInt checks if an int slice contains a specific item
func ContainsInt(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// ContainsString checks if a string slice contains a specific string (case-insensitive option)
func ContainsString(slice []string, item string, caseInsensitive bool) bool {
	for _, v := range slice {
		if caseInsensitive {
			if strings.EqualFold(v, item) {
				return true
			}
		} else {
			if v == item {
				return true
			}
		}
	}
	return false
}

// UniqueStrings removes duplicates from a string slice
func UniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// UniqueInts removes duplicates from an int slice
func UniqueInts(slice []int) []int {
	seen := make(map[int]bool)
	var result []int

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// FilterStrings filters a string slice based on a predicate function
func FilterStrings(slice []string, predicate func(string) bool) []string {
	var result []string
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// FilterInts filters an int slice based on a predicate function
func FilterInts(slice []int, predicate func(int) bool) []int {
	var result []int
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// MapStrings transforms a string slice using a mapper function
func MapStrings(slice []string, mapper func(string) string) []string {
	result := make([]string, len(slice))
	for i, item := range slice {
		result[i] = mapper(item)
	}
	return result
}

// MapInts transforms an int slice using a mapper function
func MapInts(slice []int, mapper func(int) int) []int {
	result := make([]int, len(slice))
	for i, item := range slice {
		result[i] = mapper(item)
	}
	return result
}

// ReduceStrings reduces a string slice to a single value using a reducer function
func ReduceStrings(slice []string, initial string, reducer func(string, string) string) string {
	result := initial
	for _, item := range slice {
		result = reducer(result, item)
	}
	return result
}

// ReduceInts reduces an int slice to a single value using a reducer function
func ReduceInts(slice []int, initial int, reducer func(int, int) int) int {
	result := initial
	for _, item := range slice {
		result = reducer(result, item)
	}
	return result
}

// ChunkStrings splits a string slice into chunks of specified size
func ChunkStrings(slice []string, chunkSize int) [][]string {
	if chunkSize <= 0 {
		return [][]string{slice}
	}

	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// ChunkInts splits an int slice into chunks of specified size
func ChunkInts(slice []int, chunkSize int) [][]int {
	if chunkSize <= 0 {
		return [][]int{slice}
	}

	var chunks [][]int
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// ReverseStrings reverses a string slice in place
func ReverseStrings(slice []string) {
	for i := 0; i < len(slice)/2; i++ {
		j := len(slice) - 1 - i
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// ReverseInts reverses an int slice in place
func ReverseInts(slice []int) {
	for i := 0; i < len(slice)/2; i++ {
		j := len(slice) - 1 - i
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// ShuffleStrings shuffles a string slice randomly
func ShuffleStrings(slice []string) {
	for i := range slice {
		j := i + int(GenerateRandomInt(int64(len(slice)-i)))
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// ShuffleInts shuffles an int slice randomly
func ShuffleInts(slice []int) {
	for i := range slice {
		j := i + int(GenerateRandomInt(int64(len(slice)-i)))
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Map utilities

// MergeStringMaps merges multiple string maps into one (later maps override earlier ones)
func MergeStringMaps(maps ...map[string]string) map[string]string {
	result := make(map[string]string)

	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}

	return result
}

// MergeInterfaceMaps merges multiple interface{} maps into one
func MergeInterfaceMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}

	return result
}

// StringMapKeys returns all keys from a string map
func StringMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// StringMapValues returns all values from a string map
func StringMapValues(m map[string]string) []string {
	values := make([]string, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// InvertStringMap inverts a string map (keys become values, values become keys)
func InvertStringMap(m map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[v] = k
	}
	return result
}

// File and Path utilities

// EnsureDir ensures a directory exists, creating it if necessary
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	return !os.IsNotExist(err) && info.IsDir()
}

// GetFileSize returns the size of a file in bytes
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// GetHomeDirPath returns a path relative to the user's home directory
func GetHomeDirPath(relativePath string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, relativePath), nil
}

// ExpandPath expands ~ to the user's home directory
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

// GetRelativePath returns the relative path from base to target
func GetRelativePath(base, target string) (string, error) {
	return filepath.Rel(base, target)
}

// Time utilities

// TimeAgo returns a human-readable time duration ago
func TimeAgo(t time.Time) string {
	now := time.Now()
	duration := now.Sub(t)

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case duration < 30*24*time.Hour:
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	case duration < 365*24*time.Hour:
		months := int(duration.Hours() / (24 * 30))
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	default:
		years := int(duration.Hours() / (24 * 365))
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}

// FormatDuration formats a duration in a human-readable way
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
	return fmt.Sprintf("%.1fd", d.Hours()/24)
}

// ParseDurationExtended parses extended duration formats
func ParseDurationExtended(s string) (time.Duration, error) {
	// Try standard time.ParseDuration first
	if d, err := time.ParseDuration(s); err == nil {
		return d, nil
	}

	// Handle additional formats
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*(days?|d|weeks?|w|months?|mo|years?|y)$`)
	matches := re.FindStringSubmatch(strings.ToLower(strings.TrimSpace(s)))

	if len(matches) != 3 {
		return 0, fmt.Errorf("invalid duration format: %s", s)
	}

	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, err
	}

	unit := matches[2]

	switch {
	case strings.HasPrefix(unit, "d"):
		return time.Duration(value * float64(24*time.Hour)), nil
	case strings.HasPrefix(unit, "w"):
		return time.Duration(value * float64(7*24*time.Hour)), nil
	case strings.HasPrefix(unit, "mo"):
		return time.Duration(value * float64(30*24*time.Hour)), nil
	case strings.HasPrefix(unit, "y"):
		return time.Duration(value * float64(365*24*time.Hour)), nil
	default:
		return 0, fmt.Errorf("unknown duration unit: %s", unit)
	}
}

// Number utilities

// FormatBytes formats byte count in human-readable format
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// ParseBytes parses human-readable byte format
func ParseBytes(s string) (int64, error) {
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*(B|KB|MB|GB|TB|PB)?$`)
	matches := re.FindStringSubmatch(strings.ToUpper(strings.TrimSpace(s)))

	if len(matches) != 3 {
		return 0, fmt.Errorf("invalid byte format: %s", s)
	}

	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, err
	}

	unit := matches[2]
	if unit == "" {
		unit = "B"
	}

	multipliers := map[string]int64{
		"B":  1,
		"KB": 1024,
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
		"TB": 1024 * 1024 * 1024 * 1024,
		"PB": 1024 * 1024 * 1024 * 1024 * 1024,
	}

	multiplier, exists := multipliers[unit]
	if !exists {
		return 0, fmt.Errorf("unknown byte unit: %s", unit)
	}

	return int64(value * float64(multiplier)), nil
}

// RoundToDecimalPlaces rounds a float to specified decimal places
func RoundToDecimalPlaces(f float64, decimals int) float64 {
	shift := math.Pow(10, float64(decimals))
	return math.Round(f*shift) / shift
}

// ClampInt clamps an integer between min and max values
func ClampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampFloat clamps a float between min and max values
func ClampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Random utilities

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[GenerateRandomInt(int64(len(charset)))]
	}
	return string(b)
}

// GenerateRandomInt generates a random integer between 0 and max-1
func GenerateRandomInt(max int64) int64 {
	if max <= 0 {
		return 0
	}

	// Use crypto/rand for better randomness
	bytes := make([]byte, 8)
	rand.Read(bytes)

	// Convert to int64
	var result int64
	for i, b := range bytes {
		result |= int64(b) << (8 * i)
	}

	// Make positive and constrain to max
	if result < 0 {
		result = -result
	}
	return result % max
}

// GenerateUUID generates a simple UUID-like string
func GenerateUUID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// JSON utilities

// CompactJSON compacts JSON by removing unnecessary whitespace
func CompactJSON(jsonStr string) (string, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return "", err
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// JSONToMap converts JSON string to map
func JSONToMap(jsonStr string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	return result, err
}

// MapToJSON converts map to JSON string
func MapToJSON(m map[string]interface{}) (string, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Reflection utilities

// GetTypeName returns the type name of a value
func GetTypeName(v interface{}) string {
	if v == nil {
		return "nil"
	}
	return reflect.TypeOf(v).String()
}

// IsZeroValue checks if a value is the zero value for its type
func IsZeroValue(v interface{}) bool {
	if v == nil {
		return true
	}
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}

// CopyStruct performs a deep copy of a struct using JSON marshaling
func CopyStruct(src, dst interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

// System utilities

// GetOSInfo returns information about the operating system
func GetOSInfo() map[string]string {
	return map[string]string{
		"os":   runtime.GOOS,
		"arch": runtime.GOARCH,
	}
}

// GetMemoryUsage returns current memory usage statistics
func GetMemoryUsage() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"alloc_mb":       float64(m.Alloc) / 1024 / 1024,
		"total_alloc_mb": float64(m.TotalAlloc) / 1024 / 1024,
		"sys_mb":         float64(m.Sys) / 1024 / 1024,
		"num_gc":         m.NumGC,
		"goroutines":     runtime.NumGoroutine(),
	}
}

// GetStackTrace returns the current stack trace
func GetStackTrace(skip int) []string {
	var traces []string
	for i := skip; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		traces = append(traces, fmt.Sprintf("%s:%d", filepath.Base(file), line))
	}
	return traces
}

// Sorting utilities

// SortStringSlice sorts a string slice with various options
func SortStringSlice(slice []string, ascending bool, caseInsensitive bool) {
	sort.Slice(slice, func(i, j int) bool {
		a, b := slice[i], slice[j]

		if caseInsensitive {
			a = strings.ToLower(a)
			b = strings.ToLower(b)
		}

		if ascending {
			return a < b
		}
		return a > b
	})
}

// SortStringMapByKeys sorts a string map by its keys and returns sorted key-value pairs
func SortStringMapByKeys(m map[string]string) []struct {
	Key   string
	Value string
} {
	keys := StringMapKeys(m)
	sort.Strings(keys)

	var result []struct {
		Key   string
		Value string
	}
	for _, key := range keys {
		result = append(result, struct {
			Key   string
			Value string
		}{Key: key, Value: m[key]})
	}

	return result
}

// Utility for error handling

// IgnoreError executes a function and ignores any error
func IgnoreError(fn func() error) {
	_ = fn()
}

// MustString panics if the error is not nil, otherwise returns the string value
func MustString(value string, err error) string {
	if err != nil {
		panic(err)
	}
	return value
}

// MustInt panics if the error is not nil, otherwise returns the int value
func MustInt(value int, err error) int {
	if err != nil {
		panic(err)
	}
	return value
}

// DefaultString returns the default value if the original is empty
func DefaultString(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// DefaultInt returns the default value if the original is zero
func DefaultInt(value, defaultValue int) int {
	if value == 0 {
		return defaultValue
	}
	return value
}

// Retry executes a function with retry logic
func Retry(attempts int, delay time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
		if i < attempts-1 {
			time.Sleep(delay)
		}
	}
	return fmt.Errorf("failed after %d attempts: %w", attempts, err)
}

// Debounce creates a debounced version of a function
func Debounce(delay time.Duration, fn func()) func() {
	var timer *time.Timer
	return func() {
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(delay, fn)
	}
}

// Throttle creates a throttled version of a function
func Throttle(interval time.Duration, fn func()) func() {
	var lastCall time.Time
	return func() {
		if time.Since(lastCall) >= interval {
			fn()
			lastCall = time.Now()
		}
	}
}
