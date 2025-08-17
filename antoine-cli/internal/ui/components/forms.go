// Package components provides reusable UI components for Antoine CLI
// This file implements interactive forms and input components
package components

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"antoine-cli/internal/ui/styles"
	"antoine-cli/internal/utils"
	"github.com/charmbracelet/lipgloss"
)

// InputType defines different input field types
type InputType string

const (
	InputTypeText     InputType = "text"     // Regular text input
	InputTypePassword InputType = "password" // Hidden password input
	InputTypeEmail    InputType = "email"    // Email validation
	InputTypeNumber   InputType = "number"   // Numeric input
	InputTypeURL      InputType = "url"      // URL validation
	InputTypeDate     InputType = "date"     // Date input
	InputTypeTime     InputType = "time"     // Time input
	InputTypeSelect   InputType = "select"   // Dropdown selection
	InputTypeMulti    InputType = "multi"    // Multi-select
	InputTypeTextarea InputType = "textarea" // Multi-line text
	InputTypeToggle   InputType = "toggle"   // Boolean toggle
	InputTypeSlider   InputType = "slider"   // Range slider
)

// ValidationRule defines a validation rule for form fields
type ValidationRule struct {
	Name      string
	Validator func(string) error
	Message   string
}

// FormField represents a form input field
type FormField struct {
	ID          string
	Label       string
	Type        InputType
	Value       string
	Placeholder string
	Required    bool
	Disabled    bool
	Hidden      bool
	Options     []string // For select/multi inputs
	Validations []ValidationRule
	MinLength   int
	MaxLength   int
	Pattern     string
	Style       lipgloss.Style
	HelpText    string
	Error       string
}

// FormConfig configures the form component
type FormConfig struct {
	Title       string
	Description string
	Fields      []FormField
	Width       int
	ShowHelp    bool
	ShowErrors  bool
	CompactMode bool
	SubmitText  string
	CancelText  string
}

// Form represents an interactive form component
type Form struct {
	config       FormConfig
	values       map[string]string
	errors       map[string]string
	focusedField int
	submitted    bool
	cancelled    bool
	style        lipgloss.Style
}

// NewForm creates a new form component
func NewForm(config FormConfig) *Form {
	// Set defaults
	if config.Width == 0 {
		config.Width = 80
	}
	if config.SubmitText == "" {
		config.SubmitText = "Submit"
	}
	if config.CancelText == "" {
		config.CancelText = "Cancel"
	}

	// Initialize values map
	values := make(map[string]string)
	for _, field := range config.Fields {
		values[field.ID] = field.Value
	}

	return &Form{
		config: config,
		values: values,
		errors: make(map[string]string),
		style: lipgloss.NewStyle().
			Width(config.Width).
			Padding(1, 2),
	}
}

// Render renders the form component
func (f *Form) Render() string {
	var sections []string

	// Title
	if f.config.Title != "" {
		title := styles.H2Style.Render(f.config.Title)
		sections = append(sections, title)
	}

	// Description
	if f.config.Description != "" {
		desc := styles.BodyStyle.Render(f.config.Description)
		sections = append(sections, desc)
	}

	// Form fields
	for i, field := range f.config.Fields {
		if field.Hidden {
			continue
		}

		fieldHTML := f.renderField(field, i == f.focusedField)
		sections = append(sections, fieldHTML)

		// Add spacing between fields (except in compact mode)
		if !f.config.CompactMode && i < len(f.config.Fields)-1 {
			sections = append(sections, "")
		}
	}

	// Form actions
	actions := f.renderActions()
	if actions != "" {
		sections = append(sections, "")
		sections = append(sections, actions)
	}

	content := strings.Join(sections, "\n")
	return f.style.Render(content)
}

// renderField renders a single form field
func (f *Form) renderField(field FormField, focused bool) string {
	var parts []string

	// Label
	label := f.renderLabel(field, focused)
	parts = append(parts, label)

	// Input
	input := f.renderInput(field, focused)
	parts = append(parts, input)

	// Error message
	if f.config.ShowErrors && f.errors[field.ID] != "" {
		errorMsg := styles.ErrorStyle.Render("✗ " + f.errors[field.ID])
		parts = append(parts, errorMsg)
	}

	// Help text
	if f.config.ShowHelp && field.HelpText != "" && focused {
		helpMsg := styles.BodySecondaryStyle.Render("💡 " + field.HelpText)
		parts = append(parts, helpMsg)
	}

	return strings.Join(parts, "\n")
}

