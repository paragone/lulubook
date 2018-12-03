package main

import (
	"lulubook/modules/db"
	"lulubook/service"
	"lulubook/utils"
)

func main() {
	//r := gin.Default()
	//for test
	db.DropDB()
    sp ,err:= service.NewSpider("booktxt")
    if err != nil{
		utils.Logger.Println("error " + err.Error())
		return
	}
    sp.SpiderSite("http://www.booktxt.com")
}