package backup

import (
	"container/list"
	"fmt"
	"log"

	//"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/mapstructure"
	"github.com/wade-sam/fypstoragenode/Infrastructure/Repositories/rabbit"
	"github.com/wade-sam/fypstoragenode/Infrastructure/Repositories/socket"
	"github.com/wade-sam/fypstoragenode/entity"
)

type StoragenodeData struct {
	Clients    []string `json:"clients"`
	PolicyName string   `json:"policyname"`
}

type BackupService struct {
	Channel    chan rabbit.DTO
	BackupJob  *entity.BackupRun
	RabbitRepo RabbitRepository
	SockRepo   SocketRepository
	FileRepo   FileRepository
	Storage    StorageDetails
}
type StorageDetails struct {
	BackupPath  string
	ClientPaths map[string]string
}

func NewBackupService(r RabbitRepository, s SocketRepository, f FileRepository) *BackupService {
	return &BackupService{
		RabbitRepo: r,
		SockRepo:   s,
		FileRepo:   f,
	}
}

func (service *BackupService) NewBackupRun(dto *rabbit.DTO) {
	sndto := StoragenodeData{}
	backupjob := entity.BackupRun{}
	mapstructure.Decode(dto.Data, &sndto)
	for i := range sndto.Clients {
		client, err := entity.NewClientRun(sndto.Clients[i])
		if err != nil {
			log.Println(err)
		}
		backupjob.Clients = append(backupjob.Clients, *client)
	}
	backupjob.PolicyName = sndto.PolicyName
	backupjob.CreateNameTimeProperties()
	service.BackupJob = &backupjob

	clientpaths, jobpath, err := service.FileRepo.CreateJobLayout(sndto.Clients, backupjob.ID)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Backup path: ", jobpath, "Clients: ", sndto.Clients)
	storage := StorageDetails{
		BackupPath:  jobpath,
		ClientPaths: clientpaths,
	}
	service.Storage = storage
	err = service.RabbitRepo.SendBackupSetup(jobpath)
	if err != nil {
		log.Println(err)
	}
	chn := make(chan (*socket.SockItem))
	err = service.SockRepo.Start(chn)
	service.ReceiveData(chn)
}

func (service *BackupService) ReceiveData(chn chan *socket.SockItem) {
	for msg := range chn {
		switch msg.Type {
		case "directoryscan":

			var directory = entity.Directory{}
			mapstructure.Decode(msg.Item, &directory)
			//fmt.Println("Directory Received:", directory)
			service.CreateBackupDirectoryLayout(msg.Client, &directory)
		case "filedata":
			var file = entity.File{}
			mapstructure.Decode(msg.Item, &file)

		}

	}
}

func (service *BackupService) CreateBackupDirectoryLayout(client string, directories *entity.Directory) error {
	visited := make(map[string]*entity.Directory)
	queue := list.New()
	queue.PushBack(directories)
	visited[directories.Name] = directories
	fmt.Println(client, "clients", service.Storage.ClientPaths)
	for queue.Len() > 0 {
		pop := queue.Front()
		for id, node := range pop.Value.(*entity.Directory).Folders {
			//	fmt.Sprintf("Client: %v, StoredClient: %v", client, service.Storage.ClientPaths)
			if _, ok := visited[id]; !ok {

				if val, ok := service.Storage.ClientPaths[client]; ok {
					fmt.Println("CLIENT", client)
					fmt.Println(client, "clients", service.Storage.ClientPaths)
					paths := fmt.Sprintf("%v/%v", val, node.Path)
					err := service.FileRepo.CreateDirectoryLayout(paths, node.Properties)
					if err != nil {
						fmt.Println("ERR", err)
						return err
					}

				}
				visited[id] = node
				queue.PushBack(node)
			}

		}
		queue.Remove(pop)
	}
	return nil
}
