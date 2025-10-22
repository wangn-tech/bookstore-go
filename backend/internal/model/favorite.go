package model

import "time"

type Favorite struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `json:"user_id"`
	BookID    uint64    `json:"book_id"`
	CreatedAt time.Time `json:"created_at"`

	Book *Book `json:"book,omitempty" gorm:"foreignKey:BookID"`
}

func (f *Favorite) TableName() string {
	return "favorites"
}
