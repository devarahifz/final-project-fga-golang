package controllers

import (
	"final-project/database"
	"final-project/helpers"
	"final-project/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreateSocialMedia(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(ctx)

	socialMedia := models.SocialMedia{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		err := ctx.ShouldBindJSON(&socialMedia)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
	} else {
		ctx.AbortWithStatus(http.StatusUnsupportedMediaType)
		return
	}

	socialMedia.UserID = userID

	err := db.Debug().Create(&socialMedia).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":               socialMedia.ID,
		"name":             socialMedia.Name,
		"social_media_url": socialMedia.SocialMediaURL,
		"user_id":          socialMedia.UserID,
		"created_at":       socialMedia.CreatedAt,
	})
}

func GetSocialMedia(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)

	userID := uint(userData["id"].(float64))

	var socialMedia []models.SocialMedia

	err := db.Debug().Where("user_id = ?", userID).Preload("User").Find(&socialMedia).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	mappedSocialMedia := make([]gin.H, len(socialMedia))
	for i, socialMedia := range socialMedia {
		mappedSocialMedia[i] = gin.H{
			"id":               socialMedia.ID,
			"name":             socialMedia.Name,
			"social_media_url": socialMedia.SocialMediaURL,
			"user_id":          socialMedia.UserID,
			"created_at":       socialMedia.CreatedAt,
			"updated_at":       socialMedia.UpdatedAt,
			"User": gin.H{
				"id":       socialMedia.User.ID,
				"username": socialMedia.User.Username,
				"email":    socialMedia.User.Email,
			},
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"social_media": mappedSocialMedia,
	})
}

func UpdateSocialMedia(ctx *gin.Context) {
	db := database.GetDB()
	socialMedia := ctx.MustGet("socialMedia").(models.SocialMedia)

	var updatedSocialMedia models.SocialMedia
	if err := ctx.ShouldBindJSON(&updatedSocialMedia); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	socialMedia.Name = updatedSocialMedia.Name
	socialMedia.SocialMediaURL = updatedSocialMedia.SocialMediaURL

	if err := db.Debug().Save(&socialMedia).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":               socialMedia.ID,
		"name":             socialMedia.Name,
		"social_media_url": socialMedia.SocialMediaURL,
		"user_id":          socialMedia.UserID,
		"updated_at":       socialMedia.UpdatedAt,
	})
}

func DeleteSocialMedia(ctx *gin.Context) {
	db := database.GetDB()
	socialMedia := ctx.MustGet("socialMedia").(models.SocialMedia)

	if err := db.Delete(&socialMedia).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your social media has been sucessfully deleted",
	})
}
