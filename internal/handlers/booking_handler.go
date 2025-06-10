package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/streadway/amqp"
	"github.com/thienhi/fusionstart/internal/dto"
	"github.com/thienhi/fusionstart/internal/repositories"
	"github.com/thienhi/fusionstart/internal/utils"
	"github.com/thienhi/fusionstart/internal/workers"
)

type BookingHandler interface {
	CreateBooking() gin.HandlerFunc
	GetBookingsHandler() gin.HandlerFunc
	// CancelBooking() gin.HandlerFunc
}

type bookingHandler struct {
	bookingRepository repositories.BookingRepository
	validator         *validator.Validate
}

func NewBookingHandler(bookingRepository repositories.BookingRepository) *bookingHandler {
	return &bookingHandler{
		bookingRepository: bookingRepository,
		validator:         validator.New(),
	}
}

func (b *bookingHandler) GetBookingsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		bookings, err := b.bookingRepository.GetAll()
		res := utils.Response(100, false, "Get Booking successfully", map[string]interface{}{})
		if err != nil {
			res.Code = 201
			res.Error = true
			res.Message = "Get Booking failed"
			res.Data = map[string]interface{}{"error": err.Error()}
			c.JSON(http.StatusBadRequest, res)
			return
		}
		res.Data = bookings
		c.JSON(http.StatusOK, res)
	}
}

func (b *bookingHandler) CreateBooking(ch *amqp.Channel) gin.HandlerFunc {
	return func(c *gin.Context) {
		var bookingInput dto.BookingCreateDTO
		res := utils.Response(100, false, "Create Booking successfully", map[string]interface{}{})
		if err := b.validator.Struct(bookingInput); err != nil {
			res.Code = 202
			res.Error = true
			res.Message = "Validation failed"
			res.Data = map[string]interface{}{"errors": utils.FormatValidationErrors(err)}
			c.JSON(http.StatusBadRequest, res)
			return
		}

		if err := c.BindJSON(&bookingInput); err != nil {
			res.Code = 201
			res.Error = true
			res.Message = "Bind JSON failed"
			res.Data = map[string]interface{}{"error": err.Error()}
			c.JSON(http.StatusBadRequest, res)
			return
		}

		booking, err := b.bookingRepository.Create(bookingInput)
		if err != nil {
			res.Code = 201
			res.Error = true
			res.Message = "Create Booking failed"
			res.Data = map[string]interface{}{"error": err.Error()}
			c.JSON(http.StatusBadRequest, res)
			return
		}
		errPublish := workers.PublishBookingCancelTicket(ch, booking.ID)
		if errPublish != nil {
			res.Code = 201
			res.Error = true
			res.Message = "Publish failed"
			res.Data = map[string]interface{}{"error": errPublish.Error()}
			c.JSON(http.StatusBadRequest, res)
			return
		}
		res.Data = bookingInput
		c.JSON(http.StatusOK, res)
	}
}
