package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
	// "github.com/stretchr/testify/assert"
)

// setupGitTestEnv creates a temporary test environment with config and git repo
func setupGitTestEnv(t *testing.T) (repoDir string, config *Config, cleanup func()) {
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

	return repoDir, &Config{RepoDir: RepoRoot(repoDir)}, cleanup
}

// TestInitgit tests initializing a git repository
func TestInitgit(t *testing.T) {
	repoDir, c, cleanup := setupGitTestEnv(t)
	defer cleanup()

	// Initialize git repo
	_, err := NewGitReop(c)
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
func TestGitRmFile(t *testing.T) {
	repoDir, c, cleanup := setupGitTestEnv(t)
	defer cleanup()

	r, gitpath := setupAddFile(t, repoDir, c)
	if ret, err := r.GitRmFile(gitpath); err != nil {
		t.Fatalf("GitRmFile failed: %v", err)
	} else if ret.Action != GitResultTypeRm {
		t.Errorf("Expected GitResultRm, got %v", ret)
	} else if len(ret.Files) != 1 {
		t.Errorf("Expected 1 file, got %v", ret)
	}

	if ret, err := r.GitRmFile(gitpath); err != nil {
		t.Fatalf("GitRmFile failed: %v", err)
	} else if ret.Action != GitResultTypeNochange {
		t.Errorf("Expected GitResultRm, got %v", ret)
	} else if len(ret.Files) != 0 {
		t.Errorf("Expected 1 file, got %v", ret)
	}

}

func setupAddFile(t *testing.T, repoDir string, c *Config) (*GitRepo, RepoPath) {
	r, err := NewGitReop(c)
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
	newVar := r.AbsRepo2Repo(testFile)
	_, err = r.GitAddFile(newVar)
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
	return r, newVar
}
func TestGitAddDir(t *testing.T) {
	repoDir, c, cleanup := setupGitTestEnv(t)
	defer cleanup()

	r, err := NewGitReop(c)
	if err != nil {
		t.Fatalf("NewGitReop failed: %v", err)
	}

	// Initialize git repo
	// if err := Initgit(); err != nil {
	// 	t.Fatalf("Initgit failed: %v", err)
	// }

	// Create a test file in the repo
	dir1 := filepath.Join(repoDir, "dir1")
	if err := os.MkdirAll(dir1, 0755); err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}
	testFile := filepath.Join(dir1, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	testFile2 := filepath.Join(dir1, "test2.txt")
	if err := os.WriteFile(testFile2, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add the file (using relative path from repo root)
	newVar := r.AbsRepo2Repo(dir1)
	ret, err := r.GitAddFile(newVar)
	if err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}
	if ret.Action != GitResultTypeAdd {
		t.Errorf("Expected GitResultAdd, got %v", ret)
	}

	// repo, err := r.Open()
	// if err != nil {
	// 	t.Fatalf("Failed to open repo: %v", err)
	// }

	// // Get the HEAD commit
	// ref, err := repo.Head()
	// if err != nil {
	// 	t.Fatalf("Failed to get HEAD: %v", err)
	// }

	// commit, err := repo.CommitObject(ref.Hash())
	// if err != nil {
	// 	t.Fatalf("Failed to get commit: %v", err)
	// }

	// // Verify commit message
	// if commit.Message != fmt.Sprintf("ADD %s", "test.txt") {
	// 	t.Errorf("Expected commit message 'ADD test.txt', got %q", commit.Message)
	// }

	// // Verify author
	// if commit.Author.Name != "anybakup" {
	// 	t.Errorf("Expected author 'anybakup', got %q", commit.Author.Name)
	// }
	// ret, err = r.GitAddFile(newVar)
	// if err != nil {
	// 	t.Fatalf("GitAddFile failed: %v", err)
	// }
	// if ret != GitResultNochange {
	// 	t.Errorf("Expected GitResultAdd, got %v", ret)
	// }
}

// TestGitAddFile tests adding and committing a file
func TestGitAddFile(t *testing.T) {
	repoDir, c, cleanup := setupGitTestEnv(t)
	defer cleanup()

	r, err := NewGitReop(c)
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
	newVar := r.AbsRepo2Repo(testFile)
	ret, err := r.GitAddFile(newVar)
	if err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}
	if ret.Action != GitResultTypeAdd {
		t.Errorf("Expected GitResultAdd, got %v", ret)
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
	ret, err = r.GitAddFile(newVar)
	if err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}
	if ret.Action != GitResultTypeNochange {
		t.Errorf("Expected GitResultAdd, got %v", ret)
	}
}

