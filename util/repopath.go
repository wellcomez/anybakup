package util

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type RepoPath string

func (s RepoPath) Sting() string {
	return string(s)
}

func (s RepoPath) ToAbs(repo GitRepo) string {
	r := RepoRoot(repo.root)
	return r.With(s.PlatformStyle().Sting())
}
func (s RepoPath) PlatformStyle() RepoPath {
	if runtime.GOOS == "windows" {
		return RepoPath(strings.Replace(s.Sting(), "/", "\\", -1))
	}
	return s
}
func (s RepoPath) UnixStyle() RepoPath {
	return RepoPath(strings.Replace(s.Sting(), "\\", "/", -1))
}
func (s RepoPath) ToSrc() (SrcPath, error) {
	if s.Sting() == "" {
		return SrcPath(""), nil
	}
	if filepath.IsAbs(s.Sting()) {
		return SrcPath(""), fmt.Errorf("nvalid path: %s", s.Sting())
	}
	if runtime.GOOS == "windows" {
		path := filepath.Clean(s.Sting())
		if strings.Index(path, "\\") == 1 {
			if len(path) > 1 && path[0] >= 'a' && path[0] <= 'z' || path[0] >= 'A' && path[0] <= 'Z' {
				return SrcPath(filepath.Clean(string(path[0]) + ":\\" + path[1:])), nil
			}
		}
		return SrcPath(""), fmt.Errorf("invalid path: %s",s.Sting())
	}
	return SrcPath(filepath.Clean(filepath.Join("/", s.Sting()))), nil
}
