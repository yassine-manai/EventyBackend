package models

import "github.com/uptrace/bun"

////////// THIS FILE REPRESENT STRCTS FOR USER TABLE //////////

type User struct {
	bun.BaseModel `json:"-" bun:"table:user"`
	UserID        int    `bun:"user_id,autoincrement" json:"user_id"`
	Email         string `bun:"email,pk" json:"email" binding:"required"`
	Password      string `bun:"password" json:"password" binding:"required"`
	Name          string `bun:"name" json:"name" binding:"required"`
	Is_guest      bool   `bun:"is_guest" json:"is_guest"`
	EventID       []int  `bun:"event_id" json:"event_id"`
	BookedEvents  []int  `bun:"booked_events" json:"booked_events"`
	Balance       int    `bun:"balance" json:"balance"`
}

type Login struct {
	bun.BaseModel `json:"-" bun:"table:user"`
	Email         string `bun:"email,pk" json:"email" binding:"required"`
	Password      string `bun:"password" json:"password" binding:"required"`
}
