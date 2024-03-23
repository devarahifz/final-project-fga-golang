package controllers

import (
	"final-project/database"
	"final-project/helpers"
	"final-project/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	appJSON = "application/json"
)

func UserRegister(ctx *gin.Context) {
	var newUser models.User

	db := database.GetDB()

	contentType := helpers.GetContentType(ctx)
	_, _ = db, contentType

	if contentType == appJSON {
		err := ctx.ShouldBindJSON(&newUser)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
	} else {
		ctx.AbortWithStatus(http.StatusUnsupportedMediaType)
		return
	}

	err := db.Debug().Create(&newUser).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":       newUser.ID,
		"username": newUser.Username,
		"email":    newUser.Email,
		"age":      newUser.Age,
	})
}

func UserLogin(ctx *gin.Context) {
	var user models.User

	db := database.GetDB()

	contentType := helpers.GetContentType(ctx)
	_, _ = db, contentType

	if contentType == appJSON {
		err := ctx.ShouldBindJSON(&user)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
	} else {
		ctx.AbortWithStatus(http.StatusUnsupportedMediaType)
		return
	}

	password := user.Password

	err := db.Debug().Where("email = ?", user.Email).Take(&user).Error

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "invalid email or password",
		})
		return
	}

	comparePass := helpers.ComparePass([]byte(user.Password), []byte(password))

	if !comparePass {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid password",
		})
		return
	}

	token := helpers.GenerateToken(int(user.ID), user.Email)

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func UserUpdate(ctx *gin.Context) {
	// Mendapatkan data pengguna dari konteks autentikasi
	userData, exists := ctx.Get("userData")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "User data not found in context",
		})
		return
	}

	claims, ok := userData.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to parse user data",
		})
		return
	}

	userIDFloat, ok := claims["id"].(float64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Invalid user ID data type",
		})
		return
	}

	userID := uint(userIDFloat)
	// Mendapatkan data pengguna dari database
	var user models.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to retrieve user data",
		})
		return
	}

	// Mendapatkan data pembaruan dari body permintaan
	var updateUser models.User
	if err := ctx.ShouldBindJSON(&updateUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}

	// Memperbarui email dan username pengguna
	user.Email = updateUser.Email
	user.Username = updateUser.Username
	user.UpdatedAt = time.Now()

	// Simpan perubahan ke database
	if err := database.GetDB().Save(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to update user data",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"age":        user.Age,
		"updated_at": user.UpdatedAt,
	})
}

func UserDelete(ctx *gin.Context) {
	// Mendapatkan data pengguna dari konteks autentikasi
	userData, exists := ctx.Get("userData")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "User data not found in context",
		})
		return
	}

	claims, ok := userData.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to parse user data",
		})
		return
	}

	userIDFloat, ok := claims["id"].(float64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Invalid user ID data type",
		})
		return
	}

	userID := uint(userIDFloat)
	// Mendapatkan data pengguna dari database
	var user models.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to retrieve user data",
		})
		return
	}

	// Menghapus pengguna dari database
	if err := database.GetDB().Where("id = ?", userID).Delete(&models.User{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to delete user account",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your account has been successfully deleted",
	})
}
