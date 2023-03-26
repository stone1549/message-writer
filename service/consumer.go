package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"strings"
	"sync"
	"time"
)

type Consumer interface {
	Start() error
	Stop() error
}

var mutex sync.RWMutex

type kafkaConsumer struct {
	adminClient *kafka.AdminClient
	consumer    *kafka.Consumer
	repo        MessageRepository
	run         bool
	*sync.RWMutex
}

func (c *kafkaConsumer) Start() error {
	var err error
	c.RWMutex.RLock()
	run := c.run
	c.RWMutex.RUnlock()
	for run {
		c.RWMutex.RLock()
		run = c.run
		c.RWMutex.RUnlock()
		msg, err := c.consumer.ReadMessage(time.Second)

		if err != nil && err.(kafka.Error).IsTimeout() {
			continue
		} else if err != nil {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			c.RWMutex.Lock()
			c.run = false
			c.RWMutex.Unlock()
			continue
		}

		decoder := json.NewDecoder(strings.NewReader(string(msg.Value)))
		var message KafkaMessage
		err = decoder.Decode(&message)
		if err != nil {
			fmt.Printf("Unable to decode message: %s", err.Error())
			//TODO: put error on separate topic
			_, err = c.consumer.CommitMessage(msg)
			if err != nil {
				fmt.Printf("Unable to store offset: %s", err.Error())
			}
			continue
		}

		err = c.repo.AddMessage(message)

		if err != nil {
			switch err {
			case sql.ErrNoRows:
				println("No rows found")
				fmt.Printf("Unable to store message: %s", err.Error())
				//TODO: put error on separate topic
				_, err = c.consumer.CommitMessage(msg)
				if err != nil {
					fmt.Printf("Unable to store offset: %s", err.Error())
				}
				continue
			case sql.ErrConnDone:
				println("Connection closed")
				c.RWMutex.Lock()
				c.run = false
				c.RWMutex.Unlock()
				fmt.Printf("Unable to store message: %s", err.Error())
				continue
			default:
				fmt.Printf("Unable to store message: %s", err.Error())
				//TODO: put error on separate topic
				_, err = c.consumer.CommitMessage(msg)
				if err != nil {
					fmt.Printf("Unable to store offset: %s", err.Error())
				}
				continue
			}
		}

		_, err = c.consumer.CommitMessage(msg)

		if err != nil {
			fmt.Printf("Unable to store offset: %s", err.Error())
		}

		c.RWMutex.RLock()
		run = c.run
		c.RWMutex.RUnlock()
	}

	_ = c.consumer.Close()

	return err
}

func (c *kafkaConsumer) Stop() error {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()

	c.run = false
	return nil
}

func NewConsumer(config Configuration, repo MessageRepository) (Consumer, error) {
	brokers := config.GetKafkaBrokers()
	var brokerStr = ""
	for _, broker := range brokers {
		brokerStr += broker + ","
	}
	brokerStr = brokerStr[:len(brokerStr)-1]

	ac, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": brokerStr,
	})

	if err != nil {
		return nil, err
	}

	res, err := ac.CreateTopics(context.Background(), []kafka.TopicSpecification{{
		Topic:         config.GetTopic(),
		NumPartitions: 1,
	}}, nil)

	if err != nil {
		return nil, err
	}

	for _, topic := range res {
		if topic.Error.Code() != kafka.ErrNoError && topic.Error.Code() != kafka.ErrTopicAlreadyExists {
			fmt.Printf("Failed to create topic %s: %v\n", topic.Topic, topic.Error)
			return nil, topic.Error
		} else if topic.Error.Code() == kafka.ErrTopicAlreadyExists {
			fmt.Printf("Topic %s already exists\n", topic.Topic)
		} else {
			fmt.Printf("Topic %s created\n", topic.Topic)
		}
	}
	if err != nil {
		return nil, err
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        brokerStr,
		"enable.auto.offset.store": false,
		"message.max.bytes":        1024,
		"enable.auto.commit":       false,
		"enable.idempotence":       true,
		"group.id":                 config.GetGroupId(),
	})

	if err != nil {
		return nil, err
	}

	err = consumer.SubscribeTopics([]string{config.GetTopic()}, nil)

	if err != nil {
		return nil, err
	}

	return &kafkaConsumer{consumer: consumer, adminClient: ac, run: true, RWMutex: &mutex, repo: repo}, nil
}
