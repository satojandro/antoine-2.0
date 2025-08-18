# Antoine CLI - Guía de Arquitectura y Estructura del Proyecto

## 📋 Información General

**Antoine CLI** es una herramienta de línea de comandos inteligente que actúa como "Your Ultimate Hackathon Mentor". Construida en Go con Charm libraries, integra múltiples servicios MCP (Model Context Protocol) para proporcionar búsqueda, análisis y mentoría IA para desarrolladores en hackathons.

### Tecnologías Core
- **Lenguaje**: Go 1.21+
- **CLI Framework**: Cobra + Viper
- **UI Framework**: Charm (Bubble Tea, Lip Gloss, Bubbles)
- **Protocolo**: MCP (Model Context Protocol)
- **Configuración**: YAML con Viper
- **Build System**: Make + GitHub Actions

---

## 🏗️ Arquitectura del Sistema

### Principios de Diseño
1. **Modularidad**: Separación clara entre comandos, lógica de negocio y UI
2. **Extensibilidad**: Fácil adición de nuevos comandos y servicios MCP
3. **Experiencia de Usuario**: UI rica en terminal que rivaliza con apps gráficas
4. **Performance**: Sistema de caché inteligente y operaciones asíncronas
5. **Configurabilidad**: Adaptable a diferentes entornos y preferencias

### Componentes Principales

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLI Commands  │───▶│  Core Client    │───▶│  MCP Servers    │
│   (Cobra)       │    │  (Business)     │    │  (External)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   UI Views      │    │  Cache System   │    │  Analytics      │
│   (Bubble Tea)  │    │  (Memory)       │    │  (Metrics)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

---

## 📁 Estructura del Proyecto

```
antoine-cli/
├── cmd/                           # Comandos CLI principales
│   ├── root.go                   # Comando raíz y configuración global
│   ├── search.go                 # antoine search hackathons|projects
│   ├── analyze.go                # antoine analyze repo|trends
│   ├── mentor.go                 # antoine mentor start|feedback
│   ├── trends.go                 # antoine trends --tech X
│   └── config.go                 # antoine config show|set
│
├── internal/                     # Lógica de negocio interna
│   ├── core/                     # Core business logic
│   │   ├── client.go            # Cliente principal de Antoine
│   │   ├── session.go           # Gestión de sesiones de usuario
│   │   ├── cache.go             # Sistema de caché en memoria
│   │   └── analytics.go         # Métricas y analytics
│   │
│   ├── mcp/                      # Integraciones MCP
│   │   ├── client.go            # Cliente MCP base (interface)
│   │   ├── exa.go               # Búsquedas web semánticas
│   │   ├── github.go            # Análisis de repositorios GitHub
│   │   ├── deepwiki.go          # Resúmenes inteligentes de repos
│   │   ├── e2b.go               # Ejecución de código sandboxed
│   │   ├── browserbase.go       # Navegación web automatizada
│   │   ├── firecrawl.go         # Extracción estructurada de contenido
│   │   └── utils.go             # Utilidades comunes MCP
│   │
│   ├── ui/                       # Componentes de interfaz Charm
│   │   ├── components/          # Componentes reutilizables
│   │   │   ├── header.go        # Header con ASCII art animado
│   │   │   ├── spinner.go       # Spinners personalizados
│   │   │   ├── table.go         # Tablas de datos interactivas
│   │   │   ├── progress.go      # Barras de progreso temáticas
│   │   │   └── form.go          # Formularios y inputs
│   │   │
│   │   ├── views/               # Vistas principales (Bubble Tea)
│   │   │   ├── dashboard.go     # Dashboard principal interactivo
│   │   │   ├── search.go        # Vista de búsqueda con filtros
│   │   │   ├── results.go       # Vista de resultados tabulares
│   │   │   ├── analysis.go      # Vista de análisis con progress
│   │   │   └── mentor.go        # Vista de chat interactivo
│   │   │
│   │   └── styles/              # Estilos y temas consistentes
│   │       ├── colors.go        # Paleta dorada/azul de Antoine
│   │       ├── typography.go    # Tipografía y formato de texto
│   │       └── layout.go        # Layouts base y spacing
│   │
│   ├── models/                   # Modelos de datos estructurados
│   │   ├── hackathon.go         # Hackathon, Location, Prize, etc.
│   │   ├── project.go           # Project, Team, Repository, etc.
│   │   ├── analysis.go          # AnalysisRequest, Result, Insights
│   │   └── response.go          # Respuestas estandarizadas MCP
│   │
│   ├── config/                   # Sistema de configuración
│   │   ├── config.go            # Carga y validación de config
│   │   ├── defaults.go          # Valores por defecto del sistema
│   │   └── validation.go        # Validadores de configuración
│   │
│   └── utils/                    # Utilidades generales
│       ├── logger.go            # Sistema de logging estructurado
│       ├── terminal.go          # Detección de capacidades terminal
│       └── helpers.go           # Funciones helper comunes
│
├── pkg/                          # Paquetes exportables
│   ├── ascii/                   # ASCII art y elementos visuales
│   │   ├── logo.go              # Logo de Antoine en múltiples tamaños
│   │   ├── animations.go        # Animaciones ASCII para estados
│   │   ├── banners.go           # Banners informativos dinámicos
│   │   └── components.go        # Componentes visuales reutilizables
│   │
│   └── terminal/                # Utilidades de terminal
│       ├── detector.go          # Detección de capabilities
│       ├── size.go              # Gestión de tamaño de terminal
│       └── colors.go            # Soporte y detección de colores
│
├── configs/                      # Archivos de configuración
│   ├── default.yaml             # Configuración base del sistema
│   ├── development.yaml         # Config para entorno desarrollo
│   └── production.yaml          # Config para entorno producción
│
├── scripts/                      # Scripts de automatización
│   ├── build.sh                 # Script de build multiplataforma
│   ├── install.sh               # Script de instalación automática
│   ├── dev.sh                   # Script para desarrollo local
│   └── mocks/                   # Servidores mock para testing
│
├── docs/                         # Documentación del proyecto
│   ├── README.md                # Documentación principal
│   ├── COMMANDS.md              # Referencia completa de comandos
│   ├── MCP_INTEGRATION.md       # Guía de integración MCP
│   ├── DEVELOPMENT.md           # Guía para contribuidores
│   └── API.md                   # Documentación de APIs internas
│
├── .github/                      # Workflows de GitHub Actions
│   └── workflows/
│       ├── ci.yml               # Tests, lint, build automático
│       ├── release.yml          # Release y distribución binarios
│       └── security.yml         # Scans de seguridad automáticos
│
├── go.mod                        # Dependencias Go principales
├── go.sum                        # Checksums de dependencias
├── Makefile                      # Comandos de build y desarrollo
├── Dockerfile                    # Imagen Docker para distribución
├── docker-compose.yml           # Stack completo para desarrollo
├── .golangci.yml                # Configuración del linter
└── .gitignore                   # Archivos excluidos de Git
```

