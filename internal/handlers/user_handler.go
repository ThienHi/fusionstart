package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thienhi/fusionstart/internal/dto"
	"github.com/thienhi/fusionstart/internal/repositories"
	"github.com/thienhi/fusionstart/internal/utils"
)

func RegisterController(repository repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dtoRegister dto.UserRegisterDTO
		if err := c.BindJSON(&dtoRegister); err != nil {
			res := utils.Response(100, false, "Register failed", map[string]interface{}{
				"error": err.Error(),
			})
			c.JSON(http.StatusBadRequest, res)
			return
		}
		dtoRegister.Password = utils.HashingPassword(dtoRegister.Password)
		registerUser := repository.CreateUser(dtoRegister)
		res := utils.Response(100, false, "Register successfully", map[string]interface{}{})
		if registerUser != nil {
			res.Message = "Register Failed"
			res.Error = true
			c.JSON(http.StatusBadRequest, res)
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func LoginController(repository repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		res := utils.Response(100, false, "Login successfully", nil)
		var dtoLogin dto.UserLoginDTO
		if err := c.BindJSON(&dtoLogin); err != nil {
			res.Message = "Invalid username or password"
			res.Error = true
			res.Data = "Invalid username or password"
		}
		loginUser, err := repository.FindByEmail(dtoLogin.Email)
		if err != nil {
			res.Message = "Email not found"
			res.Error = true
			res.Data = err.Error()
		}

		comparePwd := utils.ComparePassword(loginUser.Password, dtoLogin.Password)
		if !comparePwd {
			res.Message = "Password not found"
			res.Error = true
			res.Data = "Password not found"
		} else {
			token := utils.GenerateToken(loginUser.Email, loginUser.Password)
			res.Data = token
		}
		c.JSON(http.StatusOK, res)
	}
}
