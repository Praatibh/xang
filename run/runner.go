package run

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunCommand(cmd string, arg ...string) (string, error) {
	out, err := exec.Command(cmd, arg...).Output()
	if err != nil {
		return fmt.Sprintf("error: %v", err), err
	}

	return string(out), nil
}

func PrepareInteractiveCommand(input string) *exec.Cmd {
	// Clean the input command
	cleanInput := strings.TrimSpace(strings.TrimRight(input, ";"))
	
	// Create the command with proper shell execution
	cmd := exec.Command("bash", "-c", cleanInput)
	
	// Set up for interactive execution
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd
}

func PrepareEditSettingsCommand(input string) *exec.Cmd {
	cleanInput := strings.TrimSpace(strings.TrimRight(input, ";"))
	
	cmd := exec.Command("bash", "-c", cleanInput)
	
	// Set up for interactive execution
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd
}