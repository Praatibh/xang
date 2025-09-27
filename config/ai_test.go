package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAiConfig(t *testing.T) {
	t.Run("GetKey", testGetKey)
	t.Run("GetModel", testGetModel)
}

func testGetKey(t *testing.T) {
	expectedKey := "test_key"
	aiConfig := AiConfig{key: expectedKey}

	actualKey := aiConfig.GetKey()

	assert.Equal(t, expectedKey, actualKey, "The two keys should be the same.")
}

func testGetModel(t *testing.T) {
	expectedModel := "test_model"
	aiConfig := AiConfig{model: expectedModel}

	actualModel := aiConfig.GetModel()

	assert.Equal(t, expectedModel, actualModel, "The two models should be the same.")
}