package main

import (
	"log"
	"math/rand"
	"time"
	"fmt"

	"github.com/Praatibh/xang/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	asciiArt := `
░██    ░██    ░███    ░███    ░██   ░██████  
 ░██  ░██    ░██░██   ░████   ░██  ░██   ░██ 
  ░██░██    ░██  ░██  ░██░██  ░██ ░██        
   ░███    ░█████████ ░██ ░██ ░██ ░██  █████ 
  ░██░██   ░██    ░██ ░██  ░██░██ ░██     ██ 
 ░██  ░██  ░██    ░██ ░██   ░████  ░██  ░███ 
░██    ░██ ░██    ░██ ░██    ░███   ░█████░█

 Welcome to Xang! Ask me something...
`
	fmt.Print(asciiArt)

	rand.Seed(time.Now().UnixNano())

	input, err := ui.NewUIInput()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tea.NewProgram(ui.NewUi(input)).Run(); err != nil {
		log.Fatal(err)
	}
}