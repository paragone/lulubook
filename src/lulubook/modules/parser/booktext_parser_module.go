package parser

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"lulubook/dto/spider_dto"
	"lulubook/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BookTxtParser struct{
	crawledUrl map[string]bool
	books     []string
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

func (parser *BookTxtParser) ParseAllBookUrl(url string) ([]string, error){
	if parser.crawledUrl == nil{
		parser.crawledUrl = make(map[string]bool)
	}
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		utils.Logger.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		utils.Logger.Println("status code error: %d %s", res.StatusCode, res.Status)
		return nil,errors.New("status code error")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		utils.Logger.Println(err)
		return nil,err
	}
	utils.Logger.Println("SpiderSite url:" + url)
	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		href, exist := selection.Attr("href")
		if exist {
			if isLibHref(href){
				_,ok := parser.crawledUrl[href]
				if !ok {
					parser.crawledUrl[href] = true
					parser.ParseAllBookUrl("http://www.booktxt.com" + href)
				}
			}
			if isBookHref(href){
				_,ok := parser.crawledUrl[href]
				if !ok{
					parser.crawledUrl[href] = true
					parser.books = append(parser.books,"http://www.booktxt.com" + href)
					//utils.Logger.Println("append url to books : http://www.booktxt.com" + href)
				}
			}
		}
	})
	return parser.books,nil
}

func (parser *BookTxtParser)ParseBook(id string,url string) (*spider_dto.SBook, []spider_dto.SChapter, error){

	querybook := spider_dto.SBook{}
	querybook.Id = id

	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		utils.Logger.Println(err)
		return nil,nil,err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		utils.Logger.Println("status code error: %d %s", res.StatusCode, res.Status)
		return nil,nil,errors.New("status code error")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		utils.Logger.Println(err)
		return nil,nil,err
	}

	bookname := utils.GbkToUtf8(doc.Find("#info h1").Text())
	bookimg ,existed:= doc.Find("#fmimg").Attr("src")
	if existed {
		querybook.Image = bookimg
	}
	querybook.Name = bookname
	querybook.Url = url
	querybook.CreatedAt = time.Now()
	querybook.UpdatedAt = time.Now()


	var chapters []spider_dto.SChapter
	doc.Find("#list dd").Each(func (i int, contentSelection *goquery.Selection){
		pre := i - 1
		next := i + 1
		chapterid := i
		title := utils.GbkToUtf8(contentSelection.Find("a").Text())
		href, _ := contentSelection.Find("a").Attr("href")
		url := "http://www.booktxt.com"+href
		chapter := spider_dto.SChapter{Id:strconv.Itoa(chapterid), BookId:id,Title:title,Url:url, Pre:pre, Next:next}
		chapters = append(chapters, chapter)
	})
	querybook.ChapterNum = len(chapters)


	return &querybook, chapters, nil
}


func (parser *BookTxtParser)ParseChapter(chapter *spider_dto.SChapter){
	utils.Logger.Println("SpiderChapter bookid:" + chapter.BookId +" chaptername：" + chapter.Title)

	// Request the HTML page.
	res, err := http.Get(chapter.Url)
	if err != nil {
		utils.Logger.Println(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		utils.Logger.Println("status code error: %d %s", res.StatusCode, res.Status)
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		utils.Logger.Println(err)
		return
	}
	content := doc.Find("#content").Text()
	content = utils.GbkToUtf8(content)
	content = strings.Replace(content, "聽", " ", -1)
	chapter.Content = content
	chapter.CreatedAt = time.Now()
	chapter.UpdatedAt = time.Now()
}

