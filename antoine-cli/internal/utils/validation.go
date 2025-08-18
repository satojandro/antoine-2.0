// Package utils provides validation utilities for Antoine CLI
// This file implements comprehensive validation for inputs, configurations, and data
package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ValidationError represents a validation error with context
type ValidationError struct {
	Field   string `json:"field"`
	Value   string `json:"value"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (ve ValidationError) Error() string {
	if ve.Field != "" {
		return fmt.Sprintf("validation error for field '%s': %s", ve.Field, ve.Message)
	}
	return fmt.Sprintf("validation error: %s", ve.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "no validation errors"
	}
	if len(ve) == 1 {
		return ve[0].Error()
	}

	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Error())
	}
	return fmt.Sprintf("multiple validation errors: %s", strings.Join(messages, "; "))
}

// HasErrors returns true if there are validation errors
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

// Add adds a validation error
func (ve *ValidationErrors) Add(field, value, message, code string) {
	*ve = append(*ve, ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
		Code:    code,
	})
}

// Validator provides validation methods
type Validator struct {
	errors ValidationErrors
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{
		errors: make(ValidationErrors, 0),
	}
}

// Errors returns all validation errors
func (v *Validator) Errors() ValidationErrors {
	return v.errors
}

// HasErrors returns true if there are validation errors
func (v *Validator) HasErrors() bool {
	return v.errors.HasErrors()
}

// Clear clears all validation errors
func (v *Validator) Clear() {
	v.errors = make(ValidationErrors, 0)
}

// Add adds a validation error
func (v *Validator) Add(field, value, message, code string) {
	v.errors.Add(field, value, message, code)
}

// Required validates that a field is not empty
func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.Add(field, value, "field is required", "required")
	}
	return v
}

// MinLength validates minimum string length
func (v *Validator) MinLength(field, value string, min int) *Validator {
	if len(value) < min {
		v.Add(field, value, fmt.Sprintf("must be at least %d characters", min), "min_length")
	}
	return v
}

// MaxLength validates maximum string length
func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if len(value) > max {
		v.Add(field, value, fmt.Sprintf("must be no more than %d characters", max), "max_length")
	}
	return v
}

// Length validates exact string length
func (v *Validator) Length(field, value string, length int) *Validator {
	if len(value) != length {
		v.Add(field, value, fmt.Sprintf("must be exactly %d characters", length), "exact_length")
	}
	return v
}

// Email validates email format
func (v *Validator) Email(field, value string) *Validator {
	if value == "" {
		return v
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		v.Add(field, value, "must be a valid email address", "invalid_email")
	}
	return v
}

// URL validates URL format
func (v *Validator) URL(field, value string) *Validator {
	if value == "" {
		return v
	}

	parsedURL, err := url.Parse(value)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		v.Add(field, value, "must be a valid URL", "invalid_url")
	}
	return v
}

// GitHubURL validates GitHub repository URL
func (v *Validator) GitHubURL(field, value string) *Validator {
	if value == "" {
		return v
	}

	githubRegex := regexp.MustCompile(`^https://github\.com/[a-zA-Z0-9\-_.]+/[a-zA-Z0-9\-_.]+/?$`)
	if !githubRegex.MatchString(value) {
		v.Add(field, value, "must be a valid GitHub repository URL", "invalid_github_url")
	}
	return v
}

// MCPEndpoint validates MCP endpoint format
func (v *Validator) MCPEndpoint(field, value string) *Validator {
	if value == "" {
		return v
	}

	mcpRegex := regexp.MustCompile(`^mcp://[a-zA-Z0-9\-_.]+:\d+$`)
	if !mcpRegex.MatchString(value) {
		v.Add(field, value, "must be a valid MCP endpoint (mcp://host:port)", "invalid_mcp_endpoint")
	}
	return v
}

