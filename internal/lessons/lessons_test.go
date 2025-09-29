package lessons

import "testing"

func TestCatalogRegistration(t *testing.T) {
	lessons := List()
	if len(lessons) < 2 {
		t.Fatalf("expected at least 2 lessons, got %d", len(lessons))
	}

	initLesson, err := Get("init-basics")
	if err != nil {
		t.Fatalf("Get(init-basics): %v", err)
	}
	if initLesson.ID != "init-basics" {
		t.Fatalf("unexpected lesson ID: %s", initLesson.ID)
	}

	duplicate := *initLesson
	duplicate.ID = "init-basics"
	if err := Register(&duplicate); err == nil {
		t.Fatalf("expected duplicate registration to fail")
	}
}
