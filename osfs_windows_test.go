package osfs_test

import (
	"fmt"
	"io/fs"

	"github.com/l4go/osfs"
)

func ExampleOsFs_Open_win() {
	rf, err := osfs.OsRootFS.Open(".")
	if err != nil {
		return
	}
	defer rf.Close()

	rdf, ok := rf.(fs.ReadDirFile)
	if !ok {
		fmt.Println(err)
		return
	}

	dent, err := rdf.ReadDir(-1)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, d := range dent {
		if d.Name() == "C:" {
			fmt.Println(d.Name())
		}
	}

	f, err := osfs.OsRootFS.Open("C:")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	fmt.Println("opened")

	// Output:
	// C:
	// opened
}

func ExampleOsFs_Stat_win() {
	var err error
	_, err = fs.Stat(osfs.OsRootFS, ".")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("stated")

	_, err = fs.Stat(osfs.OsRootFS, "C:")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("stated")
	// Output:
	// stated
	// stated
}

func ExampleOsFs_ReadFile_win() {
	var err error
	_, err = fs.ReadFile(osfs.OsRootFS, ".")
	fmt.Println(err == nil)

	_, err = fs.ReadFile(osfs.OsRootFS, "C:")
	fmt.Println(err == nil)

	var buf []byte
	buf, err = fs.ReadFile(osfs.OsRootFS, "C:/Windows/system.ini")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(buf) > 0)

	// Output:
	// false
	// false
	// true
}

func ExampleOsFs_ReadDir_win() {
	dent, err := fs.ReadDir(osfs.OsRootFS, ".")
	if err != nil {
		return
	}
	for _, d := range dent {
		if d.Name() == "C:" {
			fmt.Println(d.Name())
		}
	}
	// Output:
	// C:
}
