package main

import (
	"fmt"

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

	for i := 0; i < 2; i++ {
		// uuid := uuid.NewString()
		uuid := "f78585d1-02ec-49b9-bf88-b08ec3f50c62"
		fmt.Println("producer "+uuid+" time : ", i+1)
		err = conn.Channel.Publish(
			"",
			"message",
			false,
			false,
			amqp.Publishing{
				MessageId:       uuid,
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

}
