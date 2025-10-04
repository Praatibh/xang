package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAiConfig(t *testing.T) {
	t.Run("GetKey", func(t *testing.T) {
		aiConfig := &AiConfig{key: "test-key"}
		assert.Equal(t, "test-key", aiConfig.GetKey())
	})

	t.Run("GetModel", func(t *testing.T) {
		aiConfig := &AiConfig{model: "test-model"}
		assert.Equal(t, "test-model", aiConfig.GetModel())
	})
}

func TestConfig(t *testing.T) {
	// Setup: Create a temporary config file for testing
	setupTestConfig := func(t *testing.T) (string, func()) {
		// Create temporary directory
		tmpDir, err := os.MkdirTemp("", "xang-test-*")
		assert.NoError(t, err)

		// Create config file
		configPath := filepath.Join(tmpDir, "xang.yaml")
		configContent := `ai:
  key: "test-api-key"
  model: "gemini-1.5-flash"
user:
  default_prompt_mode: "exec"
  preferences: []
system:
  editor: "nano"
`
		err = os.WriteFile(configPath, []byte(configContent), 0644)
		assert.NoError(t, err)

		// Set config path environment variable
		oldConfigPath := os.Getenv("XDG_CONFIG_HOME")
		os.Setenv("XDG_CONFIG_HOME", tmpDir)

		// Return cleanup function
		cleanup := func() {
			os.Setenv("XDG_CONFIG_HOME", oldConfigPath)
			os.RemoveAll(tmpDir)
		}

		return tmpDir, cleanup
	}

	t.Run("NewConfig", func(t *testing.T) {
		tmpDir, cleanup := setupTestConfig(t)
		defer cleanup()

		config, err := NewConfig()
		assert.NoError(t, err)
		assert.NotNil(t, config)
		
		// Verify config was loaded correctly
		assert.Equal(t, "test-api-key", config.GetAiConfig().GetKey())
		assert.Equal(t, "gemini-1.5-flash", config.GetAiConfig().GetModel())
		
		// Verify config file path
		expectedPath := filepath.Join(tmpDir, "xang.yaml")
		assert.Equal(t, expectedPath, config.GetSystemConfig().GetConfigFile())
	})

	t.Run("WriteConfig", func(t *testing.T) {
		tmpDir, cleanup := setupTestConfig(t)
		defer cleanup()

		// Write new config
		config, err := WriteConfig("new-api-key", true)
		assert.NoError(t, err)
		assert.NotNil(t, config)
		
		// Verify the key was updated
		assert.Equal(t, "new-api-key", config.GetAiConfig().GetKey())
		
		// Verify file was written
		configPath := filepath.Join(tmpDir, "xang.yaml")
		_, err = os.Stat(configPath)
		assert.NoError(t, err)
		
		// Read config again to verify persistence
		config2, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "new-api-key", config2.GetAiConfig().GetKey())
	})
}

func TestUserConfig(t *testing.T) {
	t.Run("GetDefaultPromptMode", func(t *testing.T) {
		userConfig := &UserConfig{defaultPromptMode: "exec"}
		assert.Equal(t, "exec", userConfig.GetDefaultPromptMode())
	})

	t.Run("GetPreferences", func(t *testing.T) {
		prefs := []string{"pref1", "pref2"}
		userConfig := &UserConfig{preferences: prefs}
		assert.Equal(t, prefs, userConfig.GetPreferences())
	})
}