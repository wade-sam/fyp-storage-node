package rabbit

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type ConsumerConfig struct {
	ExchangeName string
	ExchangeType string
	RoutingKey   string
	QueueName    string
	ConsumerName string
	MaxAttempt   int
	Interval     time.Duration
	connection   *amqp.Connection
	Channels     *Channels
}

type Channels struct {
	Config  chan string
	Backup  chan DTO
	Restore chan DTO
}

type DTO struct {
	Title string
	Data  interface{} `json:"data"`
}

func (b *Broker) Start() (*amqp.Channel, error) {

	con, err := b.Connection()
	if err != nil {
		return nil, err
	}
	chn, err := con.Channel()
	if err != nil {
		return nil, err
	}

	if err := chn.ExchangeDeclare(
		b.Consumer.ExchangeName,
		b.Consumer.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, err
	}
	if _, err := chn.QueueDeclare(
		b.Consumer.QueueName,
		true,
		false,
		false,
		false,
		amqp.Table{"x-message-ttl": 6000},
	); err != nil {
		return nil, err
	}

	if err := chn.QueueBind(
		b.Consumer.QueueName,
		b.Consumer.RoutingKey,
		b.Consumer.ExchangeName,
		false,
		nil,
	); err != nil {
		return nil, err
	}
	return chn, nil
}

func (b *Broker) Consume(channel *amqp.Channel) error {
	msgs, err := channel.Consume(
		b.Consumer.QueueName,
		b.Consumer.ConsumerName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	for msg := range msgs {
		d := DTO{}
		err = json.Unmarshal(msg.Body, &d)
		//_, err := Deserialize(msg.Body)
		if err != nil {
			log.Println("Can't deserialise message", err)
		}

		switch msg.Type {
		case "New.Backup.Job":
			b.Consumer.Channels.Backup <- d
		case "Full.Backup":
			//dto := DTO{}
			//dto.Title = msg.Type
			//dto.Data = msg.Body
			fmt.Println("DATA", d.Data)
			d.Title = msg.Type
			b.Consumer.Channels.Backup <- d
			fmt.Println("placed")
		case "Full.Restore":

		}

		fmt.Println("msg consumed")
	}
	log.Println("Exiting")
	return nil
}
