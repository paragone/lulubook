package main

import (
	"github.com/gin-gonic/gin"
	"lulubook/router"
)

func main() {
	//for test
	//db.DropDB()
	r := gin.Default()

	r = router.SetupRouter(r)

    r.Run(":80")

	/*
    sp ,err:= service.NewSpider("booktxt")
    if err != nil{
		utils.Logger.Println("error " + err.Error())
		return
	}
    sp.SpiderSite("http://www.booktxt.com")
	*/
}