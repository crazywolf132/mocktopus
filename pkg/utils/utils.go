package utils

import (
	"os/exec"
	"strings"
)

func ExecCommand(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = "."

	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}
func ExecCommandLines(args ...string) ([]string, error) {
	output, err := ExecCommand(args...)
	if err != nil {
		return nil, err
	}

	// Split the output into lines and trim whitespace from each line.ExecCommandLines
	lines := strings.Split(string(output), "\n")
	trimmedLines := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			trimmedLines = append(trimmedLines, trimmedLine)
		}
	}

	return trimmedLines, nil
}

func Unpack[T interface{}](s []T, vars ...*T) {
	for i, res := range s {
		*vars[i] = res
	}
}

func StripFileFromPath(str string) string {
	parts := strings.Split(str, "/")
	return strings.Join(parts[:len(parts)-1], "/")
}
