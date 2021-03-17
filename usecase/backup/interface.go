package backup

import (
	"github.com/wade-sam/fypstoragenode/Infrastructure/Repositories/rabbit"
	"github.com/wade-sam/fypstoragenode/entity"
)

type Usecase interface {
	//CreateBackup returns with a unique ID that resembles the backuprun
	CreateBackup(policyname string, clients []string) (string, error)
	CreateDirectoryLayout(client string, directories *entity.Directory) error
	StartBackupCopy(chn chan (rabbit.DTO))
}

type RabbitRepository interface {
	SendBackupSetup(id string) error
	SendFileMessage(name, checksum string) error
}

type SocketRepository interface {
	Start(chn chan (rabbit.DTO)) error
	End() error
}
