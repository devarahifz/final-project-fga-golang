package authorizations

import (
	"final-project/database"
	"final-project/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func SocialMediaAuth() gin.HandlerFunc {
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

		socialMediaIDStr := ctx.Param("socialMediaId")
		socialMediaID, err := strconv.ParseUint(socialMediaIDStr, 10, 64)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Invalid social media ID",
			})
			return
		}

		var socialMedia models.SocialMedia
		db := database.GetDB()
		if err := db.First(&socialMedia, socialMediaID).Error; err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Social media not found",
			})
			return
		}

		if socialMedia.UserID != userID {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "You are not allowed to update this social media",
			})
			return
		}

		ctx.Set("socialMedia", socialMedia)
		ctx.Next()
	}
}