---

## 🎨 Sistema de Diseño Visual

### Paleta de Colores (Basada en el Logo)
```go
var (
    Gold     = lipgloss.Color("#FFD700")  // Color primario del logo
    DarkBlue = lipgloss.Color("#1a1b26")  // Fondo principal
    Cyan     = lipgloss.Color("#7dcfff")  // Acentos y highlights
    Green    = lipgloss.Color("#9ece6a")  // Estados de éxito
    Orange   = lipgloss.Color("#ff9e64")  // Advertencias
    Red      = lipgloss.Color("#f7768e")  // Errores
    White    = lipgloss.Color("#ffffff")  // Texto principal
)
```

### ASCII Art Principal
```
███████╗ ███╗  ██╗ ████████╗  ██████╗  ██╗ ███╗  ██╗ ███████╗
██╔══██║ ████╗ ██║ ╚══██╔══╝ ██╔═══██╗ ██║ ████╗ ██║ ██╔════╝
███████║ ██╔██╗██║    ██║    ██║   ██║ ██║ ██╔██╗██║ █████╗
██╔══██║ ██║╚████║    ██║    ██║   ██║ ██║ ██║╚████║ ██╔══╝
██║  ██║ ██║ ╚███║    ██║    ╚██████╔╝ ██║ ██║ ╚███║ ███████╗
╚═╝  ╚═╝ ╚═╝  ╚══╝    ╚═╝     ╚═════╝  ╚═╝ ╚═╝  ╚══╝ ╚══════╝

           Your Ultimate Hackathon Mentor 🤖✨
```

