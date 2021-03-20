package entity

import (
	"fmt"
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
	Clients    []ClientRun
	Status     string
}

type ClientRun struct {
	Name   string
	Status string
	Files  *Directory
}

func NewBackupRun(policyname, Type string) (*BackupRun, error) {
	run := &BackupRun{
		PolicyName: policyname,
		Type:       Type,
	}
	return run, nil

}

func (br *BackupRun) CreateNameTimeProperties() {
	date := time.Now()
	name := fmt.Sprintf("%v-%v", br.PolicyName, date.Format("01-02-2006 15:04:05"))
	br.ID = name
	br.Date = date.Format("01-02-2006 15:04:05")
}

func (br *BackupRun) AddClient(client ClientRun) error {
	_, err := br.GetClient(client.Name)
	if err != nil {
		return ErrClientAlreadyAdded
	}
	br.Clients = append(br.Clients, client)
	return nil
}

func (br *BackupRun) GetClient(name string) (string, error) {
	for i := range br.Clients {
		if br.Clients[i].Name == name {
			return name, nil
		}
	}
	return name, ErrNotFound
}

func NewClientRun(name string) (*ClientRun, error) {
	return &ClientRun{
		Name:   name,
		Status: "In Progress",
		Files:  &Directory{},
	}, nil
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
