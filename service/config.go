package service

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	lifeCycleKey    string = "MESSAGE_WRITER_ENVIRONMENT"
	pgUrlKey        string = "MESSAGE_WRITER_PG_URL"
	groupIdKey      string = "MESSAGE_WRITER_GROUP_ID"
	topicKey        string = "MESSAGE_WRITER_TOPIC"
	kafkaBrokersKey string = "MESSAGE_WRITER_KAFKA_BROKERS"
)

// LifeCycle represents a particular application life cycle.
type LifeCycle int

const (
	// DevLifeCycle represents the development environment.
	DevLifeCycle LifeCycle = 0
	// PreProdLifeCycle represents the pre production environment.
	PreProdLifeCycle LifeCycle = iota
	// ProdLifeCycle represents the production environment.
	ProdLifeCycle LifeCycle = iota
)

func (lc LifeCycle) String() string {
	switch lc {
	case DevLifeCycle:
		return "DEV"
	case PreProdLifeCycle:
		return "PRE_PROD"
	case ProdLifeCycle:
		return "PROD"
	default:
		return ""
	}
}

// Configuration provides methods for retrieving aspects of the applications configuration.
type Configuration interface {
	// GetLifeCycle retrieves the configured life cycle.
	GetLifeCycle() LifeCycle
	// GetPgUrl retrieves the configured url string for connecting to PostgreSQL.
	GetPgUrl() string
	GetGroupId() string
	// GetTopic retrieves the configured topic.
	GetTopic() string
	// GetKafkaBrokers retrieves the configured Kafka brokers.
	GetKafkaBrokers() []string
}

type configuration struct {
	lifeCycle    LifeCycle
	pgUrl        string
	topic        string
	kafkaBrokers []string
	groupId      string
}

func (conf *configuration) GetLifeCycle() LifeCycle {
	return conf.lifeCycle
}

func (conf *configuration) GetPgUrl() string {
	return conf.pgUrl
}
func (conf *configuration) GetGroupId() string {
	return conf.groupId
}
func (conf *configuration) GetTopic() string {
	return conf.topic
}
func (conf *configuration) GetKafkaBrokers() []string {
	return conf.kafkaBrokers
}

// GetConfiguration constructs a Configuration based on environment variables.
func GetConfiguration() (Configuration, error) {
	var err error
	config := configuration{}

	lcStr := os.Getenv(lifeCycleKey)

	switch lcStr {
	case DevLifeCycle.String():
		config.lifeCycle = DevLifeCycle
	case PreProdLifeCycle.String():
		config.lifeCycle = PreProdLifeCycle
	case ProdLifeCycle.String():
		config.lifeCycle = ProdLifeCycle
	default:
		config.lifeCycle = DevLifeCycle
	}

	err = setPostgresqlConfig(&config)

	if err != nil {
		return nil, err
	}

	config.topic = os.Getenv(topicKey)

	if config.topic == "" {
		return nil, errors.New(fmt.Sprintf("No topic configured, set %s environment variable", topicKey))
	}

	brokersStr := os.Getenv(kafkaBrokersKey)

	config.kafkaBrokers = strings.Split(brokersStr, ",")
	if len(config.kafkaBrokers) == 0 {
		return nil, errors.New(fmt.Sprintf("No Kafka brokers configured, set %s environment variable", kafkaBrokersKey))
	}

	config.groupId = os.Getenv(groupIdKey)

	if config.groupId == "" {
		return nil, errors.New(fmt.Sprintf("No group id configured, set %s environment variable", groupIdKey))
	}

	return &config, nil
}

func setPostgresqlConfig(config *configuration) error {
	var err error

	config.pgUrl = os.Getenv(pgUrlKey)

	if strings.TrimSpace(config.pgUrl) == "" {
		err = errors.New(fmt.Sprintf("No PostgreSqlRepo url configured, set %s environment variable", pgUrlKey))
	}

	return err
}
