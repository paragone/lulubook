package spider

import (
	"github.com/PuerkitoBio/goquery"
	"lulubook/dto/spider_dto"
	"lulubook/modules/db"
	"lulubook/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BookTextSpider struct{
	crawledUrl map[string]bool
	bookUrl []string
}



func IsValidUrl(url string) bool{
	if strings.HasPrefix(url,"http://") && strings.Contains(url, "booktxt"){
		return true
	} else {
		return false
	}
}

func isBookHref(href string) bool{
	if(len(href) < 2){
		return false
	}
	trimHref := href[1:len(href)-1]
	if !strings.Contains(trimHref, "/"){
		firstchar := trimHref[0:1]
		if  firstchar >= "0" && firstchar <= "9" {
			return true
		}
	}
	return false
}

func isLibHref(href string) bool{
	if(len(href) < 2){
		return false
	}
	trimHref := href[1:len(href)-1]
	if !strings.Contains(trimHref, "/"){
		firstchar := trimHref[0:1]
		if  !(firstchar >= "0" && firstchar <= "9") {
			return true
		}
	}
	return false
}

func (spider *BookTextSpider)SpiderSite(url string) error {

	if spider.crawledUrl == nil{
		spider.crawledUrl = make(map[string]bool)
	}
	utils.Logger.Println("SpiderSite url:" + url)

	err := getAllBookUrl(spider, url)
	if err != nil{
		utils.Logger.Println("getAllBookUrl error" + err.Error())
		return err
	}

	utils.Logger.Println("got  books" + strconv.Itoa(len(spider.bookUrl)))

    for i,bookurl := range spider.bookUrl {
		SpiderBook(strconv.Itoa(i), bookurl)
	}

	return err
}

func getAllBookUrl(spider *BookTextSpider, url string) error{
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		utils.Logger.Println(err)
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		utils.Logger.Println("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		utils.Logger.Println(err)
		return err
	}

	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		href, exist := selection.Attr("href")
		if exist {
			if isLibHref(href){
				_,ok := spider.crawledUrl[href]
				if !ok {
					spider.crawledUrl[href] = true
					getAllBookUrl(spider,"http://www.booktxt.com" + href)
				}
			}
			if isBookHref(href){
				_,ok := spider.crawledUrl[href]
				if !ok{
					spider.crawledUrl[href] = true
					spider.bookUrl = append(spider.bookUrl,"http://www.booktxt.com" + href)
				}
			}
		}
	})
	return nil
}

func SpiderBook(id string,url string) error{
	utils.Logger.Println("SpiderBook url:" + url)
	querybook := spider_dto.SBook{}
	querybook.Id = id
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		utils.Logger.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		utils.Logger.Println("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		utils.Logger.Println(err)
	}
	bookname := utils.GbkToUtf8(doc.Find("#info h1").Text())
	querybook.Name = bookname
	querybook.Url = url

	doc.Find("#list dd").Each(func (i int, contentSelection *goquery.Selection){
		pre := i - 1
		next := i + 1
		chapterid := i
		title := utils.GbkToUtf8(contentSelection.Find("a").Text())
		href, _ := contentSelection.Find("a").Attr("href")
		url := "http://www.booktxt.com"+href
		chapter := spider_dto.SChapter{Id:strconv.Itoa(chapterid), BookId:id,Title:title,Url:url, Pre:pre, Next:next}
		querybook.Chapters = append(querybook.Chapters, chapter)
	})
	//for range创建副本
	/*
	for _,chap := range querybook.Chapters{
		SpiderChapter( &chap)
	}
	*/
	channel := make(chan struct{}, 100)
	for i:= 0; i< len(querybook.Chapters); i++{
		channel <- struct{}{}
		SpiderChapter(&querybook.Chapters[i], channel)
	}
	for i := 0; i < 100; i++{
		channel <- struct{}{}
	}
	close(channel)
	/*
    if book,err := db.ListBookByName(&querybook){

	}
	*/
	db.InsertBook(&querybook)
	return nil
}

type ChanTag struct{}

func SpiderChapter(chapter *spider_dto.SChapter, c chan struct{}){
	defer func(){<- c}()
	utils.Logger.Println("SpiderChapter bookid:" + chapter.BookId +" chaptername：" + chapter.Title)
	if  IsValidUrl(chapter.Url){
		// Request the HTML page.
		res, err := http.Get(chapter.Url)
		if err != nil {
			utils.Logger.Println(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			utils.Logger.Println("status code error: %d %s", res.StatusCode, res.Status)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			utils.Logger.Println(err)
		}
		content := doc.Find("#content").Text()
		content = utils.GbkToUtf8(content)
		content = strings.Replace(content, "聽", " ", -1)
		chapter.Content = content
		chapter.UpdatedAt = time.Now()
	}
}

