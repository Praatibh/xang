package ui

import (
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// AnimeCharacter represents the AI assistant's visual state
type AnimeCharacter struct {
	currentExpression string
	isAnimating      bool
	animationFrame   int
	lastUpdate       time.Time
	idleCounter      int // Track how long we've been idle
}

// Different expressions for the anime character
var expressions = map[string][]string{
	"idle": {
		"  ╭─────────╮  ",
		"  │  ◕   ◕  │  ",
		"  │     ‿   │  ",
		"  │   (xang)│  ",
		"  ╰─────────╯  ",
	},
	"thinking": {
		"  ╭─────────╮  ",
		"  │  ◔   ◔  │  ",
		"  │     ~   │  ",
		"  │thinking.│  ",
		"  ╰─────────╯  ",
	},
	"happy": {
		"  ╭─────────╮  ",
		"  │  ★   ★  │  ",
		"  │     ‿   │  ",
		"  │ excited!│  ",
		"  ╰─────────╯  ",
	},
	"error": {
		"  ╭─────────╮  ",
		"  │  x   x  │  ",
		"  │     ―   │  ",
		"  │  oops.. │  ",
		"  ╰─────────╯  ",
	},
	"processing": {
		"  ╭─────────╮  ",
		"  │  ●   ●  │  ",
		"  │     ―   │  ",
		"  │analyzing│  ",
		"  ╰─────────╯  ",
	},
	"success": {
		"  ╭─────────╮  ",
		"  │  ◕   ◕  │  ",
		"  │     ‿   │  ",
		"  │  done!  │  ",
		"  ╰─────────╯  ",
	},
	"curious": {
		"  ╭─────────╮  ",
		"  │  ◕   ●  │  ",
		"  │     ○   │  ",
		"  │curious? │  ",
		"  ╰─────────╯  ",
	},
	"sleepy": {
		"  ╭─────────╮  ",
		"  │  ―   ―  │  ",
		"  │     ω   │  ",
		"  │  zzz... │  ",
		"  ╰─────────╯  ",
	},
	"confused": {
		"  ╭─────────╮  ",
		"  │  ◔   ◕  │  ",
		"  │     ?   │  ",
		"  │  huh?   │  ",
		"  ╰─────────╯  ",
	},
	"working": {
		"  ╭─────────╮  ",
		"  │  ●   ●  │  ",
		"  │     ▽   │  ",
		"  │ working │  ",
		"  ╰─────────╯  ",
	},
	"celebrating": {
		"  ╭─────────╮  ",
		"  │  ☆   ☆  │  ",
		"  │     ▿   │  ",
		"  │  yay!   │  ",
		"  ╰─────────╯  ",
	},
}

// Blinking animation frames
var blinkFrames = [][]string{
	{
		"  ╭─────────╮  ",
		"  │  ◕   ◕  │  ",
		"  │     ‿   │  ",
		"  │   (xang)│  ",
		"  ╰─────────╯  ",
	},
	{
		"  ╭─────────╮  ",
		"  │  ―   ―  │  ",
		"  │     ‿   │  ",
		"  │   (xang)│  ",
		"  ╰─────────╯  ",
	},
}

// Thinking animation frames (dots cycling)
var thinkingFrames = [][]string{
	{
		"  ╭─────────╮  ",
		"  │  ◔   ◔  │  ",
		"  │     ~   │  ",
		"  │thinking.│  ",
		"  ╰─────────╯  ",
	},
	{
		"  ╭─────────╮  ",
		"  │  ◔   ◔  │  ",
		"  │     ~   │  ",
		"  │thinking..│ ",
		"  ╰─────────╯  ",
	},
	{
		"  ╭─────────╮  ",
		"  │  ◔   ◔  │  ",
		"  │     ~   │  ",
		"  │thinking...│",
		"  ╰─────────╯  ",
	},
}

// Processing animation frames (eyes moving)
var processingFrames = [][]string{
	{
		"  ╭─────────╮  ",
		"  │  ●   ●  │  ",
		"  │     ―   │  ",
		"  │analyzing│  ",
		"  ╰─────────╯  ",
	},
	{
		"  ╭─────────╮  ",
		"  │  ◐   ◐  │  ",
		"  │     ―   │  ",
		"  │analyzing│  ",
		"  ╰─────────╯  ",
	},
	{
		"  ╭─────────╮  ",
		"  │  ◑   ◑  │  ",
		"  │     ―   │  ",
		"  │analyzing│  ",
		"  ╰─────────╯  ",
	},
	{
		"  ╭─────────╮  ",
		"  │  ◒   ◒  │  ",
		"  │     ―   │  ",
		"  │analyzing│  ",
		"  ╰─────────╯  ",
	},
}

func NewAnimeCharacter() *AnimeCharacter {
	return &AnimeCharacter{
		currentExpression: "idle",
		isAnimating:      false,
		animationFrame:   0,
		lastUpdate:       time.Now(),
		idleCounter:      0,
	}
}

func (ac *AnimeCharacter) SetExpression(expr string) {
	if expr != ac.currentExpression {
		ac.currentExpression = expr
		ac.animationFrame = 0
		ac.lastUpdate = time.Now()
		ac.idleCounter = 0
	}
}

func (ac *AnimeCharacter) GetCurrentFrame() []string {
	now := time.Now()
	
	// Update animation frame every 500ms
	if now.Sub(ac.lastUpdate) > 500*time.Millisecond {
		ac.Update()
		ac.lastUpdate = now
	}
	
	// Handle different animated states
	switch ac.currentExpression {
	case "idle":
		if ac.isAnimating {
			// Blinking animation
			frameIndex := ac.animationFrame % len(blinkFrames)
			return blinkFrames[frameIndex]
		}
	case "thinking":
		// Animated thinking dots
		frameIndex := ac.animationFrame % len(thinkingFrames)
		return thinkingFrames[frameIndex]
	case "processing":
		// Animated processing eyes
		frameIndex := ac.animationFrame % len(processingFrames)
		return processingFrames[frameIndex]
	}
	
	if frames, exists := expressions[ac.currentExpression]; exists {
		return frames
	}
	return expressions["idle"]
}

func (ac *AnimeCharacter) Update() {
	ac.animationFrame++
	
	switch ac.currentExpression {
	case "idle":
		ac.idleCounter++
		
		// Random blinking when idle (5% chance each update)
		if rand.Intn(100) < 5 {
			ac.isAnimating = true
		}
		
		// Stop blinking animation after 3 frames
		if ac.isAnimating && ac.animationFrame > 3 {
			ac.isAnimating = false
			ac.animationFrame = 0
		}
		
		// After being idle for a while (30 seconds = 60 updates), get sleepy
		if ac.idleCounter > 60 {
			ac.SetExpression("sleepy")
		}
	
	case "thinking", "processing":
		// These states continuously animate
		// Animation frames loop automatically via modulo in GetCurrentFrame
		
	case "sleepy":
		// Random chance to wake up and go back to idle
		if rand.Intn(100) < 2 {
			ac.SetExpression("idle")
		}
	}
}

func (ac *AnimeCharacter) Render() string {
	frame := ac.GetCurrentFrame()
	
	var color string
	switch ac.currentExpression {
	case "idle":
		color = "86"  // Green
	case "thinking":
		color = "33"  // Blue
	case "happy":
		color = "226" // Yellow
	case "error":
		color = "196" // Red
	case "processing":
		color = "129" // Purple
	case "success":
		color = "82"  // Bright Green
	case "curious":
		color = "51"  // Cyan
	case "sleepy":
		color = "242" // Gray
	case "confused":
		color = "208" // Orange
	case "working":
		color = "141" // Light Purple
	case "celebrating":
		color = "201" // Magenta
	default:
		color = "86"
	}
	
	characterStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true).
		Margin(0, 1)
		
	return characterStyle.Render(strings.Join(frame, "\n"))
}