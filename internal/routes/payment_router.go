package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"github.com/thienhi/fusionstart/internal/handlers"
	"gorm.io/gorm"
)

func PaymentRouter(r *gin.RouterGroup, db *gorm.DB, ch *amqp.Channel) {
	r.POST("/", handlers.PaymentHookHandler(ch))
}
