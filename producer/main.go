package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/infernalslam/idenpotency-messaing-rabbitmq/pkg/rabbitmq"
	"github.com/streadway/amqp"
)

func main() {
	// connect rabbimtq & redis
	conn, err := rabbitmq.NewConnection("amqp://guest:guest@localhost:5672//")
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	// producer

	err = conn.Channel.Publish(
		"",
		"message",
		false,
		false,
		amqp.Publishing{
			MessageId:       uuid.NewString(),
			Headers:         amqp.Table{},
			ContentEncoding: "",
			Body:            []byte("test"),
			DeliveryMode:    amqp.Transient,
		},
	)
	if err != nil {
		fmt.Println("Cannot publishing : ", err)
	}
}
