package entity

import (
	"fmt"
	"path"
	"time"
)

type ClientFile struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Checksum string `json:"checksum"`
}

type BackupRun struct {
	PolicyName string
	ID         string
	Type       string
	Date       string
	Clients    map[string]*ClientRun
	Status     string
}

type ClientRun struct {
	Name         string `json:"client"`
	Status       string `json:"status"`
	Files        *Directory
	RemovePrefix int                    `json:"remove_prefix"`
	Map          map[string]interface{} `json:"backed_up_files"`
}

func NewBackupRun(policyname, Type string) (*BackupRun, error) {
	run := &BackupRun{
		PolicyName: policyname,
		Type:       Type,
		Clients:    make(map[string]*ClientRun),
	}
	return run, nil

}

func (br *BackupRun) CreateNameTimeProperties() {
	date := time.Now()
	name := fmt.Sprintf("%v-%v", br.PolicyName, date.Format("01-02-2006 15:04:05"))
	br.ID = name
	br.Date = date.Format("01-02-2006 15:04:05")
}

func (br *BackupRun) AddClient(client *ClientRun) error {
	//_, err := br.GetClient(client.Name)
	//if err != nil {
	//	return ErrClientAlreadyAdded
	//}
	br.Clients[client.Name] = client
	return nil
}

func (br *BackupRun) GetClient(name string) (*ClientRun, error) {
	for i, j := range br.Clients {
		if br.Clients[i].Name == name {
			return j, nil
		}
	}
	return nil, ErrNotFound
}

func NewClientRun(name string) (*ClientRun, error) {
	return &ClientRun{
		Name:   name,
		Status: "In Progress",
		Files:  &Directory{},
		Map:    make(map[string]interface{}),
	}, nil
}

func (f *ClientRun) AddMap(clientmap map[string]interface{}) error {
	f.Map = clientmap
	return nil
}

func (cr *BackupRun) AddRemovePrefix(client string, prefix int) (*ClientRun, error) {
	foundclient, err := cr.GetClient(client)
	if err != nil {
		return nil, nil
	}
	foundclient.RemovePrefix = prefix
	return foundclient, nil
}

func (f *ClientRun) CompileReport() error {
	var report *Directory
	//fmt.Println("FOUND MAP", f.Map)
	fmt.Println("head: ", f.Files)
	for i, j := range f.Map {
		var parent *Directory
		fmt.Println()
		if path.Dir(i) == f.Files.Path {
			fmt.Println("KEY: ", i, f.Files.Name, path.Base(i))
			report = j.(*Directory)
			continue
		} else {
			fmt.Println("MAP", i, "parent", path.Dir(i))
			//fmt.Println("FOUND MAP", f.Map[i])
			parent = f.Map[path.Dir(i)].(*Directory)
		}
		switch v := j.(type) {
		case *Directory:
			parent.Folders[v.Name] = v
		case *File:
			parent.Files = append(parent.Files, v)
		}

	}
	f.Files = report
	return nil
}
func NewClientFile(id string) (*ClientFile, error) {
	if id == "" {
		return nil, ErrFieldWasEmpty
	} else {
		return &ClientFile{
			ID: id,
		}, nil
	}
}
func (f *ClientFile) AddChecksum(checksum string) error {
	if checksum == "" {
		return ErrFieldWasEmpty
	}
	f.Checksum = checksum
	return nil
}

func (f *ClientFile) UpdateStatus(status string) error {
	if status != "client_success" || status != "client_failure" || status == "" {
		return ErrWrongStatus
	}
	f.Status = status
	return nil
}
