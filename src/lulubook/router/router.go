package router

import (
	"github.com/gin-gonic/gin"
	"lulubook/service"
)

func SetupRouter(router *gin.Engine) *gin.Engine{
	v1 := router.Group("/api/v1")

	spiderGroup := v1.Group("/spider")
	{
		spiderGroup.POST("/", service.SpiderRun)
	}
	viewGroup := v1.Group("/view")
	{
		viewGroup.GET("/",service.ListAllBook)
		viewGroup.GET("/:bookid",service.ListBook)
		//viewGroup.GET("/:bookid/:chapterid",service.ListChapter)
	}
	return router
}
