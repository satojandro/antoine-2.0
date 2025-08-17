package ascii

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Colores del tema Antoine
var (
	Gold     = lipgloss.Color("#FFD700")
	DarkBlue = lipgloss.Color("#1a1b26")
	Cyan     = lipgloss.Color("#7dcfff")
	White    = lipgloss.Color("#ffffff")
)

// Estilo para el logo principal
var logoStyle = lipgloss.NewStyle().
	Foreground(Gold).
	Bold(true)

// Estilo para el subtítulo
var subtitleStyle = lipgloss.NewStyle().
	Foreground(Cyan).
	Italic(true)

// Logo principal de Antoine en ASCII art (estilo pixelado como tu imagen)
const AntoineLogoLarge = `
███████╗ ███╗  ██╗ ████████╗  ██████╗  ██╗ ███╗  ██╗ ███████╗
██╔══██║ ████╗ ██║ ╚══██╔══╝ ██╔═══██╗ ██║ ████╗ ██║ ██╔════╝
███████║ ██╔██╗██║    ██║    ██║   ██║ ██║ ██╔██╗██║ █████╗  
██╔══██║ ██║╚████║    ██║    ██║   ██║ ██║ ██║╚████║ ██╔══╝  
██║  ██║ ██║ ╚███║    ██║    ╚██████╔╝ ██║ ██║ ╚███║ ███████╗
╚═╝  ╚═╝ ╚═╝  ╚══╝    ╚═╝     ╚═════╝  ╚═╝ ╚═╝  ╚══╝ ╚══════╝`

// Logo mediano para espacios más pequeños
const AntoineLogoMedium = `
 █████╗ ███╗  ██╗████████╗ ██████╗ ██╗███╗  ██╗███████╗
██╔══██╗████╗ ██║╚══██╔══╝██╔═══██╗██║████╗ ██║██╔════╝
███████║██╔██╗██║   ██║   ██║   ██║██║██╔██╗██║█████╗  
██╔══██║██║╚████║   ██║   ██║   ██║██║██║╚████║██╔══╝  
██║  ██║██║ ╚███║   ██║   ╚██████╔╝██║██║ ╚███║███████╗
╚═╝  ╚═╝╚═╝  ╚══╝   ╚═╝    ╚═════╝ ╚═╝╚═╝  ╚══╝╚══════╝`

// Logo pequeño para headers compactos
const AntoineLogoSmall = `
▄▀█ █▄░█ ▀█▀ █▀█ █ █▄░█ █▀▀
█▀█ █░▀█ ░█░ █▄█ █ █░▀█ ██▄`

// Logo minimalista para estados de carga
const AntoineLogoMini = `
╔═╗┌┐┌┌┬┐┌─┐┬┌┐┌┌─┐
╠═╣│││ │ │ │││││├┤ 
╩ ╩┘└┘ ┴ └─┘┴┘└┘└─┘`

// Subtítulos y taglines
const (
	SubtitleMain      = "Your Ultimate Hackathon Mentor"
	SubtitleThinking  = "AI-Powered Innovation Assistant"
	SubtitleSearching = "Discovering the Next Big Thing"
	SubtitleAnalyzing = "Deep Code Intelligence"
)

// Emojis temáticos para diferentes estados
const (
	EmojiRobot    = "🤖"
	EmojiSparkles = "✨"
	EmojiTarget   = "🎯"
	EmojiTrophy   = "🏆"
	EmojiBrain    = "🧠"
	EmojiSearch   = "🔍"
	EmojiCode     = "💻"
	EmojiRocket   = "🚀"
)

// GetLogo retorna el logo en el tamaño apropiado según el ancho de terminal
func GetLogo(terminalWidth int) string {
	switch {
	case terminalWidth >= 80:
		return logoStyle.Render(AntoineLogoLarge)
	case terminalWidth >= 60:
		return logoStyle.Render(AntoineLogoMedium)
	case terminalWidth >= 40:
		return logoStyle.Render(AntoineLogoSmall)
	default:
		return logoStyle.Render(AntoineLogoMini)
	}
}

