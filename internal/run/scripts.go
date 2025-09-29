package run

import (
	"sort"
)

// Step describes a single command invocation within a run lesson.
type Step struct {
	Command        string
	ExpectExitCode int
	ExpectStdout   []string
}

// Script represents a scripted exercise composed of multiple steps.
type Script struct {
	ID          string
	Title       string
	Description string
	Steps       []Step
}

var (
	registry   = map[string]*Script{}
	scriptList []string
)

// Register adds the provided script to the registry.
func Register(script *Script) {
	if script == nil || script.ID == "" {
		return
	}
	if _, exists := registry[script.ID]; exists {
		return
	}
	registry[script.ID] = script
	scriptList = append(scriptList, script.ID)
	sort.Strings(scriptList)
}

// Get retrieves a script by ID.
func Get(id string) (*Script, bool) {
	script, ok := registry[id]
	return script, ok
}

// List returns registered scripts in deterministic order.
func List() []*Script {
	out := make([]*Script, 0, len(scriptList))
	for _, id := range scriptList {
		out = append(out, registry[id])
	}
	return out
}

func init() {
	Register(&Script{
		ID:          "0",
		Title:       "Test your CLI",
		Description: "Confirm the CLI can execute commands in your environment.",
		Steps: []Step{
			{
				Command:        "echo \"naad karti kay?\"",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"naad karti kay?"},
			},
		},
	})

	Register(&Script{
		ID:          "1",
		Title:       "Install Git",
		Description: "Ensure git is available on your PATH.",
		Steps: []Step{
			{
				Command:        "git --version",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"git version 2"},
			},
		},
	})

	Register(&Script{
		ID:          "2",
		Title:       "Configure Git identity",
		Description: "Validate that Git global identity settings are in place.",
		Steps: []Step{
			{
				Command:        "git config get --global user.name",
				ExpectExitCode: 0,
			},
			{
				Command:        "git config get --global user.email",
				ExpectExitCode: 0,
			},
			{
				Command:        "git config get --global init.defaultBranch",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"master"},
			},
		},
	})

	Register(&Script{
		ID:          "3",
		Title:       "Initialize a repository",
		Description: "Check the current directory is the webflyx repo with a .git folder.",
		Steps: []Step{
			{
				Command:        "pwd",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"webflyx"},
			},
			{
				Command:        "ls .git",
				ExpectExitCode: 0,
				ExpectStdout: []string{
					"config",
					"description",
					"HEAD",
					"hooks",
					"info",
					"objects",
					"refs",
				},
			},
		},
	})

	Register(&Script{
		ID:          "4a",
		Title:       "Track contents.md",
		Description: "Ensure contents.md exists with expected content before staging.",
		Steps: []Step{
			{
				Command:        "git status",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"Untracked files", "contents.md"},
			},
			{
				Command:        "cat contents.md",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"# contents"},
			},
		},
	})

	Register(&Script{
		ID:          "4b",
		Title:       "Stage contents.md",
		Description: "Verify contents.md is staged prior to commit.",
		Steps: []Step{
			{
				Command:        "git status",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"Changes to be committed", "new file:", "contents.md"},
			},
		},
	})

	Register(&Script{
		ID:          "5",
		Title:       "Commit contents.md",
		Description: "Confirm the commit message A: add contents.md is recorded.",
		Steps: []Step{
			{
				Command:        "git --no-pager log -n 1",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"A:", "add contents.md", "commit"},
			},
		},
	})

	Register(&Script{
		ID:          "6a",
		Title:       "Inspect commit metadata",
		Description: "Review the temporary catfileout.txt output from git cat-file.",
		Steps: []Step{
			{
				Command:        "cat catfileout.txt",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"tree", "author", "committer"},
			},
		},
	})

	Register(&Script{
		ID:          "6b",
		Title:       "Inspect blob contents",
		Description: "Validate blobfile.txt captures blob output.",
		Steps: []Step{
			{
				Command:        "cat blobfile.txt",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"# contents"},
			},
		},
	})

	Register(&Script{
		ID:          "7",
		Title:       "Create titles.md",
		Description: "Ensure titles.md exists and commit message B: add titles.md was created.",
		Steps: []Step{
			{
				Command:        "git --no-pager log",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"B:"},
			},
			{
				Command:        "cat titles.md",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"# Titles"},
			},
		},
	})

	Register(&Script{
		ID:          "8",
		Title:       "Switch default branch",
		Description: "Check global default branch and local branch rename to main.",
		Steps: []Step{
			{
				Command:        "git config get --global init.defaultBranch",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"main"},
			},
			{
				Command:        "git branch",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"* main"},
			},
		},
	})

	Register(&Script{
		ID:          "9",
		Title:       "Create add_classics branch",
		Description: "Verify branch creation and checkout to add_classics.",
		Steps: []Step{
			{
				Command:        "git branch",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"* add_classics", "main"},
			},
		},
	})

	Register(&Script{
		ID:          "10",
		Title:       "Commit classics.csv",
		Description: "Confirm commit C: add classics.csv exists on add_classics.",
		Steps: []Step{
			{
				Command:        "git --no-pager log -n 1 --pretty=%s",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"C: add classics.csv"},
			},
			{
				Command:        "git branch",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"* add_classics"},
			},
		},
	})

	Register(&Script{
		ID:          "11a",
		Title:       "Prepare merge history",
		Description: "Inspect log graph prior to merging add_classics.",
		Steps: []Step{
			{
				Command:        "git --no-pager log --oneline --graph --all",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"A:", "B:", "C:", "D:", "|/"},
			},
		},
	})

	Register(&Script{
		ID:          "11b",
		Title:       "Merge add_classics",
		Description: "Verify the merge commit exists with parents displayed.",
		Steps: []Step{
			{
				Command:        "git --no-pager log --oneline --decorate --graph --parents",
				ExpectExitCode: 0,
				ExpectStdout:   []string{"E:", "|\\"},
			},
		},
	})
}