// Integer validates integer format and range
func (v *Validator) Integer(field, value string, min, max int) *Validator {
	if value == "" {
		return v
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		v.Add(field, value, "must be a valid integer", "invalid_integer")
		return v
	}

	if intValue < min {
		v.Add(field, value, fmt.Sprintf("must be at least %d", min), "min_value")
	}

	if intValue > max {
		v.Add(field, value, fmt.Sprintf("must be no more than %d", max), "max_value")
	}

	return v
}

// Float validates float format and range
func (v *Validator) Float(field, value string, min, max float64) *Validator {
	if value == "" {
		return v
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		v.Add(field, value, "must be a valid number", "invalid_float")
		return v
	}

	if floatValue < min {
		v.Add(field, value, fmt.Sprintf("must be at least %.2f", min), "min_value")
	}

	if floatValue > max {
		v.Add(field, value, fmt.Sprintf("must be no more than %.2f", max), "max_value")
	}

	return v
}

// Duration validates duration format
func (v *Validator) Duration(field, value string) *Validator {
	if value == "" {
		return v
	}

	_, err := time.ParseDuration(value)
	if err != nil {
		v.Add(field, value, "must be a valid duration (e.g., 5m, 1h, 30s)", "invalid_duration")
	}
	return v
}

// OneOf validates that value is one of allowed values
func (v *Validator) OneOf(field, value string, allowed []string) *Validator {
	if value == "" {
		return v
	}

	for _, allowedValue := range allowed {
		if value == allowedValue {
			return v
		}
	}

	v.Add(field, value, fmt.Sprintf("must be one of: %s", strings.Join(allowed, ", ")), "invalid_choice")
	return v
}

// Regex validates against a regular expression
func (v *Validator) Regex(field, value, pattern, message string) *Validator {
	if value == "" {
		return v
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		v.Add(field, value, "invalid validation pattern", "invalid_pattern")
		return v
	}

	if !regex.MatchString(value) {
		v.Add(field, value, message, "pattern_mismatch")
	}

	return v
}

// APIKey validates API key format
func (v *Validator) APIKey(field, value string) *Validator {
	if value == "" {
		return v
	}

	// Basic API key validation (adjust patterns as needed)
	patterns := map[string]*regexp.Regexp{
		"openai":    regexp.MustCompile(`^sk-[a-zA-Z0-9]{48}$`),
		"anthropic": regexp.MustCompile(`^sk-ant-[a-zA-Z0-9\-]{95}$`),
		"github":    regexp.MustCompile(`^ghp_[a-zA-Z0-9]{36}$`),
		"generic":   regexp.MustCompile(`^[a-zA-Z0-9\-_]{20,}$`),
	}

	// Try to match any known pattern
	for _, pattern := range patterns {
		if pattern.MatchString(value) {
			return v
		}
	}

	v.Add(field, value, "must be a valid API key format", "invalid_api_key")
	return v
}

// Technology validates technology/programming language names
func (v *Validator) Technology(field, value string) *Validator {
	if value == "" {
		return v
	}

	// List of valid technologies (expand as needed)
	validTechnologies := []string{
		"javascript", "typescript", "python", "go", "rust", "java", "c++", "c#",
		"php", "ruby", "swift", "kotlin", "dart", "scala", "clojure", "haskell",
		"react", "vue", "angular", "svelte", "nextjs", "nuxtjs", "express",
		"fastapi", "django", "flask", "spring", "gin", "echo", "fiber",
		"postgresql", "mysql", "mongodb", "redis", "elasticsearch",
		"docker", "kubernetes", "aws", "gcp", "azure", "terraform",
		"blockchain", "solidity", "web3", "ethereum", "bitcoin", "defi",
		"ai", "ml", "tensorflow", "pytorch", "opencv", "nlp",
	}

	lowerValue := strings.ToLower(value)
	for _, tech := range validTechnologies {
		if lowerValue == tech {
			return v
		}
	}

	v.Add(field, value, "must be a valid technology name", "invalid_technology")
	return v
}

