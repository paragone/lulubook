package spider

import (
	"github.com/PuerkitoBio/goquery"
	"time"
	"strings"
	"lulubook/dto/spider_dto"
	"lulubook/utils"
	"net/http"
	"log"
)

type BookTextSpider struct{

}

func (self *BookTextSpider)SpiderUrl(url string)( error){
	book := spider_dto.SBook{}
	book.Url = url
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	bookname := utils.GbkToUtf8(doc.Find("#info h1").Text())

	b, err := models.GetBookByName(bookname)
	if err != nil{
		b := models.Book{Name:bookname, CreatedAt:time.Now(), UpdatedAt:time.Now()}
		models.BookAdd(&b)
	}
	doc.Find("#list dd").Each(func (i int, contentSelection *goquery.Selection){
		if i < 9{
			return
		}
		pre := i - 9
		next := i -7
		title := utils.GbkToUtf8(contentSelection.Find("a").Text())
		href, _ := contentSelection.Find("a").Attr("href")
		chapter := SChapter{Title:title,Url:"http://www.booktxt.net"+href, Order:i - 8, Pre:pre, Next:next}
		book.Chapters = append(book.Chapters, &chapter)
		u := models.Url{Url:chapter.Url}
		models.UrlAdd(&u)
	})



	channel := make(chan struct{}, 100)
	for _, chapter := range book.Chapters{
		channel <- struct{}{}
		go SpiderChapter(b.Id, chapter, channel)
	}

	for i := 0; i < 100; i++{
		channel <- struct{}{}
	}
	close(channel)
	return nil
}

type ChanTag struct{}

func SpiderChapter(bookid int, chapter *SChapter, c chan struct{}){
	defer func(){<- c}()
	if models.IsValidUrl(chapter.Url){
		doc, err := goquery.NewDocument(chapter.Url)
		if err != nil{
			ilog.AppLog.Error("get chapter details error: ", err.Error())
			return
		}
		content := doc.Find("#content").Text()
		content = utils.GbkToUtf8(content)
		content = strings.Replace(content, "聽", " ", -1)
		ch := models.Chapter{BookId:bookid, Title:chapter.Title, Content:content,Sort:chapter.Order, Pre:chapter.Pre, Next:chapter.Next, CreatedAt:time.Now(),UpdatedAt:time.Now()}
		models.ChapterAdd(&ch)
		models.SpideredUrl(chapter.Url)
	}
}

