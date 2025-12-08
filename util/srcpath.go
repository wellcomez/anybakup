package util

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type SrcPath string

func (s SrcPath) Sting() string {
	return string(s)
}
func (s SrcPath) Repo() RepoPath {
	var rel string
	var err error

	if filepath.IsAbs(string(s)) {
		// Handle Windows drive letters properly
		if runtime.GOOS == "windows" {
			// Extract drive letter and get relative path within that drive
			disk := filepath.VolumeName(string(s))
			rel, err = filepath.Rel(disk+"\\", string(s))
			if err == nil {
				disk = strings.Replace(disk, ":", "", -1)
				rel = fmt.Sprintf("vol%s\\%s", disk, rel)
			}
		} else {
			rel, err = filepath.Rel("/", string(s))
		}
	} else {
		rel = string(s)
		err = nil
	}

	if err != nil {
		return ""
	}
	return RepoPath(rel)
}
