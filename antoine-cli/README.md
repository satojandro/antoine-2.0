> Your Ultimate Hackathon Mentor - Powered by AI and MCP Servers

Antoine CLI es una herramienta de línea de comandos inteligente que ayuda a desarrolladores a destacar en hackathons. Desde descubrir el hackathon perfecto hasta analizar proyectos ganadores y proporcionar mentoría personalizada, Antoine es tu ventaja competitiva definitiva.

## ✨ Características Principales

### 🔍 Búsqueda Inteligente
- **Hackathons**: Encuentra hackathons perfectos basados en tus intereses y habilidades
- **Proyectos**: Descubre proyectos ganadores de hackathons pasados
- **Tendencias**: Analiza tendencias tecnológicas y del mercado

### 🧠 Análisis Profundo
- **Repositorios**: Análisis completo de código, arquitectura y calidad
- **Competencia**: Evaluación comparativa con proyectos similares
- **Mercado**: Análisis de potencial comercial y oportunidades

### 🎯 Mentoría IA
- **Sesiones Interactivas**: Chat en tiempo real con Antoine
- **Feedback Personalizado**: Retroalimentación específica para tu proyecto
- **Estrategias**: Consejos basados en miles de proyectos exitosos

### 🚀 Potenciado por MCP
- **Exa**: Búsquedas web semánticas avanzadas
- **GitHub Tools**: Análisis profundo de repositorios
- **DeepWiki**: Resúmenes inteligentes de proyectos
- **E2B**: Ejecución de código y análisis personalizado
- **Browserbase**: Navegación web automatizada
- **Firecrawl**: Extracción estructurada de contenido

## 🛠️ Instalación

