package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.wndv.co/jirasak/poc-rabbit-dlx-ttl/internal/connection"
	"github.com/streadway/amqp"

	_ "github.com/joho/godotenv/autoload"
)

const (
	exchangeName = "e-2"
	queueName    = "q-2"
	routingKey   = "*.*"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	conn := connection.NewRabbitmq()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	queue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := ch.QueueBind(queue.Name, routingKey, exchangeName, false, nil); err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	go func() {
		for msg := range msgs {
			fmt.Println(msg.Headers)
			fmt.Println(msg.Redelivered)
			_ = msg.Reject(true)
		}
	}()

	closeChannel := make(chan os.Signal, 1)
	signal.Notify(closeChannel, syscall.SIGINT, syscall.SIGTERM)

	rabbitClose := make(chan *amqp.Error)
	ch.NotifyClose(rabbitClose)

	fmt.Println("running !!")
	select {
	case <-closeChannel:
		fmt.Println("ctrl + c")
	case <-rabbitClose:
		fmt.Printf("rabbitmq close")
	}
}
