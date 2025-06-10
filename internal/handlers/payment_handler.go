package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"github.com/thienhi/fusionstart/internal/dto"
	"github.com/thienhi/fusionstart/internal/workers"
)

func PaymentHookHandler(ch *amqp.Channel) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payment dto.PaymentDTO
		if err := c.BindJSON(&payment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		// can ->  publish payment event to payment worker
		err := workers.PublishPaymentTicket(ch, payment.BookingID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish payment"})
			return
		}
		c.JSON(http.StatusOK, payment)
	}
}
