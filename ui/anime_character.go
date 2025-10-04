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
    "  ╭───────────╮  ",
    "  │  ◕    ●   │  ",  // one eye bigger
    "  │      ○    │  ",
    "  │  curious? │  ",
    "  ╰───────────╯  ",
    },
	"sleepy": {
    "  ╭───────────╮  ",
    "  │  ―    ―   │  ",
    "  │      ω    │  ",
    "  │   zzz...  │  ",
    "  ╰───────────╯  ",
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

func NewAnimeCharacter() *AnimeCharacter {
	return &AnimeCharacter{
		currentExpression: "idle",
		isAnimating:      false,
		animationFrame:   0,
		lastUpdate:       time.Now(),
	}
}

func (ac *AnimeCharacter) SetExpression(expr string) {
	if expr != ac.currentExpression {
		ac.currentExpression = expr
		ac.animationFrame = 0
		ac.lastUpdate = time.Now()
	}
}

func (ac *AnimeCharacter) GetCurrentFrame() []string {
	now := time.Now()
	
	// Update animation frame every 500ms
	if now.Sub(ac.lastUpdate) > 500*time.Millisecond {
		ac.Update()
		ac.lastUpdate = now
	}
	
	if ac.isAnimating && ac.currentExpression == "idle" {
		// Blinking animation
		frameIndex := ac.animationFrame % len(blinkFrames)
		return blinkFrames[frameIndex]
	}
	
	if frames, exists := expressions[ac.currentExpression]; exists {
		return frames
	}
	return expressions["idle"]
}

func (ac *AnimeCharacter) Update() {
	ac.animationFrame++
	
	// Random blinking when idle (5% chance each update)
	if ac.currentExpression == "idle" && rand.Intn(100) < 5 {
		ac.isAnimating = true
	}
	
	// Stop blinking animation after 3 frames
	if ac.isAnimating && ac.animationFrame > 3 {
		ac.isAnimating = false
		ac.animationFrame = 0
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
	default:
		color = "86"
	}
	
	characterStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true).
		Margin(0, 1)
		
	return characterStyle.Render(strings.Join(frame, "\n"))
}