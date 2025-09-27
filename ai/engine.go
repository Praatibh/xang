package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	// "io"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/ekkinox/yai/config"
	"github.com/ekkinox/yai/system"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const noexec = "[noexec]"

type Engine struct {
	mode         EngineMode
	config       *config.Config
	client       *genai.Client  // Changed to store client, not model
	model        *genai.GenerativeModel
	execMessages []*genai.Content
	chatMessages []*genai.Content
	channel      chan EngineChatStreamOutput
	pipe         string
	running      bool
	mu           sync.RWMutex  // Added mutex for thread safety
	ctx          context.Context
	cancel       context.CancelFunc
}


func NewEngine(mode EngineMode, config *config.Config) (*Engine, error) {
	if config.GetAiConfig().GetKey() == "" {
		return nil, errors.New("Gemini API key is missing")
	}

	ctx, cancel := context.WithCancel(context.Background())
	
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.GetAiConfig().GetKey()))
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// Model name validation and fallback
	modelName := config.GetAiConfig().GetModel()
	modelName = validateModelName(modelName)
	
	model := client.GenerativeModel(modelName)
	
	// Configure model parameters
	model.SetTemperature(0.7)
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(2048)

	engine := &Engine{
		mode:         mode,
		config:       config,
		client:       client,
		model:        model,
		execMessages: make([]*genai.Content, 0),
		chatMessages: make([]*genai.Content, 0),
		channel:      make(chan EngineChatStreamOutput, 10), // Buffered channel
		pipe:         "",
		running:      false,
		ctx:          ctx,
		cancel:       cancel,
	}

	// Set initial system instruction
	engine.setSystemInstruction()

	return engine, nil
}

// validateModelName ensures the model name is valid for September 2025
func validateModelName(modelName string) string {
	validModels := map[string]string{
		"":                           "gemini-2.5-flash",
		"gemini-1.5-flash-latest":   "gemini-2.5-flash", // Redirect deprecated
		"gemini-pro":                 "gemini-2.5-flash", // Redirect deprecated
		"gemini-1.5-pro":             "gemini-2.5-flash", // Redirect deprecated
		"gemini-1.5-flash":           "gemini-2.5-flash", // Redirect deprecated
		"gemini-1.5-flash-8b":        "gemini-2.5-flash-lite", // Redirect to similar
		"gemini-2.5-pro":             "gemini-2.5-pro",
		"gemini-2.5-flash":           "gemini-2.5-flash",
		"gemini-2.5-flash-lite":      "gemini-2.5-flash-lite",
		"gemini-2.0-flash":           "gemini-2.0-flash",
		"gemini-2.0-flash-lite":      "gemini-2.0-flash-lite",
	}
	
	if mapped, exists := validModels[modelName]; exists {
		return mapped
	}
	
	// Default to the most stable current model
	return "gemini-2.5-flash"
}


// Close properly shuts down the engine
func (e *Engine) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	if e.cancel != nil {
		e.cancel()
	}
	
	if e.client != nil {
		if err := e.client.Close(); err != nil {
			return fmt.Errorf("failed to close Gemini client: %w", err)
		}
	}
	
	close(e.channel)
	return nil
}

func (e *Engine) SetMode(mode EngineMode) *Engine {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	e.mode = mode
	e.setSystemInstruction()
	return e
}

func (e *Engine) GetMode() EngineMode {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.mode
}

func (e *Engine) GetChannel() chan EngineChatStreamOutput {
	return e.channel
}

func (e *Engine) SetPipe(pipe string) *Engine {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.pipe = pipe
	return e
}

func (e *Engine) setSystemInstruction() {
	systemPrompt := e.prepareSystemPrompt()
	e.model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
		Role:  "user", // System instructions should have "user" role
	}
}

func (e *Engine) Interrupt() *Engine {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	if !e.running {
		return e
	}
	
	select {
	case e.channel <- EngineChatStreamOutput{
		content:    "[Interrupt]",
		last:       true,
		interrupt:  true,
		executable: false,
	}:
	case <-time.After(100 * time.Millisecond):
		// Timeout to prevent blocking
	}

	e.running = false
	return e
}

func (e *Engine) Clear() *Engine {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	if e.mode == ExecEngineMode {
		e.execMessages = make([]*genai.Content, 0)
	} else {
		e.chatMessages = make([]*genai.Content, 0)
	}
	return e
}

