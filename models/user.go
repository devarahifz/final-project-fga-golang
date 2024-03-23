package models

import (
	"errors"
	"final-project/helpers"
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"not null;type:varchar; uniqueIndex" valid:"required~Your username is required"`
	Email        string `gorm:"not null;type:varchar; uniqueIndex" valid:"required~Your email is required,email~Invalid email format"`
	Password     string `gorm:"not null;type:varchar" valid:"required~Your password is required,minstringlength(6)~Your password must be at least 6 characters long"`
	Age          int    `gorm:"not null;type:int" valid:"required~Your age is required"`
	Photo        []Photo
	Comments     []Comment
	SocialMedias []SocialMedia
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Age < 8 {
		return errors.New("your age must be at least 8 years old")
	}

	_, errCreate := govalidator.ValidateStruct(u)

	if errCreate != nil {
		err = errCreate
		return
	}

	u.Password = helpers.HashPassword(u.Password)
	err = nil
	return err
}
