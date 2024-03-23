package routers

import (
	"final-project/controllers"
	"final-project/middleware"
	"final-project/middleware/authorizations"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	// StartServer starts the server
	router := gin.Default()

	user := router.Group("/users")
	{
		user.POST("/register", controllers.UserRegister)
		user.POST("/login", controllers.UserLogin)
		user.PUT("/:userId", middleware.Authentication(), controllers.UserUpdate)
		user.DELETE("", middleware.Authentication(), controllers.UserDelete)
	}

	socialmedias := router.Group("/socialmedias")
	{
		socialmedias.POST("", middleware.Authentication(), controllers.CreateSocialMedia)
		socialmedias.GET("", middleware.Authentication(), controllers.GetSocialMedia)
		socialmedias.PUT("/:socialMediaId", middleware.Authentication(), authorizations.SocialMediaAuth(), controllers.UpdateSocialMedia)
		socialmedias.DELETE("/:socialMediaId", middleware.Authentication(), authorizations.SocialMediaAuth(), controllers.DeleteSocialMedia)
	}

	photos := router.Group("/photos")
	{
		photos.POST("", middleware.Authentication(), controllers.CreatePhoto)
		photos.GET("", middleware.Authentication(), controllers.GetPhoto)
		photos.PUT("/:photoId", middleware.Authentication(), authorizations.PhotoAuth(), controllers.UpdatePhoto)
		photos.DELETE("/:photoId", middleware.Authentication(), authorizations.PhotoAuth(), controllers.DeletePhoto)
	}

	comments := router.Group("/comments")
	{
		comments.POST("", middleware.Authentication(), controllers.CreateComment)
		comments.GET("", middleware.Authentication(), controllers.GetComment)
		comments.PUT("/:commentId", middleware.Authentication(), authorizations.CommentAuth(), controllers.UpdateComment)
		comments.DELETE("/:commentId", middleware.Authentication(), authorizations.CommentAuth(), controllers.DeleteComment)
	}
	return router
}
