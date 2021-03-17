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
	connection *net.Conn
}
type FileTransfer struct {
	SNFile *entity.File `json="metadata"`
	Data   *[]byte      `json="data"`
}

func NewRepository(address, port, conn_type string) *Repository {
	return &Repository{
		Address:   address,
		Port:      port,
		Conn_type: conn_type,
	}
}

func (sock *Repository) StartBackup() error {
	dstream, err := net.Listen("tcp", "8080")
	if err != nil {
		return err
	}
	defer dstream.Close()
	var connected = false

	for connected != true {
		con, err := dstream.Accept()
		if err != nil {
			return err
		}
		connected = true
		sock.connection = &con
	}
}

func (socket *Repository) ReceiveBackupData() {
	dec := gob.NewDecoder(*socket.connection)
	dto := FileTransfer{}
	err := dec.Decode(&dto)
	//i, err := bufio.NewReader(*socket.connection).ReadByte()
	fmt.Println(dto.SNFile)
	if err != nil {
		log.Println(err)
	}

	// err = json.Unmarshal(i, &dto)
	// if err != nil {
	// 	log.Println("Can't deserialise message", err)
	// }
}
