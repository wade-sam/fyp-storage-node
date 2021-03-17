package entity

type File struct {
	Name       string   `json:"filename,omitempty"`
	Path       string   `json:"path,omitempty"`
	Checksum   string   `json:"checksum,omitempty"`
	Properties []string `json:"properties,omitempty"`
	Status     string
}

func NewFile(name string) *File {
	file := &File{
		Name: name,
	}

	return file
}

func (f *File) AddChecksum(path string) error {
	f.Checksum = path
	return nil
}

func (f *File) AddProperties(perms, uid, gid string) {
	props := []string{perms, uid, gid}
	f.Properties = props
}
