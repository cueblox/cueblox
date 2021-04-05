package repo

import (
	"os"
	"path"

	"github.com/pterm/pterm"
)

func createWithPath(root string) error {
	repo := path.Join(root, "repository")
	return os.MkdirAll(repo, 0755)
}

func Create() error {
	root, err := os.Getwd()
	if err != nil {
		pterm.Error.Println("Error creating repository root directory")
		return err
	}
	return createWithPath(root)
}
