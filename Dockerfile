FROM golang:1.20

ENV MESSAGE_WRITER_ENVIRONMENT=DEV
ENV MESSAGE_WRITER_PG_URL=postgres://postgres:postgres@yapyapyap-db:5432/postgres?sslmode=disable
ENV MESSAGE_WRITER_GROUP_ID=message-writer
ENV MESSAGE_WRITER_TOPIC=store
ENV MESSAGE_WRITER_KAFKA_BROKERS=yapyapyap-kafka:9092

WORKDIR /go/src/github.com/stone1549/yapyapyap/message-writer/
COPY . .

RUN go mod tidy

CMD ["go", "run", "main.go"]