// renderLabel renders the field label
func (f *Form) renderLabel(field FormField, focused bool) string {
	label := field.Label

	if field.Required {
		label += " *"
	}

	labelStyle := styles.LabelStyle
	if focused {
		labelStyle = labelStyle.Copy().Foreground(styles.Gold)
	}

	return labelStyle.Render(label)
}

// renderInput renders the input field based on its type
func (f *Form) renderInput(field FormField, focused bool) string {
	value := f.values[field.ID]

	switch field.Type {
	case InputTypeText, InputTypeEmail, InputTypeURL:
		return f.renderTextInput(field, value, focused)
	case InputTypePassword:
		return f.renderPasswordInput(field, value, focused)
	case InputTypeNumber:
		return f.renderNumberInput(field, value, focused)
	case InputTypeDate:
		return f.renderDateInput(field, value, focused)
	case InputTypeTime:
		return f.renderTimeInput(field, value, focused)
	case InputTypeSelect:
		return f.renderSelectInput(field, value, focused)
	case InputTypeMulti:
		return f.renderMultiInput(field, value, focused)
	case InputTypeTextarea:
		return f.renderTextareaInput(field, value, focused)
	case InputTypeToggle:
		return f.renderToggleInput(field, value, focused)
	case InputTypeSlider:
		return f.renderSliderInput(field, value, focused)
	default:
		return f.renderTextInput(field, value, focused)
	}
}

// renderTextInput renders a text input field
func (f *Form) renderTextInput(field FormField, value string, focused bool) string {
	displayValue := value
	if displayValue == "" && field.Placeholder != "" {
		displayValue = field.Placeholder
		displayValue = styles.BodySecondaryStyle.Render(displayValue)
	}

	inputStyle := styles.InputStyle
	if focused {
		inputStyle = styles.InputFocusStyle
	}
	if field.Disabled {
		inputStyle = inputStyle.Copy().Foreground(styles.DarkGray)
	}

	// Add cursor if focused
	if focused && !field.Disabled {
		displayValue += "│"
	}

	width := f.config.Width - 8 // Account for padding and borders
	inputStyle = inputStyle.Width(width)

	return inputStyle.Render(displayValue)
}

// renderPasswordInput renders a password input field
func (f *Form) renderPasswordInput(field FormField, value string, focused bool) string {
	// Mask the password
	maskedValue := strings.Repeat("●", len(value))

	// Create a temporary field with masked value
	maskedField := field
	maskedField.Placeholder = "Enter password..."

	return f.renderTextInput(maskedField, maskedValue, focused)
}

// renderNumberInput renders a number input field
func (f *Form) renderNumberInput(field FormField, value string, focused bool) string {
	// Add number-specific placeholder if none provided
	if field.Placeholder == "" {
		field.Placeholder = "Enter a number..."
	}

	return f.renderTextInput(field, value, focused)
}

// renderDateInput renders a date input field
func (f *Form) renderDateInput(field FormField, value string, focused bool) string {
	// Add date-specific placeholder if none provided
	if field.Placeholder == "" {
		field.Placeholder = "YYYY-MM-DD"
	}

	return f.renderTextInput(field, value, focused)
}

// renderTimeInput renders a time input field
func (f *Form) renderTimeInput(field FormField, value string, focused bool) string {
	// Add time-specific placeholder if none provided
	if field.Placeholder == "" {
		field.Placeholder = "HH:MM"
	}

	return f.renderTextInput(field, value, focused)
}

// renderSelectInput renders a select dropdown field
func (f *Form) renderSelectInput(field FormField, value string, focused bool) string {
	var displayValue string

	if value == "" {
		displayValue = field.Placeholder
		if displayValue == "" {
			displayValue = "Select an option..."
		}
		displayValue = styles.BodySecondaryStyle.Render(displayValue)
	} else {
		displayValue = value
	}

	// Add dropdown indicator
	displayValue += " ▼"

	inputStyle := styles.InputStyle
	if focused {
		inputStyle = styles.InputFocusStyle
	}

	width := f.config.Width - 8
	inputStyle = inputStyle.Width(width)

	result := inputStyle.Render(displayValue)

	// Show options if focused
	if focused && len(field.Options) > 0 {
		optionsStr := f.renderSelectOptions(field.Options, value)
		result += "\n" + optionsStr
	}

	return result
}

