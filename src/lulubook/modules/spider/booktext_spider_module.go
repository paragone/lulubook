package spider

import (
	"github.com/PuerkitoBio/goquery"
	"lulubook/dto/spider_dto"
	"lulubook/modules/db"
	"lulubook/utils"
	"net/http"
	"strings"
	"time"
)

type BookTextSpider struct{
	crawledUrl map[string]bool
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
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		utils.Logger.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		utils.Logger.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		utils.Logger.Fatal(err)
	}

	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		href, exist := selection.Attr("href")
		if exist {
			if isLibHref(href){
				_,ok := spider.crawledUrl[href]
				if !ok {
					spider.crawledUrl[href] = true
					spider.SpiderSite("http://www.booktxt.com" + href)
				}
			}
			if isBookHref(href){
				_,ok := spider.crawledUrl[href]
				if !ok{
					spider.crawledUrl[href] = true
					SpiderBook("http://www.booktxt.com"+href)
				}
			}
		}
	})
	return err
}

func SpiderBook(url string) error{
	utils.Logger.Println("SpiderBook url:" + url)
	querybook := spider_dto.SBook{}
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		utils.Logger.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		utils.Logger.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		utils.Logger.Fatal(err)
	}
	bookname := utils.GbkToUtf8(doc.Find("#info h1").Text())
	querybook.Name = bookname


	doc.Find("#list dd").Each(func (i int, contentSelection *goquery.Selection){
		pre := i - 1
		next := i + 1
		title := utils.GbkToUtf8(contentSelection.Find("a").Text())
		href, _ := contentSelection.Find("a").Attr("href")
		url := []string{href}
		chapter := spider_dto.SChapter{Title:title,Url:url, Order:i , Pre:pre, Next:next}
		querybook.Chapters = append(querybook.Chapters, chapter)
	})

	book, err := db.ListBookByName(&querybook)
	if err != nil{
		db.InsertBook(&querybook)
	} else {
		if book != nil {
			contain := false
			for _, u := range book.Url {
				if u == url {
					contain = true
				}
			}
			if !contain {
				book.Url = append(book.Url, url)
			}
			book.Chapters = querybook.Chapters
			db.UpdateBook(book)
			querybook = *book
		} else {
			utils.Logger.Fatalf("list error ", err.Error())
			return err
		}
	}

	channel := make(chan struct{}, 100)
	for _,chapter := range querybook.Chapters {
		channel <- struct{}{}
		go SpiderChapter(querybook.Name, &chapter, channel)
	}
	for i := 0; i < 100; i++{
		channel <- struct{}{}
	}
	close(channel)
	return nil
}

type ChanTag struct{}

func SpiderChapter(bookname string, chapter *spider_dto.SChapter, c chan struct{}){
	utils.Logger.Println("SpiderChapter bookname:" + bookname +" chaptername：" + chapter.Title)
	defer func(){<- c}()
	if  IsValidUrl("http://www.booktxt.com"+chapter.Url[len(chapter.Url) - 1]){
		// Request the HTML page.
		res, err := http.Get("http://www.booktxt.com"+chapter.Url[len(chapter.Url) - 1])
		if err != nil {
			utils.Logger.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			utils.Logger.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			utils.Logger.Fatal(err)
		}
		content := doc.Find("#content").Text()
		content = utils.GbkToUtf8(content)
		content = strings.Replace(content, "聽", " ", -1)
		ch := spider_dto.SChapter{BookName:bookname, Title:chapter.Title, Content:content,Order:chapter.Order, Pre:chapter.Pre, Next:chapter.Next, CreatedAt:time.Now(),UpdatedAt:time.Now()}
		querych,err := db.ListChapterByTitle(&ch)
		if err == nil{
			if querych != nil {
				querych.Content = ch.Content
				querych.UpdatedAt = time.Now()
				db.UpdateChapter(querych)
			} else {
				utils.Logger.Fatalf("list error ", err.Error())
				return
			}
		} else {
			db.InsertChapter(&ch)
		}
	}
}

