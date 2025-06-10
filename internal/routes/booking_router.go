package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"github.com/thienhi/fusionstart/internal/handlers"
	"github.com/thienhi/fusionstart/internal/repositories"
	"gorm.io/gorm"
)

func BookingRouter(r *gin.RouterGroup, db *gorm.DB, ch *amqp.Channel) {
	bookingRepository := repositories.NewBookingRepository(db)
	bookingHandler := handlers.NewBookingHandler(bookingRepository)

	r.GET("/", bookingHandler.GetBookingsHandler())
	r.POST("/", bookingHandler.CreateBooking(ch))
}
