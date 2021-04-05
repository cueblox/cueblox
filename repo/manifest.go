package repo

import (
	"encoding/json"
	"os"
	"path"

	"github.com/devrel-blox/cueblox/config"
	"github.com/pterm/pterm"
)

// createManifest writes a JSON encoded manifest
// file with the provided namespace at the directory
// specified by `root`. Root path provided for testing
// purposes.
func create(namespace, root string) error {
	var m config.Manifest
	m.Namespace = namespace
	manifest := path.Join(root, "manifest.json")
	bb, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return os.WriteFile(manifest, bb, 0755)

}

// CreateManifest writes a JSON encoded file
// with the provided namespace to a repository's
// root directory
func CreateManifest(namespace string) error {
	root, err := os.Getwd()
	if err != nil {
		pterm.Error.Println("Error creating repository root directory")
		return err
	}
	repo := path.Join(root, "repository")
	return create(namespace, repo)
}

func get(root string) (config.Manifest, error) {
	manifest := path.Join(root, "manifest.json")
	bb, err := os.ReadFile(manifest)
	if err != nil {
		pterm.Error.Println("Error reading manifest.json")
		return config.Manifest{}, err
	}
	var m config.Manifest
	err = json.Unmarshal(bb, &m)
	if err != nil {
		pterm.Error.Println("Error parsing manifest.json")
		return config.Manifest{}, err
	}
	return m, nil
}

func GetManifest() (config.Manifest, error) {
	root, err := os.Getwd()
	if err != nil {
		pterm.Error.Println("Error getting repository root directory")
		return config.Manifest{}, err
	}
	repo := path.Join(root, "repository")
	return get(repo)
}
