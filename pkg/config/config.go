package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Firestore     FirestoreConfig     `mapstructure:"firestore"`
	Server        ServerConfig        `mapstructure:"server"`
	Scraping      ScrapingConfig      `mapstructure:"scraping"`
	AI            AIConfig            `mapstructure:"ai"`
	Notifications NotificationsConfig `mapstructure:"notifications"`
	Logging       LoggingConfig       `mapstructure:"logging"`
}

type FirestoreConfig struct {
	ProjectID       string `mapstructure:"project_id"`
	CredentialsFile string `mapstructure:"credentials_file"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type ScrapingConfig struct {
	DefaultDelay string `mapstructure:"default_delay"`
	MaxRetries   int    `mapstructure:"max_retries"`
	Timeout      string `mapstructure:"timeout"`
	UserAgent    string `mapstructure:"user_agent"`
}

type AIConfig struct {
	Provider      string `mapstructure:"provider"`
	Model         string `mapstructure:"model"`
	APIKeyEnv     string `mapstructure:"api_key_env"`
	GeminiAPIKey  string `mapstructure:"-"` // Populated from environment
	MaxTokens     int    `mapstructure:"max_tokens"`
}

type NotificationsConfig struct {
	SMTPHost        string `mapstructure:"smtp_host"`
	SMTPPort        int    `mapstructure:"smtp_port"`
	SMTPUsername    string `mapstructure:"smtp_username"`
	SMTPPasswordEnv string `mapstructure:"smtp_password_env"`
	FromEmail       string `mapstructure:"from_email"`
	ToEmail         string `mapstructure:"to_email"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/lookie")

	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("firestore.project_id", "lookie-quantum-intelligence")
	viper.SetDefault("firestore.credentials_file", "./config/lookie-dev-key.json")
	viper.SetDefault("scraping.default_delay", "10s")
	viper.SetDefault("scraping.max_retries", 3)
	viper.SetDefault("scraping.timeout", "30s")
	viper.SetDefault("ai.provider", "gemini")
	viper.SetDefault("ai.model", "gemini-pro")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %w", err)
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate required environment variables and populate from env
	if config.AI.APIKeyEnv != "" {
		apiKey := os.Getenv(config.AI.APIKeyEnv)
		if apiKey == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", config.AI.APIKeyEnv)
		}
		config.AI.GeminiAPIKey = apiKey
	}

	if config.Notifications.SMTPPasswordEnv != "" {
		if os.Getenv(config.Notifications.SMTPPasswordEnv) == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", config.Notifications.SMTPPasswordEnv)
		}
	}

	return &config, nil
}