package config

import (
	"antoine-cli/pkg/ascii"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	API       APIConfig       `mapstructure:"api"`
	MCP       MCPConfig       `mapstructure:"mcp"`
	UI        UIConfig        `mapstructure:"ui"`
	Cache     CacheConfig     `mapstructure:"cache"`
	Analytics AnalyticsConfig `mapstructure:"analytics"`
}

type APIConfig struct {
	BaseURL    string `mapstructure:"base_url"`
	Timeout    string `mapstructure:"timeout"`
	RetryCount int    `mapstructure:"retry_count"`
	RateLimit  int    `mapstructure:"rate_limit"`
}

type MCPConfig struct {
	Servers map[string]string `mapstructure:"servers"`
	Timeout string            `mapstructure:"timeout"`
}

type UIConfig struct {
	Theme      string `mapstructure:"theme"`
	Animations bool   `mapstructure:"animations"`
	ASCIIArt   bool   `mapstructure:"ascii_art"`
	Colors     bool   `mapstructure:"colors"`
}

type CacheConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	TTL     string `mapstructure:"ttl"`
	MaxSize int    `mapstructure:"max_size"`
}

type AnalyticsConfig struct {
	Enabled   bool `mapstructure:"enabled"`
	Anonymous bool `mapstructure:"anonymous"`
}

func Get() *Config {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Sprintf("Unable to decode config: %v", err))
	}
	return &cfg
}

func SetDefaults() {
	// API defaults
	viper.SetDefault("api.base_url", "https://api.antoine.ai")
	viper.SetDefault("api.timeout", "30s")
	viper.SetDefault("api.retry_count", 3)
	viper.SetDefault("api.rate_limit", 100)

	// MCP defaults
	viper.SetDefault("mcp.servers.exa", "mcp://localhost:8001")
	viper.SetDefault("mcp.servers.github", "mcp://localhost:8002")
	viper.SetDefault("mcp.servers.deepwiki", "mcp://localhost:8003")
	viper.SetDefault("mcp.servers.e2b", "mcp://localhost:8004")
	viper.SetDefault("mcp.servers.browserbase", "mcp://localhost:8005")
	viper.SetDefault("mcp.servers.firecrawl", "mcp://localhost:8006")
	viper.SetDefault("mcp.timeout", "30s")

	// UI defaults
	viper.SetDefault("ui.theme", "dark")
	viper.SetDefault("ui.animations", true)
	viper.SetDefault("ui.ascii_art", true)
	viper.SetDefault("ui.colors", true)

	// Cache defaults
	viper.SetDefault("cache.enabled", true)
	viper.SetDefault("cache.ttl", "30m")
	viper.SetDefault("cache.max_size", 1000)

	// Analytics defaults
	viper.SetDefault("analytics.enabled", true)
	viper.SetDefault("analytics.anonymous", true)
}

func Show() {
	cfg := Get()

	fmt.Println(ascii.GetBanner("Antoine Configuration", ascii.EmojiCode, 80))
	fmt.Println()

	// API Configuration
	fmt.Println("🌐 API Configuration:")
	fmt.Printf("  Base URL: %s\n", cfg.API.BaseURL)
	fmt.Printf("  Timeout: %s\n", cfg.API.Timeout)
	fmt.Printf("  Retry Count: %d\n", cfg.API.RetryCount)
	fmt.Printf("  Rate Limit: %d req/min\n", cfg.API.RateLimit)
	fmt.Println()

	// MCP Configuration
	fmt.Println("🔗 MCP Servers:")
	for name, endpoint := range cfg.MCP.Servers {
		status := "🔴 Disconnected"
		// En implementación real, verificar conexión
		if name == "exa" || name == "github" {
			status = "🟢 Connected"
		}
		fmt.Printf("  %s: %s %s\n", name, endpoint, status)
	}
	fmt.Println()

	// UI Configuration
	fmt.Println("🎨 UI Configuration:")
	fmt.Printf("  Theme: %s\n", cfg.UI.Theme)
	fmt.Printf("  Animations: %v\n", cfg.UI.Animations)
	fmt.Printf("  ASCII Art: %v\n", cfg.UI.ASCIIArt)
	fmt.Printf("  Colors: %v\n", cfg.UI.Colors)
	fmt.Println()

	// Cache Configuration
	fmt.Println("💾 Cache Configuration:")
	fmt.Printf("  Enabled: %v\n", cfg.Cache.Enabled)
	fmt.Printf("  TTL: %s\n", cfg.Cache.TTL)
	fmt.Printf("  Max Size: %d items\n", cfg.Cache.MaxSize)
	fmt.Println()

	// Analytics Configuration
	fmt.Println("📊 Analytics Configuration:")
	fmt.Printf("  Enabled: %v\n", cfg.Analytics.Enabled)
	fmt.Printf("  Anonymous: %v\n", cfg.Analytics.Anonymous)
	fmt.Println()

	fmt.Printf("Config file: %s\n", viper.ConfigFileUsed())
}

func Set(key, value string) error {
	viper.Set(key, value)

	// Guardar al archivo de configuración
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		// Crear archivo de configuración si no existe
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configFile = filepath.Join(home, ".antoine.yaml")
	}

	return viper.WriteConfigAs(configFile)
}

func Reset() error {
	SetDefaults()

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configFile := filepath.Join(home, ".antoine.yaml")
	return viper.WriteConfigAs(configFile)
}

func Validate() error {
	cfg := Get()

	// Validar configuración API
	if cfg.API.BaseURL == "" {
		return fmt.Errorf("api.base_url cannot be empty")
	}

	// Validar servidores MCP
	if len(cfg.MCP.Servers) == 0 {
		return fmt.Errorf("at least one MCP server must be configured")
	}

	return nil
}