// GetHeader crea un header completo con logo y subtítulo
func GetHeader(terminalWidth int, subtitle string) string {
	logo := GetLogo(terminalWidth)
	sub := subtitleStyle.Render(subtitle)

	// Centrar el contenido
	logoLines := strings.Split(logo, "\n")
	maxWidth := 0
	for _, line := range logoLines {
		if len(stripAnsi(line)) > maxWidth {
			maxWidth = len(stripAnsi(line))
		}
	}

	// Crear el header centrado
	var header strings.Builder
	for _, line := range logoLines {
		padding := (terminalWidth - len(stripAnsi(line))) / 2
		if padding > 0 {
			header.WriteString(strings.Repeat(" ", padding))
		}
		header.WriteString(line)
		header.WriteString("\n")
	}

	// Agregar subtítulo centrado
	subtitlePadding := (terminalWidth - len(subtitle)) / 2
	if subtitlePadding > 0 {
		header.WriteString(strings.Repeat(" ", subtitlePadding))
	}
	header.WriteString(sub)

	return header.String()
}

// GetBanner crea un banner decorativo
func GetBanner(title, emoji string, terminalWidth int) string {
	if terminalWidth < 20 {
		return fmt.Sprintf("%s %s", emoji, title)
	}

	bannerStyle := lipgloss.NewStyle().
		Foreground(Gold).
		Background(DarkBlue).
		Bold(true).
		Padding(0, 1).
		Margin(1, 0)

	content := fmt.Sprintf("%s %s %s", emoji, title, emoji)
	padding := (terminalWidth - len(content) - 4) / 2 // -4 for border chars

	if padding > 0 {
		border := strings.Repeat("═", padding)
		return bannerStyle.Render(fmt.Sprintf("╔%s %s %s╗", border, content, border))
	}

	return bannerStyle.Render(content)
}

// Animaciones ASCII para diferentes estados
var Animations = map[string][]string{
	"thinking": {
		"🤔 Thinking...",
		"💭 Pondering...",
		"🧠 Processing...",
		"💡 Ideating...",
	},
	"searching": {
		"🔍 Searching...",
		"🌐 Crawling web...",
		"📡 Fetching data...",
		"🎯 Finding matches...",
	},
	"analyzing": {
		"⚡ Analyzing...",
		"🔬 Deep diving...",
		"📊 Computing metrics...",
		"📈 Generating insights...",
	},
	"learning": {
		"📚 Learning...",
		"🎓 Studying patterns...",
		"💫 Absorbing knowledge...",
		"✨ Synthesizing...",
	},
}

// GetLoadingFrame retorna un frame de animación para el estado dado
func GetLoadingFrame(state string, frame int) string {
	if frames, exists := Animations[state]; exists {
		return frames[frame%len(frames)]
	}
	return "⏳ Loading..."
}

// Box styles para diferentes tipos de contenido
var (
	InfoBoxStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(Cyan).
		Padding(1, 2)

	SuccessBoxStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#9ece6a")).
		Padding(1, 2)

	WarningBoxStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ff9e64")).
		Padding(1, 2)

	ErrorBoxStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#f7768e")).
		Padding(1, 2)
)

// Progress bars temáticos
func GetProgressBar(current, total int, width int) string {
	if width < 10 {
		return fmt.Sprintf("%d/%d", current, total)
	}

	percentage := float64(current) / float64(total)
	filled := int(percentage * float64(width-2))
	empty := width - 2 - filled

	bar := "["
	bar += strings.Repeat("█", filled)
	bar += strings.Repeat("░", empty)
	bar += "]"

	return logoStyle.Render(bar) + fmt.Sprintf(" %d%%", int(percentage*100))
}

// stripAnsi removes ANSI escape codes for width calculation
func stripAnsi(s string) string {
	// Simple ANSI stripping - in production use a proper library
	result := ""
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		result += string(r)
	}
	return result
}

// Dashboard components
func GetDashboardStats(stats map[string]interface{}) string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(Gold).
		Padding(1, 2).
		Margin(1, 0)

	content := fmt.Sprintf(`%s Antoine Dashboard %s

🎯 Active Hackathons: %v    📊 Projects Analyzed: %v
🔍 Searches Today: %v       ⭐ Success Rate: %v%%

%s Trending Technologies:
████████████ AI/ML (89%%)
██████████   Blockchain (76%%)
████████     Web3 (65%%)

%s Recent Achievements:
• ETHGlobal Winner: "DeFi Innovation"
• 3 projects secured funding this week
• Community grew by 150 developers`,
		EmojiRocket, EmojiSparkles,
		stats["hackathons"], stats["projects"],
		stats["searches"], stats["success_rate"],
		EmojiBrain, EmojiTrophy)

	return style.Render(content)
}
