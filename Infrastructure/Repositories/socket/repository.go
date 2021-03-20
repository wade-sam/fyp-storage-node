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
	Connection *net.Conn
	Decoder    *gob.Decoder
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
	ID     string       `json:"id"`
	Type   string       `json="type"`
	Client string       `json="client"`
	Item   *entity.File `json="item"`
}

func NewRepository(address, port, conn_type string) *Repository {
	return &Repository{
		Address:   address,
		Port:      port,
		Conn_type: conn_type,
	}
}

func (sock *Repository) Start(chn chan *SockItem) error {
	port := fmt.Sprintf(":%v", sock.Port)
	dstream, err := net.Listen(sock.Conn_type, port)
	if err != nil {
		fmt.Println("1", err)
		return err
	}
	con, err := dstream.Accept()
	if err != nil {
		fmt.Println("2", err)
		return err
	}
	sock.Connection = &con
	sock.Decoder = gob.NewDecoder(*sock.Connection)

	fmt.Println("CONNECTION!!!!", sock.Connection)
	go sock.ReceiveBackupData(chn)
	return nil
}

func (sock *Repository) ReceiveBackupData(chn chan *SockItem) {
	condition := false

	dto := SockItem{}
	for condition == false {
		gob.Register(&entity.Directory{})
		err := sock.Decoder.Decode(&dto)
		//i, err := bufio.NewReader(*socket.connection).ReadByte()
		if err != nil {
			log.Println(err)
		}
		if dto.Type == "clientcomplete" {
			close(chn)
		} else {
			chn <- &dto
		}

	}

	// err = json.Unmarshal(i, &dto)
	// if err != nil {
	// 	log.Println("Can't deserialise message", err)
	// }
}
