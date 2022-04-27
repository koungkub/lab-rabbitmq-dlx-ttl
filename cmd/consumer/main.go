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
	exchangeName = "e-1"
	queueName    = "q-1"
	routingKey   = "*.*"

	dlxExchange = "e-dlx"
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

	option := amqp.Table{
		"x-dead-letter-exchange": dlxExchange,
	}

	queue, err := ch.QueueDeclare(queueName, true, false, false, false, option)
	if err != nil {
		log.Fatal(err)
	}

	if err := ch.QueueBind(queue.Name, routingKey, exchangeName, false, option); err != nil {
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
			random := rand.Intn(2)
			if random == 1 {
				fmt.Println("reject", msg.Headers)
				fmt.Println(parsing(msg.Headers))
				_ = msg.Reject(false)
			} else {
				fmt.Println("ack", msg.Headers)
				fmt.Println(parsing(msg.Headers))
				_ = msg.Ack(false)
			}
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

func parsing(table amqp.Table) int {
	death, ok := table["x-death"]
	if !ok {
		return 0
	}

	m2, ok := death.([]interface{})
	if !ok {
		return 0
	}

	if len(m2) < 1 {
		return 0
	}

	m3, ok := m2[0].(amqp.Table)
	if !ok {
		return 0
	}

	m4, ok := m3["count"].(int64)
	if !ok {
		return 0
	}

	return int(m4)
}
