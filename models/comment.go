package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"foreignKey" json:"user_id"`
	PhotoID   uint      `gorm:"foreignKey" json:"photo_id"`
	Message   string    `gorm:"not null;type:varchar" json:"message" valid:"required~Your message is required"`
	User      User      `valid:"-"`
	Photo     Photo     `valid:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(c)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return err
}
