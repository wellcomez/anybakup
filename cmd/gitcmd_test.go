package cmd

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
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
func TestHideAddFile(t *testing.T) {
	_, c, cleanup := setupTestEnv(t)
	tmpDir, err := os.MkdirTemp("", "anybakup-cmd-test-*")
	if err != nil {
		t.Error("temp file error", err)
	}
	defer cleanup()
	g := NewGitCmd("")
	g.C = c
	test1txt := filepath.Join(tmpDir, ".1.txt")
	if err := os.WriteFile(test1txt, []byte("xxx"), 0755); err != nil {
		t.Error("write file error", err)
	}
	ret := g.AddFile(test1txt)
	if ret.Err != nil {
		t.Error("add file error", ret.Err)
	}

	if testfilepath := filepath.Join(tmpDir, "a"); os.MkdirAll(testfilepath, 0755) == nil {
		if err := os.WriteFile(filepath.Join(testfilepath, "1.txt"), []byte("xxx"), 0755); err != nil {
			t.Error("write file error", err)
		}
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
	// h, _ := os.UserHomeDir()
	file := path.Join("/home/z", ".bashrc")
	ret = g.AddFile(file)
	if ret.Result != util.GitResultTypeAdd {
		t.Error("fail to add", ret.Err)
	}

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
	if err := os.WriteFile(test1txt, []byte("xxx"), 0755); err != nil {
		t.Error("write file error", err)
	}
	ret := g.AddFile(test1txt)
	if ret.Err != nil {
		t.Error("add file error", ret.Err)
	}

	if testfilepath := filepath.Join(tmpDir, "a"); os.MkdirAll(testfilepath, 0755) == nil {
		if err := os.WriteFile(filepath.Join(testfilepath, "1.txt"), []byte("xxx"), 0755); err != nil {
			t.Error("write file error", err)
		}
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
func TestAddDirTag(t *testing.T) {
	_, c, cleanup := setupTestEnv(t)
	tmpDir, err := os.MkdirTemp("", "anybakup-cmd-test-*")
	if err != nil {
		t.Error("temp file error", err)
	}
	defer os.RemoveAll(tmpDir)
	defer cleanup()

	g := NewGitCmd("")
	g.C = c
	test1txt := filepath.Join(tmpDir, "1.txt")
	if err := os.WriteFile(test1txt, []byte("xxx"), 0755); err != nil {
		t.Error("write file error", err)
	}

	test2txt := filepath.Join(tmpDir, "2.txt")
	if err := os.WriteFile(test2txt, []byte("xxx"), 0755); err != nil {
		t.Error("write file error", err)
	}
	dira := filepath.Join(tmpDir, "a")
	if err := os.MkdirAll(dira, 0755); err != nil {
		t.Error("mkdir error", err)
	}

	ret := g.AddFile(tmpDir, "tag")
	if ret.Err != nil {
		t.Error("add file error", ret.Err)
	}
	if len(ret.Files) == 0 {
		t.Error("add file error", ret.Err)
	}
	srcDir := util.SrcPath(tmpDir).Repo()
	if tag, err := GetFileTag(srcDir, g.C); err != nil {
		t.Error("add file error", ret.Err)
	} else if tag != "tag" {
		t.Error("add file error", ret.Err)
	}
	for _, f := range ret.Files {
		if tag, err := GetFileTag(f, g.C); err != nil {
			t.Error("add file error", ret.Err)
		} else if tag != "tag" {
			t.Error("add file error", ret.Err)
		}
	}
	const tagb = "bbbb"
	err = SetFileTag(srcDir, tagb, g.C)
	if err != nil {
		t.Error("add file error", ret.Err)
	}
	if tag, err := GetFileTag(srcDir, g.C); err != nil {
		t.Error("add file error", ret.Err)
	} else if tag != tagb {
		t.Error("add file error", ret.Err)
	}
	for _, f := range ret.Files {
		if tag, err := GetFileTag(f, g.C); err != nil {
			t.Error("add file error", ret.Err)
		} else if tag != tagb {
			t.Error("add file error", ret.Err)
		}
	}
}
func TestAddDir(t *testing.T) {
	_, c, cleanup := setupTestEnv(t)
	tmpDir, err := os.MkdirTemp("", "anybakup-cmd-test-*")
	if err != nil {
		t.Error("temp file error", err)
	}
	defer os.RemoveAll(tmpDir)
	defer cleanup()

	g := NewGitCmd("")
	g.C = c
	test1txt := filepath.Join(tmpDir, "1.txt")
	if err := os.WriteFile(test1txt, []byte("xxx"), 0755); err != nil {
		t.Error("write file error", err)
	}

	test2txt := filepath.Join(tmpDir, "2.txt")
	if err := os.WriteFile(test2txt, []byte("xxx"), 0755); err != nil {
		t.Error("write file error", err)
	}
	dira := filepath.Join(tmpDir, "a")
	if err := os.MkdirAll(dira, 0755); err != nil {
		t.Error("mkdir error", err)
	}

	ret := g.AddFile(tmpDir)
	if ret.Err != nil {
		t.Error("add file error", ret.Err)
	}

}

func TestGitFile(t *testing.T) {
	if runtime.GOOS != "windows" {
		return
	}
	_, c, cleanup := setupTestEnv(t)
	defer cleanup()
	g := NewGitCmd("")
	g.C = c
	const test1txt = "E:\\anybakup\\ui\\any\\lib"
	ret := g.AddFile(test1txt)
	if ret.Err != nil {
		t.Error("add file error", ret.Err)
	}
}
