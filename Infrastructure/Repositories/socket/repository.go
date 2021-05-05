package socket

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"

	"github.com/wade-sam/fypstoragenode/entity"
)

type Repository struct {
	Address    string
	Port       string
	Conn_type  string
	Connection net.Conn
	Decoder    *gob.Decoder
	Status     bool
}
type FileTransfer struct {
	SNFile *entity.File `json:"metadata"`
	Data   *[]byte      `json:"data"`
}

type DirectoryTransfer struct {
	SNDirectory *entity.Directory
}

type SockItem struct {
	Type   string      `json:"type"`
	ID     string      `json:"id"`
	Client string      `json:"client"`
	Item   interface{} `json:"item"`
}

type SockDirectory struct {
	ID     string            `json:"id"`
	Type   string            `json="type"`
	Client string            `json="client"`
	Item   *entity.Directory `json="item"`
}

type SockFile struct {
	Metadata *entity.File
	Data     []byte
}

func NewRepository(address, port, conn_type string) *Repository {
	return &Repository{
		Address:   address,
		Port:      port,
		Conn_type: conn_type,
		Status:    true,
	}
}

func (sock *Repository) Start(chn chan *SockItem) error {
	sock.Status = true
	port := fmt.Sprintf(":%v", sock.Port)
	dstream, err := net.Listen(sock.Conn_type, port)
	if err != nil {
		fmt.Println("Error Could not listen on port", err)
		return err
	}
	log.Println("Opened TCP Stream")
	con, err := dstream.Accept()
	if err != nil {
		log.Println("Error Client Not Accepted", err)
		return err
	}
	sock.Connection = con
	sock.Decoder = gob.NewDecoder(sock.Connection)

	fmt.Println("Client Connected")
	sock.ReceiveBackupData(chn)
	err = sock.Connection.Close()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Closed Succesfully", err)
	close(chn)
	dstream.Close()
	return nil
}

func (sock *Repository) End() error {
	sock.Status = false
	return nil
}

func (sock *Repository) ReceiveBackupData(chn chan *SockItem) {

	dto := SockItem{}
	for sock.Status == true {
		gob.Register(&entity.Directory{})
		gob.Register(&SockFile{})
		err := sock.Decoder.Decode(&dto)
		//i, err := bufio.NewReader(*socket.connection).ReadByte()
		if err != nil {
			log.Println(err)
		}
		if dto.Type == "clientcomplete" {
			chn <- &dto
			fmt.Println("close socket message", dto.ID)
			//return
			//condition = true
		} else {
			//fmt.Println("new socket message", dto.ID)
			chn <- &dto
			//sock.Connection.Close()
		}

	}

	// err = json.Unmarshal(i, &dto)
	// if err != nil {
	// 	log.Println("Can't deserialise message", err)
	// }
}
