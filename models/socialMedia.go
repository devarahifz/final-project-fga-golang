package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type SocialMedia struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"not null;type:varchar" json:"name" valid:"required~Your social media name is required"`
	SocialMediaURL string    `gorm:"not null;type:varchar" json:"social_media_url" valid:"required~Your social media URL is required"`
	UserID         uint      `gorm:"foreignKey" json:"userId"`
	User           User      `json:"user" valid:"-"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (s *SocialMedia) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(s)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return err
}
