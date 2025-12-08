package util

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// setupTestEnv creates a temporary test environment with config and test files
func setupTestEnv(t *testing.T) (repoDir string, cleanup func()) {
	// Create temporary directories
	tmpDir, err := os.MkdirTemp("", "anybakup-test-*")
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

// TestCopyToRepo_SingleFile tests copying a single file
func TestCopyToRepo_SingleFile(t *testing.T) {
	iswindows := runtime.GOOS == "windows"
	if iswindows {
		test_copy_window(t)
	} else {
		test_copy_linux(t)
	}
}
func test_copy_window(t *testing.T) {
	repoDir, cleanup := setupTestEnv(t)
	r := GitRepo{
		root: repoDir,
	}
	defer cleanup()

	// Create a test file
	testDir, err := os.MkdirTemp("", "test-src-*")
	if err != nil {
		t.Fatalf("Failed to create test dir: %v", err)
	}
	defer os.RemoveAll(testDir)

	testFile := filepath.Join(testDir, "test.txt")
	testContent := "Hello, World!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Copy the file
	d, err := r.CopyToRepo(SrcPath(testFile))
	if err != nil {
		t.Fatalf("CopyToRepo failed: %v", err)
	}
	dest := d.ToAbs(r)
	// Verify destination path
	expectedDest := filepath.Join(repoDir, string(r.Src2Repo(testFile)))
	if dest != expectedDest {
		t.Errorf("Expected dest %s, got %s", expectedDest, dest)
	}

	// Verify file exists and has correct content
	content, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("Failed to read copied file: %v", err)
	}
	if string(content) != testContent {
		t.Errorf("Expected content %q, got %q", testContent, string(content))
	}

	// Verify permissions
	srcInfo, _ := os.Stat(testFile)
	dstInfo, _ := os.Stat(dest)
	if srcInfo.Mode() != dstInfo.Mode() {
		t.Errorf("Permissions not preserved: src=%v, dst=%v", srcInfo.Mode(), dstInfo.Mode())
	}
}
func test_copy_linux(t *testing.T) {
	repoDir, cleanup := setupTestEnv(t)
	r := GitRepo{
		root: repoDir,
	}
	defer cleanup()

	// Create a test file
	testDir, err := os.MkdirTemp("", "test-src-*")
	if err != nil {
		t.Fatalf("Failed to create test dir: %v", err)
	}
	defer os.RemoveAll(testDir)

	testFile := filepath.Join(testDir, "test.txt")
	testContent := "Hello, World!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Copy the file
	d, err := r.CopyToRepo(SrcPath(testFile))
	if err != nil {
		t.Fatalf("CopyToRepo failed: %v", err)
	}
	dest := d.ToAbs(r)
	// Verify destination path
	expectedDest := filepath.Join(repoDir, testFile[1:])
	if dest != expectedDest {
		t.Errorf("Expected dest %s, got %s", expectedDest, dest)
	}

	// Verify file exists and has correct content
	content, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("Failed to read copied file: %v", err)
	}
	if string(content) != testContent {
		t.Errorf("Expected content %q, got %q", testContent, string(content))
	}

	// Verify permissions
	srcInfo, _ := os.Stat(testFile)
	dstInfo, _ := os.Stat(dest)
	if srcInfo.Mode() != dstInfo.Mode() {
		t.Errorf("Permissions not preserved: src=%v, dst=%v", srcInfo.Mode(), dstInfo.Mode())
	}
}

