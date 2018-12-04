package service

import (
	"github.com/gin-gonic/gin"
	"lulubook/dto/spider_dto"
	"lulubook/utils"
	"lulubook/modules/db"
)

func DbHandler(c *gin.Context){
	var req spider_dto.DbRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.Logger.Println("req error " + err.Error())
		utils.SendFailedResponse(c, utils.ErrorCodeInvalidRequest, utils.ErrorDescInvalidRequest)
		return
	}
	if req.Action == "reset" {
		db.DropDB()
	}
	utils.SendSuccessResponse(c)
	return
}