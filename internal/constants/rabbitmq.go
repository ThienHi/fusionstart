package constants

const (
	// RabbitMQ Constants of Booking Cancel Ticket
	MESSAGE_TTL                    = 15 * 60 * 1000
	BOOKING_EXCHANGE_NAME          = "booking.cancel.exchange"
	BOOKING_QUEUE_NAME             = "booking.cancel.queue"
	BOOKING_DEAD_LETTER_QUEUE_NAME = "booking.cancel.dlq"
	BOOKING_ROUTING_KEY            = "booking.cancel.routing.key"
	BOOKING_ROUTING_KEY_DLQ        = "booking.cancel.routing.key.dlq"

	// RabbitMQ Constants of Payment Hook
	PAYMENT_EXCHANGE_NAME = "payment.exchange"
	PAYMENT_QUEUE_NAME    = "payment.queue"
	PAYMENT_ROUTING_KEY   = "payment.routing.key"
)
