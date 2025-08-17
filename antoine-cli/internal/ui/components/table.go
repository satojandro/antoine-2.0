// Package components provides reusable UI components for Antoine CLI
// This file implements interactive tables for displaying data
package components

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"antoine-cli/internal/ui/styles"
	"antoine-cli/internal/utils"
	"github.com/charmbracelet/lipgloss"
)

// TableColumn defines a table column
type TableColumn struct {
	Key       string                   // Data key
	Title     string                   // Display title
	Width     int                      // Column width (0 = auto)
	Align     lipgloss.Position        // Text alignment
	Sortable  bool                     // Whether column is sortable
	Formatter func(interface{}) string // Custom formatter
}

// TableRow represents a table row with data and metadata
type TableRow struct {
	ID       string                 // Unique row identifier
	Data     map[string]interface{} // Row data
	Selected bool                   // Whether row is selected
	Style    lipgloss.Style         // Custom row styling
}

// TableConfig configures the table component
type TableConfig struct {
	Title         string
	Columns       []TableColumn
	Data          []TableRow
	Width         int
	Height        int
	ShowHeader    bool
	ShowBorders   bool
	ShowFooter    bool
	Sortable      bool
	Selectable    bool
	MultiSelect   bool
	ZebraStriping bool
	HeaderStyle   lipgloss.Style
	RowStyle      lipgloss.Style
	SelectedStyle lipgloss.Style
	EmptyMessage  string
	FooterMessage string
}

// Table represents an interactive table component
type Table struct {
	config         TableConfig
	sortColumn     string
	sortDescending bool
	selectedRows   map[string]bool
	currentPage    int
	pageSize       int
	filteredData   []TableRow
}

// NewTable creates a new table component
func NewTable(config TableConfig) *Table {
	// Set defaults
	if config.Width == 0 {
		config.Width = 120
	}
	if config.Height == 0 {
		config.Height = 20
	}
	if config.EmptyMessage == "" {
		config.EmptyMessage = "No data available"
	}
	if config.HeaderStyle.String() == "" {
		config.HeaderStyle = styles.TableHeaderStyle
	}
	if config.RowStyle.String() == "" {
		config.RowStyle = styles.TableCellStyle
	}
	if config.SelectedStyle.String() == "" {
		config.SelectedStyle = styles.TableSelectedStyle
	}

	table := &Table{
		config:       config,
		selectedRows: make(map[string]bool),
		pageSize:     config.Height - 3, // Account for header and borders
		filteredData: config.Data,
	}

	// Auto-calculate column widths if not specified
	table.calculateColumnWidths()

	return table
}

// Render renders the table component
func (t *Table) Render() string {
	if len(t.filteredData) == 0 {
		return t.renderEmpty()
	}

	var sections []string

	// Title
	if t.config.Title != "" {
		title := styles.H3Style.Render(t.config.Title)
		sections = append(sections, title)
	}

	// Header
	if t.config.ShowHeader {
		header := t.renderHeader()
		sections = append(sections, header)
	}

	// Data rows
	rows := t.renderRows()
	sections = append(sections, rows)

	// Footer
	if t.config.ShowFooter {
		footer := t.renderFooter()
		sections = append(sections, footer)
	}

	return strings.Join(sections, "\n")
}

// renderEmpty renders the empty state
func (t *Table) renderEmpty() string {
	var sections []string

	if t.config.Title != "" {
		title := styles.H3Style.Render(t.config.Title)
		sections = append(sections, title)
	}

	// Empty message
	emptyStyle := lipgloss.NewStyle().
		Width(t.config.Width).
		Align(lipgloss.Center).
		Padding(2).
		Foreground(styles.Gray)

	emptyMsg := emptyStyle.Render(t.config.EmptyMessage)
	sections = append(sections, emptyMsg)

	return strings.Join(sections, "\n")
}

