package run

import (
	"context"
	"strings"
	"testing"
)

func TestExecuteSuccess(t *testing.T) {
	script := &Script{
		ID:    "test",
		Title: "Test",
		Steps: []Step{
			{
				Command:        "echo 'hello world'",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"hello world"},
			},
		},
	}

	results, err := Execute(context.Background(), script, nil)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Passed {
		t.Fatalf("expected step to pass, failures: %v", results[0].Failures)
	}
}

func TestExecuteFailureOnMissingStdout(t *testing.T) {
	script := &Script{
		ID:    "fail",
		Title: "Fail",
		Steps: []Step{
			{
				Command:        "echo 'mismatch'",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"expected"},
			},
		},
	}

	results, err := Execute(context.Background(), script, nil)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Passed {
		t.Fatalf("expected step to fail")
	}
	if !contains(results[0].Failures, `stdout missing "expected"`) {
		t.Fatalf("missing failure about stdout: %v", results[0].Failures)
	}
}

func contains(slice []string, target string) bool {
	for _, item := range slice {
		if strings.Contains(item, target) {
			return true
		}
	}
	return false
}