// renderSelectOptions renders the options for a select field
func (f *Form) renderSelectOptions(options []string, currentValue string) string {
	var optionLines []string

	for _, option := range options {
		prefix := "  "
		optionStyle := styles.BodyStyle

		if option == currentValue {
			prefix = "▶ "
			optionStyle = styles.HighlightStyle
		}

		optionLine := optionStyle.Render(prefix + option)
		optionLines = append(optionLines, optionLine)
	}

	optionsContainer := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.Cyan).
		Padding(0, 1).
		Width(f.config.Width - 12)

	return optionsContainer.Render(strings.Join(optionLines, "\n"))
}

// renderMultiInput renders a multi-select input field
func (f *Form) renderMultiInput(field FormField, value string, focused bool) string {
	selectedValues := strings.Split(value, ",")
	if value == "" {
		selectedValues = []string{}
	}

	displayValue := fmt.Sprintf("%d selected", len(selectedValues))
	if len(selectedValues) == 0 {
		displayValue = field.Placeholder
		if displayValue == "" {
			displayValue = "Select options..."
		}
		displayValue = styles.BodySecondaryStyle.Render(displayValue)
	}

	// Add multi-select indicator
	displayValue += " ▼"

	inputStyle := styles.InputStyle
	if focused {
		inputStyle = styles.InputFocusStyle
	}

	width := f.config.Width - 8
	inputStyle = inputStyle.Width(width)

	result := inputStyle.Render(displayValue)

	// Show options if focused
	if focused && len(field.Options) > 0 {
		optionsStr := f.renderMultiOptions(field.Options, selectedValues)
		result += "\n" + optionsStr
	}

	return result
}

// renderMultiOptions renders the options for a multi-select field
func (f *Form) renderMultiOptions(options []string, selectedValues []string) string {
	var optionLines []string

	for _, option := range options {
		prefix := "☐ "
		optionStyle := styles.BodyStyle

		// Check if this option is selected
		for _, selected := range selectedValues {
			if strings.TrimSpace(selected) == option {
				prefix = "☑ "
				optionStyle = styles.AccentStyle
				break
			}
		}

		optionLine := optionStyle.Render(prefix + option)
		optionLines = append(optionLines, optionLine)
	}

	optionsContainer := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.Cyan).
		Padding(0, 1).
		Width(f.config.Width - 12)

	return optionsContainer.Render(strings.Join(optionLines, "\n"))
}

// renderTextareaInput renders a textarea input field
func (f *Form) renderTextareaInput(field FormField, value string, focused bool) string {
	lines := strings.Split(value, "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		lines = []string{field.Placeholder}
		if field.Placeholder == "" {
			lines = []string{"Enter text..."}
		}
	}

	// Limit to reasonable number of lines
	maxLines := 5
	if len(lines) > maxLines {
		lines = lines[:maxLines]
		lines = append(lines, "...")
	}

	content := strings.Join(lines, "\n")

	inputStyle := styles.InputStyle
	if focused {
		inputStyle = styles.InputFocusStyle
	}

	width := f.config.Width - 8
	height := len(lines) + 1
	inputStyle = inputStyle.Width(width).Height(height)

	return inputStyle.Render(content)
}

// renderToggleInput renders a toggle/boolean input field
func (f *Form) renderToggleInput(field FormField, value string, focused bool) string {
	isOn := value == "true" || value == "1" || strings.ToLower(value) == "yes"

	var toggle string
	var toggleStyle lipgloss.Style

	if isOn {
		toggle = "☑ ON"
		toggleStyle = styles.SuccessStyle
	} else {
		toggle = "☐ OFF"
		toggleStyle = styles.BodySecondaryStyle
	}

	if focused {
		toggleStyle = toggleStyle.Copy().Bold(true)
	}

	return toggleStyle.Render(toggle)
}

