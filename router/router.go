package router

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func SetUpRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "ok")
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	v1 := r.Group("api/v1")

	v1.POST("/sign", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)

	v1.Use(middleware.JWTAuthMiddleware(), middleware.RateLimitMiddleware(2*time.Second, 1)) //增加限流策略

	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)
		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		v1.GET("/posts", controller.GetPostListDetailHandler)
		v1.POST("/vote", controller.PostVoteController)
		v1.GET("/posts2", controller.GetPostListDetailHandlerV2)

	}
	//v1.GET("/ping", middleware.JWTAuthMiddleware(), func(context *gin.Context) {
	//	context.JSON(http.StatusOK, "pong")
	//})

	r.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusNotFound, gin.H{
			"msg": "404",
		})
	})
	return r
}
