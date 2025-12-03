package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"anybakup/cmd"
	"anybakup/util"
)

// setupTestEnv creates a temporary test environment with config and git repo
func setupTestEnv(t *testing.T) (repoDir string, cleanup func()) {
	// Create temporary directories
	tmpDir, err := os.MkdirTemp("", "anybakup-cmd-test-*")
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
func TestGitAddFile(t *testing.T) {
	r, err := util.NewGitReop()
	if err != nil {
		t.Fatalf("NewGitReop failed: %v", err)
	}

	tmpDir, _ := os.MkdirTemp("", "anybakup-git-test-*")
	// Initialize git repo
	// if err := Initgit(); err != nil {
	// 	t.Fatalf("Initgit failed: %v", err)
	// }

	// Create a test file in the repo
	dir1 := filepath.Join(tmpDir, "dir1")
	if err := os.MkdirAll(dir1, 0755); err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}
	testFile := filepath.Join(dir1, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	repodir, err := r.CopyToRepo(util.SrcPath(testFile))
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	ret, err := r.GitAddFile(repodir)
	if err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}
	if ret.Action != util.GitResultTypeAdd {
		t.Errorf("Expected GitResultAdd, got %v", ret)
	}
	if err := cmd.BakupOptAdd(testFile, repodir.Sting(), false); err != nil {
		t.Errorf("Expected GitResultAdd, got %v", ret)
	}
}

func TestGitAddDir(t *testing.T) {
	r, err := util.NewGitReop()
	if err != nil {
		t.Fatalf("NewGitReop failed: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "anybakup-git-test-*")
	// Initialize git repo
	// if err := Initgit(); err != nil {
	// 	t.Fatalf("Initgit failed: %v", err)
	// }

	// Create a test file in the repo
	dir1 := filepath.Join(tmpDir, "dir1")
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
	srcdir := util.SrcPath(dir1)
	repodir, err := r.CopyToRepo(srcdir)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	ret, err := r.GitAddFile(repodir)
	if err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}
	if ret.Action != util.GitResultTypeAdd {
		t.Errorf("Expected GitResultAdd, got %v", ret)
	}
	if err := cmd.BakupOptAdd(dir1, repodir.Sting(), false); err != nil {
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
