package git

import (
	"bytes"
	"strings"
	"os/exec"
)

func RunGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	cmd.Stderr = &out
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func Fetch() error {
	_, err := RunGit("fetch", "origin")
	return err
}
