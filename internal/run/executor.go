package run

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"
)

// StepResult captures the outcome of a single scripted command.
type StepResult struct {
	Step      Step
	Stdout    string
	Stderr    string
	ExitCode  int
	Duration  time.Duration
	Passed    bool
	Failures  []string
	ExecError error
}

// Execute runs every step within the script sequentially.
func Execute(ctx context.Context, script *Script, emitter func(StepResult)) ([]StepResult, error) {
	if script == nil {
		return nil, errors.New("run: script is nil")
	}

	results := make([]StepResult, 0, len(script.Steps))

	for _, step := range script.Steps {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}

		res := RunStep(ctx, step)

		results = append(results, res)
		if emitter != nil {
			emitter(res)
		}
	}

	return results, nil
}

func runCommand(ctx context.Context, command string) (string, string, int, error) {
	shell, args := defaultShell()
	args = append(args, command)
	cmd := exec.CommandContext(ctx, shell, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := exitStatus(err)
	return stdout.String(), stderr.String(), exitCode, err
}

// RunStep executes a single Step and returns the result.
func RunStep(ctx context.Context, step Step) StepResult {
	start := time.Now()
	stdout, stderr, exitCode, execErr := runCommand(ctx, step.Command)
	duration := time.Since(start)

	res := StepResult{
		Step:      step,
		Stdout:    stdout,
		Stderr:    stderr,
		ExitCode:  exitCode,
		Duration:  duration,
		ExecError: execErr,
	}

	if execErr == nil {
		if step.ExpectExitCode >= 0 && step.ExpectExitCode != exitCode {
			res.Failures = append(res.Failures, fmt.Sprintf("expected exit code %d, got %d", step.ExpectExitCode, exitCode))
		}
	} else {
		res.Failures = append(res.Failures, execErr.Error())
	}

	for _, expected := range step.ExpectStdout {
		if !strings.Contains(res.Stdout, expected) {
			res.Failures = append(res.Failures, fmt.Sprintf("stdout missing %q", expected))
		}
	}

	res.Passed = len(res.Failures) == 0
	return res
}

func defaultShell() (string, []string) {
	if runtime.GOOS == "windows" {
		if _, err := exec.LookPath("pwsh"); err == nil {
			return "pwsh", []string{"-NoLogo", "-NoProfile", "-Command"}
		}
		return "powershell", []string{"-NoLogo", "-NoProfile", "-Command"}
	}
	return "sh", []string{"-c"}
}

func exitStatus(err error) int {
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return -1
}