// renderSliderInput renders a slider/range input field
func (f *Form) renderSliderInput(field FormField, value string, focused bool) string {
	// Parse current value
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		floatValue = 0
	}

	// Default range 0-100
	min, max := 0.0, 100.0

	// Calculate position
	percentage := (floatValue - min) / (max - min)
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 1 {
		percentage = 1
	}

	// Create slider visualization
	sliderWidth := 20
	filledWidth := int(percentage * float64(sliderWidth))

	filled := strings.Repeat("━", filledWidth)
	handle := "●"
	empty := strings.Repeat("─", sliderWidth-filledWidth)

	slider := filled + handle + empty

	sliderStyle := styles.AccentStyle
	if focused {
		sliderStyle = styles.HighlightStyle
	}

	valueText := fmt.Sprintf(" %.1f", floatValue)

	return sliderStyle.Render(slider) + styles.BodyStyle.Render(valueText)
}

// renderActions renders form action buttons
func (f *Form) renderActions() string {
	var buttons []string

	// Submit button
	submitStyle := styles.ButtonStyle
	if f.focusedField == len(f.config.Fields) {
		submitStyle = styles.ButtonHoverStyle
	}
	submitBtn := submitStyle.Render(f.config.SubmitText)
	buttons = append(buttons, submitBtn)

	// Cancel button
	cancelStyle := styles.ButtonStyle.Copy().
		Foreground(styles.Gray).
		Background(styles.Navy)
	if f.focusedField == len(f.config.Fields)+1 {
		cancelStyle = styles.ButtonHoverStyle.Copy().
			Foreground(styles.White).
			Background(styles.Red)
	}
	cancelBtn := cancelStyle.Render(f.config.CancelText)
	buttons = append(buttons, cancelBtn)

	buttonsContainer := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(f.config.Width - 4)

	return buttonsContainer.Render(strings.Join(buttons, "  "))
}

// Form interaction methods

// SetValue sets the value of a field
func (f *Form) SetValue(fieldID, value string) {
	f.values[fieldID] = value
	f.validateField(fieldID)
}

// GetValue gets the value of a field
func (f *Form) GetValue(fieldID string) string {
	return f.values[fieldID]
}

// GetAllValues returns all form values
func (f *Form) GetAllValues() map[string]string {
	result := make(map[string]string)
	for k, v := range f.values {
		result[k] = v
	}
	return result
}

// Focus moves focus to the next field
func (f *Form) Focus() {
	maxFocus := len(f.config.Fields) + 1 // +1 for submit button, +1 for cancel
	f.focusedField = (f.focusedField + 1) % (maxFocus + 1)

	// Skip hidden or disabled fields
	for f.isFieldSkippable(f.focusedField) && f.focusedField < len(f.config.Fields) {
		f.focusedField = (f.focusedField + 1) % (maxFocus + 1)
	}
}

// FocusPrev moves focus to the previous field
func (f *Form) FocusPrev() {
	maxFocus := len(f.config.Fields) + 1
	f.focusedField--
	if f.focusedField < 0 {
		f.focusedField = maxFocus
	}

	// Skip hidden or disabled fields
	for f.isFieldSkippable(f.focusedField) && f.focusedField >= 0 {
		f.focusedField--
		if f.focusedField < 0 {
			f.focusedField = maxFocus
		}
	}
}

// isFieldSkippable checks if a field should be skipped during navigation
func (f *Form) isFieldSkippable(fieldIndex int) bool {
	if fieldIndex >= len(f.config.Fields) {
		return false // Action buttons are never skippable
	}

	field := f.config.Fields[fieldIndex]
	return field.Hidden || field.Disabled
}

// HandleInput processes user input for the focused field
func (f *Form) HandleInput(input string) bool {
	if f.focusedField >= len(f.config.Fields) {
		// Handle action buttons
		if f.focusedField == len(f.config.Fields) {
			// Submit button
			return f.Submit()
		} else if f.focusedField == len(f.config.Fields)+1 {
			// Cancel button
			f.Cancel()
			return true
		}
		return false
	}

	field := f.config.Fields[f.focusedField]
	if field.Disabled {
		return false
	}

	switch field.Type {
	case InputTypeToggle:
		return f.handleToggleInput(field.ID, input)
	case InputTypeSelect:
		return f.handleSelectInput(field.ID, input)
	case InputTypeMulti:
		return f.handleMultiInput(field.ID, input)
	default:
		return f.handleTextInput(field.ID, input)
	}
}

