package backup

import (
	"container/list"
	"fmt"
	"log"
	"path"

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

	mapstructure.Decode(dto.Data, &sndto)
	backupjob, err := entity.NewBackupRun(sndto.PolicyName, "Full")
	if err != nil {
		log.Println(err)
	}
	for i := range sndto.Clients {
		client, err := entity.NewClientRun(sndto.Clients[i])
		if err != nil {
			log.Println(err)
		}

		backupjob.AddClient(client)
	}
	backupjob.PolicyName = sndto.PolicyName
	backupjob.CreateNameTimeProperties()
	service.BackupJob = backupjob

	clientpaths, jobpath, err := service.FileRepo.CreateJobLayout(sndto.Clients, backupjob.ID)
	if err != nil {
		log.Println("Error Creating job layout", err)
	}
	storage := StorageDetails{
		BackupPath:  jobpath,
		ClientPaths: clientpaths,
	}
	service.Storage = storage

	chn := make(chan (*socket.SockItem))
	go service.SockRepo.Start(chn, len(sndto.Clients))
	err = service.RabbitRepo.SendBackupSetup(jobpath)
	if err != nil {
		service.SockRepo.End()
		close(chn)
		log.Println("Error sending backup setup to server", err)
		return
	}
	for c := 1; c <= len(sndto.Clients); c++ {
		log.Println(len(sndto.Clients))
		err = service.ReceiveData(chn)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	//service.SockRepo.End()
	//close(chn)
	// err = service.ReceiveData(chn)
	// if err != nil {
	// 	service.SockRepo.End()
	// 	close(chn)
	service.Storage = StorageDetails{}
	return
	//}

	fmt.Println("Completed backup")
	
}

func (service *BackupService) ReceiveData(chn chan *socket.SockItem) error {
	for msg := range chn {
		switch msg.Type {
		case "directoryscan":

			var directory = entity.Directory{}
			mapstructure.Decode(msg.Item, &directory)
			fmt.Println("Directory Received:")
			err := service.CreateBackupDirectoryLayout(msg.Client, &directory)
			if err != nil {
				log.Println("Error Couldn't create Backup Setup", err)
				return err
			}
			if val, ok := service.BackupJob.Clients[msg.Client]; ok {
				head := entity.NewDirectory(directory.Name)
				head.Path = directory.Path
				head.Properties = directory.Properties
				val.Files = head
			}
		case "filedata":
			var sockfile = socket.SockFile{}
			mapstructure.Decode(msg.Item, &sockfile)
			response, err := service.WriteFileSetup(msg.Client, &sockfile)
			if err != nil {
				log.Println("Error Couldn't write file: ", err)
				service.RabbitRepo.SendBackupFileMessage(response)
				//return err
			}
			if response.Status == "Success" {
				service.RabbitRepo.SendBackupFileMessage(response)
			} else {
				service.RabbitRepo.SendBackupFileMessage(response)
			}

		case "clientcomplete":
			var sockfile = socket.SockFile{}
			rabbitresponse := entity.ClientFile{}
			mapstructure.Decode(msg.Item, &sockfile)
			filerun := service.BackupJob.Clients[msg.Client]
			if _, ok := service.Storage.ClientPaths[msg.Client]; ok {
				//fmt.Println("Client location", service.Storage.BackupPath)
				err := service.FileRepo.CreateBackupReport(msg.Client, service.Storage.BackupPath, filerun.Map)
				if err != nil {
					log.Println("Error: Couldn't write backup report ", err)
					rabbitresponse.ID = "Completion"
					rabbitresponse.Status = "Failed"
					service.RabbitRepo.SendBackupFileMessage(&rabbitresponse)
					return err
				}
				rabbitresponse.ID = "Completion"
				rabbitresponse.Status = "Success"
				service.RabbitRepo.SendBackupFileMessage(&rabbitresponse)
			}

			fmt.Println("FINISHED")
			return nil

		}

	}
	return nil
}

func (service *BackupService) WriteFileSetup(client string, file *socket.SockFile) (*entity.ClientFile, error) {
	var filePath string
	msg := entity.ClientFile{
		ID: file.Metadata.Path,
	}
	runePath := []rune(file.Metadata.Path)
	_, err := service.BackupJob.GetClient(client)
	if err != nil {
		//log.Println("Error Couldn't find client",err)
		return nil, err
	}
	if val, ok := service.Storage.ClientPaths[client]; ok {
		filePath = fmt.Sprintf("%v/%v", val, string(runePath[service.BackupJob.Clients[client].RemovePrefix:]))
	}

	err = service.FileRepo.CreateFile(client, path.Clean(filePath), file)
	if err != nil {
		msg.Status = "Failed"
		log.Println("EROR REACHED", err)
		return &msg, err
	}
	msg.Status = "Success"
	msg.Checksum = file.Metadata.Checksum
	service.BackupJob.Clients[client].Map[file.Metadata.Path] = file.Metadata
	return &msg, nil
}

func (service *BackupService) CreateBackupDirectoryLayout(client string, directories *entity.Directory) error {
	lenPath := len(directories.Path)
	_, err := service.BackupJob.AddRemovePrefix(client, lenPath)
	if err != nil {
		return err
	}
	visited := make(map[string]*entity.Directory)
	queue := list.New()
	queue.PushBack(directories)
	visited["/"] = directories

	for queue.Len() > 0 {
		pop := queue.Front()
		for _, node := range pop.Value.(*entity.Directory).Folders {
			if _, ok := visited[node.Path]; !ok {

				if val, ok := service.Storage.ClientPaths[client]; ok {
					runePath := []rune(node.Path)
					paths := fmt.Sprintf("%v/%v", val, string(runePath[lenPath:]))
					err := service.FileRepo.CreateDirectoryLayout(path.Clean(paths), node.Properties)
					if err != nil {
						//fmt.Println("ERR", err)
						return err
					}

					n := entity.NewDirectory(directories.Name)
					n.Path = node.Path
					n.Properties = node.Properties
					service.BackupJob.Clients[client].Map[n.Path] = n
					//fmt.Println(n.Path)
				}
				visited[node.Path] = node
				queue.PushBack(node)
			}

		}
		queue.Remove(pop)
	}
	return nil
}
