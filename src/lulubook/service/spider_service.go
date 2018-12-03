package service

import (
	"errors"
	"lulubook/modules/spider"
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
