package utils

import (
	"log"
	"os"
)

var Logger *log.Logger

func init(){
	if CONFIG.LogPath != ""{
		file, err := os.Create(CONFIG.LogPath)
		if err != nil{
			panic(err.Error())
		}

		Logger = log.New(file, "", log.LstdFlags|log.Llongfile)
	} else {
		Logger = log.New(os.Stdout, "", log.LstdFlags)
	}

}