// renderHeader renders the table header
func (t *Table) renderHeader() string {
	var headerCells []string

	for _, col := range t.config.Columns {
		title := col.Title

		// Add sort indicator
		if t.config.Sortable && col.Sortable && t.sortColumn == col.Key {
			if t.sortDescending {
				title += " ↓"
			} else {
				title += " ↑"
			}
		}

		// Style the header cell
		cellStyle := t.config.HeaderStyle.Copy().
			Width(col.Width).
			Align(col.Align)

		cell := cellStyle.Render(title)
		headerCells = append(headerCells, cell)
	}

	header := lipgloss.JoinHorizontal(lipgloss.Top, headerCells...)

	if t.config.ShowBorders {
		borderStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(styles.Gold).
			Width(t.config.Width)
		header = borderStyle.Render(header)
	}

	return header
}

// renderRows renders the table data rows
func (t *Table) renderRows() string {
	var rowStrings []string

	// Calculate pagination
	startIdx := t.currentPage * t.pageSize
	endIdx := startIdx + t.pageSize
	if endIdx > len(t.filteredData) {
		endIdx = len(t.filteredData)
	}

	pageData := t.filteredData[startIdx:endIdx]

	for i, row := range pageData {
		rowStr := t.renderRow(row, i)
		rowStrings = append(rowStrings, rowStr)
	}

	return strings.Join(rowStrings, "\n")
}

// renderRow renders a single table row
func (t *Table) renderRow(row TableRow, index int) string {
	var cells []string

	// Determine row style
	rowStyle := t.config.RowStyle
	if row.Selected || t.selectedRows[row.ID] {
		rowStyle = t.config.SelectedStyle
	} else if t.config.ZebraStriping && index%2 == 1 {
		rowStyle = rowStyle.Copy().Background(styles.Navy)
	}

	// Custom row style override
	if row.Style.String() != "" {
		rowStyle = row.Style
	}

	for _, col := range t.config.Columns {
		value := row.Data[col.Key]

		// Format the value
		var displayValue string
		if col.Formatter != nil {
			displayValue = col.Formatter(value)
		} else {
			displayValue = t.formatValue(value)
		}

		// Truncate if necessary
		if len(displayValue) > col.Width-2 {
			displayValue = utils.TruncateString(displayValue, col.Width-2)
		}

		// Style the cell
		cellStyle := rowStyle.Copy().
			Width(col.Width).
			Align(col.Align)

		cell := cellStyle.Render(displayValue)
		cells = append(cells, cell)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, cells...)
}

// renderFooter renders the table footer with pagination info
func (t *Table) renderFooter() string {
	var footerParts []string

	// Pagination info
	totalRows := len(t.filteredData)
	startIdx := t.currentPage*t.pageSize + 1
	endIdx := (t.currentPage + 1) * t.pageSize
	if endIdx > totalRows {
		endIdx = totalRows
	}

	if totalRows > 0 {
		paginationInfo := fmt.Sprintf("Showing %d-%d of %d", startIdx, endIdx, totalRows)
		footerParts = append(footerParts, paginationInfo)
	}

	// Sort info
	if t.sortColumn != "" {
		sortInfo := fmt.Sprintf("Sorted by %s", t.sortColumn)
		if t.sortDescending {
			sortInfo += " (desc)"
		} else {
			sortInfo += " (asc)"
		}
		footerParts = append(footerParts, sortInfo)
	}

	// Selection info
	selectedCount := len(t.selectedRows)
	if selectedCount > 0 {
		selectionInfo := fmt.Sprintf("%d selected", selectedCount)
		footerParts = append(footerParts, selectionInfo)
	}

	// Custom footer message
	if t.config.FooterMessage != "" {
		footerParts = append(footerParts, t.config.FooterMessage)
	}

	footerText := strings.Join(footerParts, " • ")

	footerStyle := lipgloss.NewStyle().
		Width(t.config.Width).
		Align(lipgloss.Center).
		Foreground(styles.Gray).
		Padding(0, 1)

	if t.config.ShowBorders {
		footerStyle = footerStyle.
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(styles.Gold)
	}

	return footerStyle.Render(footerText)
}

