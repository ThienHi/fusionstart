package workers

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"github.com/thienhi/fusionstart/internal/constants"
	"github.com/thienhi/fusionstart/internal/dto"
)

func PublishBookingCancelTicket(ch *amqp.Channel, bookingID uint) error {
	job := dto.BookingMsgQueueDTO{BookingID: bookingID, Event: "booking_cancel_ticket"}
	body, _ := json.Marshal(job)

	err := ch.Publish(
		constants.BOOKING_EXCHANGE_NAME,
		constants.BOOKING_ROUTING_KEY,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Print("Failed to publish booking:", err)
	}
	log.Print("Success to publish booking")
	return err
}

func PublishPaymentTicket(ch *amqp.Channel, bookingID uint) error {
	job := dto.PaymentMsgQueueDTO{BookingID: bookingID}
	body, _ := json.Marshal(job)

	err := ch.Publish(
		constants.PAYMENT_EXCHANGE_NAME,
		constants.PAYMENT_ROUTING_KEY,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Print("Failed to publish payment:", err)
	}
	log.Print("Success to publish payment")
	return err
}