// TestGitAddFile_MultipleFiles tests adding multiple files
func TestGitAddFile_MultipleFiles(t *testing.T) {
	repoDir, c, cleanup := setupGitTestEnv(t)
	defer cleanup()

	r, err := NewGitReop(c)
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
		if _, err := r.GitAddFile(RepoPath(filename)); err != nil {
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
	repoDir, c, cleanup := setupGitTestEnv(t)
	defer cleanup()

	r, err := NewGitReop(c)
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
	_, err = r.GitAddFile("subdir/test.txt")
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
	repoDir, c, cleanup := setupGitTestEnv(t)
	defer cleanup()

	r, err := NewGitReop(c)
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

	_, err = r.GitAddFile("test.txt")
	if err != nil {
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
	repoDir, c, cleanup := setupGitTestEnv(t)
	defer cleanup()

	r, err := NewGitReop(c)
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
	_, err = r.GitAddFile("file1.txt")
	if err != nil {
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
	repoDir, c, cleanup := setupGitTestEnv(t)
	defer cleanup()

	r, err := NewGitReop(c)
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
	_, err = r.GitAddFile("test.txt")
	if err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	// Second commit
	if err := os.WriteFile(testFile, []byte("version 2"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}
	_, err = r.GitAddFile("test.txt")
	if err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	// Third commit
	if err := os.WriteFile(testFile, []byte("version 3"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}
	_, err = r.GitAddFile("test.txt")
	if err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	// Get change history
	changes, err := r.GitLogFile("test.txt")
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

	for _, c := range changes {
		tmpfile := filepath.Join(os.TempDir(), "test.txt."+c.Commit)
		if _, err := r.GitViewFile("test.txt", c.Commit, tmpfile); err != nil {
			t.Fatalf("GitViewFile failed: %v", err)
		}

	}
}

// TestGitChangesFile_NoCommits tests getting history for a file with no commits
func TestGitChangesFile_NoCommits(t *testing.T) {
	repoDir, c, cleanup := setupGitTestEnv(t)
	defer cleanup()

	r, err := NewGitReop(c)
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
	_, err = r.GitAddFile("file1.txt")
	if err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	// Try to get history for a file that was never committed
	changes, err := r.GitLogFile("nonexistent.txt")
	if err == nil {
		t.Fatalf("GitChangesFile failed: %v", err)
	}

	// Should indicate no commits found
	if len(changes) != 0 {
		t.Error("Expected commit history, got empty string")
	}
}

// TestGitStatusFile tests checking git status
func TestGitStatusFile(t *testing.T) {
	repoDir, c, cleanup := setupGitTestEnv(t)
	defer cleanup()

	// Initialize git repo
	// if err := Initgit(); err != nil {
	// 	t.Fatalf("Initgit failed: %v", err)
	// }

	const newConst = "test.txt"
	// Create and add a file
	testFile := filepath.Join(repoDir, newConst)
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	r, _ := NewGitReop(c)
	// Add the file
	if _, err := r.GetStateStage(testFile); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}
	gitapth := r.AbsRepo2Repo(testFile)
	if _, err := r.GitAddFile(gitapth); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}
	// if s, err := GetState(testFile); err != nil {
	// 	t.Fatalf("GitAddFile failed: %v", err)
	// } else if s != GitUnmodified {
	// 	t.Fatalf("GitAddFile failed: %v", err)
	// }

	if err := os.WriteFile(testFile, []byte("test content sss"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if s, err := r.GetStateWorkTree(testFile); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	} else if s != GitModified {
		t.Fatalf("GitAddFile failed: %v", err)
	}
	if _, err := r.GitAddFile(gitapth); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}
	if s, err := r.GetStateWorkTree(testFile); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	} else if s != GitUntracked {
		t.Fatalf("GitAddFile failed: %v", err)
	}
	if s, err := r.GetStateStage(testFile); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	} else if s != GitUntracked {
		t.Fatalf("GitAddFile failed: %v", err)
	}
	if _, err := r.GitAddFile(gitapth); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}
	if s, err := r.GetStateWorkTree(testFile); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	} else if s != GitUntracked {
		t.Fatalf("GitAddFile failed: %v", err)
	}
	if s, err := r.GetStateStage(testFile); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	} else if s != GitUntracked {
		t.Fatalf("GitAddFile failed: %v", err)
	}
	// if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
	// 	t.Fatalf("Failed to create test file: %v", err)
	// }

	// if s, err := GetStateStage(testFile, nil); err != nil {
	// 	t.Fatalf("GitAddFile failed: %v", err)
	// } else if s != GitUntracked {
	// 	t.Fatalf("GitAddFile failed: %v", err)
	// }

	if _, err := r.GitRmFile(gitapth); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}
	if s, err := r.GetStateWorkTree(testFile); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	} else if s != GitUntracked {
		t.Fatalf("GitAddFile failed: %v", err)
	}
	if s, err := r.GetStateStage(testFile); err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	} else if s != GitUntracked {
		t.Fatalf("GitAddFile failed: %v", err)
	}
}

