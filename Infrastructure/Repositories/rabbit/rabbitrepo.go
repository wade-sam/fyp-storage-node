package rabbit

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

type BrokerConfig struct {
	Schema         string
	Username       string
	Password       string
	Host           string
	Port           string
	VHost          string
	ConnectionName string
}

type Broker struct {
	config           BrokerConfig
	BrokerConnection *amqp.Connection
}
type DTO struct {
	Title string
	Data  interface{} `json:"data"`
}

func NewBroker(config BrokerConfig) *Broker {
	return &Broker{
		config: config,
	}
}

func (r *Broker) Connect() error {
	if r.BrokerConnection == nil || r.BrokerConnection.IsClosed() {
		conn, err := amqp.Dial(fmt.Sprintf("%s://%s:%s@%s:%s/%s",
			r.config.Schema,
			r.config.Username,
			r.config.Password,
			r.config.Host,
			r.config.Port,
			r.config.VHost,
		))
		if err != nil {
			return err
		}
		r.BrokerConnection = conn
	}
	return nil
}
func (r *Broker) Channel() (*amqp.Channel, error) {
	chn, err := r.BrokerConnection.Channel()
	if err != nil {
		return nil, err
	}
	return chn, nil
}

func (r *Broker) Connection() (*amqp.Connection, error) {
	if r.BrokerConnection == nil || r.BrokerConnection.IsClosed() {
		return nil, errors.New("connection isnt open")
	}
	return r.BrokerConnection, nil
}
