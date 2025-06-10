package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"github.com/thienhi/fusionstart/internal/configs"
	"github.com/thienhi/fusionstart/internal/constants"
)

func ConnectRabbitMQ(config *configs.Config) (*amqp.Connection, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", config.RabbitMQ.User, config.RabbitMQ.Password, config.RabbitMQ.Host, config.RabbitMQ.Port)
	con, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}
	return con, nil
}

func CloseConnectionRabbitMQ(con *amqp.Connection) error {
	return con.Close()
}

func SetupRabbitMQ(rbmq *amqp.Connection) {
	// setup booking queue
	channelBooking, err := rbmq.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	SetupBookingQueue(channelBooking)
	defer channelBooking.Close()

	// // setup payment queue
	channelPayment, err := rbmq.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	SetupPaymentQueue(channelPayment)
	defer channelPayment.Close()
}

func DeclareQueue(ch *amqp.Channel, queueName string, args *amqp.Table) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		*args,     // arguments
	)
	if err != nil {
		return amqp.Queue{}, err
	}
	return q, nil
}

func SetupBookingQueue(ch *amqp.Channel) (amqp.Queue, error) {
	ch.ExchangeDeclare(constants.BOOKING_EXCHANGE_NAME, "direct", true, false, false, false, nil)
	// declare queue - dead letter queue
	args := amqp.Table{
		"x-dead-letter-exchange":    constants.BOOKING_EXCHANGE_NAME,
		"x-dead-letter-routing-key": constants.BOOKING_ROUTING_KEY_DLQ,
		"x-message-ttl":             int32(constants.MESSAGE_TTL),
	}
	q, err := DeclareQueue(ch, constants.BOOKING_QUEUE_NAME, &args)
	if err != nil {
		log.Printf("Failed to declare Booking Queue dead letter: %s", err)
		return amqp.Queue{}, err
	}
	// Declare DLQ
	ch.QueueBind(constants.BOOKING_QUEUE_NAME, constants.BOOKING_ROUTING_KEY, constants.BOOKING_EXCHANGE_NAME, false, nil)
	ch.QueueDeclare(constants.BOOKING_DEAD_LETTER_QUEUE_NAME, true, false, false, false, nil)
	ch.QueueBind(constants.BOOKING_DEAD_LETTER_QUEUE_NAME, constants.BOOKING_ROUTING_KEY_DLQ, constants.BOOKING_EXCHANGE_NAME, false, nil)

	return q, nil
}

func SetupPaymentQueue(ch *amqp.Channel) (amqp.Queue, error) {
	// declare queue - dead letter queue
	args := amqp.Table{}
	q, err := DeclareQueue(ch, constants.PAYMENT_QUEUE_NAME, &args)
	if err != nil {
		log.Printf("Failed to declare Payment Queue: %s", err)
		return amqp.Queue{}, err
	}
	// Declare DLX and DLQ
	ch.ExchangeDeclare(constants.PAYMENT_EXCHANGE_NAME, "direct", true, false, false, false, nil)
	ch.QueueBind(constants.PAYMENT_QUEUE_NAME, constants.PAYMENT_ROUTING_KEY, constants.PAYMENT_EXCHANGE_NAME, false, nil)

	return q, nil
}