// TestGitViewFile tests checking out a file from a specific commit
func TestGitViewFile(t *testing.T) {
	repoDir, c, cleanup := setupGitTestEnv(t)
	defer cleanup()

	r, err := NewGitReop(c)
	if err != nil {
		t.Fatalf("NewGitReop failed: %v", err)
	}

	// Create and commit a file with initial content
	testFile := filepath.Join(repoDir, "test.txt")
	initialContent := "version 1 content"
	if err := os.WriteFile(testFile, []byte(initialContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = r.GitAddFile("test.txt")
	if err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	// Get the commit hash of the first version
	repo, err := r.Open()
	if err != nil {
		t.Fatalf("Failed to open repo: %v", err)
	}

	ref, err := repo.Head()
	if err != nil {
		t.Fatalf("Failed to get HEAD: %v", err)
	}
	firstCommitHash := ref.Hash().String()

	// Modify the file and commit again
	modifiedContent := "version 2 content - modified"
	if err := os.WriteFile(testFile, []byte(modifiedContent), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	_, err = r.GitAddFile("test.txt")
	if err != nil {
		t.Fatalf("GitAddFile failed for second commit: %v", err)
	}

	// Create output directory
	outputDir := filepath.Join(repoDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}

	// Test 1: Checkout the first version using commit hash
	outputFile := filepath.Join(outputDir, "test_v1.txt")
	resultPath, err := r.GitViewFile(RepoPath("test.txt"), firstCommitHash, outputFile)
	if err != nil {
		t.Fatalf("GitViewFile failed: %v", err)
	}

	if resultPath != outputFile {
		t.Errorf("Expected result path %s, got %s", outputFile, resultPath)
	}

	// Verify the content matches the first version
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if string(content) != initialContent {
		t.Errorf("Expected content %q, got %q", initialContent, string(content))
	}

	// Test 2: Checkout to a subdirectory that doesn't exist yet
	outputFile2 := filepath.Join(outputDir, "subdir", "nested", "test.txt")
	resultPath2, err := r.GitViewFile(RepoPath("test.txt"), firstCommitHash, outputFile2)
	if err != nil {
		t.Fatalf("GitViewFile failed for nested directory: %v", err)
	}

	if resultPath2 != outputFile2 {
		t.Errorf("Expected result path %s, got %s", outputFile2, resultPath2)
	}

	// Verify the file exists and has correct content
	content2, err := os.ReadFile(outputFile2)
	if err != nil {
		t.Fatalf("Failed to read nested output file: %v", err)
	}

	if string(content2) != initialContent {
		t.Errorf("Expected content %q in nested file, got %q", initialContent, string(content2))
	}
}

// TestGitViewFile_InvalidCommit tests error handling for invalid commit hash
func TestGitViewFile_InvalidCommit(t *testing.T) {
	repoDir, c, cleanup := setupGitTestEnv(t)
	defer cleanup()

	r, err := NewGitReop(c)
	if err != nil {
		t.Fatalf("NewGitReop failed: %v", err)
	}

	// Create and commit a file
	testFile := filepath.Join(repoDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = r.GitAddFile("test.txt")
	if err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	// Try to checkout with an invalid commit hash
	outputFile := filepath.Join(repoDir, "output.txt")
	invalidHash := "0000000000000000000000000000000000000000"
	_, err = r.GitViewFile(RepoPath("test.txt"), invalidHash, outputFile)
	if err == nil {
		t.Error("Expected error for invalid commit hash, got nil")
	}

	if !strings.Contains(err.Error(), "failed to get commit") {
		t.Errorf("Expected error message about commit, got: %v", err)
	}
}

// TestGitViewFile_FileNotInCommit tests error handling when file doesn't exist in commit
func TestGitViewFile_FileNotInCommit(t *testing.T) {
	repoDir, c, cleanup := setupGitTestEnv(t)
	defer cleanup()

	r, err := NewGitReop(c)
	if err != nil {
		t.Fatalf("NewGitReop failed: %v", err)
	}

	// Create and commit a file
	testFile := filepath.Join(repoDir, "existing.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = r.GitAddFile("existing.txt")
	if err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	// Get the commit hash
	repo, err := r.Open()
	if err != nil {
		t.Fatalf("Failed to open repo: %v", err)
	}

	ref, err := repo.Head()
	if err != nil {
		t.Fatalf("Failed to get HEAD: %v", err)
	}
	commitHash := ref.Hash().String()

	// Try to checkout a file that doesn't exist in this commit
	outputFile := filepath.Join(repoDir, "output.txt")
	_, err = r.GitViewFile(RepoPath("nonexistent.txt"), commitHash, outputFile)
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected error message about file not found, got: %v", err)
	}
}

// TestCleanEmptyDir tests the CleanEmptyDir function with various scenarios.
func TestCleanEmptyDir(t *testing.T) {
	reporoot, c, clean := setupGitTestEnv(t)
	defer clean()
	// 使用 afero 模拟文件系统
	// tmpDir, err := os.MkdirTemp("", "anybakup-git-test-*")
	// if err!=nil {
	// 	t.Fatal(err)
	// }
	repo, err := NewGitReop(c)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("Non-existent path", func(t *testing.T) {
		// 测试路径不存在的情况
		n, err := repo.CleanEmptyDir(filepath.Join(reporoot, "nonexistent"))
		if n != 0 {
			t.Fatal("Expected non-zero number of deleted files, got zero", err)
		}
	})

	t.Run("Non-empty directory", func(t *testing.T) {
		// 创建一个非空目录
		dirPath := filepath.Join(reporoot, "nonempty")
		os.MkdirAll(dirPath, 0755)
		os.WriteFile(filepath.Join(dirPath, "file.txt"), []byte("content"), 0644)

		n, err := repo.CleanEmptyDir(dirPath)
		if n != 0 {
			t.Fatal("Expected non-zero number of deleted files, got zero", err)
		}
	})

	t.Run("Empty directory", func(t *testing.T) {
		// 创建一个空目录
		dirPath := filepath.Join(reporoot, "empty")
		os.MkdirAll(dirPath, 0755)

		n, err := repo.CleanEmptyDir(dirPath)
		if n != 1 {
			t.Fatal("Expected non-zero number of deleted files, got zero", err)
		}
	})
	t.Run("Empty directory 2", func(t *testing.T) {
		// 创建一个空目录
		dirPath := filepath.Join(reporoot, "empty")
		os.MkdirAll(dirPath, 0755)

		if dirPath := filepath.Join(reporoot, "empty", "empty2"); os.MkdirAll(dirPath, 0755) == nil {

		}

		n, err := repo.CleanEmptyDir(dirPath)
		if n != 2 {
			t.Fatal("Expected non-zero number of deleted files, got zero", err)
		}
	})
	t.Run("childe not empty", func(t *testing.T) {
		// 创建一个空目录
		dirPath := filepath.Join(reporoot, "empty")
		os.MkdirAll(dirPath, 0755)

		if dirPath := filepath.Join(reporoot, "empty", "empty2"); os.MkdirAll(dirPath, 0755) == nil {
			if eerr := os.WriteFile(filepath.Join(dirPath, "1.txt"), []byte("content"), 0644); eerr != nil {
				t.Fatal(eerr)
			}
		}

		n, err := repo.CleanEmptyDir(dirPath)
		if n != 0 {
			t.Fatal("Expected non-zero number of deleted files, got zero", err)
		}
	})
	// t.Run("File path", func(t *testing.T) {
	// 	// 创建一个文件
	// 	filePath := filepath.Join(reporoot, "file.txt")
	// 	os.WriteFile(filePath, []byte("content"), 0644)

	// 	err := repo.CleanEmptyDir(filePath)
	// 	assert.NoError(t, err)
	// 	if _, err = os.Stat(filePath); err != nil {
	// 		t.Fatal(err)
	// 	}
	// })

}
