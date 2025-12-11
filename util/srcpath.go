package util

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type SrcPath string

func (s SrcPath) String() string {
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
			if disk != "" {
				rel, err = filepath.Rel(disk+"\\", string(s))
				if err == nil {
					// Get just the drive letter (without colon)
					driveLetter := disk[:len(disk)-1]
					rel = fmt.Sprintf("%s\\%s", driveLetter, rel)
				}
			} else {
				// Handle paths like C:\path where VolumeName returns ""
				if len(string(s)) > 1 && string(s)[1] == ':' {
					driveLetter := string(s)[0]
					rel, err = filepath.Rel(string(s)[:2], string(s))
					if err == nil {
						rel = fmt.Sprintf("%c%s", driveLetter, rel[1:]) // Skip the leading backslash
					}
				}
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