### Componentes Visuales
- **Spinners Temáticos**: Diferentes animaciones por tipo de operación
- **Progress Bars**: Barras de progreso con símbolos relevantes
- **Tablas Interactivas**: Datos organizados con navegación por teclado
- **Banners Dinámicos**: Headers contextuales por comando
- **Chat Interface**: Vista de conversación para mentoría IA

---

## 🔗 Integración de Servicios MCP

### Arquitectura MCP
```go
type MCPClient interface {
    Connect(endpoint string) error
    Call(ctx context.Context, method string, params interface{}) (*MCPResponse, error)
    Subscribe(event string, handler EventHandler) error
    Disconnect() error
    Health() error
}
```

### Servicios Integrados

#### 1. **Exa** - Búsquedas Web Semánticas
```go
// Búsqueda de hackathons
hackathons, err := client.SearchHackathons(ctx, "AI blockchain", filters)

// Búsqueda de proyectos
projects, err := client.SearchProjects(ctx, "DeFi gaming", filters)

// Análisis de tendencias
trends, err := client.SearchTrends(ctx, []string{"Solidity", "Rust"}, "6months")
```

#### 2. **GitHub Tools** - Análisis de Repositorios
```go
// Análisis completo de repositorio
analysis, err := client.AnalyzeRepository(ctx, repoURL, options)

// Información básica de repo
repo, err := client.GetRepositoryInfo(ctx, repoURL)

// Lectura de archivos específicos
content, err := client.ReadFile(ctx, repoURL, "README.md")
```

#### 3. **DeepWiki** - Resúmenes Inteligentes
```go
// Overview rápido de proyecto
overview, err := client.GenerateOverview(ctx, repoURL)

// Documentación automática
docs, err := client.GenerateDocumentation(ctx, repoURL, sections)
```

#### 4. **E2B** - Ejecución de Código
```go
// Ejecución de scripts personalizados
result, err := client.ExecuteCode(ctx, &CodeExecutionRequest{
    Code:     pythonScript,
    Language: "python",
    Timeout:  60,
})

// Análisis de datos personalizado
analysis, err := client.RunAnalysis(ctx, analysisScript, data)
```

### Configuración MCP
```yaml
mcp:
  servers:
    exa: "mcp://localhost:8001"
    github: "mcp://localhost:8002"
    deepwiki: "mcp://localhost:8003"
    e2b: "mcp://localhost:8004"
    browserbase: "mcp://localhost:8005"
    firecrawl: "mcp://localhost:8006"
  timeout: "30s"
  retry_count: 3
```

---

## 🚀 Comandos y Funcionalidades

### Estructura de Comandos
```
antoine
├── search                        # Búsqueda de contenido
│   ├── hackathons               # Buscar hackathons
│   └── projects                 # Buscar proyectos
├── analyze                       # Análisis profundo
│   ├── repo [url]               # Analizar repositorio
│   ├── trends                   # Analizar tendencias
│   └── compare                  # Comparar proyectos
├── mentor                        # Mentoría IA
│   ├── start                    # Sesión interactiva
│   ├── feedback                 # Feedback de proyecto
│   └── ideate                   # Generación de ideas
├── trends                        # Análisis de mercado
│   ├── tech                     # Tendencias tecnológicas
│   └── hackathons               # Tendencias de hackathons
└── config                        # Configuración
    ├── show                     # Mostrar configuración
    ├── set [key] [value]        # Establecer valor
    └── reset                    # Restaurar defaults
```

### Ejemplos de Uso

#### Búsqueda Avanzada
```bash
# Hackathons específicos
antoine search hackathons \
  --tech "AI,Blockchain,Gaming" \
  --location "online" \
  --prize-min 10000 \
  --date-from "2025-03-01" \
  --format interactive

# Proyectos por categoría
antoine search projects \
  --hackathon "ETHGlobal" \
  --category "DeFi,Infrastructure" \
  --sort "popularity" \
  --limit 20
```

#### Análisis Profundo
```bash
# Análisis completo de repositorio
antoine analyze repo https://github.com/user/project \
  --depth deep \
  --include-dependencies \
  --focus "security,architecture,performance" \
  --generate-report

# Análisis de tendencias de mercado
antoine analyze trends \
  --tech "Solidity,Rust,Move" \
  --timeframe "1year" \
  --include-metrics \
  --market "web3"
```

#### Mentoría Interactiva
```bash
# Sesión completa de mentoría
antoine mentor start

# Feedback específico y rápido
antoine mentor feedback \
  --project-url "https://github.com/user/project" \
  --focus "architecture,scalability" \
  --quick
```

