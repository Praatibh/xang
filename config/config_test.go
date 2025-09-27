package config

import (
	"os"
	"strings"
	"testing"

	"github.com/ekkinox/yai/system"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("NewConfig", testNewConfig)
	t.Run("WriteConfig", testWriteConfig)
}

func setupViper(t *testing.T) {
	t.Helper()
	system := system.Analyse()

	viper.SetConfigName(strings.ToLower(system.GetApplicationName()))
	viper.AddConfigPath("/tmp/")
	viper.Set(gemini_key, "test_key")
	viper.Set(gemini_model, "gemini-pro")
	viper.Set(user_default_prompt_mode, "exec")
	viper.Set(user_preferences, "test_preferences")

	require.NoError(t, viper.SafeWriteConfigAs("/tmp/yai.json"))
}

func cleanup(t *testing.T) {
	t.Helper()
	require.NoError(t, os.Remove("/tmp/yai.json"))
}

func testNewConfig(t *testing.T) {
	setupViper(t)
	defer cleanup(t)

	cfg, err := NewConfig()
	require.NoError(t, err)

	assert.Equal(t, "test_key", cfg.GetAiConfig().GetKey())
	assert.Equal(t, "gemini-pro", cfg.GetAiConfig().GetModel())
	assert.Equal(t, "exec", cfg.GetUserConfig().GetDefaultPromptMode())
	assert.Equal(t, "test_preferences", cfg.GetUserConfig().GetPreferences())

	assert.NotNil(t, cfg.GetSystemConfig())
}

func testWriteConfig(t *testing.T) {
	setupViper(t)
	defer cleanup(t)

	cfg, err := WriteConfig("new_test_key", false)
	require.NoError(t, err)

	assert.Equal(t, "new_test_key", cfg.GetAiConfig().GetKey())
	assert.Equal(t, "gemini-pro", cfg.GetAiConfig().GetModel())
	assert.Equal(t, "exec", cfg.GetUserConfig().GetDefaultPromptMode())
	assert.Equal(t, "test_preferences", cfg.GetUserConfig().GetPreferences())

	assert.NotNil(t, cfg.GetSystemConfig())

	assert.Equal(t, "new_test_key", viper.GetString(gemini_key))
	assert.Equal(t, "gemini-pro", viper.GetString(gemini_model))
	assert.Equal(t, "exec", viper.GetString(user_default_prompt_mode))
	assert.Equal(t, "test_preferences", viper.GetString(user_preferences))
}