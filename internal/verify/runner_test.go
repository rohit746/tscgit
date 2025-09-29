package verify

import (
	"context"
	"errors"
	"testing"

	"github.com/rohit746/tscgit/internal/gitutil"
	"github.com/rohit746/tscgit/internal/lessons"
)

type testEmitter struct {
	results []Result
}

func (t *testEmitter) Emit(res Result) {
	t.results = append(t.results, res)
}

func TestRunInvokesChecks(t *testing.T) {
	lesson := &lessons.Lesson{
		ID:    "test",
		Title: "Test Lesson",
		Checks: []lessons.Check{
			{
				ID:    "check-1",
				Title: "Pass",
				Verify: func(context.Context, *gitutil.Repository) lessons.CheckResult {
					return lessons.CheckResult{Passed: true, Message: "ok"}
				},
			},
			{
				ID:    "check-2",
				Title: "Fail",
				Verify: func(context.Context, *gitutil.Repository) lessons.CheckResult {
					return lessons.CheckResult{Passed: false, Message: "nope"}
				},
			},
		},
	}

	repo := &gitutil.Repository{}
	emitter := &testEmitter{}

	results, err := Run(context.Background(), lesson, repo, emitter)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if len(emitter.results) != 2 {
		t.Fatalf("expected emitter to receive 2 results, got %d", len(emitter.results))
	}
	if !results[0].Outcome.Passed || results[1].Outcome.Passed {
		t.Fatalf("unexpected pass/fail outcomes: %+v", results)
	}
}

func TestRunRequiresLessonAndRepo(t *testing.T) {
	if _, err := Run(context.Background(), nil, &gitutil.Repository{}, nil); !errors.Is(err, ErrLessonNil) {
		t.Fatalf("expected nil lesson error, got %v", err)
	}
	if _, err := Run(context.Background(), &lessons.Lesson{}, nil, nil); !errors.Is(err, ErrRepositoryNil) {
		t.Fatalf("expected nil repo error, got %v", err)
	}
}