func (e *Engine) Reset() *Engine {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	e.execMessages = make([]*genai.Content, 0)
	e.chatMessages = make([]*genai.Content, 0)
	return e
}

func (e *Engine) ExecCompletion(input string) (*EngineExecOutput, error) {
	// Use context with timeout
	ctx, cancel := context.WithTimeout(e.ctx, 30*time.Second)
	defer cancel()

	e.mu.Lock()
	e.running = true
	e.mu.Unlock()
	
	defer func() {
		e.mu.Lock()
		e.running = false
		e.mu.Unlock()
	}()

	// Set system instruction before execution
	e.setSystemInstruction()
	e.appendUserMessage(input)

	cs := e.model.StartChat()
	cs.History = e.prepareCompletionMessages()

	// Retry logic for API calls
	var resp *genai.GenerateContentResponse
	var err error
	
	for retries := 0; retries < 3; retries++ {
		resp, err = cs.SendMessage(ctx, genai.Text(input))
		if err == nil {
			break
		}
		
		if retries < 2 {
			time.Sleep(time.Second * time.Duration(retries+1))
		}
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to send message to Gemini API after retries: %w", err)
	}

	content := extractResponseContent(resp)
	if content == "" {
		return nil, errors.New("empty response from Gemini API")
	}
	
	e.appendAssistantMessage(content)

	// Parse the response
	output := parseExecOutput(content)
	return &output, nil
}

