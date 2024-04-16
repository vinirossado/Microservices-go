package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
)

type Consumer struct {
	Conn      *amqp.Connection
	QueueName string
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		Conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.Conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.Conn.Channel()
	if err != nil {
		return err
	}

	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {

		}
	}(ch)

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, t := range topics {
		ch.QueueBind(q.Name, t, "logs_topic", false, nil)
		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()
	fmt.Printf("Waiting for messages on [Exchange, Queue] [logs_topic, %s]...\n", q.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	fmt.Printf("Received message: %s\n", payload)

	switch payload.Name {
	case "log", "event":
		err := logEvent(payload)
		if err != nil {
			fmt.Println(err)
		}

	case "auth":
	//Authenticate
	default:
		err := logEvent(payload)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest(http.MethodPost, logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    //name
		false, //durable?
		false, //auto-delete?
		true,  //exclusive?
		false, //no-wait?
		nil,   //arguments
	)
}
