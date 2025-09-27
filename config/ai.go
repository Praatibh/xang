package config

const (
	gemini_key         = "GEMINI_KEY"
	gemini_model       = "GEMINI_MODEL"
)

type AiConfig struct {
	key         string
	model       string
}

func (c AiConfig) GetKey() string {
	return c.key
}

func (c AiConfig) GetModel() string {
	return c.model
}