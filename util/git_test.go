package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// setupGitTestEnv creates a temporary test environment with config and git repo
func setupGitTestEnv(t *testing.T) (repoDir string, cleanup func()) {
	// Create temporary directories
	tmpDir, err := os.MkdirTemp("", "anybakup-git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	repoDir = filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("Failed to create repo dir: %v", err)
	}

	// Create config directory and file
	configDir := filepath.Join(tmpDir, ".config", "anybakup")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	configFile := filepath.Join(configDir, "config.yaml")
	configContent := "repodir: " + repoDir + "\n"
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Set HOME to temp directory so config is loaded from there
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)

	cleanup = func() {
		os.Setenv("HOME", oldHome)
		os.RemoveAll(tmpDir)
	}

	return repoDir, cleanup
}

// TestInitgit tests initializing a git repository
func TestInitgit(t *testing.T) {
	repoDir, cleanup := setupGitTestEnv(t)
	defer cleanup()

	// Initialize git repo
	_, err := NewGitReop()
	if err != nil {
		t.Fatalf("NewGitReop failed: %v", err)
	}

	// Verify .git directory exists
	gitDir := filepath.Join(repoDir, ".git")
	if info, err := os.Stat(gitDir); err != nil || !info.IsDir() {
		t.Fatalf("Git directory not created: %v", err)
	}

	// Verify we can open the repo
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		t.Fatalf("Failed to open git repo: %v", err)
	}

	// Verify it's a valid git repository
	repo.Head()
}

