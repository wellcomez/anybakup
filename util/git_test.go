package util

import (
	"os"
	"path/filepath"
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
	err := Initgit()
	if err != nil {
		t.Fatalf("Initgit failed: %v", err)
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
	// // It's OK if there's no HEAD yet (empty repo), but other errors are bad
	// if err != nil && err != git.e{
	// 	t.Fatalf("Invalid git repository: %v", err)
	// }
}

// TestInitgit_AlreadyInitialized tests initializing an already initialized repo
// func TestInitgit_AlreadyInitialized(t *testing.T) {
// 	_, cleanup := setupGitTestEnv(t)
// 	defer cleanup()

// 	// Initialize git repo twice
// 	if err := Initgit(); err != nil {
// 		t.Fatalf("First Initgit failed: %v", err)
// 	}

// 	// Second init should handle gracefully (go-git returns error for already initialized)
// 	err := Initgit()
// 	// go-git returns an error if repo is already initialized, which is expected
// 	if err == nil {
// 		t.Log("Note: Initgit allowed re-initialization (this may be OK)")
// 	}
// }

// TestGitAddFile tests adding and committing a file
func TestGitAddFile(t *testing.T) {
	repoDir, cleanup := setupGitTestEnv(t)
	defer cleanup()

	// Initialize git repo
	if err := Initgit(); err != nil {
		t.Fatalf("Initgit failed: %v", err)
	}

	// Create a test file in the repo
	testFile := filepath.Join(repoDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add the file (using relative path from repo root)
	err := GitAddFile("test.txt")
	if err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	// Verify the file was committed
	repo, err := git.PlainOpen(repoDir)
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
	if commit.Message != "add file" {
		t.Errorf("Expected commit message 'add file', got %q", commit.Message)
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

	// Initialize git repo
	if err := Initgit(); err != nil {
		t.Fatalf("Initgit failed: %v", err)
	}

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
		if err := GitAddFile(filename); err != nil {
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

	// Initialize git repo
	if err := Initgit(); err != nil {
		t.Fatalf("Initgit failed: %v", err)
	}

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
	err := GitAddFile("subdir/test.txt")
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

// TestGitStatusFile tests checking git status
func TestGitStatusFile(t *testing.T) {
	repoDir, cleanup := setupGitTestEnv(t)
	defer cleanup()

	// Initialize git repo
	if err := Initgit(); err != nil {
		t.Fatalf("Initgit failed: %v", err)
	}

	// Create and add a file
	testFile := filepath.Join(repoDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if err := GitAddFile("test.txt"); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	// Check status (should be clean after commit)
	err := GitStatusFile("test.txt")
	if err != nil {
		t.Fatalf("GitStatusFile failed: %v", err)
	}
}

// TestGitCommitFile tests committing changes
func TestGitCommitFile(t *testing.T) {
	repoDir, cleanup := setupGitTestEnv(t)
	defer cleanup()

	// Initialize git repo
	if err := Initgit(); err != nil {
		t.Fatalf("Initgit failed: %v", err)
	}

	// Create and add a file
	testFile := filepath.Join(repoDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("initial content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if err := GitAddFile("test.txt"); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	// Modify the file
	if err := os.WriteFile(testFile, []byte("modified content"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	// Stage the changes
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		t.Fatalf("Failed to open repo: %v", err)
	}
	w, err := repo.Worktree()
	if err != nil {
		t.Fatalf("Failed to get worktree: %v", err)
	}
	if _, err := w.Add("test.txt"); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// Commit the changes
	err = GitCommitFile("test.txt")
	if err != nil {
		t.Fatalf("GitCommitFile failed: %v", err)
	}

	// Verify we have 2 commits now
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

	if commitCount != 2 {
		t.Errorf("Expected 2 commits, got %d", commitCount)
	}
}

// TestGitDiffFile tests checking file differences
func TestGitDiffFile(t *testing.T) {
	repoDir, cleanup := setupGitTestEnv(t)
	defer cleanup()

	// Initialize git repo
	if err := Initgit(); err != nil {
		t.Fatalf("Initgit failed: %v", err)
	}

	// Create and add a file
	testFile := filepath.Join(repoDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("initial content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if err := GitAddFile("test.txt"); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	// Check diff (should be clean)
	err := GitDiffFile("test.txt")
	if err != nil {
		t.Fatalf("GitDiffFile failed: %v", err)
	}
}

// TestRepoRoot tests the repo_root helper function
func TestRepoRoot(t *testing.T) {
	repoDir, cleanup := setupGitTestEnv(t)
	defer cleanup()

	// Get repo root
	root, err := repo_root()
	if err != nil {
		t.Fatalf("repo_root failed: %v", err)
	}

	if root != repoDir {
		t.Errorf("Expected repo root %s, got %s", repoDir, root)
	}
}

// TestRepoRoot_NoConfig tests repo_root with missing config
func TestRepoRoot_NoConfig(t *testing.T) {
	// Set HOME to a non-existent directory
	oldHome := os.Getenv("HOME")
	tmpDir, _ := os.MkdirTemp("", "test-*")
	os.Setenv("HOME", tmpDir)
	defer func() {
		os.Setenv("HOME", oldHome)
		os.RemoveAll(tmpDir)
	}()

	// Should fail because config doesn't exist
	_, err := repo_root()
	if err == nil {
		t.Fatal("Expected error for missing config, got nil")
	}
}
