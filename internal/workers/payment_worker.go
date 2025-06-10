package workers

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
	"github.com/thienhi/fusionstart/internal/constants"
	"github.com/thienhi/fusionstart/internal/dto"
	"github.com/thienhi/fusionstart/internal/models"
	"gorm.io/gorm"
)

func PaymentTicketConsumer(ch *amqp.Channel, db *gorm.DB) {
	msgs, err := ch.Consume(constants.PAYMENT_QUEUE_NAME, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register payment consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			var (
				paymentMsg dto.PaymentMsgQueueDTO
				booking    models.Booking
			)
			if err := json.Unmarshal(msg.Body, &paymentMsg); err != nil {
				log.Fatalf("Invalid message: %s", err)
			}

			err := db.Transaction(func(tx *gorm.DB) error {
				// Update Booking
				if err := tx.First(&booking, "id = ?", paymentMsg.BookingID).Error; err != nil {
					log.Printf("Booking not found: %v", paymentMsg.BookingID)
					return nil
				}
				if booking.Status != constants.BookingStatusPending {
					log.Printf("Booking already processed: %v", paymentMsg.BookingID)
					return nil
				}
				confirmedAt := time.Now()
				updateErr := tx.Model(&booking).Updates(map[string]interface{}{
					"status":       constants.BookingStatusConfirmed,
					"updated_at":   time.Now(),
					"comfirmed_at": &confirmedAt,
				}).Error
				if updateErr != nil {
					log.Printf("Error updating booking: %v", err)
					return updateErr
				}

				log.Printf("Booking comfirmed: %v", paymentMsg.BookingID)
				return nil
			})
			if err != nil {
				log.Printf("Error processing message: %v", err)
			}
		}
	}()
}