---

## 📊 Modelos de Datos

### Modelo Principal: Hackathon
```go
type Hackathon struct {
    ID              string            `json:"id"`
    Name            string            `json:"name"`
    Description     string            `json:"description"`
    URL             string            `json:"url"`
    StartDate       time.Time         `json:"start_date"`
    EndDate         time.Time         `json:"end_date"`
    Location        Location          `json:"location"`
    Technologies    []string          `json:"technologies"`
    Categories      []string          `json:"categories"`
    PrizePool       PrizeInfo         `json:"prize_pool"`
    Organizer       Organizer         `json:"organizer"`
    Difficulty      string            `json:"difficulty"`
    TeamSize        TeamSizeInfo      `json:"team_size"`
    Status          string            `json:"status"`
    ParticipantCount int              `json:"participant_count"`
    ProjectCount    int               `json:"project_count"`
    Tags            []string          `json:"tags"`
    Metadata        map[string]interface{} `json:"metadata"`
}
```

### Modelo Principal: Project
```go
type Project struct {
    ID              string              `json:"id"`
    Name            string              `json:"name"`
    Description     string              `json:"description"`
    HackathonID     string              `json:"hackathon_id"`
    Team            Team                `json:"team"`
    Repository      Repository          `json:"repository"`
    Technologies    []string            `json:"technologies"`
    Categories      []string            `json:"categories"`
    Awards          []Award             `json:"awards"`
    Metrics         ProjectMetrics      `json:"metrics"`
    Innovation      InnovationScore     `json:"innovation"`
    TechnicalDepth  TechnicalAnalysis   `json:"technical_depth"`
    MarketPotential MarketAnalysis      `json:"market_potential"`
    Status          string              `json:"status"`
    SubmissionDate  time.Time           `json:"submission_date"`
}
```

### Modelo de Análisis
```go
type AnalysisResult struct {
    ID              string                 `json:"id"`
    Type            string                 `json:"type"`
    Status          string                 `json:"status"`
    Progress        int                    `json:"progress"`
    Results         interface{}            `json:"results"`
    Summary         string                 `json:"summary"`
    Insights        []Insight              `json:"insights"`
    Recommendations []Recommendation       `json:"recommendations"`
    StartTime       time.Time              `json:"start_time"`
    EndTime         *time.Time             `json:"end_time,omitempty"`
    Duration        time.Duration          `json:"duration"`
}
```

---

## ⚙️ Sistema de Configuración

### Archivo de Configuración Principal
```yaml
# ~/.antoine.yaml
api:
  base_url: "https://api.antoine.ai"
  timeout: "30s"
  retry_count: 3
  rate_limit: 100

mcp:
  servers:
    exa: "mcp://localhost:8001"
    github: "mcp://localhost:8002"
    deepwiki: "mcp://localhost:8003"
    e2b: "mcp://localhost:8004"
  timeout: "30s"

ui:
  theme: "dark"                    # dark, light, minimal
  animations: true
  ascii_art: true
  colors: true

cache:
  enabled: true
  ttl: "30m"
  max_size: 1000

analytics:
  enabled: true
  anonymous: true
```

### Variables de Entorno
```bash
# API Configuration
export ANTOINE_API_BASE_URL="https://api.antoine.ai"
export ANTOINE_API_KEY="your-api-key"

# MCP Configuration
export ANTOINE_MCP_EXA_ENDPOINT="mcp://localhost:8001"
export ANTOINE_MCP_GITHUB_ENDPOINT="mcp://localhost:8002"

# UI Configuration
export ANTOINE_UI_THEME="dark"
export ANTOINE_NO_COLOR="false"
export ANTOINE_NO_ANIMATIONS="false"

# Debug
export ANTOINE_DEBUG="false"
export ANTOINE_LOG_LEVEL="info"
```

---

## 🔧 Sistema de Build y Deploy

### Makefile Principal
```makefile
# Variables principales
BINARY_NAME=antoine
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags="-X 'main.Version=${VERSION}'"

# Comandos principales
build:           # Build para plataforma actual
build-all:       # Build multiplataforma
test:            # Ejecutar tests
test-coverage:   # Tests con cobertura
lint:            # Linting con golangci-lint
fmt:             # Formateo de código
security:        # Scan de seguridad
install:         # Instalar en sistema
clean:           # Limpiar artifacts
docker-build:    # Build imagen Docker
```

