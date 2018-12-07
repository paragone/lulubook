package router

import (
	"github.com/gin-gonic/gin"
	"lulubook/service"
)

func SetupRouter(router *gin.Engine) *gin.Engine{
	router.Static("/web", "web")
	v1 := router.Group("/api/v1")
	dbGroup := v1.Group("db")
	{
		dbGroup.POST("/", service.DbHandler)
	}

	spiderGroup := v1.Group("/spider")
	{
		spiderGroup.POST("/", service.SpiderRun)
		spiderGroup.GET("/verify", service.SpiderVerify)
	}
	viewGroup := v1.Group("/view")
	{
		viewGroup.GET("/",service.ListAllBook)
		viewGroup.GET("/:bookid",service.ListBook)
		viewGroup.GET("/:bookid/:chapterid",service.ListChapter)
	}
	return router
}
