package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyToRepo copies a file or directory from src (absolute path) to the repository directory,
// preserving the full path structure.
// For example:
//   - /a/b/c.file -> /repo/a/b/c.file (file)
//   - /a/b -> /repo/a/b (directory, copied recursively)
func CopyToRepo(src string) (string, error) {
	conf := Config{}
	if err := conf.Load(); err != nil {
		return "", fmt.Errorf("copytorepo error load config: %v", err)
	}

	// Verify source exists and get its info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return "", fmt.Errorf("copytorepo error stat src: %v", err)
	}

	// Create destination path by appending src path (without leading /) to repo dir
	dest := filepath.Join(conf.RepoDir, src[1:])

	// Copy based on whether src is a file or directory
	if srcInfo.IsDir() {
		err = copyDir(src, dest)
	} else {
		err = copyFile(src, dest)
	}

	if err != nil {
		return "", fmt.Errorf("copytorepo error copying: %v", err)
	}

	return dest, nil
}

// copyFile copies a single file from src to dst, preserving permissions
func copyFile(src, dst string) error {
	// Get source file info for permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create destination directory if it doesn't exist
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy file contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// Preserve permissions
	if err := os.Chmod(dst, srcInfo.Mode()); err != nil {
		return err
	}

	return nil
}

// copyDir recursively copies a directory from src to dst
func copyDir(src, dst string) error {
	// Get source directory info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create destination directory with same permissions
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// Read directory entries
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}
