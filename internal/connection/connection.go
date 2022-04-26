package connection

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

func NewRabbitmq() *amqp.Connection {
	connection, err := amqp.Dial(os.Getenv("MESSAGE_BROKER_URI"))
	if err != nil {
		panic(fmt.Errorf("can not dial to RabbitMQ server, err: %s", err.Error()))
	}

	return connection
}
