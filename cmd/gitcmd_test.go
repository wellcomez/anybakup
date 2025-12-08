package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"anybakup/util"
)

// setupTestEnv creates a temporary test environment with config and git repo
func setupTestEnv(t *testing.T) (repoDir string, config *util.Config, cleanup func()) {
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

	return repoDir, &util.Config{RepoDir: util.RepoRoot(repoDir)}, cleanup
}

// TestAddFile tests adding a file to the repository
func TestAddFile(t *testing.T) {
	_, c, cleanup := setupTestEnv(t)
	tmpDir, err := os.MkdirTemp("", "anybakup-cmd-test-*")
	if err != nil {
		t.Error("temp file error", err)
	}
	defer cleanup()
	g := NewGitCmd("")
	g.C = c
	test1txt := filepath.Join(tmpDir, "1.txt")
	os.WriteFile(test1txt, []byte("xxx"), 0755)
	ret := g.AddFile(test1txt)
	if ret.Err != nil {
		t.Error("add file error", ret.Err)
	}

	if testfilepath := filepath.Join(tmpDir, "a"); os.MkdirAll(testfilepath, 0755) == nil {
		os.WriteFile(filepath.Join(testfilepath, "1.txt"), []byte("xxx"), 0755)
		ret := g.AddFile(testfilepath)
		if ret.Err != nil {
			t.Error("add file error", ret.Err)
		}
		if err := g.RmFileAbs(testfilepath); err != nil {
			t.Error("rm file error", err)
		}
		if err := g.RmFileAbs(test1txt); err != nil {
			t.Error("rm file error", err)
		}

		ret = g.AddFile(testfilepath)
		if ret.Err != nil {
			t.Error("add file error", ret.Err)
		}
	}

}
