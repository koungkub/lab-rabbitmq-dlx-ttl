package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"git.wndv.co/jirasak/poc-rabbit-dlx-ttl/internal/connection"

	"github.com/streadway/amqp"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	conn := connection.NewRabbitmq()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	option := amqp.Table{
		"x-dead-letter-exchange": "e2",
		"x-message-ttl":          10000,
	}

	queue, err := ch.QueueDeclare("q-dlx", true, false, false, false, option)
	if err != nil {
		log.Fatal(err)
	}

	if err := ch.QueueBind(queue.Name, "*.*", "e-dlx", false, option); err != nil {
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
			fmt.Println(string(msg.Body), msg.Headers)
			_ = msg.Reject(false)
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
