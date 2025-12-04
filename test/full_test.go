package cmd

import (
	"fmt"
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
	if len(ret.Files) != 1 {
		t.Errorf("Expected GitResultAdd, got %v", ret.Files)
	}

	for _, v := range ret.Files {
		if err := cmd.BakupOptAdd(testFile, v, false, false); err != nil {
			t.Errorf("Expected GitResultAdd, got %v", err)
		}
	}

	if err := os.WriteFile(testFile, []byte("test content xx"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	repodir, err = r.CopyToRepo(util.SrcPath(testFile))
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	ret, err = r.GitAddFile(repodir)
	if err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}

	for _, v := range ret.Files {
		if err := cmd.BakupOptAdd(testFile, v, false, false); err != nil {
			t.Errorf("Expected GitResultAdd, got %v", err)
		}
	}
}

func TestGitAddDir(t *testing.T) {
	repo, clean := setupTestEnv(t)
	defer clean()
	fmt.Printf("t: %v\n", repo)
	r, err := util.NewGitReop()
	if err != nil {
		t.Fatalf("NewGitReop failed: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "anybakup-git-test-*")
	if err != nil {
		t.Fatalf(" create tmepdir dir: %v", err)
		return
	}
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
	if err := cmd.BakupOptAdd(dir1, repodir.Sting(), false, false); err != nil {
		t.Errorf("Expected GitResultAdd, got %v", ret)
	}
	for _, v := range ret.Files {
		if err := cmd.BakupOptAdd(fmt.Sprintf("/%v", v), v, true, true); err != nil {
			t.Errorf("Expected GitResultAdd, got %v", ret)
		}
	}

	testFile2 = filepath.Join(dir1, "test2.txt")
	if err := os.WriteFile(testFile2, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add the file (using relative path from repo root)
	srcdir = util.SrcPath(dir1)
	repodir, err = r.CopyToRepo(srcdir)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	ret, err = r.GitAddFile(repodir)
	if err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}
	if ret.Action != util.GitResultTypeNochange {
		t.Errorf("Expected GitResultAdd, got %v", ret)
	}

	testFile2 = filepath.Join(dir1, "test2.txt")
	if err := os.WriteFile(testFile2, []byte("test contentxxx"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add the file (using relative path from repo root)
	srcdir = util.SrcPath(dir1)
	repodir, err = r.CopyToRepo(srcdir)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	ret, err = r.GitAddFile(repodir)
	if err != nil {
		t.Fatalf("GitAddFile failed: %v", err)
	}
	if ret.Action != util.GitResultTypeAdd {
		t.Errorf("Expected GitResultAdd, got %v", ret)
	}

	if err := cmd.BakupOptAdd(dir1, repodir.Sting(), false, false); err != nil {
		t.Errorf("Expected GitResultAdd, got %v", ret)
	}
	for _, v := range ret.Files {
		if err := cmd.BakupOptAdd(fmt.Sprintf("/%v", v), v, true, true); err != nil {
			t.Errorf("Expected GitResultAdd, got %v", ret)
		}
	}
}
func TestGitRmDir(t *testing.T) {
	repo, clean := setupTestEnv(t)
	defer clean()
	fmt.Printf("t: %v\n", repo)
	r, err := util.NewGitReop()
	if err != nil {
		t.Fatalf("NewGitReop failed: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "anybakup-git-test-*")
	if err != nil {
		t.Fatalf(" create tmepdir dir: %v", err)
		return
	}
	dir1 := filepath.Join(tmpDir, "dir1")
	if err := os.MkdirAll(dir1, 0755); err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}

	// Add the file (using relative path from repo root)
	repodir := setupAddDir(t, dir1, r, "aaa")
	///change file content
	setupAddDir(t, dir1, r, "bbb")

	rmret, err := r.GitRmFile(repodir)
	if err != nil {
		t.Errorf("GitRmFile failed: %v", err)
	}
	if rmret.Action != util.GitResultTypeRm {
		t.Errorf("Expected GitResultRm, got %v", rmret)
	}
	if len(rmret.Files) != 2 {
		t.Errorf("Expected 2 file, got %v", len(rmret.Files))
	}
	if err := cmd.BakupOptRm(dir1); err != nil {
		t.Errorf("Expected GitResultRm, got %v", rmret)
	}
	for _, v := range rmret.Files {
		if err := cmd.BakupOptRm(fmt.Sprintf("/%v", v)); err != nil {
			t.Errorf("Expected GitResultRm, got %v", rmret)
		}
	}
}

func setupAddDir(t *testing.T, dir1 string, r *util.GitRepo, content string) util.RepoPath {
	testFile := filepath.Join(dir1, "test.txt")
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	testFile2 := filepath.Join(dir1, "test2.txt")
	if err := os.WriteFile(testFile2, []byte(content), 0644); err != nil {
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
	if err := cmd.BakupOptAdd(dir1, repodir.Sting(), false, false); err != nil {
		t.Errorf("Expected GitResultAdd, got %v", ret)
	}
	for _, v := range ret.Files {
		if err := cmd.BakupOptAdd(fmt.Sprintf("/%v", v), v, true, true); err != nil {
			t.Errorf("Expected GitResultAdd, got %v", ret)
		}
	}
	return repodir
}
