package util

import (
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
	if runtime.GOOS == "windows" {
		return RepoPath(strings.Replace(s.Sting(), "\\", "/", -1))
	}
	return s
}
func (s RepoPath) ToSrc() SrcPath {
	if s.Sting() == "" {
		return SrcPath("")
	}
	if runtime.GOOS == "windows" {
		path := s.Sting()

		// Handle drive letter format (e.g., "c\path\to\file" -> "c:\path\to\file")
		if len(path) > 1 && path[1] == '\\' {
			// Path already has backslash after drive letter (e.g., "c:\path")
			return SrcPath(path)
		}
		if len(path) > 2 && path[1] == '/' {
			// Path has forward slash after drive letter (e.g., "c/path" -> "c:\path")
			return SrcPath(string(path[0]) + ":\\" + path[2:])
		}
		if len(path) > 1 && path[0] >= 'a' && path[0] <= 'z' || path[0] >= 'A' && path[0] <= 'Z' {
			// Single drive letter with path (e.g., "c\path" -> "c:\path")
			return SrcPath(string(path[0]) + ":\\" + path[1:])
		}

		return SrcPath(path)
	}
	return SrcPath(filepath.Join("/", s.Sting()))
}