// handleTextInput handles text input for text-based fields
func (f *Form) handleTextInput(fieldID, input string) bool {
	currentValue := f.values[fieldID]

	switch input {
	case "backspace":
		if len(currentValue) > 0 {
			f.SetValue(fieldID, currentValue[:len(currentValue)-1])
		}
	case "delete":
		f.SetValue(fieldID, "")
	default:
		// Regular character input
		if len(input) == 1 {
			newValue := currentValue + input
			field := f.getField(fieldID)
			if field != nil && field.MaxLength > 0 && len(newValue) > field.MaxLength {
				return false // Don't allow input beyond max length
			}
			f.SetValue(fieldID, newValue)
		}
	}

	return true
}

// handleToggleInput handles toggle input
func (f *Form) handleToggleInput(fieldID, input string) bool {
	if input == " " || input == "enter" {
		currentValue := f.values[fieldID]
		isOn := currentValue == "true" || currentValue == "1"
		f.SetValue(fieldID, fmt.Sprintf("%t", !isOn))
		return true
	}
	return false
}

// handleSelectInput handles select dropdown input
func (f *Form) handleSelectInput(fieldID, input string) bool {
	field := f.getField(fieldID)
	if field == nil || len(field.Options) == 0 {
		return false
	}

	switch input {
	case "up":
		currentIndex := f.findOptionIndex(fieldID, f.values[fieldID])
		if currentIndex > 0 {
			f.SetValue(fieldID, field.Options[currentIndex-1])
		}
		return true
	case "down":
		currentIndex := f.findOptionIndex(fieldID, f.values[fieldID])
		if currentIndex < len(field.Options)-1 {
			f.SetValue(fieldID, field.Options[currentIndex+1])
		}
		return true
	case "enter":
		// Confirm selection
		return true
	}

	return false
}

// handleMultiInput handles multi-select input
func (f *Form) handleMultiInput(fieldID, input string) bool {
	field := f.getField(fieldID)
	if field == nil || len(field.Options) == 0 {
		return false
	}

	// For simplicity, we'll just toggle the first option
	// In a real implementation, you'd track which option is highlighted
	if input == " " {
		currentValues := strings.Split(f.values[fieldID], ",")
		if f.values[fieldID] == "" {
			currentValues = []string{}
		}

		option := field.Options[0] // Simplified - would be highlighted option

		// Toggle option
		found := false
		var newValues []string
		for _, val := range currentValues {
			if strings.TrimSpace(val) == option {
				found = true
				// Skip this value (remove it)
			} else if strings.TrimSpace(val) != "" {
				newValues = append(newValues, strings.TrimSpace(val))
			}
		}

		if !found {
			newValues = append(newValues, option)
		}

		f.SetValue(fieldID, strings.Join(newValues, ","))
		return true
	}

	return false
}

// Validation methods

// validateField validates a single field
func (f *Form) validateField(fieldID string) {
	field := f.getField(fieldID)
	if field == nil {
		return
	}

	value := f.values[fieldID]

	// Clear previous error
	delete(f.errors, fieldID)

	// Required validation
	if field.Required && strings.TrimSpace(value) == "" {
		f.errors[fieldID] = "This field is required"
		return
	}

	// Skip other validations if field is empty and not required
	if strings.TrimSpace(value) == "" {
		return
	}

	// Length validations
	if field.MinLength > 0 && len(value) < field.MinLength {
		f.errors[fieldID] = fmt.Sprintf("Must be at least %d characters", field.MinLength)
		return
	}

	if field.MaxLength > 0 && len(value) > field.MaxLength {
		f.errors[fieldID] = fmt.Sprintf("Must be no more than %d characters", field.MaxLength)
		return
	}

	// Pattern validation
	if field.Pattern != "" {
		if matched, _ := regexp.MatchString(field.Pattern, value); !matched {
			f.errors[fieldID] = "Invalid format"
			return
		}
	}

	// Type-specific validations
	switch field.Type {
	case InputTypeEmail:
		if !f.isValidEmail(value) {
			f.errors[fieldID] = "Invalid email address"
		}
	case InputTypeURL:
		if !f.isValidURL(value) {
			f.errors[fieldID] = "Invalid URL"
		}
	case InputTypeNumber:
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			f.errors[fieldID] = "Must be a valid number"
		}
	case InputTypeDate:
		if !f.isValidDate(value) {
			f.errors[fieldID] = "Invalid date format (YYYY-MM-DD)"
		}
	case InputTypeTime:
		if !f.isValidTime(value) {
			f.errors[fieldID] = "Invalid time format (HH:MM)"
		}
	}

	// Custom validations
	for _, rule := range field.Validations {
		if err := rule.Validator(value); err != nil {
			f.errors[fieldID] = rule.Message
			if rule.Message == "" {
				f.errors[fieldID] = err.Error()
			}
			return
		}
	}
}