### Desde Release Binario
\`\`\`bash
# macOS/Linux
curl -fsSL https://get.antoine.ai | sh

# O descarga manual desde GitHub Releases
wget https://github.com/username/antoine-cli/releases/latest/download/antoine-linux-amd64
chmod +x antoine-linux-amd64
sudo mv antoine-linux-amd64 /usr/local/bin/antoine
\`\`\`

### Desde Código Fuente
\`\`\`bash
git clone https://github.com/username/antoine-cli.git
cd antoine-cli
make build
make install
\`\`\`

### Con Go
\`\`\`bash
go install github.com/username/antoine-cli/cmd/antoine@latest
\`\`\`

## 🚀 Uso Rápido

### Dashboard Principal
\`\`\`bash
antoine
\`\`\`

### Búsqueda de Hackathons
\`\`\`bash
# Búsqueda básica
antoine search hackathons

# Búsqueda avanzada
antoine search hackathons --tech "AI,Blockchain" --location "online" --prize-min 10000

# Búsqueda interactiva
antoine search hackathons --format interactive
\`\`\`

### Análisis de Proyectos
\`\`\`bash
# Análisis básico
antoine analyze repo https://github.com/user/project

# Análisis profundo
antoine analyze repo https://github.com/user/project \
--depth deep \
--include-dependencies \
--focus "security,performance"

# Análisis de tendencias
antoine analyze trends --tech "Solidity" --timeframe "1year"
\`\`\`

### Mentoría Interactiva
\`\`\`bash
# Sesión de mentoría
antoine mentor start

# Feedback específico
antoine mentor feedback --project-url "https://github.com/user/project"
\`\`\`

## ⚙️ Configuración

### Configuración Inicial
\`\`\`bash
# Ver configuración actual
antoine config show

# Configurar API key
antoine config set api-key "tu-api-key"

# Configurar servidor MCP personalizado
antoine config set mcp.servers.custom "mcp://localhost:8080"
\`\`\`

### Archivo de Configuración
\`\`\`yaml
# ~/.antoine.yaml
api:
base_url: "https://api.antoine.ai"
timeout: 30s

mcp:
servers:
exa: "mcp://localhost:8001"
github: "mcp://localhost:8002"
deepwiki: "mcp://localhost:8003"

ui:
theme: "dark"
animations: true
ascii_art: true
\`\`\`

## 📊 Ejemplos de Uso

### Prepararse para ETHGlobal
\`\`\`bash
# 1. Buscar próximos hackathons de Ethereum
antoine search hackathons --tech "Ethereum,Solidity" --org "ETHGlobal"

# 2. Analizar proyectos ganadores pasados
antoine search projects --hackathon "ETHGlobal" --category "DeFi" --sort "prize"

# 3. Obtener ideas y mentoría
antoine mentor start
> "Quiero participar en ETHGlobal con un proyecto DeFi. ¿Qué ideas innovadoras me recomiendas?"
\`\`\`

### Mejorar tu Proyecto Actual
\`\`\`bash
# 1. Análisis completo del proyecto
antoine analyze repo https://github.com/tu-usuario/tu-proyecto \
--depth deep \
--generate-report

# 2. Comparar con competencia
antoine analyze trends --tech "React,Node.js" --timeframe "6months"

# 3. Obtener feedback específico
antoine mentor feedback \
--project-url "https://github.com/tu-usuario/tu-proyecto" \
--focus "architecture,scalability"
\`\`\`

## 🎨 Interfaz de Usuario

Antoine CLI utiliza [Charm](https://charm.sh/) para crear una experiencia visual excepcional:

- **ASCII Art**: Logo animado y banners temáticos
- **Colores**: Paleta dorada/azul consistente con el branding
- **Tablas**: Datos organizados y fáciles de leer
- **Spinners**: Indicadores de progreso elegantes
- **Formularios**: Inputs interactivos y responsivos

### Temas
- `dark` (por defecto): Fondo oscuro con acentos dorados
- `light`: Fondo claro para terminales claras
- `minimal`: Sin colores ni animaciones

## 🔧 Desarrollo

### Requisitos
- Go 1.21+
- Make
- Git

### Configuración del Entorno
\`\`\`bash
git clone https://github.com/username/antoine-cli.git
cd antoine-cli

# Instalar dependencias
make deps

# Ejecutar en modo desarrollo
make dev

# Ejecutar tests
make test

# Generar cobertura
make test-coverage
\`\`\`

### Estructura del Proyecto
\`\`\`
antoine-cli/
├── cmd/              # Comandos CLI
├── internal/         # Lógica de negocio
│   ├── core/        # Cliente principal
│   ├── mcp/         # Integraciones MCP
│   ├── ui/          # Componentes UI
│   └── models/      # Modelos de datos
├── pkg/             # Paquetes exportables
└── configs/         # Configuraciones
\`\`\`

### Contribuir
1. Fork el repositorio
2. Crea una rama feature: \`git checkout -b feature/nueva-funcionalidad\`
3. Commit tus cambios: \`git commit -am 'Agregar nueva funcionalidad'\`
4. Push a la rama: \`git push origin feature/nueva-funcionalidad\`
5. Crea un Pull Request

## 📖 Documentación

- [Guía de Comandos](docs/COMMANDS.md)
- [Integración MCP](docs/MCP_INTEGRATION.md)
- [Guía de Desarrollo](docs/DEVELOPMENT.md)
- [API Reference](docs/API.md)

## 🤝 Comunidad

- [Discord](https://discord.gg/antoine-cli)
- [Twitter](https://twitter.com/antoine_ai)
- [GitHub Discussions](https://github.com/username/antoine-cli/discussions)

## 📄 Licencia

MIT License - ver [LICENSE](LICENSE) para detalles.

## 🙏 Agradecimientos

Antoine CLI está construido sobre increíbles proyectos de código abierto:

- [Charm](https://charm.sh/) - Herramientas TUI excepcionales
- [Cobra](https://github.com/spf13/cobra) - Framework CLI poderoso
- [Viper](https://github.com/spf13/viper) - Gestión de configuración
- Todos los proyectos MCP que potencian las capacidades de Antoine

---

**¡Hecho con ❤️  para la comunidad de hackathons!**

> "El futuro pertenece a aquellos que construyen con propósito" - Antoine
`
