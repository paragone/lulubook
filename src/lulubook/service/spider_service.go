package service

import (
	"errors"
	"lulubook/modules/spider"
	"github.com/gin-gonic/gin"
	"lulubook/dto/spider_dto"
	"lulubook/utils"
)

type Spider interface{
	SpiderSite(url string) error
}

func NewSpider(from string) (Spider, error){
	switch from{
	case "booktxt":
		return new(spider.BookTextSpider), nil
	default:
		return nil, errors.New("暂不支持该种爬虫")
	}
}

func SpiderRun(c *gin.Context){
	var req spider_dto.SpiderRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.Logger.Fatalf("req error " + err.Error())
		utils.SendFailedResponse(c, utils.ErrorCodeInvalidRequest, utils.ErrorDescInvalidRequest)
		return
	}
	if req.Action == "start" {
		spider,err := NewSpider(req.Name)
		if err != nil {
			utils.Logger.Fatalf("start error" + err.Error())
			utils.SendFailedResponse(c, utils.ErrorCodeFailed, utils.ErrorDescFaild + err.Error())
			return
		}
		go spider.SpiderSite(req.Url)
	}
	utils.SendSuccessResponse(c)
	return
}