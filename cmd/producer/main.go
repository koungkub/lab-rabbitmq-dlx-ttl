package main

import (
	"log"

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

	if err := ch.ExchangeDeclare("e-1", "topic", true, false, false, false, amqp.Table{}); err != nil {
		log.Fatal(err)
	}

	if err := ch.Publish("e-1", "order.messenger", false, false, amqp.Publishing{
		Body: []byte("hiw-kanom"),
	}); err != nil {
		log.Printf("%s\n", err)
	}
}
