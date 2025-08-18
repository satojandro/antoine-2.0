package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"antoine-cli/internal/config"
	"antoine-cli/internal/mcp"
	"antoine-cli/internal/models"
)

type AntoineClient struct {
	config    *config.Config
	mcp       *MCPManager
	cache     *CacheManager
	session   *SessionManager
	analytics *AnalyticsManager
	mu        sync.RWMutex
}

type MCPManager struct {
	exa      *mcp.ExaClient
	github   *mcp.GitHubClient
	deepwiki *mcp.DeepWikiClient
	e2b      *mcp.E2BClient
	//browserbase *mcp.BrowserbaseClient
	//firecrawl   *mcp.FirecrawlClient
}

func NewAntoineClient(cfg *config.Config) (*AntoineClient, error) {
	mcpManager := &MCPManager{
		exa:      mcp.NewExaClient(),
		github:   mcp.NewGitHubClient(),
		deepwiki: mcp.NewDeepWikiClient(),
		e2b:      mcp.NewE2BClient(),
		// browserbase: mcp.NewBrowserbaseClient(),
		// firecrawl: mcp.NewFirecrawlClient(),
	}

	// Conectar a los servidores MCP
	if err := mcpManager.Connect(cfg); err != nil {
		return nil, fmt.Errorf("failed to connect to MCP servers: %w", err)
	}

	client := &AntoineClient{
		config:    cfg,
		mcp:       mcpManager,
		cache:     NewCacheManager(),
		session:   NewSessionManager(),
		analytics: NewAnalyticsManager(),
	}

	return client, nil
}

func (m *MCPManager) Connect(cfg *config.Config) error {
	// Conectar Exa
	if serverConfig, ok := cfg.MCP.Servers["exa"]; ok {
		if err := m.exa.Connect(serverConfig.Endpoint); err != nil {
			return fmt.Errorf("failed to connect to Exa: %w", err)
		}
	}

	//// Conectar GitHub
	//if endpoint, ok := cfg.MCP.Servers["github"]; ok {
	//	if err := m.github.Connect(endpoint); err != nil {
	//		return fmt.Errorf("failed to connect to GitHub: %w", err)
	//	}
	//}

	// Conectar GitHub
	if serverConfig, ok := cfg.MCP.Servers["github"]; ok {
		if err := m.github.Connect(serverConfig.Endpoint); err != nil {
			return fmt.Errorf("failed to connect to GitHub: %w", err)
		}
	}

	// Conectar DeepWiki
	if serverConfig, ok := cfg.MCP.Servers["deepwiki"]; ok {
		if err := m.deepwiki.Connect(serverConfig.Endpoint); err != nil {
			return fmt.Errorf("failed to connect to DeepWiki: %w", err)
		}
	}

	// Conectar E2B
	if serverConfig, ok := cfg.MCP.Servers["e2b"]; ok {
		if err := m.e2b.Connect(serverConfig.Endpoint); err != nil {
			return fmt.Errorf("failed to connect to E2B: %w", err)
		}
	}

	return nil
}

// SearchHackathons busca hackathons usando múltiples fuentes
func (c *AntoineClient) SearchHackathons(ctx context.Context, query string, filters map[string]interface{}) ([]*models.Hackathon, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Verificar caché primero
	cacheKey := fmt.Sprintf("hackathons:%s:%v", query, filters)
	if cached, found := c.cache.Get(cacheKey); found {
		if hackathons, ok := cached.([]*models.Hackathon); ok {
			return hackathons, nil
		}
	}

	// Buscar usando Exa
	hackathons, err := c.mcp.exa.SearchHackathons(ctx, query, filters)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// Guardar en caché
	c.cache.Set(cacheKey, hackathons, 30*time.Minute)

	// Registrar métricas
	c.analytics.RecordSearch("hackathons", len(hackathons))

	return hackathons, nil
}

// SearchProjects busca proyectos de hackathons
func (c *AntoineClient) SearchProjects(ctx context.Context, query string, filters map[string]interface{}) ([]*models.Project, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cacheKey := fmt.Sprintf("projects:%s:%v", query, filters)
	if cached, found := c.cache.Get(cacheKey); found {
		if projects, ok := cached.([]*models.Project); ok {
			return projects, nil
		}
	}

	projects, err := c.mcp.exa.SearchProjects(ctx, query, filters)
	if err != nil {
		return nil, fmt.Errorf("project search failed: %w", err)
	}

	c.cache.Set(cacheKey, projects, 30*time.Minute)
	c.analytics.RecordSearch("projects", len(projects))

	return projects, nil
}

// AnalyzeRepository analiza un repositorio de GitHub
func (c *AntoineClient) AnalyzeRepository(ctx context.Context, repoURL string, options *models.AnalysisOptions) (*models.AnalysisResult, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Primero obtener overview rápido con DeepWiki
	overview, err := c.mcp.deepwiki.GenerateOverview(ctx, repoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate overview: %w", err)
	}

	// Luego análisis profundo con GitHub tools
	analysis, err := c.mcp.github.AnalyzeRepository(ctx, repoURL, options)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze repository: %w", err)
	}

	// Combinar resultados
	analysis.Summary = overview

	c.analytics.RecordAnalysis("repository", repoURL)

	return analysis, nil
}

// GetTrends obtiene tendencias de tecnologías
func (c *AntoineClient) GetTrends(ctx context.Context, technologies []string, timeframe string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cacheKey := fmt.Sprintf("trends:%v:%s", technologies, timeframe)
	if cached, found := c.cache.Get(cacheKey); found {
		return cached, nil
	}

	trends, err := c.mcp.exa.SearchTrends(ctx, technologies, timeframe)
	if err != nil {
		return nil, fmt.Errorf("failed to get trends: %w", err)
	}

	c.cache.Set(cacheKey, trends, 1*time.Hour)
	c.analytics.RecordTrends(technologies)

	return trends, nil
}

// ExecuteAnalysisScript ejecuta un script de análisis personalizado
func (c *AntoineClient) ExecuteAnalysisScript(ctx context.Context, script string, data interface{}) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result, err := c.mcp.e2b.RunAnalysis(ctx, script, data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute analysis: %w", err)
	}

	return result, nil
}

// Health verifica el estado de todos los servicios
func (c *AntoineClient) Health(ctx context.Context) map[string]bool {
	status := make(map[string]bool)

	status["exa"] = c.mcp.exa.Health() == nil
	status["github"] = c.mcp.github.Health() == nil
	status["deepwiki"] = c.mcp.deepwiki.Health() == nil
	status["e2b"] = c.mcp.e2b.Health() == nil
	status["cache"] = c.cache.Health() == nil

	return status
}

// Close cierra todas las conexiones
func (c *AntoineClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var errors []error

	if err := c.mcp.exa.Disconnect(); err != nil {
		errors = append(errors, err)
	}

	if err := c.mcp.github.Disconnect(); err != nil {
		errors = append(errors, err)
	}

	if err := c.mcp.deepwiki.Disconnect(); err != nil {
		errors = append(errors, err)
	}

	if err := c.mcp.e2b.Disconnect(); err != nil {
		errors = append(errors, err)
	}

	if err := c.cache.Close(); err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors during close: %v", errors)
	}

	return nil
}
