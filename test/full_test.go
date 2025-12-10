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
func setupTestEnv(t *testing.T) (ret cmd.GitCmd, cleanup func()) {
	// Create temporary directories
	tmpDir, err := os.MkdirTemp("", "anybakup-cmd-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	repoDir := filepath.Join(tmpDir, "repo")
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
	re := cmd.GitCmd{C: &util.Config{RepoDir: util.RepoRoot(repoDir)}}
	return re, cleanup
}
func TestGitAddFile(t *testing.T) {
	g, clean := setupTestEnv(t)
	defer clean()
	repo := g.C.RepoDir
	fmt.Printf("t: %v\n", repo)
	r, err := util.NewGitReop(g.C)
	if err != nil {
		t.Fatalf("NewGitReop failed: %v", err)
	}

	tmpDir, _ := os.MkdirTemp("", "anybakup-git-test-*")
	defer os.RemoveAll(tmpDir)
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
		if err := cmd.BakupOptAdd(testFile, v, false, false, g); err != nil {
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
		if err := cmd.BakupOptAdd(testFile, v, false, false, g); err != nil {
			t.Errorf("Expected GitResultAdd, got %v", err)
		}
	}
}

func TestGitAddDir(t *testing.T) {
	g, clean := setupTestEnv(t)
	repo := g.C.RepoDir
	defer clean()
	fmt.Printf("t: %v\n", repo)
	r, err := util.NewGitReop(g.C)
	if err != nil {
		t.Fatalf("NewGitReop failed: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "anybakup-git-test-*")
	if err != nil {
		t.Fatalf(" create tmepdir dir: %v", err)
		return
	}
	defer os.RemoveAll(tmpDir)
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
	if err := cmd.BakupOptAdd(dir1, repodir, false, false, g); err != nil {
		t.Errorf("Expected GitResultAdd, got %v", ret)
	}
	for _, v := range ret.Files {
		if err := cmd.BakupOptAdd(fmt.Sprintf("/%v", v), v, true, true, g); err != nil {
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

	if err := cmd.BakupOptAdd(dir1, repodir, false, false, g); err != nil {
		t.Errorf("Expected GitResultAdd, got %v", ret)
	}
	for _, v := range ret.Files {
		if err := cmd.BakupOptAdd(fmt.Sprintf("/%v", v), v, true, true, g); err != nil {
			t.Errorf("Expected GitResultAdd, got %v", ret)
		}
	}
	_, err = g.GetFileLog(repodir)
	if err == nil {
		t.Errorf("Expected GitResultAdd, got %v", ret)
	}
	// if len(logs) != 1 {
	// 	t.Errorf("Expected GitResultAdd, got %v", ret)
	// }
	for _, v := range ret.Files {
		logs, err := g.GetFileLog(v)
		if err != nil {
			t.Errorf("Expected GitResultAdd, got %v", ret)
		}
		if len(logs) != 2 {
			t.Errorf("Expected GitResultAdd, got %v", ret)
		}
		for _, h := range logs {
			target := filepath.Join(tmpDir, h.Commit)
			if err := g.GetFile(v, h.Commit, target); err != nil {
				t.Errorf("Expected GitResultAdd, got %v", ret)
			}
		}
	}

}
func TestGitRmDir(t *testing.T) {
	g, clean := setupTestEnv(t)
	repo := g.C.RepoDir
	defer clean()
	fmt.Printf("t: %v\n", repo)
	r, err := util.NewGitReop(g.C)
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
	repodir := setupAddDir(t, dir1, g, r, "aaa")
	///change file content
	setupAddDir(t, dir1, g, r, "bbb")

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
	if err := cmd.BakupOptRm(repodir, g.C); err != nil {
		t.Errorf("Expected GitResultRm, got %v", rmret)
	}
	for _, v := range rmret.Files {
		if err := cmd.BakupOptRm(v, g.C); err != nil {
			t.Errorf("Expected GitResultRm, got %v", rmret)
		}
	}
}

func setupAddDir(t *testing.T, dir1 string, g cmd.GitCmd, r *util.GitRepo, content string) util.RepoPath {
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
	if err := cmd.BakupOptAdd(dir1, repodir, false, false, g); err != nil {
		t.Errorf("Expected GitResultAdd, got %v", ret)
	}
	for _, v := range ret.Files {
		if err := cmd.BakupOptAdd(string(v.ToSrc()), v, true, true, g); err != nil {
			t.Errorf("Expected GitResultAdd, got %v", ret)
		}
	}
	return repodir
}
func TestAddFolder(t *testing.T) {
	g, cleanup := setupTestEnv(t)
	defer cleanup()
	tmpDir, err := os.MkdirTemp("", "anybakup-cmd-test-*")
	if err != nil {
		t.Error("temp file error", err)
	}
	defer os.RemoveAll(tmpDir)
	test1txt := filepath.Join(tmpDir, "1.txt")
	if err := os.WriteFile(test1txt, []byte("xxx"), 0755); err != nil {
		t.Error("write file error", err)
	}
	test2txt := filepath.Join(tmpDir, "2.txt")
	if err := os.WriteFile(test2txt, []byte("xxx"), 0755); err != nil {
		t.Error("write file error", err)
	}
	ret := g.AddFile(tmpDir)

	if ret.Err != nil {
		t.Error("add file error", ret.Err)
	}
	logs, err := cmd.GetAllOpt(g.C)
	if err != nil {
		t.Error("get log", err)
	}
	for _, v := range logs {
		fmt.Println(v.DestFile)
	}

	if err := os.WriteFile(test2txt, []byte("xxxssss"), 0755); err != nil {
		t.Error("write file error", err)
	}
	ret = g.AddFile(test2txt)
	if ret.Err != nil {
		t.Error("add file error", ret.Err)
	}
	logs, err = cmd.GetAllOpt(g.C)
	if err != nil {
		t.Error("get log", err)
	}
	for _, v := range logs {
		fmt.Print(v.DestFile)
	}

}

func TestSqlGetRoot(t *testing.T) {
	g, cleanup := setupTestEnv(t)
	defer cleanup()
	tmpDir, err := os.MkdirTemp("", "anybakup-cmd-test-*")
	if err != nil {
		t.Error("temp file error", err)
	}
	defer os.RemoveAll(tmpDir)
	test1 := util.SrcPath(filepath.Join(tmpDir, "1"))
	if err := os.WriteFile(test1.Sting(), []byte("xxx"), 0644); err != nil {
		t.Error("write file error", err)
	}
	test2 := util.SrcPath(filepath.Join(tmpDir, "2"))
	if err := os.WriteFile(test2.Sting(), []byte("xxx"), 0644); err != nil {
		t.Error("write file error", err)
	}
	if ret := g.AddFile(tmpDir); ret.Err != nil {
		t.Error("add file error", err)
	}
	if ret, err := cmd.GetRepoRoot(test1.Repo(), test1.Sting(), g.C); err != nil {
		t.Error("get repo root error", err)
	} else if ret == nil {
		t.Error("get repo root error ret==nil", ret)
	}
	if ret, err := cmd.GetRepoRoot(test1.Repo(), "z:\\zzz", g.C); err == nil {
		t.Error("get repo root error", err)
	} else if ret != nil {
		t.Error("get repo root error ret!=nil", ret)
	}

}
