package cmd

import (
	"os"
	"path/filepath"
	"testing"

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

// TestAddFile tests adding a file to the repository
func TestAddFile(t *testing.T) {
	repoDir, cleanup := setupTestEnv(t)
	defer cleanup()

	// Initialize repo
	repo, err := util.NewGitReop()
	if err != nil {
		t.Fatalf("NewGitReop failed: %v", err)
	}

	// Create a test file outside the repo
	testFile := filepath.Join(os.TempDir(), "test_add_file.txt")
	testContent := "test content for add"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	// Test AddFile
	result := AddFile(testFile)
	if result.Err != nil {
		t.Fatalf("AddFile failed: %v", result.Err)
	}

	if result.Result != util.GitResultTypeAdd {
		t.Errorf("Expected GitResultAdd, got %v", result.Result)
	}

	if result.Dest == "" {
		t.Error("Expected non-empty dest path")
	}

	// Verify file was copied to repo
	copiedFile := filepath.Join(repoDir, result.Dest)
	content, err := os.ReadFile(copiedFile)
	if err != nil {
		t.Fatalf("Failed to read copied file: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("Expected content %q, got %q", testContent, string(content))
	}

	// Test adding the same file again (should result in nochange)
	result2 := AddFile(testFile)
	if result2.Err != nil {
		t.Fatalf("AddFile second time failed: %v", result2.Err)
	}

	if result2.Result != util.GitResultTypeNochange {
		t.Errorf("Expected GitResultNochange, got %v", result2.Result)
	}

	_ = repo
}

// TestGetFileLog tests getting file commit history
func TestGetFileLog(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	// Create and add a test file
	testFile := filepath.Join(os.TempDir(), "test_log_file.txt")
	defer os.Remove(testFile)

	// First version
	if err := os.WriteFile(testFile, []byte("version 1"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	result1 := AddFile(testFile)
	if result1.Err != nil {
		t.Fatalf("AddFile failed: %v", result1.Err)
	}

	// Second version
	if err := os.WriteFile(testFile, []byte("version 2"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}
	result2 := AddFile(testFile)
	if result2.Err != nil {
		t.Fatalf("AddFile second time failed: %v", result2.Err)
	}

	// Third version
	if err := os.WriteFile(testFile, []byte("version 3"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}
	result3 := AddFile(testFile)
	if result3.Err != nil {
		t.Fatalf("AddFile third time failed: %v", result3.Err)
	}

	// Get log
	logs, err := GetFileLog(testFile)
	if err != nil {
		t.Fatalf("GetFileLog failed: %v", err)
	}

	if len(logs) != 3 {
		t.Errorf("Expected 3 log entries, got %d", len(logs))
	}

	// Verify log entries have required fields
	for i, log := range logs {
		if log.Commit == "" {
			t.Errorf("Log entry %d has empty commit hash", i)
		}
		if log.Author == "" {
			t.Errorf("Log entry %d has empty author", i)
		}
		if log.Date == "" {
			t.Errorf("Log entry %d has empty date", i)
		}
	}
}

// TestGetFile_WithCommit tests getting a file from a specific commit
func TestGetFile_WithCommit(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	// Create and add a test file
	testFile := filepath.Join(os.TempDir(), "test_get_file.txt")
	defer os.Remove(testFile)

	// First version
	version1Content := "version 1 content"
	if err := os.WriteFile(testFile, []byte(version1Content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	result1 := AddFile(testFile)
	if result1.Err != nil {
		t.Fatalf("AddFile failed: %v", result1.Err)
	}

	// Get the first commit hash
	logs, err := GetFileLog(testFile)
	if err != nil {
		t.Fatalf("GetFileLog failed: %v", err)
	}
	if len(logs) == 0 {
		t.Fatal("Expected at least one log entry")
	}
	firstCommit := logs[0].Commit

	// Second version
	version2Content := "version 2 content - modified"
	if err := os.WriteFile(testFile, []byte(version2Content), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}
	result2 := AddFile(testFile)
	if result2.Err != nil {
		t.Fatalf("AddFile second time failed: %v", result2.Err)
	}

	// Get the first version using commit hash
	outputFile := filepath.Join(os.TempDir(), "output_v1.txt")
	defer os.Remove(outputFile)

	err = GetFile(testFile, firstCommit, outputFile)
	if err != nil {
		t.Fatalf("GetFile failed: %v", err)
	}

	// Verify content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if string(content) != version1Content {
		t.Errorf("Expected content %q, got %q", version1Content, string(content))
	}
}

// TestGetFile_WithoutCommit tests getting a file from HEAD (no commit specified)
func TestGetFile_WithoutCommit(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	// Create and add a test file
	testFile := filepath.Join(os.TempDir(), "test_get_head.txt")
	defer os.Remove(testFile)

	// First version
	if err := os.WriteFile(testFile, []byte("version 1"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	result1 := AddFile(testFile)
	if result1.Err != nil {
		t.Fatalf("AddFile failed: %v", result1.Err)
	}

	// Second version (this will be HEAD)
	headContent := "version 2 - HEAD"
	if err := os.WriteFile(testFile, []byte(headContent), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}
	result2 := AddFile(testFile)
	if result2.Err != nil {
		t.Fatalf("AddFile second time failed: %v", result2.Err)
	}

	// Get the file from HEAD (empty commit string)
	outputFile := filepath.Join(os.TempDir(), "output_head.txt")
	defer os.Remove(outputFile)

	err := GetFile(testFile, "", outputFile)
	if err != nil {
		t.Fatalf("GetFile with empty commit failed: %v", err)
	}

	// Verify content matches HEAD
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if string(content) != headContent {
		t.Errorf("Expected HEAD content %q, got %q", headContent, string(content))
	}
}

// TestGetFile_ToNestedDirectory tests getting a file to a nested directory
func TestGetFile_ToNestedDirectory(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	// Create and add a test file
	testFile := filepath.Join(os.TempDir(), "test_nested.txt")
	defer os.Remove(testFile)

	testContent := "test content"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	result := AddFile(testFile)
	if result.Err != nil {
		t.Fatalf("AddFile failed: %v", result.Err)
	}

	// Get file to a nested directory that doesn't exist
	outputDir := filepath.Join(os.TempDir(), "nested", "deep", "directory")
	outputFile := filepath.Join(outputDir, "output.txt")
	defer os.RemoveAll(filepath.Join(os.TempDir(), "nested"))

	err := GetFile(testFile, "", outputFile)
	if err != nil {
		t.Fatalf("GetFile to nested directory failed: %v", err)
	}

	// Verify file exists and has correct content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("Expected content %q, got %q", testContent, string(content))
	}
}
