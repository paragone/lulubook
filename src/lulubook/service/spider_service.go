package service

import (
	"github.com/gin-gonic/gin"
	"lulubook/dto/spider_dto"
	"lulubook/modules/spider"
	"lulubook/utils"
)


func SpiderRun(c *gin.Context){
	var req spider_dto.SpiderRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.Logger.Println("req error " + err.Error())
		utils.SendFailedResponse(c, utils.ErrorCodeInvalidRequest, utils.ErrorDescInvalidRequest)
		return
	}
	if req.Action == "start" {
		spider,err := spider.CreateSpider(req.Name)
		if err != nil {
			utils.Logger.Println("start error" + err.Error())
			utils.SendFailedResponse(c, utils.ErrorCodeFailed, utils.ErrorDescFaild + err.Error())
			return
		}
		go spider.CrawlSite(req.Url)
	}
	utils.SendSuccessResponse(c)
	return
}