// TestCopyToRepo_Directory tests copying a directory recursively
func TestCopyToRepo_Directory(t *testing.T) {
	repoDir, cleanup := setupTestEnv(t)
	r := GitRepo{
		root: repoDir,
	}
	defer cleanup()

	// Create a test directory structure
	testDir, err := os.MkdirTemp("", "test-src-*")
	if err != nil {
		t.Fatalf("Failed to create test dir: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Create subdirectories and files
	subDir := filepath.Join(testDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	file1 := filepath.Join(testDir, "file1.txt")
	file2 := filepath.Join(subDir, "file2.txt")

	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatalf("Failed to write file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
		t.Fatalf("Failed to write file2: %v", err)
	}

	// Copy the directory
	gitpath, err := r.CopyToRepo(SrcPath(testDir))
	if err != nil {
		t.Fatalf("CopyToRepo failed: %v", err)
	}
	dest := gitpath.ToAbs(r)
	// Verify destination path
	expectedDest := filepath.Join(repoDir, string(SrcPath(testDir).Repo()))
	if dest != expectedDest {
		t.Errorf("Expected dest %s, got %s", expectedDest, dest)
	}

	// Verify directory exists
	if info, err := os.Stat(dest); err != nil || !info.IsDir() {
		t.Fatalf("Destination directory not created properly")
	}

	// Verify files exist with correct content
	destFile1 := filepath.Join(dest, "file1.txt")
	destFile2 := filepath.Join(dest, "subdir", "file2.txt")

	content1, err := os.ReadFile(destFile1)
	if err != nil {
		t.Fatalf("Failed to read copied file1: %v", err)
	}
	if string(content1) != "content1" {
		t.Errorf("Expected content1, got %q", string(content1))
	}

	content2, err := os.ReadFile(destFile2)
	if err != nil {
		t.Fatalf("Failed to read copied file2: %v", err)
	}
	if string(content2) != "content2" {
		t.Errorf("Expected content2, got %q", string(content2))
	}
}

// TestCopyToRepo_NonExistentSource tests error handling for non-existent source
func TestCopyToRepo_NonExistentSource(t *testing.T) {
	dir, cleanup := setupTestEnv(t)
	r := GitRepo{
		root: dir,
	}
	defer cleanup()

	// Try to copy a non-existent file
	_, err := r.CopyToRepo("/non/existent/file.txt")
	if err == nil {
		t.Fatal("Expected error for non-existent source, got nil")
	}
}

// TestCopyToRepo_NestedDirectories tests copying deeply nested directories
func TestCopyToRepo_NestedDirectories(t *testing.T) {
	repoDir, cleanup := setupTestEnv(t)
	r := GitRepo{root: repoDir}
	defer cleanup()

	// Create a deeply nested directory structure
	testDir, err := os.MkdirTemp("", "test-src-*")
	if err != nil {
		t.Fatalf("Failed to create test dir: %v", err)
	}
	defer os.RemoveAll(testDir)

	deepDir := filepath.Join(testDir, "a", "b", "c", "d")
	if err := os.MkdirAll(deepDir, 0755); err != nil {
		t.Fatalf("Failed to create deep dir: %v", err)
	}

	deepFile := filepath.Join(deepDir, "deep.txt")
	if err := os.WriteFile(deepFile, []byte("deep content"), 0644); err != nil {
		t.Fatalf("Failed to write deep file: %v", err)
	}

	// Copy the directory
	_, err = r.CopyToRepo(SrcPath(testDir))
	if err != nil {
		t.Fatalf("CopyToRepo failed: %v", err)
	}

	// Verify deep file exists
	destDeepFile := filepath.Join(repoDir, string(SrcPath(testDir).Repo()), "a", "b", "c", "d", "deep.txt")
	content, err := os.ReadFile(destDeepFile)
	if err != nil {
		t.Fatalf("Failed to read deep file: %v", err)
	}
	if string(content) != "deep content" {
		t.Errorf("Expected 'deep content', got %q", string(content))
	}
}

// TestCopyFile tests the copyFile helper function
func TestCopyFile(t *testing.T) {
	// Create temporary source file
	tmpDir, err := os.MkdirTemp("", "copy-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	srcFile := filepath.Join(tmpDir, "source.txt")
	dstFile := filepath.Join(tmpDir, "dest", "destination.txt")

	testContent := "test content"
	if err := os.WriteFile(srcFile, []byte(testContent), 0600); err != nil {
		t.Fatalf("Failed to write source file: %v", err)
	}

	// Copy the file
	if err := copyFile(srcFile, dstFile); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	// Verify content
	content, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	if string(content) != testContent {
		t.Errorf("Expected %q, got %q", testContent, string(content))
	}

	// Verify permissions
	srcInfo, _ := os.Stat(srcFile)
	dstInfo, _ := os.Stat(dstFile)
	if srcInfo.Mode() != dstInfo.Mode() {
		t.Errorf("Permissions not preserved: src=%v, dst=%v", srcInfo.Mode(), dstInfo.Mode())
	}
}

// TestCopyDir tests the copyDir helper function
func TestCopyDir(t *testing.T) {
	// Create temporary source directory
	tmpDir, err := os.MkdirTemp("", "copy-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	srcDir := filepath.Join(tmpDir, "source")
	dstDir := filepath.Join(tmpDir, "destination")

	// Create source structure
	if err := os.MkdirAll(filepath.Join(srcDir, "sub"), 0755); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "file.txt"), []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "sub", "file2.txt"), []byte("content2"), 0644); err != nil {
		t.Fatalf("Failed to write file2: %v", err)
	}

	// Copy the directory
	if err := copyDir(srcDir, dstDir); err != nil {
		t.Fatalf("copyDir failed: %v", err)
	}

	// Verify files exist
	if _, err := os.Stat(filepath.Join(dstDir, "file.txt")); err != nil {
		t.Errorf("file.txt not copied: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dstDir, "sub", "file2.txt")); err != nil {
		t.Errorf("sub/file2.txt not copied: %v", err)
	}
}