// formatValue formats a value for display
func (t *Table) formatValue(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case int, int32, int64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%.2f", v)
	case bool:
		if v {
			return "✓"
		}
		return "✗"
	case time.Time:
		return v.Format("2006-01-02 15:04")
	case time.Duration:
		return utils.FormatDuration(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// calculateColumnWidths auto-calculates column widths
func (t *Table) calculateColumnWidths() {
	availableWidth := t.config.Width
	autoColumns := 0

	// First pass: account for fixed-width columns
	for i, col := range t.config.Columns {
		if col.Width > 0 {
			availableWidth -= col.Width
		} else {
			autoColumns++
		}
	}

	// Second pass: distribute remaining width
	if autoColumns > 0 {
		autoWidth := availableWidth / autoColumns
		for i, col := range t.config.Columns {
			if col.Width == 0 {
				t.config.Columns[i].Width = autoWidth
			}
		}
	}
}

// Sorting methods

// SortBy sorts the table by the specified column
func (t *Table) SortBy(columnKey string, descending bool) {
	t.sortColumn = columnKey
	t.sortDescending = descending

	sort.Slice(t.filteredData, func(i, j int) bool {
		valueI := t.filteredData[i].Data[columnKey]
		valueJ := t.filteredData[j].Data[columnKey]

		result := t.compareValues(valueI, valueJ)
		if descending {
			return result > 0
		}
		return result < 0
	})
}

// ToggleSort toggles sort order for a column
func (t *Table) ToggleSort(columnKey string) {
	if t.sortColumn == columnKey {
		t.sortDescending = !t.sortDescending
	} else {
		t.sortColumn = columnKey
		t.sortDescending = false
	}
	t.SortBy(columnKey, t.sortDescending)
}

// compareValues compares two values for sorting
func (t *Table) compareValues(a, b interface{}) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	switch va := a.(type) {
	case string:
		if vb, ok := b.(string); ok {
			return strings.Compare(va, vb)
		}
	case int:
		if vb, ok := b.(int); ok {
			return va - vb
		}
	case int64:
		if vb, ok := b.(int64); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case float64:
		if vb, ok := b.(float64); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case time.Time:
		if vb, ok := b.(time.Time); ok {
			if va.Before(vb) {
				return -1
			} else if va.After(vb) {
				return 1
			}
			return 0
		}
	}

	// Fallback to string comparison
	return strings.Compare(fmt.Sprintf("%v", a), fmt.Sprintf("%v", b))
}

// Selection methods

// SelectRow selects a row by ID
func (t *Table) SelectRow(rowID string) {
	if !t.config.Selectable {
		return
	}

	if !t.config.MultiSelect {
		// Clear previous selections
		t.selectedRows = make(map[string]bool)
	}

	t.selectedRows[rowID] = true
}

// DeselectRow deselects a row by ID
func (t *Table) DeselectRow(rowID string) {
	delete(t.selectedRows, rowID)
}

// ToggleRow toggles row selection
func (t *Table) ToggleRow(rowID string) {
	if t.selectedRows[rowID] {
		t.DeselectRow(rowID)
	} else {
		t.SelectRow(rowID)
	}
}

// SelectAll selects all visible rows
func (t *Table) SelectAll() {
	if !t.config.Selectable || !t.config.MultiSelect {
		return
	}

	for _, row := range t.filteredData {
		t.selectedRows[row.ID] = true
	}
}

// DeselectAll deselects all rows
func (t *Table) DeselectAll() {
	t.selectedRows = make(map[string]bool)
}

// GetSelectedRows returns the IDs of selected rows
func (t *Table) GetSelectedRows() []string {
	var selected []string
	for rowID := range t.selectedRows {
		selected = append(selected, rowID)
	}
	return selected
}

// Pagination methods

// NextPage moves to the next page
func (t *Table) NextPage() {
	maxPage := (len(t.filteredData) - 1) / t.pageSize
	if t.currentPage < maxPage {
		t.currentPage++
	}
}

// PrevPage moves to the previous page
func (t *Table) PrevPage() {
	if t.currentPage > 0 {
		t.currentPage--
	}
}

// SetPage sets the current page
func (t *Table) SetPage(page int) {
	maxPage := (len(t.filteredData) - 1) / t.pageSize
	if page >= 0 && page <= maxPage {
		t.currentPage = page
	}
}

// GetCurrentPage returns the current page number
func (t *Table) GetCurrentPage() int {
	return t.currentPage
}

// GetTotalPages returns the total number of pages
func (t *Table) GetTotalPages() int {
	if len(t.filteredData) == 0 {
		return 0
	}
	return (len(t.filteredData)-1)/t.pageSize + 1
}

// Filtering methods

// Filter filters the table data based on a predicate
func (t *Table) Filter(predicate func(TableRow) bool) {
	var filtered []TableRow
	for _, row := range t.config.Data {
		if predicate(row) {
			filtered = append(filtered, row)
		}
	}
	t.filteredData = filtered
	t.currentPage = 0 // Reset to first page
}

// FilterByColumn filters by a specific column value
func (t *Table) FilterByColumn(columnKey string, value interface{}) {
	t.Filter(func(row TableRow) bool {
		return row.Data[columnKey] == value
	})
}

// FilterByText filters by text search across all columns
func (t *Table) FilterByText(searchText string) {
	searchText = strings.ToLower(searchText)

	t.Filter(func(row TableRow) bool {
		for _, value := range row.Data {
			text := strings.ToLower(t.formatValue(value))
			if strings.Contains(text, searchText) {
				return true
			}
		}
		return false
	})
}

// ClearFilter removes all filters
func (t *Table) ClearFilter() {
	t.filteredData = t.config.Data
	t.currentPage = 0
}

// Data manipulation methods

// AddRow adds a new row to the table
func (t *Table) AddRow(row TableRow) {
	t.config.Data = append(t.config.Data, row)
	t.filteredData = append(t.filteredData, row)
}

// UpdateRow updates an existing row
func (t *Table) UpdateRow(rowID string, data map[string]interface{}) {
	for i, row := range t.config.Data {
		if row.ID == rowID {
			for key, value := range data {
				t.config.Data[i].Data[key] = value
			}
			break
		}
	}

	for i, row := range t.filteredData {
		if row.ID == rowID {
			for key, value := range data {
				t.filteredData[i].Data[key] = value
			}
			break
		}
	}
}

// RemoveRow removes a row by ID
func (t *Table) RemoveRow(rowID string) {
	// Remove from main data
	for i, row := range t.config.Data {
		if row.ID == rowID {
			t.config.Data = append(t.config.Data[:i], t.config.Data[i+1:]...)
			break
		}
	}

	// Remove from filtered data
	for i, row := range t.filteredData {
		if row.ID == rowID {
			t.filteredData = append(t.filteredData[:i], t.filteredData[i+1:]...)
			break
		}
	}

	// Remove from selection
	delete(t.selectedRows, rowID)
}

// SetData replaces all table data
func (t *Table) SetData(data []TableRow) {
	t.config.Data = data
	t.filteredData = data
	t.selectedRows = make(map[string]bool)
	t.currentPage = 0
}

// Export methods

// ToCSV exports table data to CSV format
func (t *Table) ToCSV() string {
	var lines []string

	// Header
	var headers []string
	for _, col := range t.config.Columns {
		headers = append(headers, col.Title)
	}
	lines = append(lines, strings.Join(headers, ","))

	// Data rows
	for _, row := range t.filteredData {
		var values []string
		for _, col := range t.config.Columns {
			value := t.formatValue(row.Data[col.Key])
			// Escape quotes and wrap in quotes if contains comma
			if strings.Contains(value, ",") || strings.Contains(value, "\"") {
				value = "\"" + strings.ReplaceAll(value, "\"", "\"\"") + "\""
			}
			values = append(values, value)
		}
		lines = append(lines, strings.Join(values, ","))
	}

	return strings.Join(lines, "\n")
}

// ToJSON exports table data to JSON format
func (t *Table) ToJSON() (string, error) {
	var data []map[string]interface{}

	for _, row := range t.filteredData {
		rowData := make(map[string]interface{})
		for _, col := range t.config.Columns {
			rowData[col.Key] = row.Data[col.Key]
		}
		data = append(data, rowData)
	}

	return utils.PrettyPrintJSON(data)
}

// Utility methods

// GetRowCount returns the total number of rows
func (t *Table) GetRowCount() int {
	return len(t.filteredData)
}

// GetRow returns a row by ID
func (t *Table) GetRow(rowID string) *TableRow {
	for _, row := range t.filteredData {
		if row.ID == rowID {
			return &row
		}
	}
	return nil
}

// SetFooterMessage sets a custom footer message
func (t *Table) SetFooterMessage(message string) {
	t.config.FooterMessage = message
}

// SetTitle sets the table title
func (t *Table) SetTitle(title string) {
	t.config.Title = title
}

// Resize resizes the table
func (t *Table) Resize(width, height int) {
	t.config.Width = width
	t.config.Height = height
	t.pageSize = height - 3
	t.calculateColumnWidths()
}

// TableBuilder provides a fluent interface for building tables
type TableBuilder struct {
	config TableConfig
}

// NewTableBuilder creates a new table builder
func NewTableBuilder() *TableBuilder {
	return &TableBuilder{
		config: TableConfig{
			ShowHeader:    true,
			ShowBorders:   true,
			ShowFooter:    true,
			Sortable:      true,
			Selectable:    true,
			ZebraStriping: true,
		},
	}
}

// Title sets the table title
func (tb *TableBuilder) Title(title string) *TableBuilder {
	tb.config.Title = title
	return tb
}

// Columns sets the table columns
func (tb *TableBuilder) Columns(columns ...TableColumn) *TableBuilder {
	tb.config.Columns = columns
	return tb
}

// Data sets the table data
func (tb *TableBuilder) Data(data []TableRow) *TableBuilder {
	tb.config.Data = data
	return tb
}

// Width sets the table width
func (tb *TableBuilder) Width(width int) *TableBuilder {
	tb.config.Width = width
	return tb
}

// Height sets the table height
func (tb *TableBuilder) Height(height int) *TableBuilder {
	tb.config.Height = height
	return tb
}

// ShowHeader controls header visibility
func (tb *TableBuilder) ShowHeader(show bool) *TableBuilder {
	tb.config.ShowHeader = show
	return tb
}

// ShowBorders controls border visibility
func (tb *TableBuilder) ShowBorders(show bool) *TableBuilder {
	tb.config.ShowBorders = show
	return tb
}

// Sortable controls sorting capability
func (tb *TableBuilder) Sortable(sortable bool) *TableBuilder {
	tb.config.Sortable = sortable
	return tb
}

// Selectable controls row selection
func (tb *TableBuilder) Selectable(selectable bool) *TableBuilder {
	tb.config.Selectable = selectable
	return tb
}

// MultiSelect enables multi-row selection
func (tb *TableBuilder) MultiSelect(multiSelect bool) *TableBuilder {
	tb.config.MultiSelect = multiSelect
	return tb
}

// ZebraStriping controls alternating row colors
func (tb *TableBuilder) ZebraStriping(zebra bool) *TableBuilder {
	tb.config.ZebraStriping = zebra
	return tb
}

// EmptyMessage sets the empty state message
func (tb *TableBuilder) EmptyMessage(message string) *TableBuilder {
	tb.config.EmptyMessage = message
	return tb
}

// Build creates the table component
func (tb *TableBuilder) Build() *Table {
	return NewTable(tb.config)
}
