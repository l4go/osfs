package osfs

import (
	"errors"
	"io"
	"io/fs"
	"math/bits"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/sys/windows"
)

var ErrIsDirectory = errors.New("is a directory")

var drivePathRe = regexp.MustCompile(`^[a-zA-Z]:$`)
var windowsPathRe = regexp.MustCompile(`^[a-zA-Z]:`)

func fsToOsPath(fs_name string) (string, error) {
	if !fs.ValidPath(fs_name) {
		return "", fs.ErrInvalid
	}
	if fs_name == "." {
		return "", fs.ErrNotExist
	}

	if !windowsPathRe.MatchString(fs_name) {
		return "", fs.ErrNotExist
	}

	win_name := strings.ReplaceAll(fs_name, "/", `\`)
	if drivePathRe.MatchString(win_name) {
		win_name = win_name + `\`
	}

	return win_name, nil
}

func (f OsFS) Open(fs_name string) (fs.File, error) {
	if !fs.ValidPath(fs_name) {
		return nil, &fs.PathError{Op: "open", Path: fs_name, Err: fs.ErrInvalid}
	}

	if fs_name == "." {
		return newRootFile()
	}
	if !windowsPathRe.MatchString(fs_name) {
		return nil, &fs.PathError{Op: "open", Path: fs_name, Err: fs.ErrNotExist}
	}

	os_name, err := fsToOsPath(fs_name)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: fs_name, Err: err}
	}
	return os.Open(os_name)
}

func (f OsFS) Stat(fs_name string) (fs.FileInfo, error) {
	if !fs.ValidPath(fs_name) {
		return nil, &fs.PathError{Op: "stat", Path: fs_name, Err: fs.ErrInvalid}
	}

	if fs_name == "." {
		return dummyDirInfo("."), nil
	}

	if !windowsPathRe.MatchString(fs_name) {
		return nil, &fs.PathError{Op: "stat", Path: fs_name, Err: fs.ErrNotExist}
	}
	if drivePathRe.MatchString(fs_name) {
		drv := strings.ToUpper(fs_name)[0]
		if !found_drive(drv) {
			return nil, &fs.PathError{Op: "stat", Path: fs_name, Err: fs.ErrNotExist}
		}

		return dummyDirInfo(fs_name), nil
	}

	os_name, err := fsToOsPath(fs_name)
	if err != nil {
		return nil, &fs.PathError{Op: "stat", Path: fs_name, Err: err}
	}
	return os.Stat(os_name)
}

func (f OsFS) ReadFile(fs_name string) ([]byte, error) {
	if !fs.ValidPath(fs_name) {
		return nil, &fs.PathError{Op: "readfile", Path: fs_name, Err: fs.ErrInvalid}
	}

	if fs_name == "." {
		return nil, &fs.PathError{Op: "readfile", Path: fs_name, Err: ErrIsDirectory}
	}

	if !windowsPathRe.MatchString(fs_name) {
		return nil, &fs.PathError{Op: "readfile", Path: fs_name, Err: fs.ErrNotExist}
	}
	if drivePathRe.MatchString(fs_name) {
		return nil, &fs.PathError{Op: "readfile", Path: fs_name, Err: ErrIsDirectory}
	}

	os_name, err := fsToOsPath(fs_name)
	if err != nil {
		return nil, &fs.PathError{Op: "readfile", Path: fs_name, Err: err}
	}
	return os.ReadFile(os_name)
}

func (f OsFS) ReadDir(fs_name string) ([]fs.DirEntry, error) {
	if !fs.ValidPath(fs_name) {
		return nil, &fs.PathError{Op: "readdir", Path: fs_name, Err: fs.ErrInvalid}
	}

	if fs_name == "." {
		dent, err := get_drives()
		if err != nil {
			return nil, &fs.PathError{Op: "readddir", Path: ".", Err: err}
		}

		return dent, nil
	}

	if !windowsPathRe.MatchString(fs_name) {
		return nil, &fs.PathError{Op: "readdir", Path: fs_name, Err: fs.ErrNotExist}
	}

	os_name, err := fsToOsPath(fs_name)
	if err != nil {
		return nil, &fs.PathError{Op: "readdir", Path: fs_name, Err: err}
	}
	return os.ReadDir(os_name)
}

type rootFile struct {
	drives []fs.DirEntry
	offset int
}

func newRootFile() (*rootFile, error) {
	dent, err := get_drives()
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: ".", Err: err}
	}

	return &rootFile{drives: dent, offset: 0}, nil
}
func (rf *rootFile) Close() error {
	return nil
}

func (rf *rootFile) Stat() (fs.FileInfo, error) {
	return dummyDirInfo("."), nil
}

func (rf *rootFile) Read([]byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: ".", Err: ErrIsDirectory}
}

func (rf *rootFile) ReadDir(count int) ([]fs.DirEntry, error) {
	n := len(rf.drives) - rf.offset
	if n == 0 {
		if count <= 0 {
			return nil, nil
		}
		return nil, io.EOF
	}

	if count > 0 && n > count {
		n = count
	}

	list := make([]fs.DirEntry, n)
	for i := range list {
		list[i] = rf.drives[rf.offset+i]
	}
	rf.offset += n
	return list, nil
}

type dummyDirInfo string

func (dfi dummyDirInfo) Name() string {
	return string(dfi)
}
func (dfi dummyDirInfo) Size() int64 {
	return 0
}
func (dfi dummyDirInfo) Mode() fs.FileMode {
	return fs.ModeDir
}
func (dfi dummyDirInfo) ModTime() time.Time {
	return time.Time{}
}
func (dfi dummyDirInfo) IsDir() bool {
	return true
}
func (dfi dummyDirInfo) Sys() any {
	return nil
}

type driveDirEntry string

func (dde driveDirEntry) Name() string {
	return string(dde)
}
func (dde driveDirEntry) IsDir() bool {
	return true
}
func (dde driveDirEntry) Type() fs.FileMode {
	return fs.ModeDir
}
func (dde driveDirEntry) Info() (fs.FileInfo, error) {
	return dummyDirInfo(dde), nil
}

func found_drive(drv byte) bool {
	drv_bit, err := windows.GetLogicalDrives()
	if err != nil {
		return false
	}
	if drv < byte('A') || drv > byte('Z') {
		return false
	}

	return (drv_bit>>(drv-byte('A')))&0x1 != 0
}

func get_drives() ([]fs.DirEntry, error) {
	drv_bit, err := windows.GetLogicalDrives()
	if err != nil {
		return nil, err
	}

	drives := make([]fs.DirEntry, 0, bits.OnesCount32(drv_bit))
	for d := byte('A'); d <= byte('Z'); d++ {
		if drv_bit&0x1 != 0 {
			drives = append(drives, driveDirEntry(string([]byte{d, ':'})))
		}
		drv_bit >>= 1
	}

	return drives, err
}
