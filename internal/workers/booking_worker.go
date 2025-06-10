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

func BookingCancelTicketConsumer(ch *amqp.Channel, db *gorm.DB) {
	msgs, err := ch.Consume(constants.BOOKING_DEAD_LETTER_QUEUE_NAME, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			var (
				bookingMsg dto.BookingMsgQueueDTO
				booking    models.Booking
				event      models.Event
			)
			if err := json.Unmarshal(msg.Body, &bookingMsg); err != nil {
				log.Fatalf("Invalid message: %s", err)
			}

			err := db.Transaction(func(tx *gorm.DB) error {
				// Update Booking
				if err := tx.First(&booking, "id = ?", bookingMsg.BookingID).Error; err != nil {
					log.Printf("Booking not found: %v", bookingMsg.BookingID)
					return nil
				}
				if booking.Status != constants.BookingStatusPending {
					log.Printf("Booking already processed: %v", bookingMsg.BookingID)
					return nil
				}
				cancelledAt := time.Now()
				updateErr := tx.Model(&booking).Updates(map[string]interface{}{
					"status":       constants.BookingStatusCancelled,
					"updated_at":   time.Now(),
					"cancelled_at": &cancelledAt,
				}).Error
				if updateErr != nil {
					log.Printf("Error updating booking: %v", err)
					return updateErr
				}

				// Update Event
				if err := tx.First(&event, "id = ?", booking.EventID).Error; err != nil {
					log.Printf("Event not found: %v", bookingMsg.BookingID)
					return nil
				}
				updateErr = tx.Model(&event).Updates(map[string]interface{}{
					"available_tickets": event.AvailableTickets + booking.Quantity,
					"updated_at":        time.Now(),
				}).Error
				if updateErr != nil {
					log.Printf("Error updating event: %v", err)
					return updateErr
				}

				log.Printf("Booking cancelled: %v", bookingMsg.BookingID)
				return nil
			})
			if err != nil {
				log.Printf("Error processing message: %v", err)
			}
		}
	}()
}
