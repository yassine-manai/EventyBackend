package models

import "github.com/uptrace/bun"

////////// THIS FILE REPRESENT STRCTS FOR EVENT TABLE //////////

type Event struct {
	bun.BaseModel `json:"-" bun:"table:event"`
	EventID       int    `bun:"event_id,autoincrement,pk" json:"event_id" `
	Title         string `bun:"title" json:"title" binding:"required"`
	StartDate     string `bun:"start_date" json:"start_date" binding:"required"`
	EndDate       string `bun:"end_date" json:"end_date" binding:"required"`
	Location      string `bun:"location" json:"location" binding:"required"`
	Image         string `bun:"image,type:bytea" json:"image"`
	Category      int    `bun:"category" json:"category" binding:"required"`
	Capacity      int    `bun:"capacity" json:"capacity" binding:"required"`
	IsArchived    bool   `bun:"isArchived" json:"isArchived"`
	Price         int    `bun:"price" json:"price" binding:"required"`
	UserID        []int  `bun:"user_id,array" json:"user_id" `
}

type EventNoBind struct {
	bun.BaseModel `json:"-" bun:"table:event"`
	EventID       int    `bun:"event_id,autoincrement,pk" json:"event_id" `
	Title         string `bun:"title" json:"title"`
	StartDate     string `bun:"start_date" json:"start_date"`
	EndDate       string `bun:"end_date" json:"end_date"`
	Location      string `bun:"location" json:"location"`
	Image         string `bun:"image,type:bytea" json:"image"`
	Capacity      int    `bun:"capacity" json:"capacity"`
	IsArchived    bool   `bun:"isArchived" json:"isArchived"`
	Category      int    `bun:"category" json:"category"`
	Price         int    `bun:"price" json:"price"`
	UserID        []int  `bun:"user_id,array" json:"user_id" `
}
