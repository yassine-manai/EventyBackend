package models

import "github.com/uptrace/bun"

type Category struct {
	bun.BaseModel `json:"-" bun:"table:category"`
	CategoryID    int    `bun:"category_id,autoincrement,pk" json:"category_id"`
	CategoryName  string `bun:"category_name" json:"category_name" binding:"required"`
}

type CategoryNoBind struct {
	bun.BaseModel `json:"-" bun:"table:category"`
	CategoryID    int    `bun:"category_id,autoincrement,pk" json:"category_id"`
	CategoryName  string `bun:"category_name" json:"category_name"`
}
