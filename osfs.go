package osfs

import (
	"io/fs"
)

type OsFS struct {
}

var OsRootFS fs.FS = OsFS{}
