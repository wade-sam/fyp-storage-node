package rabbit

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
	"github.com/wade-sam/fypstoragenode/entity"
)

type ProducerConfig struct {
	ExchangeName string
	ExchangeType string
	RoutingKey   string
	MaxAttempt   int
	Interval     time.Duration
	Connection   *amqp.Connection
}

type Producer struct {
	Producer ProducerConfig
}

func NewProducer(p ProducerConfig) *Producer {
	return &Producer{
		Producer: p,
	}
}
func (p *Producer) Publish(Type string, body *DTO) error {
	channel, err := p.Producer.Connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	data, err := json.Marshal(body.Data)
	//dto, err := Serialize(body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(Type, "", p.Producer.RoutingKey)
	err = channel.Publish(
		p.Producer.ExchangeName,
		p.Producer.RoutingKey,
		false,
		false,
		amqp.Publishing{
			Type:        Type,
			ContentType: "encoding/json",
			Body:        []byte(data),
		},
	)
	if err != nil {
		return err
	}
	r := &entity.Directory{}
	err = json.Unmarshal([]byte(data), &r)
	fmt.Println("Sent message back")
	return nil
}

func (p *Producer) SendBackupSetup(id string) error {
	dto := DTO{}
	dto.Data = id
	err := p.Publish("StorageNode.Job", &dto)
	if err != nil {
		return err
	}
	return nil
}
