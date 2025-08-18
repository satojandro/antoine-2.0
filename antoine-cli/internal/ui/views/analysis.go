package views

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"antoine-cli/internal/core"
	"antoine-cli/internal/models"
	"antoine-cli/pkg/ascii"
)

type AnalysisView struct {
	client *core.AntoineClient
}

type AnalysisOptions struct {
	RepoURL             string
	Depth               string
	IncludeDependencies bool
	GenerateReport      bool
	Focus               string
	Tech                []string
	Timeframe           string
	Metrics             bool
	Format              string
}

func NewAnalysisView(client *core.AntoineClient) *AnalysisView {
	return &AnalysisView{client: client}
}

type analysisModel struct {
	spinner  spinner.Model
	progress progress.Model
	repoURL  string
	options  *AnalysisOptions
	result   *models.AnalysisResult
	loading  bool
	err      error
	step     string
	width    int
	height   int
	client   *core.AntoineClient
}

type analysisCompleteMsg struct {
	result *models.AnalysisResult
	err    error
}

type analysisProgressMsg struct {
	step     string
	progress float64
}

func (av *AnalysisView) AnalyzeRepository(options *AnalysisOptions) {
	if options.Format == "json" || options.Format == "yaml" {
		av.analyzeRepositoryNonInteractive(options)
		return
	}

	model := av.createAnalysisModel(options)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func (av *AnalysisView) AnalyzeTrends(options *AnalysisOptions) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	trends, err := av.client.GetTrends(ctx, options.Tech, options.Timeframe)
	if err != nil {
		fmt.Printf("Error analyzing trends: %v\n", err)
		return
	}

	fmt.Printf("📈 Technology Trends (%s)\n\n", options.Timeframe)
	fmt.Printf("Results: %+v\n", trends)
}

func (av *AnalysisView) createAnalysisModel(options *AnalysisOptions) analysisModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(ascii.Gold)

	prog := progress.New(progress.WithDefaultGradient())

	return analysisModel{
		spinner:  s,
		progress: prog,
		repoURL:  options.RepoURL,
		options:  options,
		loading:  true,
		step:     "Initializing analysis...",
		client:   av.client,
	}
}

func (m analysisModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.performAnalysis(),
	)
}

func (m analysisModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = msg.Width - 4

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case analysisCompleteMsg:
		m.loading = false
		m.result = msg.result
		m.err = msg.err

	case analysisProgressMsg:
		m.step = msg.step
		cmd = m.progress.SetPercent(msg.progress)
		cmds = append(cmds, cmd)

	case spinner.TickMsg:
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case progress.FrameMsg:
		progressModel, progressCmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		cmds = append(cmds, progressCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m analysisModel) View() string {
	var s strings.Builder

	// Header
	header := ascii.GetBanner("Repository Analysis", ascii.EmojiBrain, m.width)
	s.WriteString(header)
	s.WriteString("\n\n")

	if m.loading {
		s.WriteString(fmt.Sprintf("🔬 Analyzing: %s\n\n", m.repoURL))
		s.WriteString(fmt.Sprintf("%s %s\n\n", m.spinner.View(), m.step))
		s.WriteString(m.progress.View())
		s.WriteString("\n\n")
	} else if m.err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f7768e"))
		s.WriteString(errorStyle.Render(fmt.Sprintf("❌ Analysis failed: %v", m.err)))
		s.WriteString("\n\n")
	} else if m.result != nil {
		s.WriteString(m.renderResults())
	}

	// Footer
	helpStyle := lipgloss.NewStyle().Foreground(ascii.Cyan)
	if m.loading {
		s.WriteString(helpStyle.Render("Analysis in progress... Press 'ctrl+c' to quit"))
	} else {
		s.WriteString(helpStyle.Render("Press 'q' to quit"))
	}

	return s.String()
}

