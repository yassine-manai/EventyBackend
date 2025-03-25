package models

type BookEventRequest struct {
	EventID int `json:"event_id" binding:"required"`
	UserID  int `json:"user_id" binding:"required"`
}
