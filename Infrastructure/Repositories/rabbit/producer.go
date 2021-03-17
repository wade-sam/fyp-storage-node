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
	connection   *amqp.Connection
}

func (b *Broker) Publish(Type string, body *DTO) error {
	channel, err := b.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	data, err := json.Marshal(body.Data)
	//dto, err := Serialize(body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(Type, "", b.Producer.RoutingKey)
	err = channel.Publish(
		b.Producer.ExchangeName,
		b.Producer.RoutingKey,
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
