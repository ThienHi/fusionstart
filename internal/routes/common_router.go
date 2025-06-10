package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"github.com/thienhi/fusionstart/internal/middleware"
	"gorm.io/gorm"
)

func SetupRouter(router *gin.Engine, db *gorm.DB, ch *amqp.Channel) {
	// // health check
	// router.GET("/health", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	// })

	api := router.Group("/api")
	// define router for user
	user_router := api.Group("/user")
	UserRouter(user_router, db)

	// define router for event
	event_router := api.Group("/event", middleware.Athentication())
	EventRouter(event_router, db)

	// define router for booking
	booking_router := api.Group("/booking", middleware.Athentication())
	BookingRouter(booking_router, db, ch)

	// define router for payment
	payment_router := api.Group("/payment")
	PaymentRouter(payment_router, db, ch)
}
