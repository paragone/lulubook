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
		return nil, errors.New("系统暂未处理该类型的配置文件")
	}
}

func SpiderRun(c *gin.Context) spider_dto.SpiderResponse{
	var req spider_dto.SpiderRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.Logger.Fatalf("req error ", err.Error())
		utils.SendFailedResponse(c, utils.ErrorCodeInvalidRequest, utils.ErrorDescInvalidRequest)
	}

}