package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Photo struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"not null;type:varchar" json:"title" valid:"required~Your title is required"`
	Caption   string    `gorm:";type:varchar" json:"caption"`
	PhotoURL  string    `gorm:"not null;type:varchar" json:"photo_url" valid:"required~Your photo URL is required"`
	UserID    uint      `gorm:"foreignKey" json:"user_id"`
	User      User      `json:"user" valid:"-"`
	Comments  []Comment `json:"comments,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (p *Photo) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(p)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return err
}
