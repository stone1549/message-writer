# message-writer

Service that consumes chat messages from Kafka and writes them to the database. Written in Go.

---

## Prereqs
You probably want to checkout the monorepo [yapyapyap](https://www.github.com/stone1549/yapyapyap) instead

## Configuration

### Environment Variables

| Variable                     | Description                                    | Possible Values      |
|------------------------------|------------------------------------------------|----------------------|
| MESSAGE_WRITER_ENVIRONMENT   | Controls log levels and configuration defaults | DEV, PRE_PROD, PROD  |
| MESSAGE_WRITER_GROUP_ID      | Sets the group id the consumer will provide    | string               |
| MESSAGE_WRITER_TOPIC         | Sets the topic the consumer will read from     | string               |  
| MESSAGE_WRITER_KAFKA_BROKERS | A comma separated list of kafka brokers        | host:9092,host2:9092 |
| MESSAGE_WRITER_PG_URL        | Full connection string for PG                  | string               |

## Run

```go run main.go```