package config

import (
	"fmt"
	"strings"

	"github.com/ekkinox/yai/system"
	"github.com/spf13/viper"
)

type Config struct {
	ai     AiConfig
	user   UserConfig
	system *system.Analysis
}

func (c *Config) GetAiConfig() AiConfig {
	return c.ai
}

func (c *Config) GetUserConfig() UserConfig {
	return c.user
}

func (c *Config) GetSystemConfig() *system.Analysis {
	return c.system
}

func NewConfig() (*Config, error) {
	system := system.Analyse()

	viper.SetConfigName(strings.ToLower(system.GetApplicationName()))
	viper.AddConfigPath(fmt.Sprintf("%s/.config/", system.GetHomeDirectory()))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return &Config{
		ai: AiConfig{
			key:   viper.GetString(gemini_key),
			model: viper.GetString(gemini_model),
		},
		user: UserConfig{
			defaultPromptMode: viper.GetString(user_default_prompt_mode),
			preferences:       viper.GetString(user_preferences),
		},
		system: system,
	}, nil
}

func WriteConfig(key string, write bool) (*Config, error) {
	if key == "" {
		return nil, fmt.Errorf("API key cannot be empty")
	}

	system := system.Analyse()

	// ai defaults - use the most stable model name
	viper.Set(gemini_key, key)
	viper.Set(gemini_model, "gemini-1.5-flash")

	// user defaults
	viper.SetDefault(user_default_prompt_mode, "exec")
	viper.SetDefault(user_preferences, "")

	if write {
		err := viper.WriteConfigAs(system.GetConfigFile())
		if err != nil {
			return nil, fmt.Errorf("failed to write config file: %w", err)
		}
	}

	return NewConfig()
}