func (m analysisModel) performAnalysis() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()

		// Simular progreso del análisis
		steps := []string{
			"Fetching repository information...",
			"Generating overview...",
			"Analyzing code structure...",
			"Evaluating dependencies...",
			"Computing metrics...",
			"Generating insights...",
		}

		for i, step := range steps {
			time.Sleep(500 * time.Millisecond) // Simular trabajo
			//progress := float64(i+1) / float64(len(steps))
			// En una implementación real, esto se enviaría a través de un canal
			_ = step // Para evitar el error de variable no utilizada
			_ = i
		}

		// Crear opciones de análisis
		analysisOptions := &models.AnalysisOptions{
			Depth:               m.options.Depth,
			IncludeDependencies: m.options.IncludeDependencies,
		}

		if m.options.Focus != "" {
			analysisOptions.Focus = strings.Split(m.options.Focus, ",")
		}

		//result, err := av.client.AnalyzeRepository(ctx, m.repoURL, analysisOptions)
		result, err := m.client.AnalyzeRepository(ctx, m.repoURL, analysisOptions)
		return analysisCompleteMsg{result: result, err: err}
	}
}

func (m analysisModel) renderResults() string {
	var s strings.Builder

	result := m.result

	// Resumen ejecutivo
	summaryStyle := ascii.InfoBoxStyle
	s.WriteString(summaryStyle.Render(fmt.Sprintf("📋 Analysis Summary\n\n%s", result.Summary)))
	s.WriteString("\n\n")

	// Insights clave
	if len(result.Insights) > 0 {
		s.WriteString("💡 Key Insights:\n\n")
		for _, insight := range result.Insights[:min(3, len(result.Insights))] {
			confidenceBar := ascii.GetProgressBar(int(insight.Confidence*100), 100, 20)
			s.WriteString(fmt.Sprintf("• %s\n  %s\n  Confidence: %s\n\n",
				insight.Title,
				insight.Description,
				confidenceBar))
		}
	}

	// Recomendaciones principales
	if len(result.Recommendations) > 0 {
		s.WriteString("🎯 Top Recommendations:\n\n")
		for i, rec := range result.Recommendations[:min(3, len(result.Recommendations))] {
			priorityEmoji := "🔵"
			if rec.Priority == "high" {
				priorityEmoji = "🔴"
			} else if rec.Priority == "medium" {
				priorityEmoji = "🟡"
			}

			s.WriteString(fmt.Sprintf("%d. %s %s\n   %s\n   Impact: %s | Effort: %s\n\n",
				i+1, priorityEmoji, rec.Title, rec.Description, rec.Impact, rec.Effort))
		}
	}

	// Estadísticas de tiempo
	s.WriteString(fmt.Sprintf("⏱️  Analysis completed in %v\n", result.Duration))

	return s.String()
}

func (av *AnalysisView) analyzeRepositoryNonInteractive(options *AnalysisOptions) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	fmt.Printf("🔬 Analyzing repository: %s\n", options.RepoURL)

	analysisOptions := &models.AnalysisOptions{
		Depth:               options.Depth,
		IncludeDependencies: options.IncludeDependencies,
	}

	if options.Focus != "" {
		analysisOptions.Focus = strings.Split(options.Focus, ",")
	}

	result, err := av.client.AnalyzeRepository(ctx, options.RepoURL, analysisOptions)
	if err != nil {
		fmt.Printf("❌ Analysis failed: %v\n", err)
		return
	}

	// Output según formato
	switch options.Format {
	case "json":
		// Implementar output JSON
		fmt.Printf("Analysis completed: %s\n", result.Summary)
	case "yaml":
		// Implementar output YAML
		fmt.Printf("Analysis completed: %s\n", result.Summary)
	default:
		fmt.Printf("\n📋 Analysis Results:\n")
		fmt.Printf("Summary: %s\n\n", result.Summary)

		if len(result.Insights) > 0 {
			fmt.Printf("💡 Key Insights:\n")
			for _, insight := range result.Insights {
				fmt.Printf("• %s: %s\n", insight.Title, insight.Description)
			}
		}
	}
}
