//go:build !windows

package osfs

import (
	"io/fs"
	"os"
	"path"
)

func fsToOsPath(fs_name string) (string, error) {
	if !fs.ValidPath(fs_name) {
		return "", fs.ErrInvalid
	}

	if fs_name == "." {
		return "/", nil
	}

	return path.Clean(path.Join("/", fs_name)), nil
}

func (f OsFS) Open(fs_name string) (fs.File, error) {
	os_name, err := fsToOsPath(fs_name)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: fs_name, Err: err}
	}
	return os.Open(os_name)
}

func (f OsFS) Stat(fs_name string) (fs.FileInfo, error) {
	os_name, err := fsToOsPath(fs_name)
	if err != nil {
		return nil, &fs.PathError{Op: "stat", Path: fs_name, Err: err}
	}
	return os.Stat(os_name)
}

func (f OsFS) ReadFile(fs_name string) ([]byte, error) {
	os_name, err := fsToOsPath(fs_name)
	if err != nil {
		return nil, &fs.PathError{Op: "readfile", Path: fs_name, Err: err}
	}
	return os.ReadFile(os_name)
}

func (f OsFS) ReadDir(fs_name string) ([]fs.DirEntry, error) {
	os_name, err := fsToOsPath(fs_name)
	if err != nil {
		return nil, &fs.PathError{Op: "readdir", Path: fs_name, Err: err}
	}
	return os.ReadDir(os_name)
}
