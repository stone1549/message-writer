package service

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

const (
	insertMessage = "INSERT INTO message (id, user_id, content, location, client_id, sent_at, received_at) VALUES ($1, $2, $3, ST_GeomFromText($4), $5, $6, $7) RETURNING created_at"
)

// MessageRepository represents a data source through which users can be managed.
type MessageRepository interface {
	AddMessage(message KafkaMessage) error
}

type postgresqlMessageRepository struct {
	db *sql.DB
}

func (p *postgresqlMessageRepository) AddMessage(message KafkaMessage) error {
	row := p.db.QueryRow(
		insertMessage,
		message.Payload.Id,
		message.Payload.Message.Sender.Id,
		message.Payload.Message.Content,
		fmt.Sprintf("POINT (%f %f)", message.Payload.Message.Location.Long, message.Payload.Message.Location.Lat),
		message.Payload.Message.ClientId,
		message.Payload.Message.SentAt,
		message.Payload.Message.ReceivedAt,
	)

	var createdAt time.Time
	err := row.Scan(&createdAt)

	if err != nil {
		return err
	}

	return nil
}

func makePostgresqlRespository(db *sql.DB) (MessageRepository, error) {
	return &postgresqlMessageRepository{db}, nil
}

func NewMessageRepository(config Configuration) (MessageRepository, error) {
	var err error
	var repo MessageRepository
	//var db *sql.DB
	db, err := sql.Open("postgres", config.GetPgUrl())

	if err != nil {
		return nil, err
	}
	repo, err = makePostgresqlRespository(db)

	if err != nil {
		return nil, err
	}

	return repo, err
}
