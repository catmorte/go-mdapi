package command

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func RunCommand(command string) (string, error) {
	// Execute the command and capture its output
	cmd := exec.Command("bash", "-c", command) // Use bash to handle piping and stderr
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
