package verify

import (
	"context"
	"errors"
	"time"

	"github.com/rohit746/tscgit/internal/gitutil"
	"github.com/rohit746/tscgit/internal/lessons"
)

var (
	ErrLessonNil     = errors.New("verify: lesson is nil")
	ErrRepositoryNil = errors.New("verify: repository is nil")
)

// Result captures the outcome of executing a lesson check.
type Result struct {
	Check    lessons.Check
	Outcome  lessons.CheckResult
	Duration time.Duration
}

// Emitter receives results as they are produced.
type Emitter interface {
	Emit(Result)
}

// Run executes each check in the provided lesson sequentially. The optional
// emitter is invoked after every check with its result. Execution halts if the
// context is cancelled.
func Run(ctx context.Context, lesson *lessons.Lesson, repo *gitutil.Repository, emitter Emitter) ([]Result, error) {
	if lesson == nil {
		return nil, ErrLessonNil
	}
	if repo == nil {
		return nil, ErrRepositoryNil
	}

	results := make([]Result, 0, len(lesson.Checks))

	for _, check := range lesson.Checks {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}

		start := time.Now()
		outcome := check.Verify(ctx, repo)
		res := Result{Check: check, Outcome: outcome, Duration: time.Since(start)}
		results = append(results, res)
		if emitter != nil {
			emitter.Emit(res)
		}
	}

	return results, nil
}
