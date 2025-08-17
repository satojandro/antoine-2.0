package views

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"antoine-cli/internal/core"
	"antoine-cli/internal/models"
	"antoine-cli/pkg/ascii"
)

type SearchView struct {
	client *core.AntoineClient
}

type SearchOptions struct {
	Tech      string
	Location  string
	PrizeMin  string
	DateFrom  string
	DateTo    string
	Online    bool
	Hackathon string
	Category  string
	Sort      string
	Format    string
}

func NewSearchView(client *core.AntoineClient) *SearchView {
	return &SearchView{client: client}
}

type searchModel struct {
	searchInput textinput.Model
	spinner     spinner.Model
	table       table.Model
	results     interface{}
	loading     bool
	err         error
	searchType  string
	options     *SearchOptions
	width       int
	height      int
}

type searchCompleteMsg struct {
	results interface{}
	err     error
}

func (sv *SearchView) SearchHackathons(options *SearchOptions) {
	if options.Format == "json" || options.Format == "yaml" {
		sv.searchHackathonsNonInteractive(options)
		return
	}

	model := sv.createSearchModel("hackathons", options)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func (sv *SearchView) SearchProjects(options *SearchOptions) {
	if options.Format == "json" || options.Format == "yaml" {
		sv.searchProjectsNonInteractive(options)
		return
	}

	model := sv.createSearchModel("projects", options)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func (sv *SearchView) createSearchModel(searchType string, options *SearchOptions) searchModel {
	// Configurar input de búsqueda
	ti := textinput.New()
	ti.Placeholder = "Enter search terms..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	// Configurar spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(ascii.Gold)

	// Configurar tabla
	columns := []table.Column{}
	if searchType == "hackathons" {
		columns = []table.Column{
			{Title: "Name", Width: 30},
			{Title: "Start Date", Width: 12},
			{Title: "Prize", Width: 10},
			{Title: "Location", Width: 15},
			{Title: "Tech", Width: 20},
		}
	} else {
		columns = []table.Column{
			{Title: "Project", Width: 25},
			{Title: "Hackathon", Width: 20},
			{Title: "Tech", Width: 20},
			{Title: "Awards", Width: 15},
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s1 := table.DefaultStyles()
	s1.Header = s1.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(ascii.Gold).
		BorderBottom(true).
		Bold(false)
	s1.Selected = s1.Selected.
		Foreground(ascii.DarkBlue).
		Background(ascii.Gold).
		Bold(false)
	t.SetStyles(s1)

	return searchModel{
		searchInput: ti,
		spinner:     s,
		table:       t,
		searchType:  searchType,
		options:     options,
		loading:     false,
	}
}

func (m searchModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

func (m searchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if !m.loading {
				m.loading = true
				query := m.searchInput.Value()
				return m, tea.Batch(
					m.spinner.Tick,
					m.performSearch(query),
				)
			}
		case "esc":
			if m.loading {
				m.loading = false
				return m, nil
			}
		}

	case searchCompleteMsg:
		m.loading = false
		m.results = msg.results
		m.err = msg.err

		if msg.err == nil {
			m.updateTable()
		}

	case spinner.TickMsg:
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	if !m.loading {
		m.searchInput, cmd = m.searchInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	if len(m.table.Rows()) > 0 {
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m searchModel) View() string {
	var s strings.Builder

	// Header
	header := ascii.GetBanner(fmt.Sprintf("Search %s", strings.Title(m.searchType)), ascii.EmojiSearch, m.width)
	s.WriteString(header)
	s.WriteString("\n\n")

	if m.loading {
		s.WriteString(fmt.Sprintf("%s Searching... This may take a moment\n\n", m.spinner.View()))
	} else {
		// Input de búsqueda
		s.WriteString("🔍 Search Query:\n")
		s.WriteString(m.searchInput.View())
		s.WriteString("\n\n")

		if m.err != nil {
			errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f7768e"))
			s.WriteString(errorStyle.Render(fmt.Sprintf("❌ Error: %v", m.err)))
			s.WriteString("\n\n")
		}

		// Tabla de resultados
		if len(m.table.Rows()) > 0 {
			s.WriteString("📊 Results:\n")
			s.WriteString(m.table.View())
			s.WriteString("\n\n")
		}
	}

	// Footer con ayuda
	helpStyle := lipgloss.NewStyle().Foreground(ascii.Cyan)
	if m.loading {
		s.WriteString(helpStyle.Render("Press 'esc' to cancel • Press 'ctrl+c' to quit"))
	} else {
		s.WriteString(helpStyle.Render("Press 'enter' to search • Press 'q' to quit"))
	}

	return s.String()
}

func (m searchModel) performSearch(query string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		var results interface{}
		var err error

		// Convertir opciones a filtros
		filters := make(map[string]interface{})
		if m.options.Tech != "" {
			filters["technologies"] = strings.Split(m.options.Tech, ",")
		}
		if m.options.Location != "" {
			filters["location"] = m.options.Location
		}
		if m.options.Online {
			filters["online"] = true
		}

		if m.searchType == "hackathons" {
			// Buscar hackathons
			hackathons, searchErr := sv.client.SearchHackathons(ctx, query, filters)
			results = hackathons
			err = searchErr
		} else {
			// Buscar proyectos
			if m.options.Hackathon != "" {
				filters["hackathon"] = m.options.Hackathon
			}
			if m.options.Category != "" {
				filters["category"] = strings.Split(m.options.Category, ",")
			}

			projects, searchErr := sv.client.SearchProjects(ctx, query, filters)
			results = projects
			err = searchErr
		}

		return searchCompleteMsg{results: results, err: err}
	}
}

func (m *searchModel) updateTable() {
	var rows []table.Row

	if m.searchType == "hackathons" {
		if hackathons, ok := m.results.([]*models.Hackathon); ok {
			for _, h := range hackathons {
				prize := fmt.Sprintf("$%d", h.PrizePool.Total)
				if h.PrizePool.Total == 0 {
					prize = "N/A"
				}

				tech := strings.Join(h.Technologies[:min(3, len(h.Technologies))], ", ")
				if len(h.Technologies) > 3 {
					tech += "..."
				}

				rows = append(rows, table.Row{
					h.Name,
					h.StartDate.Format("2006-01-02"),
					prize,
					h.Location.City,
					tech,
				})
			}
		}
	} else {
		if projects, ok := m.results.([]*models.Project); ok {
			for _, p := range projects {
				tech := strings.Join(p.Technologies[:min(3, len(p.Technologies))], ", ")
				if len(p.Technologies) > 3 {
					tech += "..."
				}

				awards := "None"
				if len(p.Awards) > 0 {
					awards = p.Awards[0].Position
				}

				rows = append(rows, table.Row{
					p.Name,
					p.HackathonName,
					tech,
					awards,
				})
			}
		}
	}

	m.table.SetRows(rows)
}

func (sv *SearchView) searchHackathonsNonInteractive(options *SearchOptions) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filters := make(map[string]interface{})
	if options.Tech != "" {
		filters["technologies"] = strings.Split(options.Tech, ",")
	}

	hackathons, err := sv.client.SearchHackathons(ctx, "", filters)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Output según formato
	switch options.Format {
	case "json":
		// Implementar output JSON
		fmt.Printf("Found %d hackathons\n", len(hackathons))
	case "yaml":
		// Implementar output YAML
		fmt.Printf("Found %d hackathons\n", len(hackathons))
	default:
		// Tabla simple
		for _, h := range hackathons {
			fmt.Printf("%-30s %-12s %-10s %-15s\n",
				h.Name,
				h.StartDate.Format("2006-01-02"),
				fmt.Sprintf("$%d", h.PrizePool.Total),
				h.Location.City)
		}
	}
}

func (sv *SearchView) searchProjectsNonInteractive(options *SearchOptions) {
	// Implementación similar para proyectos
	fmt.Println("Project search not implemented in non-interactive mode")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
