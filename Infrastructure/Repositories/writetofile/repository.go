package writetofile

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"github.com/wade-sam/fypstoragenode/Infrastructure/Repositories/socket"
	"github.com/wade-sam/fypstoragenode/entity"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type FileRepo struct {
	BackupLocation string
}

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

func NewFileRepo(location string) *FileRepo {
	return &FileRepo{
		BackupLocation: location,
	}
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

func (f *FileRepo) CreateDirectoryLayout(paths string, permissions []string) error {
	//fmt.Println("CLEANED", path.Clean(paths), "NORMAL", paths)
	perms, err := strconv.ParseInt(permissions[0], 0, 32)
	if err != nil {

		return err
	}
	//fmt.Println(permissions)
	os.Mkdir(paths, os.FileMode(perms))
	os.Chmod(paths, os.FileMode(perms))
	uid, err := strconv.Atoi(permissions[1])
	if err != nil {

		return err
	}
	gid, err := strconv.Atoi(permissions[2])
	err = os.Chown(paths, uid, gid)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileRepo) CreateFile(client, path string, file *socket.SockFile) error {
	//fmt.Println("Create File", path)
	filePlacement, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = filePlacement.Write(file.Data)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileRepo) CreateJobLayout(clients []string, policyname string) (map[string]string, string, error) {
	clientsLocation := make(map[string]string)
	policyLocation := f.BackupLocation + "/" + policyname
	reportLocation := policyLocation + "/reports"
	//fmt.Println("policy location", policyLocation)
	err := os.Mkdir(policyLocation, 666)
	if err != nil {
		return nil, "", err
	}
	err = os.Mkdir(reportLocation, 666)
	if err != nil {
		return nil, "", err
	}

	for _, j := range clients {
		clientLocation := policyLocation + "/" + j
		//fmt.Println("client location", clientLocation)
		err = os.Mkdir(clientLocation, 666)
		if err != nil {
			return nil, "", err
		}
		clientsLocation[j] = clientLocation
	}
	return clientsLocation, policyLocation, nil
}

func (f *FileRepo) CreateBackupReport(client, location string, files map[string]interface{}) error {
	output, err := json.MarshalIndent(files, "", "   ")
	if err != nil {
		fmt.Println(entity.ErrCouldNotMarshallJSON)
	}

	filename := fmt.Sprintf("%v_report.gzip", client)
	path := fmt.Sprintf("%v/reports/%v", location, filename)
	//fmt.Println("PATH", path)
	file, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		fmt.Println(entity.ErrCouldNotMarshallJSON)
	}
	di, err := os.Create(path)
	q := gzip.NewWriter(di)
	_, err = q.Write([]byte(file))
	if err != nil {
		return entity.ErrCouldNotWriteToFile
	}
	q.Close()
	return nil
}

func (f *FileRepo) GetPreviousBackupResult() (map[string]*entity.ClientFile, error) {
	var files map[string]*entity.ClientFile
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
