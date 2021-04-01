package main

import (
	//"github.com/wade-sam/fypstoragenode/entity"

	"log"
	"time"

	"github.com/wade-sam/fypstoragenode/Infrastructure/Repositories/rabbit"
	"github.com/wade-sam/fypstoragenode/Infrastructure/Repositories/writetofile"
	"github.com/wade-sam/fypstoragenode/api/handler"

	//"github.com/wade-sam/fypstoragenode/api/handler"
	"github.com/wade-sam/fypstoragenode/Infrastructure/Repositories/socket"
	"github.com/wade-sam/fypstoragenode/usecase/backup"
	"github.com/wade-sam/fypstoragenode/usecase/configuration"
)

func main() {
	wtf := writetofile.NewFileRepo("/home/sam/backup")
	config_service := configuration.NewConfigurationService(wtf)
	conn_name, err := config_service.GetStorageNode()

	if err != nil {
		log.Fatal(err)
	}
	broker_config, err := wtf.GetRabbitDetails()
	if err != nil {
		log.Fatal(err)
	}

	// channs := Channels{
	// 	Config:  make(chan string),
	// 	Backup:  make(chan rabbit.DTO),
	// 	Restore: make(chan rabbit.DTO),
	// }

	BrokerConfig := rabbit.BrokerConfig{
		Schema:         broker_config.Schema,
		Username:       broker_config.Username,
		Password:       broker_config.Password,
		Host:           broker_config.Host,
		Port:           broker_config.Port,
		VHost:          broker_config.VHost,
		ConnectionName: conn_name,
	}

	broker := rabbit.NewBroker(BrokerConfig)
	err = broker.Connect()
	if err != nil {
		log.Fatal(err)
	}

	consumerConf := handler.ConsumerConfig{
		ExchangeName: "main",
		ExchangeType: "direct",
		RoutingKey:   conn_name,
		QueueName:    conn_name,
		ConsumerName: conn_name,
		MaxAttempt:   60,
		Interval:     1 * time.Second,
		Connection:   broker.BrokerConnection,
	}

	producerConf := rabbit.ProducerConfig{
		ExchangeName: "main",
		ExchangeType: "direct",
		MaxAttempt:   60,
		Interval:     1 * time.Second,
		RoutingKey:   "backupserver",
		Connection:   broker.BrokerConnection,
	}

	producer := rabbit.NewProducer(producerConf)
	socket_repo := socket.NewRepository("localhost", "8080", "tcp")

	backup_service := backup.NewBackupService(producer, socket_repo, wtf)

	consumer_repo := handler.NewConsumerRepo(consumerConf, *backup_service)
	consumer_chan, err := consumer_repo.Start()

	if err != nil {
		log.Fatal("ERR", err)
	}

	consumer_repo.Consume(consumer_chan)
}
