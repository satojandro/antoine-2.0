# Antoine CLI

<div align="center">

```
███████╗ ███╗  ██╗ ████████╗  ██████╗  ██╗ ███╗  ██╗ ███████╗
██╔══██║ ████╗ ██║ ╚══██╔══╝ ██╔═══██╗ ██║ ████╗ ██║ ██╔════╝
███████║ ██╔██╗██║    ██║    ██║   ██║ ██║ ██╔██╗██║ █████╗  
██╔══██║ ██║╚████║    ██║    ██║   ██║ ██║ ██║╚████║ ██╔══╝  
██║  ██║ ██║ ╚███║    ██║    ╚██████╔╝ ██║ ██║ ╚███║ ███████╗
╚═╝  ╚═╝ ╚═╝  ╚══╝    ╚═╝     ╚═════╝  ╚═╝ ╚═╝  ╚══╝ ╚══════╝
```

**Your Ultimate Hackathon Mentor** 🤖✨

*AI-powered CLI tool that helps developers excel in hackathons through intelligent search, deep analysis, and personalized mentorship.*

[![Go Report Card](https://goreportcard.com/badge/github.com/username/antoine-cli)](https://goreportcard.com/report/github.com/username/antoine-cli)
[![Release](https://img.shields.io/github/release/username/antoine-cli.svg)](https://github.com/username/antoine-cli/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org)

[Features](#-features) •
[Installation](#-installation) •
[Quick Start](#-quick-start) •
[Documentation](#-documentation) •
[Contributing](#-contributing)

</div>

---

## 🎯 What is Antoine CLI?

Antoine CLI is an intelligent command-line tool designed to give developers a competitive edge in hackathons. From discovering the perfect hackathon to analyzing winning projects and providing personalized AI mentorship, Antoine transforms how you approach hackathon success.

**Built with Go and powered by MCP (Model Context Protocol) servers, Antoine delivers a rich terminal experience using Charm libraries.**

## ✨ Features

### 🔍 **Intelligent Discovery**
- **Smart Hackathon Search**: Find hackathons matching your skills, interests, and goals
- **Project Research**: Discover winning projects from past hackathons with detailed insights
- **Trend Analysis**: Stay ahead with real-time technology and market trend analysis

### 🧠 **Deep Code Analysis**
- **Repository Intelligence**: Comprehensive code quality, architecture, and security analysis
- **Competitive Research**: Compare your project against similar winning submissions
- **Market Potential**: Evaluate commercial viability and market opportunities

### 🎯 **AI-Powered Mentorship**
- **Interactive Chat**: Real-time conversations with Antoine AI mentor
- **Personalized Feedback**: Project-specific recommendations and improvements
- **Strategy Guidance**: Data-driven advice based on thousands of successful projects

### 🚀 **MCP Integration**
Antoine leverages multiple specialized services through the Model Context Protocol:

- **[Exa](https://exa.ai)**: Advanced semantic web search
- **GitHub Tools**: Deep repository analysis and insights
- **DeepWiki**: Intelligent project summarization
- **E2B**: Sandboxed code execution and custom analysis
- **Browserbase**: Automated web navigation
- **Firecrawl**: Structured content extraction

## 🛠️ Installation

### Quick Install (Recommended)
```bash
# macOS/Linux
curl -fsSL https://get.antoine.ai | sh

# Or via Homebrew
brew install antoine-cli
```

### Download Binary
```bash
# Linux x64
wget https://github.com/username/antoine-cli/releases/latest/download/antoine-linux-amd64
chmod +x antoine-linux-amd64
sudo mv antoine-linux-amd64 /usr/local/bin/antoine

# macOS (Intel)
wget https://github.com/username/antoine-cli/releases/latest/download/antoine-darwin-amd64

# macOS (Apple Silicon)
wget https://github.com/username/antoine-cli/releases/latest/download/antoine-darwin-arm64
```

### Build from Source
```bash
git clone https://github.com/username/antoine-cli.git
cd antoine-cli
make build && make install
```

### Go Install
```bash
go install github.com/username/antoine-cli/cmd/antoine@latest
```

## 🚀 Quick Start

### Launch Interactive Dashboard
```bash
antoine
```

### Discover Hackathons
```bash
# Basic search
antoine search hackathons

# Advanced filtering
antoine search hackathons \
  --tech "AI,Blockchain,Web3" \
  --location "online" \
  --prize-min 10000 \
  --format interactive

# Trending hackathons
antoine search hackathons --trending
```

### Analyze Projects
```bash
# Quick analysis
antoine analyze repo https://github.com/user/awesome-project

# Deep dive analysis
antoine analyze repo https://github.com/user/awesome-project \
  --depth deep \
  --include-deps \
  --focus "security,performance,scalability" \
  --generate-report

# Technology trends
antoine analyze trends --tech "Solidity,Rust" --timeframe "1year"
```

### Get AI Mentorship
```bash
# Start interactive session
antoine mentor start

# Quick project feedback
antoine mentor feedback \
  --project https://github.com/user/project \
  --focus "architecture,market-fit"

# Ideation session
antoine mentor ideate --theme "DeFi" --difficulty "beginner"
```

## ⚙️ Configuration

### Initial Setup
```bash
# View current configuration
antoine config show

# Set API credentials
antoine config set api.key "your-api-key"

# Configure custom MCP server
antoine config set mcp.servers.custom "mcp://localhost:8080"

# Change UI theme
antoine config set ui.theme "dark"  # dark, light, minimal
```

### Configuration File
Create `~/.antoine.yaml`:

```yaml
# API Configuration
api:
  base_url: "https://api.antoine.ai"
  key: "${ANTOINE_API_KEY}"
  timeout: "30s"
  retry_count: 3

# MCP Server Endpoints
mcp:
  servers:
    exa: "mcp://localhost:8001"
    github: "mcp://localhost:8002"
    deepwiki: "mcp://localhost:8003"
    e2b: "mcp://localhost:8004"
  timeout: "30s"

# UI Preferences
ui:
  theme: "dark"                    # dark, light, minimal
  animations: true
  ascii_art: true
  colors: true

# Performance
cache:
  enabled: true
  ttl: "30m"
  max_size: 1000

# Analytics (anonymous)
analytics:
  enabled: true
  anonymous: true
```

### Environment Variables
```bash
export ANTOINE_API_KEY="your-api-key"
export ANTOINE_UI_THEME="dark"
export ANTOINE_DEBUG="false"
export ANTOINE_LOG_LEVEL="info"
```

## 📋 Use Cases

### 🏆 Preparing for ETHGlobal
```bash
# 1. Discover upcoming Ethereum hackathons
antoine search hackathons --tech "Ethereum,Solidity" --org "ETHGlobal"

# 2. Research winning DeFi projects
antoine search projects \
  --hackathon "ETHGlobal" \
  --category "DeFi" \
  --sort "prize" \
  --limit 10

# 3. Get project ideas and strategy
antoine mentor start
# > "I want to build a DeFi project for ETHGlobal. What innovative ideas do you recommend?"
```

### 🔧 Improving Your Current Project
```bash
# 1. Comprehensive project analysis
antoine analyze repo https://github.com/yourname/your-project \
  --depth deep \
  --generate-report \
  --output analysis-report.md

# 2. Competitive landscape analysis
antoine analyze trends \
  --tech "React,Node.js,GraphQL" \
  --timeframe "6months" \
  --include-market-data

# 3. Get targeted feedback
antoine mentor feedback \
  --project https://github.com/yourname/your-project \
  --focus "scalability,user-experience,monetization"
```

### 🎯 Market Research
```bash
# Technology adoption trends
antoine analyze trends \
  --tech "AI,Machine Learning,LLM" \
  --timeframe "1year" \
  --include-jobs-data

# Hackathon landscape analysis
antoine search hackathons \
  --tech "AI" \
  --date-from "2024-01-01" \
  --format "csv" \
  --output hackathons-2024.csv
```

## 🎨 Terminal Experience

Antoine CLI delivers a premium terminal experience using [Charm](https://charm.sh/) libraries:

### Visual Features
- **🎨 Rich ASCII Art**: Animated logo and themed banners
- **🌈 Color Themes**: Consistent golden/blue palette matching brand identity
- **📊 Interactive Tables**: Sortable, filterable data views
- **⚡ Smart Spinners**: Context-aware loading indicators
- **📝 Responsive Forms**: Adaptive input components

### Themes
- **`dark`** (default): Dark background with golden accents
- **`light`**: Light theme for bright terminals
- **`minimal`**: No colors or animations for maximum compatibility

### Screenshots
```bash
# View the interactive dashboard
antoine

# Example output:
╔══════════════════════════════════════════════════════════════════════════════╗
║                               🤖 ANTOINE CLI ✨                              ║
║                        Your Ultimate Hackathon Mentor                        ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ 🎯 Active Hackathons: 23     📊 Projects Analyzed: 1,247                    ║
║ 🔍 Searches Today: 45        ⭐ Success Rate: 94%                           ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

## 🏗️ Architecture

Antoine CLI follows modern Go application architecture principles:

```
antoine-cli/
├── cmd/                    # CLI commands (Cobra)
├── internal/
│   ├── core/              # Business logic
│   ├── mcp/               # MCP integrations
│   ├── ui/                # Terminal UI (Bubble Tea)
│   │   ├── components/    # Reusable UI components
│   │   ├── views/         # Main application views
│   │   └── styles/        # Consistent styling
│   └── models/            # Data models
├── pkg/                   # Exportable packages
├── configs/               # Configuration files
└── docs/                  # Documentation
```

### Key Technologies
- **[Go 1.21+](https://golang.org)**: Core language
- **[Cobra](https://github.com/spf13/cobra)**: CLI framework
- **[Viper](https://github.com/spf13/viper)**: Configuration management
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)**: Terminal UI framework
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)**: Styling and layout
- **[MCP](https://spec.modelcontextprotocol.io/)**: Model Context Protocol

## 🧪 Development

### Prerequisites
- Go 1.21+
- Make
- Git

### Setup Development Environment
```bash
git clone https://github.com/username/antoine-cli.git
cd antoine-cli

# Install dependencies
make deps

# Run in development mode
make dev

# Run tests
make test

# Generate coverage report
make test-coverage

# Lint code
make lint
```

### Build Commands
```bash
make build              # Build for current platform
make build-all          # Build for all platforms
make docker-build       # Build Docker image
make install            # Install to system
make clean              # Clean build artifacts
```

### Running Tests
```bash
make test               # Run all tests
make test-unit          # Run unit tests only
make test-integration   # Run integration tests
make test-e2e          # Run end-to-end tests
make test-coverage      # Generate coverage report
```

## 📖 Documentation

- **[Command Reference](docs/COMMANDS.md)** - Complete command documentation
- **[MCP Integration Guide](docs/MCP_INTEGRATION.md)** - How to integrate MCP servers
- **[Development Guide](docs/DEVELOPMENT.md)** - Contributing and development setup
- **[API Reference](docs/API.md)** - Internal API documentation
- **[Configuration Guide](docs/CONFIGURATION.md)** - Advanced configuration options

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Quick Contribution Steps
1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/amazing-feature`
3. **Commit** your changes: `git commit -m 'Add amazing feature'`
4. **Push** to the branch: `git push origin feature/amazing-feature`
5. **Open** a Pull Request

### Development Principles
- Write tests for new features
- Follow Go best practices
- Maintain backward compatibility
- Update documentation
- Use conventional commits

## 🌟 Community

Join our growing community of hackathon enthusiasts and developers:

- **[Discord](https://discord.gg/antoine-cli)** - Real-time chat and support
- **[GitHub Discussions](https://github.com/username/antoine-cli/discussions)** - Ideas and feedback
- **[Twitter](https://twitter.com/antoine_ai)** - Updates and announcements
- **[Reddit](https://reddit.com/r/AntoineAI)** - Community discussions

## 📊 Roadmap

### 🎯 Current (v1.0)
- ✅ Core CLI framework
- ✅ MCP integration
- ✅ Basic search and analysis
- ✅ Terminal UI with Charm

### 🚀 Next (v1.1)
- 🔄 Advanced AI mentorship
- 🔄 Project comparison features
- 🔄 Web dashboard companion
- 🔄 Plugin system

### 🌟 Future (v2.0)
- 📋 IDE integrations (VSCode, JetBrains)
- 📋 GitHub App integration
- 📋 Discord bot
- 📋 Mobile companion app

## 📄 License

Antoine CLI is released under the [MIT License](LICENSE).

## 🙏 Acknowledgments

Antoine CLI is built on the shoulders of amazing open-source projects:

- **[Charm](https://charm.sh/)** - Exceptional TUI tools that make terminals beautiful
- **[Cobra](https://github.com/spf13/cobra)** - Powerful CLI framework for Go
- **[Viper](https://github.com/spf13/viper)** - Complete configuration solution
- **[Go](https://golang.org)** - The language that powers Antoine
- All the MCP server implementations that enable Antoine's capabilities

Special thanks to the hackathon community for inspiring this project.

---

<div align="center">

**Built with ❤️ for the global hackathon community**

*"The future belongs to those who build with purpose"* - Antoine

[![Made with Go](https://img.shields.io/badge/Made%20with-Go-00ADD8?logo=go)](https://golang.org)
[![Powered by Charm](https://img.shields.io/badge/Powered%20by-Charm-FF69B4)](https://charm.sh)
[![MCP Compatible](https://img.shields.io/badge/MCP-Compatible-00D4AA)](https://spec.modelcontextprotocol.io/)

</div>
