package lessons

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/rohit746/tscgit/internal/gitutil"
)

// CheckFunc executes a single verification step against the given repository.
type CheckFunc func(ctx context.Context, repo *gitutil.Repository) CheckResult

// CheckResult represents the outcome of a CheckFunc.
type CheckResult struct {
	Passed  bool
	Message string
	Err     error
}

// Check captures a single verification step for a lesson.
type Check struct {
	ID          string
	Title       string
	Description string
	Verify      CheckFunc
}

// Lesson bundles together a series of Checks.
type Lesson struct {
	ID          string
	Title       string
	Description string
	Checks      []Check
}

var (
	catalog       = map[string]*Lesson{}
	lessonOrder   []string
	errLessonDNE  = errors.New("lesson not found")
	errCatalogSet = errors.New("lesson already registered")
)

// Register adds a new lesson to the catalog. It should typically be invoked in
// an init function to keep lesson discovery modular.
func Register(lesson *Lesson) error {
	if lesson == nil {
		return errors.New("lesson is nil")
	}
	if lesson.ID == "" {
		return errors.New("lesson ID is required")
	}
	if _, exists := catalog[lesson.ID]; exists {
		return fmt.Errorf("%w: %s", errCatalogSet, lesson.ID)
	}
	catalog[lesson.ID] = lesson
	lessonOrder = append(lessonOrder, lesson.ID)
	sort.Strings(lessonOrder)
	return nil
}

// List returns the catalog of lessons in a deterministic order.
func List() []*Lesson {
	out := make([]*Lesson, 0, len(lessonOrder))
	for _, id := range lessonOrder {
		out = append(out, catalog[id])
	}
	return out
}

// Get retrieves a lesson by ID.
func Get(id string) (*Lesson, error) {
	lesson, ok := catalog[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", errLessonDNE, id)
	}
	return lesson, nil
}

// Must ensures lesson registration succeeds.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
