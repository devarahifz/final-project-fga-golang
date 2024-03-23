package controllers

import (
	"final-project/database"
	"final-project/helpers"
	"final-project/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreatePhoto(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(ctx)

	photo := models.Photo{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		err := ctx.ShouldBindJSON(&photo)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
	} else {
		ctx.AbortWithStatus(http.StatusUnsupportedMediaType)
		return
	}

	photo.UserID = userID

	err := db.Debug().Create(&photo).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         photo.ID,
		"title":      photo.Title,
		"caption":    photo.Caption,
		"photo_url":  photo.PhotoURL,
		"user_id":    photo.UserID,
		"created_at": photo.CreatedAt,
	})
}

func GetPhoto(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)

	userID := uint(userData["id"].(float64))

	var photos []models.Photo

	if err := db.Where("user_id = ?", userID).Preload("User").Find(&photos).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	mappedPhotos := make([]gin.H, len(photos))
	for i, photo := range photos {
		mappedPhotos[i] = gin.H{
			"id":         photo.ID,
			"title":      photo.Title,
			"caption":    photo.Caption,
			"photo_url":  photo.PhotoURL,
			"user_id":    photo.UserID,
			"created_at": photo.CreatedAt,
			"updated_at": photo.UpdatedAt,
			"User": gin.H{
				"username": photo.User.Username,
				"email":    photo.User.Email,
			},
		}
	}
	ctx.JSON(http.StatusOK, mappedPhotos)
}

func UpdatePhoto(ctx *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(ctx)

	photo := ctx.MustGet("photo").(models.Photo)

	updatePhoto := models.Photo{}

	if contentType == appJSON {
		err := ctx.ShouldBindJSON(&updatePhoto)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
	} else {
		ctx.AbortWithStatus(http.StatusUnsupportedMediaType)
		return
	}

	photo.Title = updatePhoto.Title
	photo.Caption = updatePhoto.Caption
	photo.PhotoURL = updatePhoto.PhotoURL

	err := db.Debug().Save(&photo).Error

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         photo.ID,
		"title":      photo.Title,
		"caption":    photo.Caption,
		"photo_url":  photo.PhotoURL,
		"user_id":    photo.UserID,
		"updated_at": photo.UpdatedAt,
	})
}

func DeletePhoto(ctx *gin.Context) {
	db := database.GetDB()
	photo := ctx.MustGet("photo").(models.Photo)

	if err := db.Delete(&photo).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your photo has been successfully deleted",
	})
}