### Distribución Multiplataforma
```bash
# Plataformas soportadas
- linux/amd64
- linux/arm64
- darwin/amd64 (Intel Mac)
- darwin/arm64 (Apple Silicon)
- windows/amd64
```

### CI/CD con GitHub Actions
```yaml
# Workflows automáticos
- ci.yml:          # Tests, lint, build en PRs
- release.yml:     # Release automático en tags
- security.yml:    # Scans de seguridad programados
```

---

## 🧪 Testing y Calidad

### Estrategia de Testing
```go
// Tests unitarios por paquete
internal/core/client_test.go
internal/mcp/exa_test.go
internal/ui/views/search_test.go

// Tests de integración
tests/integration/mcp_integration_test.go
tests/integration/cli_integration_test.go

// Tests end-to-end
tests/e2e/search_workflow_test.go
tests/e2e/analysis_workflow_test.go
```

### Herramientas de Calidad
- **golangci-lint**: Linting completo con 40+ linters
- **gosec**: Análisis de seguridad del código
- **govulncheck**: Detección de vulnerabilidades
- **CodeQL**: Análisis estático de seguridad
- **Codecov**: Cobertura de tests automática

---

## 🚀 Deployment y Distribución

### Métodos de Instalación

#### 1. Script de Instalación Automática
```bash
curl -fsSL https://get.antoine.ai | sh
```

#### 2. Binarios Pre-compilados
```bash
# GitHub Releases con binarios firmados
wget https://github.com/org/antoine-cli/releases/latest/download/antoine-linux-amd64
```

#### 3. Package Managers
```bash
# Homebrew (macOS/Linux)
brew install antoine-cli

# Scoop (Windows)
scoop install antoine-cli

# APT (Ubuntu/Debian)
apt install antoine-cli
```

#### 4. Docker
```bash
# Imagen oficial
docker run --rm -it antoine/cli:latest

# Con configuración persistente
docker run --rm -it -v ~/.antoine.yaml:/home/antoine/.antoine.yaml antoine/cli
```

#### 5. Go Install
```bash
go install github.com/org/antoine-cli/cmd/antoine@latest
```

---

## 🔒 Seguridad y Privacidad

### Gestión de Credenciales
```go
// Almacenamiento seguro con keyring del sistema
type CredentialManager struct {
    keyring   keyring.Keyring
    encrypted map[string][]byte
}

// APIs keys encriptadas localmente
// Rotación automática de tokens
// Sin almacenamiento de datos sensibles
```

### Medidas de Seguridad
- **Encriptación AES-256** para credenciales locales
- **HTTPS obligatorio** para todas las conexiones
- **Validación de certificados** SSL/TLS
- **Rate limiting** integrado para APIs
- **Timeouts configurables** para prevenir ataques
- **Análisis de vulnerabilidades** automático en CI/CD

---

## 📈 Performance y Optimización

### Sistema de Caché
```go
type CacheManager struct {
    memory    *ristretto.Cache     // Caché en memoria rápido
    disk      *badger.DB           // Caché persistente en disco
    ttl       time.Duration        // Time-to-live configurable
}

// Estrategias de caché por tipo de contenido:
// - Hackathons: 30 minutos
// - Proyectos: 1 hora
// - Análisis: 24 horas
// - Tendencias: 6 horas
```

### Optimizaciones
- **Conexiones MCP reutilizables** con pool de conexiones
- **Requests paralelos** cuando es posible
- **Paginación automática** para grandes datasets
- **Lazy loading** de datos no críticos
- **Compresión** de respuestas grandes

---

## 🤝 Guía para Contribuidores

### Setup de Desarrollo
```bash
# 1. Clonar repositorio
git clone https://github.com/org/antoine-cli.git
cd antoine-cli

# 2. Instalar dependencias
make deps

# 3. Configurar entorno de desarrollo
cp configs/development.yaml ~/.antoine.yaml
make dev

# 4. Ejecutar tests
make test
make lint
```

### Estructura para Nuevas Features

#### Nuevo Comando
```bash
# 1. Crear archivo de comando
cmd/nuevo_comando.go

# 2. Implementar lógica de negocio
internal/core/nuevo_servicio.go

# 3. Crear vista si es necesaria
internal/ui/views/nueva_vista.go

# 4. Añadir tests
cmd/nuevo_comando_test.go
internal/core/nuevo_servicio_test.go

# 5. Documentar
docs/NUEVO_COMANDO.md
```

