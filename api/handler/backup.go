package handler

import (
	"github.com/wade-sam/fypstoragenode/Infrastructure/Repositories/rabbit"
	"github.com/wade-sam/fypstoragenode/usecase/backup"
	"github.com/wade-sam/fypstoragenode/usecase/configuration"
)

func BackupHandler(service backup.Usecase, configservice configuration.Usecase, b *rabbit.Broker, chn chan rabbit.DTO) {
	for msg := range chn {

	}
}
