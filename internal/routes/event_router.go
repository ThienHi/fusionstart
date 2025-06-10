package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thienhi/fusionstart/internal/handlers"
	"github.com/thienhi/fusionstart/internal/repositories"
	"gorm.io/gorm"
)

func EventRouter(r *gin.RouterGroup, db *gorm.DB) {
	eventRepository := repositories.NewEventRepository(db)
	eventHandlers := handlers.NewEventHandlers(eventRepository)

	r.GET("/", eventHandlers.GetEventHandler())
	r.GET("/:id", eventHandlers.GetDetailEventHandler())
	r.POST("/", eventHandlers.CreateEventHandler())
	r.PUT("/:id", eventHandlers.UpdateEventHandler())
	r.DELETE("/:id", eventHandlers.DeleteEventHandler())
}