// validateAll validates all form fields
func (f *Form) validateAll() bool {
	for _, field := range f.config.Fields {
		f.validateField(field.ID)
	}
	return len(f.errors) == 0
}

// Validation helper methods
func (f *Form) isValidEmail(email string) bool {
	return utils.IsValidEmail(email)
}

func (f *Form) isValidURL(url string) bool {
	return utils.IsValidURL(url)
}

func (f *Form) isValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func (f *Form) isValidTime(timeStr string) bool {
	_, err := time.Parse("15:04", timeStr)
	return err == nil
}

// Utility methods

// getField returns a field by ID
func (f *Form) getField(fieldID string) *FormField {
	for i, field := range f.config.Fields {
		if field.ID == fieldID {
			return &f.config.Fields[i]
		}
	}
	return nil
}

// findOptionIndex finds the index of an option in a select field
func (f *Form) findOptionIndex(fieldID, value string) int {
	field := f.getField(fieldID)
	if field == nil {
		return -1
	}

	for i, option := range field.Options {
		if option == value {
			return i
		}
	}
	return -1
}

// Form lifecycle methods

// Submit submits the form if validation passes
func (f *Form) Submit() bool {
	if f.validateAll() {
		f.submitted = true
		return true
	}
	return false
}

// Cancel cancels the form
func (f *Form) Cancel() {
	f.cancelled = true
}

// IsSubmitted returns true if the form was submitted
func (f *Form) IsSubmitted() bool {
	return f.submitted
}

// IsCancelled returns true if the form was cancelled
func (f *Form) IsCancelled() bool {
	return f.cancelled
}

// HasErrors returns true if the form has validation errors
func (f *Form) HasErrors() bool {
	return len(f.errors) > 0
}

// GetErrors returns all validation errors
func (f *Form) GetErrors() map[string]string {
	result := make(map[string]string)
	for k, v := range f.errors {
		result[k] = v
	}
	return result
}

// Reset resets the form to its initial state
func (f *Form) Reset() {
	f.values = make(map[string]string)
	f.errors = make(map[string]string)
	f.focusedField = 0
	f.submitted = false
	f.cancelled = false

	// Restore default values
	for _, field := range f.config.Fields {
		f.values[field.ID] = field.Value
	}
}

// FormBuilder provides a fluent interface for building forms
type FormBuilder struct {
	config FormConfig
}

// NewFormBuilder creates a new form builder
func NewFormBuilder() *FormBuilder {
	return &FormBuilder{
		config: FormConfig{
			Width:      80,
			ShowHelp:   true,
			ShowErrors: true,
			SubmitText: "Submit",
			CancelText: "Cancel",
		},
	}
}

// Title sets the form title
func (fb *FormBuilder) Title(title string) *FormBuilder {
	fb.config.Title = title
	return fb
}

// Description sets the form description
func (fb *FormBuilder) Description(description string) *FormBuilder {
	fb.config.Description = description
	return fb
}

// Width sets the form width
func (fb *FormBuilder) Width(width int) *FormBuilder {
	fb.config.Width = width
	return fb
}

// AddField adds a field to the form
func (fb *FormBuilder) AddField(field FormField) *FormBuilder {
	fb.config.Fields = append(fb.config.Fields, field)
	return fb
}

