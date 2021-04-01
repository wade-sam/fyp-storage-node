package backup

import (
	"github.com/wade-sam/fypstoragenode/Infrastructure/Repositories/rabbit"
	"github.com/wade-sam/fypstoragenode/Infrastructure/Repositories/socket"
	"github.com/wade-sam/fypstoragenode/entity"
)

type Usecase interface {
	//CreateBackup returns with a unique ID that resembles the backuprun
	NewBackupRun(*StoragenodeData) (string, error)
	CreateBackupDirectoryLayout(client string, directories *entity.Directory) error
	StartBackupCopy(chn chan (rabbit.DTO))
}

type RabbitRepository interface {
	SendBackupSetup(id string) error
	SendBackupFileMessage(client_file *entity.ClientFile) error
}

type SocketRepository interface {
	Start(chn chan (*socket.SockItem)) error
	End() error
}

type FileRepository interface {
	CreateJobLayout(clients []string, policyname string) (map[string]string, string, error)
	CreateDirectoryLayout(path string, permissions []string) error
	CreateFile(client, path string, file *socket.SockFile) error
	CreateBackupReport(client, path string, files map[string]interface{}) error
	//ClientFile(path string, perms []string) error
}
