// package run

// import (
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"strings"
// )

// func RunCommand(cmd string, arg ...string) (string, error) {
// 	out, err := exec.Command(cmd, arg...).Output()
// 	if err != nil {
// 		return fmt.Sprintf("error: %v", err), err
// 	}

// 	return string(out), nil
// }

// func PrepareInteractiveCommand(input string) *exec.Cmd {
// 	// Clean the input command
// 	cleanInput := strings.TrimSpace(strings.TrimRight(input, ";"))
	
// 	// Create the command with proper shell execution
// 	cmd := exec.Command("bash", "-c", cleanInput)
	
// 	// Set up for interactive execution
// 	cmd.Stdin = os.Stdin
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
	
// 	return cmd
// }

// func PrepareEditSettingsCommand(input string) *exec.Cmd {
// 	cleanInput := strings.TrimSpace(strings.TrimRight(input, ";"))
	
// 	cmd := exec.Command("bash", "-c", cleanInput)
	
// 	// Set up for interactive execution
// 	cmd.Stdin = os.Stdin
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
	
// 	return cmd
// }

package run

import (
    "fmt"
    "os/exec"
)

// RunCommand executes a command and returns its combined stdout and stderr output.
func RunCommand(name string, args ...string) (string, error) {
    cmd := exec.Command(name, args...)
    out, err := cmd.CombinedOutput()
    return string(out), err
}

// PrepareInteractiveCommand wraps cmdStr in a bash -c invocation that prints
// a leading and trailing newline around the command output.
func PrepareInteractiveCommand(cmdStr string) *exec.Cmd {
    fullCmd := fmt.Sprintf("echo \"\n\";%s; echo \"\n\";", cmdStr)
    return exec.Command("bash", "-c", fullCmd)
}

// PrepareEditSettingsCommand wraps cmdStr in a bash -c invocation that prints
// a trailing newline after the command completes.
func PrepareEditSettingsCommand(cmdStr string) *exec.Cmd {
    fullCmd := fmt.Sprintf("%s; echo \"\n\";", cmdStr)
    return exec.Command("bash", "-c", fullCmd)
}
