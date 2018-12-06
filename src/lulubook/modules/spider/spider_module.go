package spider

import (
	"errors"
	"lulubook/dto/spider_dto"
	"lulubook/modules/db"
	"lulubook/modules/parser"
	"lulubook/utils"
	"strconv"
	"time"
)

type Parser interface{
	ParseAllBookUrl(url string) ([]string,error)
	ParseBook(id string,url string) (*spider_dto.SBook, []spider_dto.SChapter, error)
	ParseChapter(chapter *spider_dto.SChapter)
}

type Spider struct{
	//bookUrl    []string
	books       []spider_dto.SBook
	parser     Parser
}

func CreateSpider(from string) (*Spider, error){

	switch from{
	case "booktxt":{
		spider := new(Spider)
		spider.parser = new(parser.BookTxtParser)
		return spider,nil
	}
	default:
		return nil, errors.New("暂不支持该种爬虫")
	}
}


func (spider *Spider)CrawlSite(url string) error {
	books, err := spider.parser.ParseAllBookUrl(url)
	if err != nil{
		utils.Logger.Println("ParseAllBookUrl error" + err.Error())
		return err
	}

	utils.Logger.Println("got  books" + strconv.Itoa(len(books)))
	var req spider_dto.SListCommon
	req.Offset = 0
	req.Limited = 0
	booklist, err := db.ListAllBook(&req)
	if err != nil {
		utils.Logger.Println("this is  no  book in db")
		spider.createSite(books)
	} else {
		spider.updateSite(booklist, books)
	}


	return err
}
func (spider *Spider)createSite(books []string) {
	for i, bookurl := range books {
		spider.CrawlBook(strconv.Itoa(i), bookurl)
	}
}
func (spider *Spider)updateSite(booklist []spider_dto.SBook,books []string) {
	has := make([]string, 0)
	for _, bookurl := range booklist {
		has = append(has, bookurl.Url)
	}
	//将未包含的加入末尾
	for _, todo := range books{
		contain := false
		for _, bookurl := range has{
			if todo == bookurl{
				contain = true
			}
		}
		if !contain {
			has = append(has, todo)
		}
	}
	for i, bookurl := range has {
		spider.CrawlBook(strconv.Itoa(i), bookurl)
	}
}
func (spider *Spider)CrawlBook(id string, url string)  {
	book,chapters,err := spider.parser.ParseBook(id, url)
	if err == nil{
		querybook, err := db.ListBookById(book)
		if err != nil{
			spider.createBook(book, chapters)
		} else {
			spider.updateBook(querybook, chapters)
		}

	}
}

func (spider *Spider)createBook(book *spider_dto.SBook,chapters []spider_dto.SChapter)  {
	channel := make(chan struct{}, 100)
	for i:= 0; i< len(chapters); i++{
		channel <- struct{}{}
		go spider.CrawlChapter(&chapters[i], channel)
	}
	for i := 0; i < 100; i++{
		channel <- struct{}{}
	}
	close(channel)
	db.InsertBook(book)
}

func (spider *Spider)updateBook(querybook *spider_dto.SBook, chapters []spider_dto.SChapter)  {
	//0123
	//0123456789
	channel := make(chan struct{}, 100)
	for i:= querybook.ChapterNum; i< len(chapters); i++{
		channel <- struct{}{}
		go spider.CrawlChapter(&chapters[i], channel)
	}
	for i := 0; i < 100; i++{
		channel <- struct{}{}
	}
	close(channel)
	querybook.ChapterNum = len(chapters)
	querybook.UpdatedAt  = time.Now()
	db.UpdateBook(querybook)
}

func (spider *Spider)CrawlChapter(chapter *spider_dto.SChapter, c chan struct{})  {
	defer func(){<- c}()
	spider.parser.ParseChapter(chapter)
	db.InsertChapter(chapter)
}


