#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)"
BIN_DIR="${ROOT_DIR}/bin"
mkdir -p "${BIN_DIR}"

TSGIT_BIN="${BIN_DIR}/tscgit"
go build -o "${TSGIT_BIN}" "${ROOT_DIR}/cmd/tscgit"

WORKDIR="$(mktemp -d -t tscgit-lessons-XXXXXX)"
trap 'rm -rf "${WORKDIR}"' EXIT

export GIT_CONFIG_GLOBAL="${WORKDIR}/gitconfig"
cat >"${GIT_CONFIG_GLOBAL}" <<'EOF'
[user]
    name = Test Student
    email = student@example.com
[init]
    defaultBranch = master
EOF

export GIT_AUTHOR_NAME="Test Student"
export GIT_AUTHOR_EMAIL="student@example.com"
export GIT_COMMITTER_NAME="${GIT_AUTHOR_NAME}"
export GIT_COMMITTER_EMAIL="${GIT_AUTHOR_EMAIL}"

REPO_DIR="${WORKDIR}/webflyx"
mkdir -p "${REPO_DIR}"
cd "${REPO_DIR}"

git init
echo "# Webflyx Practice" >README.md

"${TSGIT_BIN}" run 0
"${TSGIT_BIN}" run 1
"${TSGIT_BIN}" run 2
"${TSGIT_BIN}" run 3

echo "# contents" >contents.md
"${TSGIT_BIN}" run 4a

git add contents.md
"${TSGIT_BIN}" run 4b

git commit -m "A: add contents.md"
"${TSGIT_BIN}" run 5

git cat-file commit HEAD >catfileout.txt
BLOB_ID="$(git rev-parse HEAD:contents.md)"
git cat-file blob "${BLOB_ID}" >blobfile.txt
"${TSGIT_BIN}" run 6a
"${TSGIT_BIN}" run 6b

echo "# Titles" >titles.md
git add titles.md
git commit -m "B: add titles.md"
"${TSGIT_BIN}" run 7

git config --global init.defaultBranch main
git branch -M main
"${TSGIT_BIN}" run 8

git switch -c add_classics
"${TSGIT_BIN}" run 9

echo "title,year" >classics.csv
echo "Metropolis,1927" >>classics.csv
git add classics.csv
git commit -m "C: add classics.csv"
"${TSGIT_BIN}" run 10

git switch main
echo "" >>contents.md
echo "## Classics" >>contents.md
git add contents.md
git commit -m "D: update contents.md"
"${TSGIT_BIN}" run 11a

git merge add_classics --no-ff -m "E: merge add_classics"
"${TSGIT_BIN}" run 11b

git branch -D "C:" >/dev/null 2>&1 || true
git switch -c feature/lesson-branch
echo "- Added new section" >>contents.md
git add contents.md
git commit -m "[branch] add lesson notes"

"${TSGIT_BIN}" verify init-basics
"${TSGIT_BIN}" verify branch-basics

echo "\nAll lessons and verifications completed successfully."