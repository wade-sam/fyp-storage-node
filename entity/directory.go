package entity

import "encoding/json"

type Directory struct {
	Name       string                `json:"name,omitempty"`
	Path       string                `json:"path"`
	Properties []string              `json:"properties,omitempty"`
	Files      []*File               `json:"files,omitempty"`
	Folders    map[string]*Directory `json:"folders,omitempty"`
}

func NewDirectory(name string) *Directory {
	//	children := make(map[string]*Directory)
	return &Directory{
		Name:       name,
		Properties: []string{},
		Files:      []*File{},
		Folders:    map[string]*Directory{},
	}
}
func (d *Directory) String() string {
	j, _ := json.Marshal(d)
	return string(j)
}

func (d *Directory) AddChildDirectory(dir *Directory) error {
	for _, i := range d.Folders {
		if dir.Name == i.Name {
			return ErrChildAlreadyExists
		}
	}
	d.Folders[dir.Name] = dir
	return nil
}
func (d *Directory) AddProperties(permissions, UID, GID string) {
	d.Properties = append(d.Properties, permissions, UID, GID)
}

func (d *Directory) AddFile(file *File) {
	d.Files = append(d.Files, file)
}

func (d *Directory) GetFile(file string) (*File, error) {
	for _, v := range d.Files {
		if v.Name == file {
			return v, nil
		}
	}
	return nil, ErrFileNotFound
}

func (d *Directory) Validate() error {
	if d.Name == "" || len(d.Properties) != 3 {
		return ErrUnsuccesfulValidationOfDirectory
	}
	return nil
}
