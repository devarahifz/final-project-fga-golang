package authorizations

import (
	"final-project/database"
	"final-project/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func PhotoAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userData, exists := ctx.Get("userData")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "User data not found in token",
			})
			return
		}

		userID := uint(userData.(jwt.MapClaims)["id"].(float64))

		photoIDStr := ctx.Param("photoId")
		photoID, err := strconv.ParseUint(photoIDStr, 10, 64)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Invalid photo ID",
			})
			return
		}

		var photo models.Photo
		db := database.GetDB()
		if err := db.First(&photo, photoID).Error; err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Photo not found",
			})
			return
		}

		if photo.UserID != userID {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "You are not allowed to update this photo",
			})
			return
		}

		ctx.Set("photo", photo)
		ctx.Next()
	}
}
