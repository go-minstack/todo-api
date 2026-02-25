package todo_entities

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	Done        bool `gorm:"default:false"`
}
