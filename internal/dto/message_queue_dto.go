package dto

type BookingMsgQueueDTO struct {
	Event     string `json:"event" binding:"required"`
	BookingID uint   `json:"booking_id" binding:"required"`
}

type PaymentMsgQueueDTO struct {
	BookingID uint `json:"booking_id" binding:"required"`
}
