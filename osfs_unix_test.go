//go:build !windows

package osfs_test

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/l4go/osfs"
)

func ExampleOsFS_Open_unix() {
	var err error
	_, err = osfs.OsRootFS.Open("/")
	fmt.Println(errors.Unwrap(err))
	_, err = osfs.OsRootFS.Open(".")
	fmt.Println(err == nil)
	_, err = osfs.OsRootFS.Open("etc/hosts")
	fmt.Println(err == nil)
	// Output:
	// invalid argument
	// true
	// true
}

func ExampleOsFS_Stat_unix() {
	var err error
	_, err = fs.Stat(osfs.OsRootFS, "/")
	fmt.Println(errors.Unwrap(err))
	_, err = fs.Stat(osfs.OsRootFS, ".")
	fmt.Println(err == nil)
	_, err = osfs.OsRootFS.Open("etc/hosts")
	fmt.Println(err == nil)
	// Output:
	// invalid argument
	// true
	// true
}

func ExampleOsFS_ReadFile_unix() {
	var err error
	_, err = fs.Stat(osfs.OsRootFS, "/")
	fmt.Println(errors.Unwrap(err))
	_, err = fs.ReadFile(osfs.OsRootFS, ".")
	fmt.Println(err == nil)
	_, err = fs.ReadFile(osfs.OsRootFS, "etc/hosts")
	fmt.Println(err == nil)
	// Output:
	// invalid argument
	// false
	// true
}

func ExampleOsFS_ReadDir_unix() {
	var err error
	_, err = fs.ReadDir(osfs.OsRootFS, "/")
	fmt.Println(errors.Unwrap(err))
	_, err = fs.ReadDir(osfs.OsRootFS, ".")
	fmt.Println(err == nil)
	_, err = fs.ReadDir(osfs.OsRootFS, "etc/hosts")
	fmt.Println(err == nil)
	// Output:
	// invalid argument
	// true
	// false
}
