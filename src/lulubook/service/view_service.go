package service

import (
	"github.com/gin-gonic/gin"
	"lulubook/dto/spider_dto"
	"lulubook/modules/db"
	"lulubook/utils"
	"net/http"
)
func ListAllBook(c *gin.Context) {
	var req spider_dto.SListCommon
	err := c.ShouldBind(&req)
	if err != nil {
		utils.Logger.Println("req error " + err.Error())
		utils.SendFailedResponse(c, utils.ErrorCodeInvalidRequest, utils.ErrorDescInvalidRequest)
		return
	}
	res, err := db.ListAllBook(&req)
	if err != nil {
		utils.Logger.Println("req error " + err.Error())
		utils.SendFailedResponse(c, utils.ErrorCodeInvalidResponse, utils.ErrorDescInvalidResponse)
		return
	}
	c.JSON(http.StatusOK, &res)
	return
}

func ListBook(c *gin.Context) {
	var req spider_dto.SListCommon
	req.Id = c.Param("bookid")
	err := c.ShouldBind(&req)
	if err != nil {
		utils.Logger.Println("req error " + err.Error())
		utils.SendFailedResponse(c, utils.ErrorCodeInvalidRequest, utils.ErrorDescInvalidRequest)
		return
	}
	res, err := db.ListBookById(&req)
	if err != nil {
		utils.Logger.Println("req error " + err.Error())
		utils.SendFailedResponse(c, utils.ErrorCodeInvalidResponse, utils.ErrorDescInvalidResponse)
		return
	}
	c.JSON(http.StatusOK, &res)
	return
}
/*
func ListChapter(c *gin.Context){
	var req spider_dto.SListCommon
	req.Id = c.Param("bookid")
	req.Chapterid = c.Param("chapterid")
	err := c.ShouldBind(&req)
	if err != nil {
		utils.Logger.Println("req error " + err.Error())
		utils.SendFailedResponse(c, utils.ErrorCodeInvalidRequest, utils.ErrorDescInvalidRequest)
		return
	}

	book, err := db.ListBookById(&req)
	if err != nil {
		utils.Logger.Println("req error " + err.Error())
		utils.SendFailedResponse(c, utils.ErrorCodeInvalidResponse, utils.ErrorDescInvalidResponse)
		return
	}

	res, err := db.ListChapterById()
	c.JSON(http.StatusOK, &res)
	return
}
*/