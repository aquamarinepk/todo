package am

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
)

// DebugFS prints the tree structure of the given embedded filesystem.
func DebugFS(efs embed.FS, root string) error {
	return fs.WalkDir(efs, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(root, path)
		fmt.Println(relPath)
		return nil
	})
}
