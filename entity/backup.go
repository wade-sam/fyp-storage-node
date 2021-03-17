package entity

type FileDTO struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Checksum string `json:"checksum"`
}
type FileTransfer struct {
	SNFile *entity.File
	Data   *[]byte
}

func NewFileDTO(id string) (*FileDTO, error) {
	if id == "" {
		return nil, ErrFieldWasEmpty
	} else {
		return &FileDTO{
			ID: id,
		}, nil
	}
}
func (f *FileDTO) AddChecksum(checksum string) error {
	if checksum == "" {
		return ErrFieldWasEmpty
	}
	f.Checksum = checksum
	return nil
}

func (f *FileDTO) UpdateStatus(status string) error {
	if status != "client_success" || status != "client_failure" || status == "" {
		return ErrWrongStatus
	}
	f.Status = status
	return nil
}
