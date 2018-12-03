package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ErrorCodeSuccess  = 0
	StatusReconnecting = 1
	StatusProcessing   = 2
	StatusNoNewVersion = 3

	ErrorCodeFailed         = 10000
	ErrorCodeInvalidParam   = 10001
	ErrorCodeInvalidRequest = 10002
	ErrorCodeInvalidFormat  = 10003
	ErrorCodeInvalidResponse = 10005
	ErrorCodeNotSupportYet = 10006
8

)


const (
	ErrorDescSuccess       = "Success"
	StatusDescReconnecting = "reconnecting"
	StatusDescProcessing   = "processing"

	ErrorDescFaild          = "failed: "
	ErrorDescInvalidParam   = "invalid param "
	ErrorDescInvalidRequest = "invalid request"
	ErrorDescInvalidFormat  = "invalid format"
	ErrorDescInvalidResponse = "invalid response"
	ErrorDescNotSupportYet = "do not support yet:"

)

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


func SendSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"error_code": ErrorCodeSuccess,
		"error_desc": ErrorDescSuccess},
	)
}



func SendFailedResponse(c *gin.Context, errcode int, errdesc string) {
	c.JSON(http.StatusOK, gin.H{
		"error_code": errcode,
		"error_desc": errdesc})
}