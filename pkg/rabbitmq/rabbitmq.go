package rabbitmq

import "github.com/streadway/amqp"

type Connection struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewConnection(uri string) (*Connection, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	return &Connection{Conn: conn, Channel: channel}, nil
}
