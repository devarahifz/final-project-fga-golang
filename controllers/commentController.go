package controllers

import (
	"final-project/database"
	"final-project/helpers"
	"final-project/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreateComment(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(ctx)

	comment := models.Comment{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		err := ctx.ShouldBindJSON(&comment)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
	} else {
		ctx.AbortWithStatus(http.StatusUnsupportedMediaType)
		return
	}

	comment.UserID = userID

	err := db.Debug().Create(&comment).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         comment.ID,
		"message":    comment.Message,
		"photo_id":   comment.PhotoID,
		"user_id":    comment.UserID,
		"created_at": comment.CreatedAt,
	})
}

func GetComment(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)

	userID := uint(userData["id"].(float64))

	var comments []models.Comment

	if err := db.Where("user_id = ?", userID).Preload("User").Preload("Photo").Find(&comments).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}
	mappedComments := make([]gin.H, len(comments))
	for i, comment := range comments {
		mappedComments[i] = gin.H{
			"id":         comment.ID,
			"message":    comment.Message,
			"photo_id":   comment.PhotoID,
			"user_id":    comment.UserID,
			"updated_at": comment.UpdatedAt,
			"created_at": comment.CreatedAt,
			"User": gin.H{
				"id":       comment.User.ID,
				"username": comment.User.Username,
				"email":    comment.User.Email,
			},
			"Photo": gin.H{
				"id":        comment.Photo.ID,
				"title":     comment.Photo.Title,
				"caption":   comment.Photo.Caption,
				"photo_url": comment.Photo.PhotoURL,
				"user_id":   comment.Photo.UserID,
			},
		}
	}
	ctx.JSON(http.StatusOK, mappedComments)
}

func UpdateComment(ctx *gin.Context) {
	db := database.GetDB()
	comment := ctx.MustGet("comment").(models.Comment)

	var updatedComment models.Comment
	if err := ctx.ShouldBindJSON(&updatedComment); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	comment.Message = updatedComment.Message

	if err := db.Debug().Save(&comment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         comment.ID,
		"message":    comment.Message,
		"photo_id":   comment.PhotoID,
		"user_id":    comment.UserID,
		"created_at": comment.CreatedAt,
	})
}

func DeleteComment(ctx *gin.Context) {
	db := database.GetDB()
	comment := ctx.MustGet("comment").(models.Comment)

	if err := db.Debug().Delete(&comment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your comment has been successfully deleted",
	})
}
