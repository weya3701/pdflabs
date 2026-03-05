package filelist

import (
	"io/fs"
	"path/filepath"
)

func GetCurrentFiles(root string) []string {
	var vfiles []string = []string{}

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		vfiles = append(vfiles, path)
		return nil

	})

	return vfiles
}
