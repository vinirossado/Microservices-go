package event

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Emitter struct {
	Connection *amqp.Connection
}

func (e *Emitter) setup() error {
	channel, err := e.Connection.Channel()
	if err != nil {
		return err
	}

	defer func(channel *amqp.Channel) {
		err := channel.Close()
		if err != nil {

		}
	}(channel)

	return declareExchange(channel)
}

func (e *Emitter) Push(event, severity string) error {
	channel, err := e.Connection.Channel()
	if err != nil {
		return err
	}

	defer func(channel *amqp.Channel) {
		err := channel.Close()
		if err != nil {

		}
	}(channel)

	log.Println("Pushing to Channel")
	ctx := context.Background()

	err = channel.PublishWithContext(
		ctx,
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func NewEventEmmiter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		Connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