// AddTextField adds a text field
func (fb *FormBuilder) AddTextField(id, label, placeholder string, required bool) *FormBuilder {
	field := FormField{
		ID:          id,
		Label:       label,
		Type:        InputTypeText,
		Placeholder: placeholder,
		Required:    required,
	}
	return fb.AddField(field)
}

// AddPasswordField adds a password field
func (fb *FormBuilder) AddPasswordField(id, label string, required bool) *FormBuilder {
	field := FormField{
		ID:       id,
		Label:    label,
		Type:     InputTypePassword,
		Required: required,
	}
	return fb.AddField(field)
}

// AddSelectField adds a select field
func (fb *FormBuilder) AddSelectField(id, label string, options []string, required bool) *FormBuilder {
	field := FormField{
		ID:       id,
		Label:    label,
		Type:     InputTypeSelect,
		Options:  options,
		Required: required,
	}
	return fb.AddField(field)
}

// AddToggleField adds a toggle field
func (fb *FormBuilder) AddToggleField(id, label string, defaultValue bool) *FormBuilder {
	value := "false"
	if defaultValue {
		value = "true"
	}

	field := FormField{
		ID:    id,
		Label: label,
		Type:  InputTypeToggle,
		Value: value,
	}
	return fb.AddField(field)
}

// CompactMode enables compact mode
func (fb *FormBuilder) CompactMode(compact bool) *FormBuilder {
	fb.config.CompactMode = compact
	return fb
}

// SubmitText sets the submit button text
func (fb *FormBuilder) SubmitText(text string) *FormBuilder {
	fb.config.SubmitText = text
	return fb
}

// CancelText sets the cancel button text
func (fb *FormBuilder) CancelText(text string) *FormBuilder {
	fb.config.CancelText = text
	return fb
}

// Build creates the form component
func (fb *FormBuilder) Build() *Form {
	return NewForm(fb.config)
}

// Predefined forms for common use cases

// SearchForm creates a search form
func SearchForm() *Form {
	return NewFormBuilder().
		Title("Search Configuration").
		AddTextField("query", "Search Query", "Enter search terms...", true).
		AddSelectField("type", "Search Type", []string{"hackathons", "projects", "all"}, true).
		AddSelectField("sort", "Sort By", []string{"relevance", "date", "popularity"}, false).
		AddToggleField("include_past", "Include Past Events", false).
		SubmitText("Search").
		Build()
}

// ConfigForm creates a configuration form
func ConfigForm() *Form {
	return NewFormBuilder().
		Title("Antoine Configuration").
		Description("Configure your Antoine CLI settings").
		AddSelectField("theme", "UI Theme", []string{"dark", "light", "minimal"}, true).
		AddSelectField("log_level", "Log Level", []string{"debug", "info", "warn", "error"}, true).
		AddToggleField("animations", "Enable Animations", true).
		AddToggleField("analytics", "Enable Analytics", true).
		AddTextField("api_key", "API Key", "Enter your API key...", false).
		SubmitText("Save Configuration").
		Build()
}

// CredentialsForm creates a credentials form
func CredentialsForm() *Form {
	return NewFormBuilder().
		Title("API Credentials").
		Description("Configure your API credentials for external services").
		AddTextField("openai_key", "OpenAI API Key", "sk-...", false).
		AddTextField("github_token", "GitHub Token", "ghp_...", false).
		AddTextField("exa_key", "Exa API Key", "Enter Exa API key...", false).
		SubmitText("Save Credentials").
		Build()
}

// AnalysisForm creates an analysis configuration form
func AnalysisForm() *Form {
	return NewFormBuilder().
		Title("Repository Analysis").
		AddTextField("repo_url", "Repository URL", "https://github.com/user/repo", true).
		AddSelectField("depth", "Analysis Depth", []string{"shallow", "medium", "deep", "comprehensive"}, true).
		AddField(FormField{
			ID:      "focus",
			Label:   "Focus Areas",
			Type:    InputTypeMulti,
			Options: []string{"architecture", "security", "performance", "documentation", "testing"},
		}).
		AddToggleField("include_deps", "Include Dependencies", true).
		AddToggleField("security_scan", "Security Scan", true).
		SubmitText("Start Analysis").
		Build()
}
