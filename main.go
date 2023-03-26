package main

import (
	"fmt"
	"github.com/stone1549/yapyapyap/message-writer/service"
	"os"
	"os/signal"
)

func main() {
	config, err := service.GetConfiguration()

	if err != nil {
		panic(fmt.Sprintf("Unable to load configuration: %s", err.Error()))
	}

	repo, err := service.NewMessageRepository(config)

	if err != nil {
		panic(fmt.Sprintf("Unable to create repository: %s", err.Error()))
	}
	consumer, err := service.NewConsumer(config, repo)

	if err != nil {
		panic(fmt.Sprintf("Unable to create consumer: %s", err.Error()))
	}

	consumerChan := make(chan error, 1)

	go func() {
		consumerChan <- consumer.Start()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	select {
	case err := <-consumerChan:
		if err != nil {
			panic(fmt.Sprintf("Consumer error: %s", err.Error()))
		}
		_ = consumer.Stop()
	case <-c:
		_ = consumer.Stop()
	}
}
