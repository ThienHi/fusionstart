package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thienhi/fusionstart/internal/handlers"
	"github.com/thienhi/fusionstart/internal/repositories"
	"gorm.io/gorm"
)

func UserRouter(r *gin.RouterGroup, db *gorm.DB) {
	userRepository := repositories.NewUserRepository(db)
	r.POST("/register", handlers.RegisterController(userRepository))
	r.POST("/login", handlers.LoginController(userRepository))
}