// TestGitAddFile tests adding and committing a file
func TestGitAddFile(t *testing.T) {
	repoDir, cleanup := setupGitTestEnv(t)
	defer cleanup()

	r, err := NewGitReop()
	if err != nil {
		t.Fatalf("NewGitReop failed: %v", err)
	}

	// Initialize git repo
	// if err := Initgit(); err != nil {
	// 	t.Fatalf("Initgit failed: %v", err)
	// }

	// Create a test file in the repo
	testFile := filepath.Join(repoDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add the file (using relative path from repo root)
	err = r.GitAddFile("test.txt")
	if err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	repo, err := r.Open()
	if err != nil {
		t.Fatalf("Failed to open repo: %v", err)
	}

	// Get the HEAD commit
	ref, err := repo.Head()
	if err != nil {
		t.Fatalf("Failed to get HEAD: %v", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		t.Fatalf("Failed to get commit: %v", err)
	}

	// Verify commit message
	if commit.Message != fmt.Sprintf("ADD %s", "test.txt") {
		t.Errorf("Expected commit message 'ADD test.txt', got %q", commit.Message)
	}

	// Verify author
	if commit.Author.Name != "anybakup" {
		t.Errorf("Expected author 'anybakup', got %q", commit.Author.Name)
	}
}

// TestGitAddFile_MultipleFiles tests adding multiple files
func TestGitAddFile_MultipleFiles(t *testing.T) {
	repoDir, cleanup := setupGitTestEnv(t)
	defer cleanup()

	r, err := NewGitReop()
	if err != nil {
		t.Fatalf("NewGitReop failed: %v", err)
	}

	// Initialize git repo
	// if err := Initgit(); err != nil {
	// 	t.Fatalf("Initgit failed: %v", err)
	// }

	// Create multiple test files
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, filename := range files {
		testFile := filepath.Join(repoDir, filename)
		if err := os.WriteFile(testFile, []byte("content of "+filename), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Add all files
	for _, filename := range files {
		if err := r.GitAddFile(filename); err != nil {
			t.Fatalf("GitAddFile failed for %s: %v", filename, err)
		}
	}

	// Verify we have 3 commits
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		t.Fatalf("Failed to open repo: %v", err)
	}

	commitCount := 0
	iter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		t.Fatalf("Failed to get log: %v", err)
	}

	err = iter.ForEach(func(c *object.Commit) error {
		commitCount++
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to iterate commits: %v", err)
	}

	if commitCount != 3 {
		t.Errorf("Expected 3 commits, got %d", commitCount)
	}
}

// TestGitAddFile_Subdirectory tests adding a file in a subdirectory
func TestGitAddFile_Subdirectory(t *testing.T) {
	repoDir, cleanup := setupGitTestEnv(t)
	defer cleanup()

	r, err := NewGitReop()
	if err != nil {
		t.Fatalf("NewGitReop failed: %v", err)
	}

	// Initialize git repo
	// if err := Initgit(); err != nil {
	// 	t.Fatalf("Initgit failed: %v", err)
	// }

	// Create a subdirectory and file
	subDir := filepath.Join(repoDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	testFile := filepath.Join(subDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add the file (using relative path from repo root)
	err = r.GitAddFile("subdir/test.txt")
	if err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	// Verify the file was committed
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		t.Fatalf("Failed to open repo: %v", err)
	}

	ref, err := repo.Head()
	if err != nil {
		t.Fatalf("Failed to get HEAD: %v", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		t.Fatalf("Failed to get commit: %v", err)
	}

	// Verify the file is in the commit
	tree, err := commit.Tree()
	if err != nil {
		t.Fatalf("Failed to get tree: %v", err)
	}

	_, err = tree.File("subdir/test.txt")
	if err != nil {
		t.Fatalf("File not found in commit tree: %v", err)
	}
}

// TestGitDiffFile tests checking file differences
func TestGitDiffFile(t *testing.T) {
	repoDir, cleanup := setupGitTestEnv(t)
	defer cleanup()

	r, err := NewGitReop()
	if err != nil {
		t.Fatalf("NewGitReop failed: %v", err)
	}

	// // Initialize git repo
	// if err := Initgit(); err != nil {
	// 	t.Fatalf("Initgit failed: %v", err)
	// }

	// Create and add a file
	testFile := filepath.Join(repoDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("initial content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if err := r.GitAddFile("test.txt"); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	// Check diff (should be empty - no changes)
	diff, err := r.GitDiffFile("test.txt")
	if err != nil {
		t.Fatalf("GitDiffFile failed: %v", err)
	}
	if diff != "" {
		t.Errorf("Expected empty diff for unchanged file, got: %s", diff)
	}

	// Modify the file
	if err := os.WriteFile(testFile, []byte("modified content"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	// Check diff (should show changes)
	diff, err = r.GitDiffFile("test.txt")
	if err != nil {
		t.Fatalf("GitDiffFile failed: %v", err)
	}
	if diff == "" {
		t.Error("Expected diff output for modified file, got empty string")
	}
	if !strings.Contains(diff, "test.txt") {
		t.Errorf("Expected diff to contain filename, got: %s", diff)
	}
}

// TestGitDiffFile_NewFile tests diff for a new file not in HEAD
func TestGitDiffFile_NewFile(t *testing.T) {
	repoDir, cleanup := setupGitTestEnv(t)
	defer cleanup()

	r, err := NewGitReop()
	if err != nil {
		t.Fatalf("NewGitReop failed: %v", err)
	}

	// // Initialize git repo
	// if err := Initgit(); err != nil {
	// 	t.Fatalf("Initgit failed: %v", err)
	// }

	// Create initial file to have a commit
	testFile1 := filepath.Join(repoDir, "file1.txt")
	if err := os.WriteFile(testFile1, []byte("file1"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := r.GitAddFile("file1.txt"); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	// Create a new file not yet committed
	testFile2 := filepath.Join(repoDir, "newfile.txt")
	if err := os.WriteFile(testFile2, []byte("new content"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	// Check diff for new file
	diff, err := r.GitDiffFile("newfile.txt")
	if err != nil {
		t.Fatalf("GitDiffFile failed: %v", err)
	}
	if !strings.Contains(diff, "newfile.txt") {
		t.Errorf("Expected diff to contain filename for new file, got: %s", diff)
	}
}

// TestGitChangesFile tests getting commit history for a file
func TestGitChangesFile(t *testing.T) {
	repoDir, cleanup := setupGitTestEnv(t)
	defer cleanup()

	r, err := NewGitReop()
	if err != nil {
		t.Fatalf("NewGitReop failed: %v", err)
	}

	// // Initialize git repo
	// if err := Initgit(); err != nil {
	// 	t.Fatalf("Initgit failed: %v", err)
	// }

	// Create and commit a file multiple times
	testFile := filepath.Join(repoDir, "test.txt")

	// First commit
	if err := os.WriteFile(testFile, []byte("version 1"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := r.GitAddFile("test.txt"); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	// Second commit
	if err := os.WriteFile(testFile, []byte("version 2"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}
	if err := r.GitAddFile("test.txt"); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	// Third commit
	if err := os.WriteFile(testFile, []byte("version 3"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}
	if err := r.GitAddFile("test.txt"); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	// Get change history
	changes, err := r.GitChangesFile("test.txt")
	if err != nil {
		t.Fatalf("GitChangesFile failed: %v", err)
	}

	// Verify we have commit information
	if len(changes) == 0 {
		t.Error("Expected commit history, got empty string")
	}

	// Should contain commit markers
	commitCount := len(changes)
	if commitCount != 3 {
		t.Errorf("Expected 3 commits in history, found %d", commitCount)
	}

	// Should contain author information
	for _, c := range changes {

		if !strings.Contains(c.Author, "anybakup") {
			t.Error("Expected author information in changes")
		}

		// Should contain date information
		if c.Date == "" {
			t.Error("Expected date information in changes")
		}
	}
}

// TestGitChangesFile_NoCommits tests getting history for a file with no commits
func TestGitChangesFile_NoCommits(t *testing.T) {
	repoDir, cleanup := setupGitTestEnv(t)
	defer cleanup()

	r, err := NewGitReop()
	if err != nil {
		t.Fatalf("NewGitReop failed: %v", err)
	}

	// // Initialize git repo
	// if err := Initgit(); err != nil {
	// 	t.Fatalf("Initgit failed: %v", err)
	// }

	// Create a file to have at least one commit
	testFile1 := filepath.Join(repoDir, "file1.txt")
	if err := os.WriteFile(testFile1, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := r.GitAddFile("file1.txt"); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	// Try to get history for a file that was never committed
	changes, err := r.GitChangesFile("nonexistent.txt")
	if err != nil {
		t.Fatalf("GitChangesFile failed: %v", err)
	}

	// Should indicate no commits found
	if len(changes) == 0 {
		t.Error("Expected commit history, got empty string")
	}
}

// TestGitStatusFile tests checking git status
func TestGitStatusFile(t *testing.T) {
	repoDir, cleanup := setupGitTestEnv(t)
	defer cleanup()

	r, err := NewGitReop()
	if err != nil {
		t.Fatalf("NewGitReop failed: %v", err)
	}

	// Initialize git repo
	// if err := Initgit(); err != nil {
	// 	t.Fatalf("Initgit failed: %v", err)
	// }

	// Create and add a file
	testFile := filepath.Join(repoDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if err := r.GitAddFile("test.txt"); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	// Check status
	s := GitStatusFile{}
	err = s.CheckStatus()
	if err != nil {
		t.Fatalf("GitStatusFile.CheckStatus failed: %v", err)
	}

}