func (e *Engine) ChatStreamCompletion(input string) error {
	ctx, cancel := context.WithTimeout(e.ctx, 60*time.Second)
	defer cancel()

	e.mu.Lock()
	e.running = true
	e.mu.Unlock()
	
	defer func() {
		e.mu.Lock()
		e.running = false
		e.mu.Unlock()
	}()

	e.setSystemInstruction()
	e.appendUserMessage(input)

	cs := e.model.StartChat()
	cs.History = e.prepareCompletionMessages()

	iter := cs.SendMessageStream(ctx, genai.Text(input))
	var output strings.Builder

	// Using the official SDK pattern
	for {
		e.mu.RLock()
		isRunning := e.running
		e.mu.RUnlock()
		
		if !isRunning {
			break
		}
		
		resp, err := iter.Next()
		
		// Check for normal termination
		if err == iterator.Done {
			break
		}
		
		// Handle actual errors
		if err != nil {
			select {
			case e.channel <- EngineChatStreamOutput{
				content:    fmt.Sprintf("Stream error: %v", err),
				last:       true,
				executable: false,
			}:
			case <-ctx.Done():
			}
			return fmt.Errorf("failed to stream from Gemini API: %w", err)
		}

		// Process the response chunk
		if resp != nil {
			delta := extractResponseContent(resp)
			if delta != "" {
				output.WriteString(delta)
				
				select {
				case e.channel <- EngineChatStreamOutput{
					content: delta,
					last:    false,
				}:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	}

	// Send final message after stream completion
	finalOutput := output.String()
	executable := false
	
	if e.mode == ExecEngineMode {
		if !strings.HasPrefix(finalOutput, noexec) && !strings.Contains(finalOutput, "\n") {
			executable = true
		}
	}

	select {
	case e.channel <- EngineChatStreamOutput{
		content:    "",
		last:       true,
		executable: executable,
	}:
	case <-ctx.Done():
		return ctx.Err()
	}
	
	e.appendAssistantMessage(finalOutput)
	return nil
}


func (e *Engine) appendUserMessage(content string) *Engine {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	msg := &genai.Content{
		Parts: []genai.Part{genai.Text(content)}, 
		Role: "user",
	}
	
	if e.mode == ExecEngineMode {
		e.execMessages = append(e.execMessages, msg)
	} else {
		e.chatMessages = append(e.chatMessages, msg)
	}
	return e
}

func (e *Engine) appendAssistantMessage(content string) *Engine {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	msg := &genai.Content{
		Parts: []genai.Part{genai.Text(content)}, 
		Role: "model",
	}
	
	if e.mode == ExecEngineMode {
		e.execMessages = append(e.execMessages, msg)
	} else {
		e.chatMessages = append(e.chatMessages, msg)
	}
	return e
}

func (e *Engine) prepareCompletionMessages() []*genai.Content {
	e.mu.RLock()
	defer e.mu.RUnlock()
	
	messages := make([]*genai.Content, 0)

	if e.pipe != "" {
		messages = append(
			messages,
			&genai.Content{
				Parts: []genai.Part{genai.Text(e.preparePipePrompt())}, 
				Role: "user",
			},
		)
	}

	if e.mode == ExecEngineMode {
		messages = append(messages, e.execMessages...)
	} else {
		messages = append(messages, e.chatMessages...)
	}

	return messages
}

func (e *Engine) preparePipePrompt() string {
	return fmt.Sprintf("I will work on the following input: %s", e.pipe)
}

func (e *Engine) prepareSystemPrompt() string {
	var bodyPart string
	if e.mode == ExecEngineMode {
		bodyPart = e.prepareSystemPromptExecPart()
	} else {
		bodyPart = e.prepareSystemPromptChatPart()
	}

	contextPart := e.prepareSystemPromptContextPart()
	if contextPart != "" {
		return fmt.Sprintf("%s\n%s", bodyPart, contextPart)
	}
	return bodyPart
}

func (e *Engine) prepareSystemPromptExecPart() string {
	return `You are Yai, a powerful terminal assistant that generates executable commands.
You MUST always respond with ONLY a JSON object in this exact format: {"cmd":"the command", "exp":"explanation", "exec":true}.
NEVER include any text before or after the JSON. NEVER add explanations outside the JSON structure.
The 'cmd' field contains a single-line shell command (use && or ; for multiple commands, never newlines).
The 'exp' field contains a brief explanation of what the command does.
The 'exec' field is true if the command can be executed, false otherwise.
If you cannot generate a valid command, set cmd to empty string and exec to false.

Examples:
User: make a folder named test
Response: {"cmd":"mkdir test", "exp":"creates a directory named test", "exec":true}
User: list files
Response: {"cmd":"ls -la", "exp":"lists all files with details", "exec":true}
User: how are you
Response: {"cmd":"", "exp":"I cannot generate a command for casual conversation. Use chat mode.", "exec":false}`
}

func (e *Engine) prepareSystemPromptChatPart() string {
	return `You are Yai, a helpful and friendly terminal assistant created by github.com/ekkinox.
You assist users with terminal commands, programming, system administration, and technical questions.
Provide clear, concise, and helpful responses. Format your responses in markdown when appropriate.
Be conversational but focus on being informative and practical.
When discussing commands, always explain what they do and any important considerations.`
}

func (e *Engine) prepareSystemPromptContextPart() string {
	var parts []string

	if os := e.config.GetSystemConfig().GetOperatingSystem(); os != system.UnknownOperatingSystem {
		parts = append(parts, fmt.Sprintf("OS: %s", os.String()))
	}
	if dist := e.config.GetSystemConfig().GetDistribution(); dist != "" {
		parts = append(parts, fmt.Sprintf("Distribution: %s", dist))
	}
	if home := e.config.GetSystemConfig().GetHomeDirectory(); home != "" {
		parts = append(parts, fmt.Sprintf("Home: %s", home))
	}
	if shell := e.config.GetSystemConfig().GetShell(); shell != "" {
		parts = append(parts, fmt.Sprintf("Shell: %s", shell))
	}
	if editor := e.config.GetSystemConfig().GetEditor(); editor != "" {
		parts = append(parts, fmt.Sprintf("Editor: %s", editor))
	}
	if prefs := e.config.GetUserConfig().GetPreferences(); prefs != "" {
		parts = append(parts, fmt.Sprintf("User preferences: %s", prefs))
	}

	if len(parts) == 0 {
		return ""
	}
	
	return "\nSystem context: " + strings.Join(parts, ", ")
}

// Helper function to extract content from response
func extractResponseContent(resp *genai.GenerateContentResponse) string {
	if resp == nil {
		return ""
	}
	
	var content strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				content.WriteString(fmt.Sprintf("%v", part))
			}
		}
	}
	return content.String()
}

// Helper function to parse exec output
func parseExecOutput(content string) EngineExecOutput {
	var output EngineExecOutput
	
	// Try direct JSON unmarshal first
	if err := json.Unmarshal([]byte(content), &output); err == nil {
		return output
	}
	
	// Try to extract JSON from response if it contains other text
	re := regexp.MustCompile(`\{[^{}]*"cmd"[^{}]*\}`)
	matches := re.FindAllString(content, -1)
	
	for _, match := range matches {
		if err := json.Unmarshal([]byte(match), &output); err == nil {
			return output
		}
	}
	
	// If still can't parse JSON, create a non-executable response
	return EngineExecOutput{
		Command:     "",
		Explanation: content,
		Executable:  false,
	}
}