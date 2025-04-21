package command

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func RunCommand(command string) (string, error) {
	// Split the command into individual arguments
	args := strings.Fields(command)

	// Create a new command with the args and use bash for interpreting them
	cmd := exec.Command("bash", "-c", strings.Join(args, " "))
	// cmd := exec.Command("bash", "-c", command) // Use bash to handle piping and stderr
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	// Run the command
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run command: %v, stderr: %s", err, errOut.String())
	}

	// Return the output
	return strings.TrimRight(out.String(), "\n"), nil
}
