package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/thienhi/fusionstart/internal/dto"
	"github.com/thienhi/fusionstart/internal/repositories"
	"github.com/thienhi/fusionstart/internal/utils"
)

type EventHandlers interface {
	GetEventHandler() gin.HandlerFunc
	GetDetailEventHandler() gin.HandlerFunc
	CreateEventHandler() gin.HandlerFunc
	UpdateEventHandler() gin.HandlerFunc
	DeleteEventHandler() gin.HandlerFunc
}

type eventHandlers struct {
	repository repositories.EventRepository
	validator  *validator.Validate
}

func NewEventHandlers(repository repositories.EventRepository) *eventHandlers {
	return &eventHandlers{
		repository: repository,
		validator:  validator.New(),
	}
}

func (e *eventHandlers) GetEventHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		events, err := e.repository.GetAll()
		res := utils.Response(100, false, "Get Events successfully", map[string]interface{}{})
		if err != nil {
			res.Code = 200
			res.Error = true
			res.Message = "Get Events failed"
			res.Data = map[string]interface{}{"error": err.Error()}
			c.JSON(http.StatusBadRequest, res)
			return
		}
		res.Data = events
		c.JSON(http.StatusOK, res)
	}
}

func (e *eventHandlers) GetDetailEventHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		event, err := e.repository.FindById(uint(id))
		res := utils.Response(100, false, "Get Detail Event successfully", map[string]interface{}{})
		if err != nil {
			res.Code = 200
			res.Error = true
			res.Message = "Get Detail Event failed"
			res.Data = err.Error()
			c.JSON(http.StatusBadRequest, res)
			return
		}
		res.Data = event
		c.JSON(http.StatusOK, res)
	}
}

func (e *eventHandlers) UpdateEventHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		res := utils.Response(100, false, "Update Event successfully", map[string]interface{}{})
		if err != nil {
			res.Code = 203
			res.Error = true
			res.Message = "Invalid ID"
			res.Data = map[string]interface{}{"error": err.Error()}
			c.JSON(http.StatusBadRequest, res)
			return
		}
		var event dto.EventUpdateDTO
		if err := c.BindJSON(&event); err != nil {
			res.Code = 201
			res.Error = true
			res.Message = "Bind JSON failed"
			res.Data = map[string]interface{}{"error": err.Error()}
			c.JSON(http.StatusBadRequest, res)
			return
		}
		// Validate struct using validator
		if err := e.validator.Struct(event); err != nil {
			res.Code = 202
			res.Error = true
			res.Message = "Validation failed"
			res.Data = map[string]interface{}{"errors": utils.FormatValidationErrors(err)}
			c.JSON(http.StatusBadRequest, res)
			return
		}
		fmt.Println(event, " =================== ", id, uint(id))
		updateEvent, err := e.repository.Update(uint(id), event)
		if err != nil {
			res.Code = 202
			res.Error = true
			res.Message = "Update Event failed"
			res.Data = map[string]interface{}{"error": err.Error()}
			c.JSON(http.StatusBadRequest, res)
			return
		}
		res.Data = updateEvent
		c.JSON(http.StatusOK, res)
	}
}

func (e *eventHandlers) DeleteEventHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		res := utils.Response(100, false, "Delete Event successfully", map[string]interface{}{})
		if err != nil {
			res.Code = 203
			res.Error = true
			res.Message = "Invalid ID"
			res.Data = map[string]interface{}{"error": err.Error()}
			c.JSON(http.StatusBadRequest, res)
		}
		if err := e.repository.Delete(uint(id)); err != nil {
			res.Code = 204
			res.Error = true
			res.Message = "Delete Event failed"
			res.Data = map[string]interface{}{"error": err.Error()}
			c.JSON(http.StatusBadRequest, res)
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func (e *eventHandlers) CreateEventHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var event dto.EventCreateDTO
		res := utils.Response(100, false, "Create Event successfully", map[string]interface{}{})
		if err := c.BindJSON(&event); err != nil {
			res.Code = 201
			res.Error = true
			res.Message = "Bind JSON failed"
			res.Data = map[string]interface{}{"error": err.Error()}
			c.JSON(http.StatusBadRequest, res)
			return
		}
		// Validate struct using validator
		if err := e.validator.Struct(event); err != nil {
			res.Code = 202
			res.Error = true
			res.Message = "Validation failed"
			res.Data = map[string]interface{}{"errors": utils.FormatValidationErrors(err)}
			c.JSON(http.StatusBadRequest, res)
			return
		}
		if event.Datetime.Before(time.Now()) {
			res.Code = 201
			res.Error = true
			res.Message = "Event date must be in the future"
			res.Data = map[string]interface{}{"errors": map[string]string{"datetime": "Event date must be in the future"}}
			c.JSON(http.StatusBadRequest, res)
			return
		}
		if err := e.repository.Create(event); err != nil {
			res.Code = 205
			res.Error = true
			res.Message = "Create Event failed"
			res.Data = map[string]interface{}{"error": err.Error()}
			c.JSON(http.StatusBadRequest, res)
			return
		}
		res.Data = event
		c.JSON(http.StatusOK, res)
	}
}
