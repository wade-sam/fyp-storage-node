package main

import (
	//"github.com/wade-sam/fypstoragenode/entity"
	"log"
	"time"

	"github.com/wade-sam/fypstoragenode/Infrastructure/Repositories/rabbit"

	"github.com/wade-sam/fypstoragenode/Infrastructure/Repositories/writetofile"
	"github.com/wade-sam/fypstoragenode/usecase/configuration"
)

func main() {
	wtf := writetofile.NewFileRepo()
	config_service := configuration.NewConfigurationService(wtf)
	conn_name, err := config_service.GetStorageNode()
	if err != nil {
		log.Fatal(err)
	}
	broker_config, err := wtf.GetRabbitDetails()
	if err != nil {
		log.Fatal(err)
	}

	channs := rabbit.Channels{
		Config:  make(chan string),
		Backup:  make(chan rabbit.DTO),
		Restore: make(chan rabbit.DTO),
	}

	BrokerConfig := rabbit.BrokerConfig{
		Schema:         broker_config.Schema,
		Username:       broker_config.Username,
		Password:       broker_config.Password,
		Host:           broker_config.Host,
		Port:           broker_config.Port,
		VHost:          broker_config.VHost,
		ConnectionName: conn_name,
	}

	consumerConf := rabbit.ConsumerConfig{
		ExchangeName: "main",
		ExchangeType: "direct",
		RoutingKey:   conn_name,
		QueueName:    conn_name,
		ConsumerName: conn_name,
		MaxAttempt:   60,
		Interval:     1 * time.Second,
		Channels:     &channs,
	}

	producerConf := rabbit.ProducerConfig{
		ExchangeName: "main",
		ExchangeType: "direct",
		MaxAttempt:   60,
		Interval:     1 * time.Second,
		RoutingKey:   "backupserver",
	}

	broker := rabbit.NewBroker(BrokerConfig, producerConf, consumerConf)
	err = broker.Connect()
	if err != nil {
		log.Fatal(err)
	}
	consumer_chan, err := broker.Start()
	if err != nil {
		log.Fatal("ERR", err)
	}
	broker.Consume(consumer_chan)
}
