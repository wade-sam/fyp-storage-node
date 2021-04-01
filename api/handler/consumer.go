package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
	"github.com/wade-sam/fypstoragenode/Infrastructure/Repositories/rabbit"
	"github.com/wade-sam/fypstoragenode/usecase/backup"
)

type ConsumerConfig struct {
	ExchangeName string
	ExchangeType string
	RoutingKey   string
	QueueName    string
	ConsumerName string
	MaxAttempt   int
	Interval     time.Duration
	Connection   *amqp.Connection
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
type Consumer struct {
	Consumer ConsumerConfig
	Backup   backup.BackupService
}

func NewConsumerRepo(c ConsumerConfig, b backup.BackupService) *Consumer {
	return &Consumer{
		Consumer: c,
		Backup:   b,
	}
}

func (c *Consumer) Start() (*amqp.Channel, error) {

	con, err := c.Connection()
	if err != nil {
		return nil, err
	}
	chn, err := con.Channel()
	if err != nil {
		return nil, err
	}

	if err := chn.ExchangeDeclare(
		c.Consumer.ExchangeName,
		c.Consumer.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, err
	}
	if _, err := chn.QueueDeclare(
		c.Consumer.QueueName,
		true,
		false,
		false,
		false,
		amqp.Table{"x-message-ttl": 6000},
	); err != nil {
		return nil, err
	}

	if err := chn.QueueBind(
		c.Consumer.QueueName,
		c.Consumer.RoutingKey,
		c.Consumer.ExchangeName,
		false,
		nil,
	); err != nil {
		return nil, err
	}
	return chn, nil
}

func (c *Consumer) Consume(channel *amqp.Channel) error {
	msgs, err := channel.Consume(
		c.Consumer.QueueName,
		c.Consumer.ConsumerName,
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
		fmt.Println("new Backup Job")
		d := rabbit.DTO{}
		err = json.Unmarshal(msg.Body, &d)
		fmt.Println(msg.Type)
		//_, err := Deserialize(msg.Body)
		if err != nil {
			log.Println("Can't deserialise message", err)
		}

		switch msg.Type {
		case "New.Backup.Job":
			d.Title = msg.Type

			go c.Backup.NewBackupRun(&d)

			//c.Consumer.Channels.Backup <- d
		case "Full.Backup":
			//dto := DTO{}
			//dto.Title = msg.Type
			//dto.Data = msg.Body
			fmt.Println("DATA", d.Data)

			fmt.Println("placed")
		case "Full.Restore":

		}

		fmt.Println("msg consumed")
	}
	log.Println("Exiting")
	return nil
}

func (c *Consumer) Connection() (*amqp.Connection, error) {
	if c.Consumer.Connection == nil || c.Consumer.Connection.IsClosed() {
		return nil, errors.New("connection isnt open")
	}
	return c.Consumer.Connection, nil
}
