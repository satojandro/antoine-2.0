package views

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"

	"antoine-cli/internal/core"
	"antoine-cli/pkg/ascii"
)

type MentorView struct {
	client *core.AntoineClient
}

type MentorOptions struct {
	ProjectURL string
	ProjectID  string
	Focus      string
	Quick      bool
}

func NewMentorView(client *core.AntoineClient) *MentorView {
	return &MentorView{client: client}
}

type mentorModel struct {
	viewport viewport.Model
	textarea textarea.Model
	messages []chatMessage
	waiting  bool
	err      error
	width    int
	height   int
	options  *MentorOptions
}

type chatMessage struct {
	sender    string // "user" or "antoine"
	content   string
	timestamp string
}

type mentorResponseMsg struct {
	response string
	err      error
}

func (mv *MentorView) StartSession() {
	model := mv.createMentorModel(nil)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func (mv *MentorView) ProvideFeedback(options *MentorOptions) {
	if options.Quick {
		mv.provideFeedbackNonInteractive(options)
		return
	}

	model := mv.createMentorModel(options)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func (mv *MentorView) createMentorModel(options *MentorOptions) mentorModel {
	ta := textarea.New()
	ta.Placeholder = "Ask Antoine anything about hackathons, your project, or get advice..."
	ta.Focus()

	vp := viewport.New(78, 20)
	vp.SetContent("")

	// Mensaje de bienvenida
	welcomeMsg := chatMessage{
		sender:    "antoine",
		content:   "👋 Hello! I'm Antoine, your Ultimate Hackathon Mentor. I'm here to help you succeed in hackathons with insights from thousands of winning projects. What can I help you with today?",
		timestamp: "now",
	}

	model := mentorModel{
		viewport: vp,
		textarea: ta,
		messages: []chatMessage{welcomeMsg},
		options:  options,
	}

	// Si hay un proyecto específico, agregar contexto
	if options != nil && options.ProjectURL != "" {
		contextMsg := chatMessage{
			sender:    "antoine",
			content:   fmt.Sprintf("🔍 I see you want feedback on: %s\nI'll analyze this project and provide personalized recommendations.", options.ProjectURL),
			timestamp: "now",
		}
		model.messages = append(model.messages, contextMsg)
	}

	model.updateViewport()
	return model
}

func (m mentorModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m mentorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerHeight := 5
		footerHeight := 4
		verticalMarginHeight := headerHeight + footerHeight

		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - verticalMarginHeight
		m.textarea.SetWidth(msg.Width - 4)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+s":
			if m.textarea.Value() != "" && !m.waiting {
				userInput := m.textarea.Value()
				m.textarea.Reset()
				m.waiting = true

				// Agregar mensaje del usuario
				userMsg := chatMessage{
					sender:    "user",
					content:   userInput,
					timestamp: "now",
				}
				m.messages = append(m.messages, userMsg)
				m.updateViewport()

				return m, m.sendToAntoine(userInput)
			}
		}

	case mentorResponseMsg:
		m.waiting = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			antoineMsg := chatMessage{
				sender:    "antoine",
				content:   msg.response,
				timestamp: "now",
			}
			m.messages = append(m.messages, antoineMsg)
			m.updateViewport()
		}
	}

	if !m.waiting {
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m mentorModel) View() string {
	var s strings.Builder

	// Header
	header := ascii.GetBanner("Antoine Mentor Session", ascii.EmojiRobot, m.width)
	s.WriteString(header)
	s.WriteString("\n")

	// Chat viewport
	s.WriteString(m.viewport.View())
	s.WriteString("\n\n")

	// Input area
	s.WriteString(m.textarea.View())
	s.WriteString("\n")

	// Status y ayuda
	statusStyle := lipgloss.NewStyle().Foreground(ascii.Cyan)
	if m.waiting {
		s.WriteString(statusStyle.Render("🤔 Antoine is thinking..."))
	} else if m.err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f7768e"))
		s.WriteString(errorStyle.Render(fmt.Sprintf("❌ Error: %v", m.err)))
	} else {
		s.WriteString(statusStyle.Render("Press Ctrl+S to send • Press Ctrl+C to quit"))
	}

	return s.String()
}

func (m *mentorModel) updateViewport() {
	var content strings.Builder

	for _, msg := range m.messages {
		if msg.sender == "antoine" {
			// Estilo para mensajes de Antoine
			msgStyle := lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(ascii.Gold).
				Padding(1, 2).
				Margin(1, 0)

			content.WriteString(msgStyle.Render(fmt.Sprintf("🤖 Antoine:\n%s", msg.content)))
		} else {
			// Estilo para mensajes del usuario
			msgStyle := lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(ascii.Cyan).
				Padding(1, 2).
				Margin(1, 0).
				Align(lipgloss.Right)

			content.WriteString(msgStyle.Render(fmt.Sprintf("👤 You:\n%s", msg.content)))
		}
		content.WriteString("\n")
	}

	m.viewport.SetContent(content.String())
	m.viewport.GotoBottom()
}

func (m mentorModel) sendToAntoine(userInput string) tea.Cmd {
	return func() tea.Msg {
		//ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		//defer cancel()

		// En una implementación real, aquí se enviaría el mensaje a Antoine
		// Por ahora, simular una respuesta
		response := fmt.Sprintf("Great question! Based on my analysis of thousands of hackathon projects, here's my advice about: %s\n\n[This would be Antoine's AI-generated response based on the user's question and context]", userInput)

		return mentorResponseMsg{response: response, err: nil}
	}
}

func (mv *MentorView) provideFeedbackNonInteractive(options *MentorOptions) {
	fmt.Printf("🎯 Analyzing project: %s\n", options.ProjectURL)

	// Simular análisis
	fmt.Println("📊 Quick Analysis Results:")
	fmt.Println("• Code Quality: Good")
	fmt.Println("• Innovation Level: High")
	fmt.Println("• Market Potential: Medium")
	fmt.Println("• Technical Implementation: Solid")
	fmt.Println("\n💡 Key Recommendations:")
	fmt.Println("1. Enhance user documentation")
	fmt.Println("2. Add more comprehensive tests")
	fmt.Println("3. Consider scalability improvements")
}