#### Nueva Integración MCP
```bash
# 1. Implementar cliente MCP
internal/mcp/nuevo_servicio.go

# 2. Añadir al cliente principal
internal/core/client.go

# 3. Configurar conexión
configs/default.yaml

# 4. Añadir tests de integración
tests/integration/nuevo_servicio_test.go
```

### Estándares de Código
- **gofmt** para formateo automático
- **golangci-lint** debe pasar sin errores
- **Test coverage** mínimo del 80%
- **Documentación** obligatoria para APIs públicas
- **Conventional commits** para mensajes de commit

---

## 📚 Referencias y Recursos

### Dependencias Principales
```go
// CLI y Configuración
github.com/spf13/cobra v1.8.0          // Framework CLI
github.com/spf13/viper v1.17.0         // Gestión de configuración

// UI Terminal
github.com/charmbracelet/bubbletea v0.24.2   // TUI framework
github.com/charmbracelet/lipgloss v0.9.1     // Styling para terminal
github.com/charmbracelet/bubbles v0.17.1     // Componentes UI

// Utilidades
golang.org/x/term v0.15.0              // Terminal utilities
gopkg.in/yaml.v3 v3.0.1               // YAML parsing
```

### Documentación Externa
- [Charm Libraries](https://charm.sh/) - Framework de UI terminal
- [Cobra CLI](https://cobra.dev/) - Framework para CLIs en Go
- [MCP Specification](https://spec.modelcontextprotocol.io/) - Protocolo MCP
- [Go Best Practices](https://go.dev/doc/effective_go) - Mejores prácticas Go

### Arquitecturas de Referencia
- [Kubernetes CLI (kubectl)](https://github.com/kubernetes/kubectl) - CLI complejo bien estructurado
- [GitHub CLI (gh)](https://github.com/cli/cli) - Integración con APIs externas
- [Charm Gum](https://github.com/charmbracelet/gum) - Excelente UI terminal

---

## 🎯 Roadmap y Futuras Expansiones

### Fase 1: Core MVP (4 semanas)
- [x] Estructura básica del proyecto
- [x] Comandos principales (search, analyze, mentor)
- [x] Cliente MCP base con mocks
- [x] UI básica con Charm libraries
- [x] Sistema de configuración

### Fase 2: Integración Real (3 semanas)
- [ ] Conectar con servicios MCP reales
- [ ] Implementar análisis de repositorios completo
- [ ] Sistema de caché robusto
- [ ] Mentoría IA funcional

### Fase 3: Features Avanzadas (4 semanas)
- [ ] Dashboard interactivo avanzado
- [ ] Análisis comparativo de proyectos
- [ ] Sistema de recomendaciones personalizadas
- [ ] Integración con calendarios de hackathons

### Fase 4: Ecosystem (Ongoing)
- [ ] Plugin system para extensibilidad
- [ ] API REST para integración externa
- [ ] Web dashboard complementario
- [ ] Mobile companion app

### Posibles Extensiones
- **Antoine IDE Extension**: Plugin para VSCode/JetBrains
- **Antoine GitHub App**: Integración nativa con GitHub
- **Antoine Discord Bot**: Bot para comunidades de hackathons
- **Antoine Analytics Dashboard**: Dashboard web para organizadores

---

## 📞 Contacto y Soporte

### Información del Proyecto
- **Repository**: https://github.com/org/antoine-cli
- **Documentation**: https://docs.antoine.ai
- **Website**: https://antoine.ai

### Comunidad
- **Discord**: https://discord.gg/antoine
- **Twitter**: https://twitter.com/antoine_ai
- **Reddit**: https://reddit.com/r/AntoineAI

### Desarrollo
- **Issues**: GitHub Issues para bugs y feature requests
- **Discussions**: GitHub Discussions para preguntas generales
- **Contributing**: Ver CONTRIBUTING.md para guía de contribución

---

*Esta guía sirve como referencia completa para el desarrollo, mantenimiento y extensión del proyecto Antoine CLI. Mantenerla actualizada es responsabilidad de todos los contribuidores.*

**Version**: 1.0
**Last Updated**: 2025-01-15
**Next Review**: 2025-02-15
