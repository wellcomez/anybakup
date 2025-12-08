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
func (s RepoPath) UnitxStyle() RepoPath {
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
		//volx/ replace to c:\
		path := s.Sting()
		// Check if path starts with "vol" prefix (e.g., "volc\path\to\file")
		if strings.HasPrefix(path, "vol") && len(path) > 3 {
			// Extract drive letter (the character after "vol")
			driveLetter := string(path[3])
			// Get the rest of the path after "vol<letter>\"
			restOfPath := ""
			if len(path) > 5 && (path[4] == '\\' || path[4] == '/') {
				restOfPath = path[5:]
			} else if len(path) > 4 {
				restOfPath = path[4:]
			}
			// Reconstruct as "X:\path\to\file"
			return SrcPath(driveLetter + ":\\" + restOfPath)
		}
		return SrcPath(path)
	}
	return SrcPath(filepath.Join("/", s.Sting()))
}
