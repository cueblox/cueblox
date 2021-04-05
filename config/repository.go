package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/otiai10/copy"
	"github.com/pterm/pterm"
)

// A repository is a combination of metadata
// and files in a specific folder structure
type Repository struct {
	Meta
	root     string
	Versions []*Version
}

func NewRepository(namespace string) error {
	r := &Repository{}
	r.Namespace = namespace

	pterm.Info.Println("Creating repository root directory")
	err := r.Create()
	if err != nil {
		return err
	}

	pterm.Info.Println("Creating manifest file")
	return r.CreateManifest()
}
func Open() (*Repository, error) {
	root, err := os.Getwd()
	if err != nil {
		pterm.Error.Println("Error creating repository root directory")
		return nil, err
	}
	repo := path.Join(root, "repository")
	return openWithPath(root, repo)
}
func openWithPath(root, path string) (*Repository, error) {
	r := &Repository{
		root: root,
	}
	m, err := r.getManifest(path)
	if err != nil {
		return nil, err
	}
	r.Namespace = m.Namespace
	//r.Versions
	err = r.LoadVersions()
	if err != nil {
		return nil, err
	}

	return r, nil
}
func (r *Repository) createWithPath(root string) error {
	repo := path.Join(root, "repository")
	return os.MkdirAll(repo, 0755)
}

func (r *Repository) Create() error {
	root, err := os.Getwd()
	if err != nil {
		pterm.Error.Println("Error creating repository root directory")
		return err
	}
	return r.createWithPath(root)
}

// createManifest writes a JSON encoded manifest
// file with the provided namespace at the directory
// specified by `root`. Root path provided for testing
// purposes.
func (r *Repository) createManifest(root string) error {
	err := r.LoadVersions()
	if err != nil {
		return err
	}
	var m Manifest
	m.Namespace = r.Namespace
	m.Versions = r.SchemaVersions()
	m.Schemas = make(map[string]*DisplayVersion, len(m.Versions))
	for _, v := range r.Versions {
		dv := &DisplayVersion{
			Name:      v.Name(),
			Namespace: v.Namespace(),
			Schemas:   v.Schemas,
		}
		m.Schemas[dv.Name] = dv
	}
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
func (r *Repository) CreateManifest() error {
	root, err := os.Getwd()
	if err != nil {
		pterm.Error.Println("Error creating repository root directory")
		return err
	}
	repo := path.Join(root, "repository")
	return r.createManifest(repo)
}

func (r *Repository) getManifest(root string) (Manifest, error) {
	manifest := path.Join(root, "manifest.json")
	bb, err := os.ReadFile(manifest)
	if err != nil {
		pterm.Error.Println("Error reading manifest.json")
		return Manifest{}, err
	}
	var m Manifest
	err = json.Unmarshal(bb, &m)
	if err != nil {

		pterm.Error.Println("Error parsing manifest.json")
		return Manifest{}, err
	}
	return m, nil
}

func (r *Repository) GetManifest() (Manifest, error) {
	root, err := os.Getwd()
	if err != nil {
		pterm.Error.Println("Error getting repository root directory")
		return Manifest{}, err
	}
	repo := path.Join(root, "repository")
	return r.getManifest(repo)
}

func (r *Repository) LoadVersions() error {
	var vv []*Version
	root, err := os.Getwd()
	if err != nil {
		pterm.Error.Println("Error creating repository root directory")
		return err
	}
	repo := path.Join(root, "repository")

	err = filepath.WalkDir(repo, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			pterm.Error.Printf("failure accessing a path %q: %v\n", path, err)
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if len(name) > 1 {
				ver, err := strconv.Atoi(name[1:])
				if err != nil {
					pterm.Info.Printf("skipping non-version directory: %+v \n", name)
					return nil
				}
				v := &Version{
					RootNamespace: r.Namespace,
					Index:         ver,
				}
				err = r.loadSchema(v)

				if err != nil {
					pterm.Info.Printf("error loading schemas", name)
					return err
				}
				vv = append(vv, v)
			}
			return nil
		}
		return nil
	})
	r.Versions = vv
	return err
}

func (r *Repository) loadSchema(v *Version) error {
	root, err := os.Getwd()
	if err != nil {
		pterm.Error.Println("Error getting repository root directory")
		return err
	}
	ver := path.Join(root, "repository", v.Name())

	err = filepath.WalkDir(ver, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			pterm.Error.Printf("failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !d.IsDir() {
			name := d.Name()
			schema := &Schema{
				Namespace: v.Namespace(),
				Name:      name,
			}
			v.Schemas = append(v.Schemas, schema)

		}

		return nil
	})
	return err
}

func (r *Repository) NewVersion() (string, error) {
	count := len(r.Versions)
	next := count + 1
	v := &Version{
		RootNamespace: r.Namespace,
		Index:         next,
	}
	r.Versions = append(r.Versions, v)
	vdir := path.Join(r.root, "repository", v.Name())
	err := os.MkdirAll(vdir, 0755)
	if err != nil {
		return "", err
	}
	schema := path.Join(vdir, "schema.cue")
	if next == 1 {
		os.WriteFile(schema, schemaTemplate, 0755)
	} else {
		prev := fmt.Sprintf("%s%d", "v", count)
		prevDir := path.Join(r.root, "repository", prev)
		err := copy.Copy(prevDir, vdir)
		if err != nil {
			return "", err
		}
	}

	repo := path.Join(r.root, "repository")
	err = r.createManifest(repo)
	return v.Name(), err
}

func (r *Repository) SchemaVersions() []string {
	var versions []string
	for _, v := range r.Versions {
		versions = append(versions, v.Name())
	}
	return versions
}
func (r *Repository) SchemaList() []string {
	var schemas []string
	for _, v := range r.Versions {
		for _, s := range v.Schemas {
			schemas = append(schemas, fmt.Sprintf("%s/%s", v.Namespace(), s.Name))
		}
	}
	return schemas
}

type Manifest struct {
	Namespace string                     `json:"namespace"`
	Versions  []string                   `json:"versions"`
	Schemas   map[string]*DisplayVersion `json:"version_schemas"`
}

type DisplayVersion struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Schemas   []*Schema `json:"schemas"`
}

type Version struct {
	RootNamespace string
	Index         int
	Schemas       []*Schema
}

func (v Version) Name() string {
	return fmt.Sprintf("v%d", v.Index)
}

func (v Version) Path() string {
	return fmt.Sprintf("v%d", v.Index)
}
func (v Version) Namespace() string {
	return fmt.Sprintf("%s/v%d", v.RootNamespace, v.Index)
}

type Meta struct {
	Name      string
	Namespace string
}

//go:embed schema.cue
var schemaTemplate []byte
