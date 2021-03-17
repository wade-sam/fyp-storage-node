package writetofile

import (
	"compress/gzip"
	"container/list"
	"fmt"
	"io/ioutil"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/wade-sam/fypstoragenode/entity"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type FileRepo struct{}

type FileStruct struct {
	BackupServer  string `json:"backupserver"`
	StorageNode   string `json:"storagenode"`
	RabbitDetails *RabbitConfig
}

type RabbitConfig struct {
	Schema         string `json:"schema"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Host           string `json:"host"`
	Port           string `json:"port"`
	VHost          string `json:"vhost"`
	ConnectionName string `json:"conname"`
}

func NewFileRepo() *FileRepo {
	return &FileRepo{}
}

func ReadInJsonFile() (*FileStruct, error) {
	var file FileStruct
	jsonFile, err := os.Open("/home/sam/Documents/fyp-storage_node/Infrastructure/Repositories/writetofile/config.json")
	if err != nil {
		return nil, entity.ErrFileNotFound
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, entity.ErrCouldNotUnMarshallJSON
	}
	json.Unmarshal(byteValue, &file)
	return &file, nil
}

func WriteJsonFile(file *FileStruct) error {
	outputFile, err := json.MarshalIndent(file, "", "	")
	if err != nil {
		return entity.ErrCouldNotMarshallJSON
	}
	err = ioutil.WriteFile("/home/sam/Documents/fyp-storage_node/Infrastructure/Repositories/writetofile/config.json", outputFile, 0775)
	if err != nil {
		return entity.ErrCouldNotWriteToFile
	}
	return nil
}

func (f *FileRepo) GetStorageNode() (string, error) {
	file, err := ReadInJsonFile()
	if err != nil {
		return "", err
	}
	storagenode := file.StorageNode
	fmt.Println(storagenode)
	return storagenode, nil
}

func (f *FileRepo) GetRabbitDetails() (*RabbitConfig, error) {
	file, err := ReadInJsonFile()
	if err != nil {
		return nil, err
	}
	rabbit := file.RabbitDetails
	return rabbit, nil
}

func (f *FileRepo) SetStorageNode(ip string) error {
	file, err := ReadInJsonFile()
	if err != nil {
		return err
	}
	file.StorageNode = ip
	err = WriteJsonFile(file)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileRepo) CreateDirectoryLayout(client string, directories *entity.Directory) error {
	visited := make(map[string]*entity.Directory)
	queue := list.New()
	queue.PushBack(n)
	visited[n.Name] = n
	return nil
}

func (f *FileRepo) CreateBackupResult(files map[string]*entity.FileDTO) error {
	output, err := json.MarshalIndent(files, "", "   ")
	if err != nil {
		fmt.Println(entity.ErrCouldNotMarshallJSON)
	}

	filename := "backup_config.gzip"
	file, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		fmt.Println(entity.ErrCouldNotMarshallJSON)
	}
	di, err := os.Create(filename)
	q := gzip.NewWriter(di)
	_, err = q.Write([]byte(file))
	if err != nil {
		return entity.ErrCouldNotWriteToFile
	}
	q.Close()
	return nil
}

func (f *FileRepo) GetPreviousBackupResult() (map[string]*entity.FileDTO, error) {
	var files map[string]*entity.FileDTO
	fi, err := os.Open("backup_config.gzip")
	if err != nil {
		return nil, err
	}
	reader, err := gzip.NewReader(fi)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(reader).Decode(&files)
	if err != nil {
		return nil, err
	}
	return files, nil
}
