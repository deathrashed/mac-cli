package services

import (
	"context"
	"os"
	"os/exec"
)

func GumAvailable() bool {
	_, err := exec.LookPath("gum")
	return err == nil
}

func GumConfirm(ctx context.Context, message string) (bool, error) {
	if !GumAvailable() {
		return true, nil
	}
	cmd := exec.CommandContext(ctx, "gum", "confirm", message)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok && exit.ExitCode() == 1 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
