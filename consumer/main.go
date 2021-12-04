package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/infernalslam/idenpotency-messaing-rabbitmq/pkg/cache"
	"github.com/infernalslam/idenpotency-messaing-rabbitmq/pkg/rabbitmq"
	"github.com/streadway/amqp"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	redis   *redis.Client
	done    chan error
}

func (c Consumer) Shutdown() error {
	if err := c.channel.Cancel(*consumerTag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed : %s", err)
	}
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error : ", err)
	}
	defer fmt.Println("AMQP shutdown OK")
	return <-c.done
}

var lifetime = flag.Duration("lifetime", 500*time.Second, "lifetime of process before shutdown (0s=infinite)")
var consumerTag = flag.String("consumer-tag", "message-consumer", "AMQP consumer tag (should not be blank)")

func main() {
	rcof := cache.RedisConf{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
	redisMessage, err := cache.NewConnection(rcof)
	if err != nil {
		fmt.Println("error redis : ", err)
		return
	}

	c, err := Consume(redisMessage)
	if err != nil {
		fmt.Print("Cannot start consume ....")
		return
	}

	if *lifetime > 0 {
		log.Printf("running for %s", *lifetime)
		time.Sleep(*lifetime)
	} else {
		log.Printf("running forever")
		select {}
	}

	fmt.Println("Shut down server ....")

	if err := c.Shutdown(); err != nil {
		fmt.Println("error cannot shutdown .... ", err)
	}

}

func Consume(redis *redis.Client) (*Consumer, error) {
	conn, err := rabbitmq.NewConnection("amqp://guest:guest@localhost:5672//")
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	c := &Consumer{
		conn:    conn.Conn,
		channel: conn.Channel,
		done:    make(chan error),
		redis:   redis,
	}

	go func() {
		fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	consumer, err := c.channel.Consume(
		"message",
		*consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println("Error Consumer : ", err)
	}

	go handle(consumer, c)

	return c, nil
}

func handle(consumer <-chan amqp.Delivery, c *Consumer) {
	for d := range consumer {
		// check messageID
		// messageIDFromRedis := redis.Get(context.Background(), d.MessageId).Result()
		redisMessageID, err := c.redis.Get(context.Background(), d.MessageId).Result()
		if err != nil && err != redis.Nil {
			fmt.Println("Error get redis conusmer ... ", err)
		}

		if redisMessageID == d.MessageId {
			fmt.Println("Same messageID : ", redisMessageID)
			continue
		}
		// fmt.Printf("%+v\n\n", d)
		fmt.Println(d.MessageId)
		// process ....

		err = c.redis.Set(context.Background(), d.MessageId, d.MessageId, 180*time.Second).Err()
		if err != nil {
			fmt.Println("Error set redis conusmer ... ", err)
		}
		d.Ack(true)
	}
	fmt.Println("handle closing consumer ....")
	c.done <- nil
}