// DateRange validates date range format
func (v *Validator) DateRange(field, value string) *Validator {
	if value == "" {
		return v
	}

	// Support formats like "2023-01-01,2023-12-31" or "2023-01-01..2023-12-31"
	separators := []string{",", "..", " to ", "-"}

	var parts []string
	for _, sep := range separators {
		if strings.Contains(value, sep) {
			parts = strings.Split(value, sep)
			break
		}
	}

	if len(parts) != 2 {
		v.Add(field, value, "must be a valid date range (start,end or start..end)", "invalid_date_range")
		return v
	}

	// Validate each date
	for _, part := range parts {
		part = strings.TrimSpace(part)
		_, err := time.Parse("2006-01-02", part)
		if err != nil {
			v.Add(field, value, fmt.Sprintf("invalid date in range: %s", part), "invalid_date")
			return v
		}
	}

	return v
}

// FileSize validates file size format (e.g., "10MB", "1GB")
func (v *Validator) FileSize(field, value string) *Validator {
	if value == "" {
		return v
	}

	fileSizeRegex := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*(B|KB|MB|GB|TB)$`)
	if !fileSizeRegex.MatchString(strings.ToUpper(value)) {
		v.Add(field, value, "must be a valid file size (e.g., 10MB, 1GB)", "invalid_file_size")
	}

	return v
}

// ConfigValidation provides specialized validation for Antoine configuration
type ConfigValidation struct {
	*Validator
}

// NewConfigValidation creates a new config validator
func NewConfigValidation() *ConfigValidation {
	return &ConfigValidation{
		Validator: NewValidator(),
	}
}

// ValidateLogLevel validates logging level
func (cv *ConfigValidation) ValidateLogLevel(level string) *ConfigValidation {
	validLevels := []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}
	cv.OneOf("log_level", level, validLevels)
	return cv
}

// ValidateLogFormat validates logging format
func (cv *ConfigValidation) ValidateLogFormat(format string) *ConfigValidation {
	validFormats := []string{"text", "json"}
	cv.OneOf("log_format", format, validFormats)
	return cv
}

// ValidateLogOutput validates logging output
func (cv *ConfigValidation) ValidateLogOutput(output string) *ConfigValidation {
	validOutputs := []string{"stdout", "stderr", "file"}
	cv.OneOf("log_output", output, validOutputs)
	return cv
}

// ValidateCacheType validates cache type
func (cv *ConfigValidation) ValidateCacheType(cacheType string) *ConfigValidation {
	validTypes := []string{"memory", "disk", "hybrid"}
	cv.OneOf("cache_type", cacheType, validTypes)
	return cv
}

// ValidateUITheme validates UI theme
func (cv *ConfigValidation) ValidateUITheme(theme string) *ConfigValidation {
	validThemes := []string{"dark", "light", "minimal"}
	cv.OneOf("ui_theme", theme, validThemes)
	return cv
}

// ValidateMentorModel validates AI model name
func (cv *ConfigValidation) ValidateMentorModel(model string) *ConfigValidation {
	validModels := []string{"gpt-3.5-turbo", "gpt-4", "gpt-4-turbo", "claude-3-sonnet", "claude-3-opus"}
	cv.OneOf("mentor_model", model, validModels)
	return cv
}

// ValidateUpdateChannel validates update channel
func (cv *ConfigValidation) ValidateUpdateChannel(channel string) *ConfigValidation {
	validChannels := []string{"stable", "beta", "alpha"}
	cv.OneOf("update_channel", channel, validChannels)
	return cv
}

// SearchValidation provides specialized validation for search parameters
type SearchValidation struct {
	*Validator
}

// NewSearchValidation creates a new search validator
func NewSearchValidation() *SearchValidation {
	return &SearchValidation{
		Validator: NewValidator(),
	}
}

// ValidateSearchQuery validates search query
func (sv *SearchValidation) ValidateSearchQuery(query string) *SearchValidation {
	sv.Required("query", query).
		MinLength("query", query, 2).
		MaxLength("query", query, 200)
	return sv
}

// ValidateSortBy validates sort field
func (sv *SearchValidation) ValidateSortBy(sortBy, context string) *SearchValidation {
	var validFields []string

	switch context {
	case "hackathons":
		validFields = []string{"start_date", "end_date", "prize_pool", "popularity", "name"}
	case "projects":
		validFields = []string{"popularity", "recent", "stars", "forks", "name", "created_at"}
	default:
		validFields = []string{"name", "created_at", "updated_at"}
	}

	sv.OneOf("sort_by", sortBy, validFields)
	return sv
}

// ValidateSortOrder validates sort order
func (sv *SearchValidation) ValidateSortOrder(sortOrder string) *SearchValidation {
	validOrders := []string{"asc", "desc"}
	sv.OneOf("sort_order", sortOrder, validOrders)
	return sv
}

// ValidateLimit validates result limit
func (sv *SearchValidation) ValidateLimit(limit string) *SearchValidation {
	sv.Integer("limit", limit, 1, 1000)
	return sv
}

// ValidateTechnologies validates a list of technologies
func (sv *SearchValidation) ValidateTechnologies(technologies []string) *SearchValidation {
	for _, tech := range technologies {
		sv.Technology("technology", tech)
	}
	return sv
}

// ValidatePrizeRange validates prize range
func (sv *SearchValidation) ValidatePrizeRange(minPrize, maxPrize string) *SearchValidation {
	if minPrize != "" {
		sv.Integer("min_prize", minPrize, 0, 10000000)
	}
	if maxPrize != "" {
		sv.Integer("max_prize", maxPrize, 0, 10000000)
	}
	return sv
}

// Additional validation types

// AnalysisValidation provides specialized validation for analysis parameters
type AnalysisValidation struct {
	*Validator
}

// NewAnalysisValidation creates a new analysis validator
func NewAnalysisValidation() *AnalysisValidation {
	return &AnalysisValidation{
		Validator: NewValidator(),
	}
}

// ValidateRepositoryURL validates repository URL
func (av *AnalysisValidation) ValidateRepositoryURL(repoURL string) *AnalysisValidation {
	av.Required("repository_url", repoURL)
	av.GitHubURL("repository_url", repoURL)
	return av
}

// ValidateAnalysisDepth validates analysis depth
func (av *AnalysisValidation) ValidateAnalysisDepth(depth string) *AnalysisValidation {
	validDepths := []string{"shallow", "normal", "deep"}
	av.OneOf("analysis_depth", depth, validDepths)
	return av
}

// ValidateAnalysisTimeout validates analysis timeout
func (av *AnalysisValidation) ValidateAnalysisTimeout(timeout string) *AnalysisValidation {
	av.Duration("analysis_timeout", timeout)
	return av
}

// ValidateAnalysisFocus validates analysis focus areas
func (av *AnalysisValidation) ValidateAnalysisFocus(focus []string) *AnalysisValidation {
	validFocus := []string{"security", "performance", "architecture", "dependencies", "documentation"}
	for _, f := range focus {
		av.OneOf("analysis_focus", f, validFocus)
	}
	return av
}

// MentorValidation provides specialized validation for mentor parameters
type MentorValidation struct {
	*Validator
}

// NewMentorValidation creates a new mentor validator
func NewMentorValidation() *MentorValidation {
	return &MentorValidation{
		Validator: NewValidator(),
	}
}

// ValidatePersonality validates mentor personality
func (mv *MentorValidation) ValidatePersonality(personality string) *MentorValidation {
	validPersonalities := []string{"helpful", "casual", "professional", "technical"}
	mv.OneOf("personality", personality, validPersonalities)
	return mv
}

// ValidateExpertiseLevel validates expertise level
func (mv *MentorValidation) ValidateExpertiseLevel(level string) *MentorValidation {
	validLevels := []string{"beginner", "intermediate", "advanced"}
	mv.OneOf("expertise_level", level, validLevels)
	return mv
}

// ValidateSessionTimeout validates session timeout
func (mv *MentorValidation) ValidateSessionTimeout(timeout string) *MentorValidation {
	mv.Duration("session_timeout", timeout)
	return mv
}

// Utility validation functions

// ValidateIP validates IP address format
func ValidateIP(field, ip string) ValidationError {
	if ip == "" {
		return ValidationError{}
	}

	ipRegex := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	if !ipRegex.MatchString(ip) {
		return ValidationError{
			Field:   field,
			Value:   ip,
			Message: "must be a valid IP address",
			Code:    "invalid_ip",
		}
	}

	// Validate each octet is between 0-255
	parts := strings.Split(ip, ".")
	for _, part := range parts {
		num, _ := strconv.Atoi(part)
		if num > 255 {
			return ValidationError{
				Field:   field,
				Value:   ip,
				Message: "IP address octets must be between 0 and 255",
				Code:    "invalid_ip_range",
			}
		}
	}

	return ValidationError{}
}

// ValidateHostname validates hostname format
func ValidateHostname(field, hostname string) ValidationError {
	if hostname == "" {
		return ValidationError{}
	}

	hostnameRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	if !hostnameRegex.MatchString(hostname) {
		return ValidationError{
			Field:   field,
			Value:   hostname,
			Message: "must be a valid hostname",
			Code:    "invalid_hostname",
		}
	}

	if len(hostname) > 253 {
		return ValidationError{
			Field:   field,
			Value:   hostname,
			Message: "hostname cannot exceed 253 characters",
			Code:    "hostname_too_long",
		}
	}

	return ValidationError{}
}

// ValidateSemVer validates semantic version format
func ValidateSemVer(field, version string) ValidationError {
	if version == "" {
		return ValidationError{}
	}

	semVerRegex := regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	if !semVerRegex.MatchString(version) {
		return ValidationError{
			Field:   field,
			Value:   version,
			Message: "must be a valid semantic version (e.g., 1.0.0, 2.1.3-beta.1)",
			Code:    "invalid_semver",
		}
	}

	return ValidationError{}
}

// ValidateHexColor validates hexadecimal color format
func ValidateHexColor(field, color string) ValidationError {
	if color == "" {
		return ValidationError{}
	}

	hexColorRegex := regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)
	if !hexColorRegex.MatchString(color) {
		return ValidationError{
			Field:   field,
			Value:   color,
			Message: "must be a valid hex color (e.g., #FF0000, #f00)",
			Code:    "invalid_hex_color",
		}
	}

	return ValidationError{}
}

// Quick validation helper functions

// IsValidEmail checks if email is valid
func IsValidEmail(email string) bool {
	validator := NewValidator()
	validator.Email("email", email)
	return !validator.HasErrors()
}

// IsValidURL checks if URL is valid
func IsValidURL(url string) bool {
	validator := NewValidator()
	validator.URL("url", url)
	return !validator.HasErrors()
}

// IsValidGitHubURL checks if GitHub URL is valid
func IsValidGitHubURL(url string) bool {
	validator := NewValidator()
	validator.GitHubURL("url", url)
	return !validator.HasErrors()
}

// IsValidMCPEndpoint checks if MCP endpoint is valid
func IsValidMCPEndpoint(endpoint string) bool {
	validator := NewValidator()
	validator.MCPEndpoint("endpoint", endpoint)
	return !validator.HasErrors()
}

// IsValidDuration checks if duration string is valid
func IsValidDuration(duration string) bool {
	validator := NewValidator()
	validator.Duration("duration", duration)
	return !validator.HasErrors()
}
