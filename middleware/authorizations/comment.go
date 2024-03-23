package authorizations

import (
	"final-project/database"
	"final-project/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CommentAuth() gin.HandlerFunc {
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

		commentIDStr := ctx.Param("commentId")
		commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Invalid comment ID",
			})
			return
		}

		var comment models.Comment
		db := database.GetDB()
		if err := db.First(&comment, commentID).Error; err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Comment not found",
			})
			return
		}

		if comment.UserID != userID {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "You are not allowed to update this comment",
			})
			return
		}

		ctx.Set("comment", comment)
		ctx.Next()
	}
}
