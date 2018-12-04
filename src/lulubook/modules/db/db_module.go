package db

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"lulubook/dto/spider_dto"
	"lulubook/utils"
)
const(
	DBNAME = "lulubook"
	BOOKCOLLECTION = "books"
)

func DropDB() error{
	err := utils.DropMongoDB(DBNAME)
    if err != nil{
    	utils.Logger.Println("drop fail" + err.Error())
	}
	return err
}
func InsertBook(book *spider_dto.SBook) error{
	err := utils.WithMongoDBCollection(DBNAME, BOOKCOLLECTION, func(c *mgo.Collection) error{
		err := c.Insert(book)
		
		if err != nil {
			utils.Logger.Println("InsertBook error" + err.Error())
		}
		return err
	})

	return err
}

func UpdateBook(book *spider_dto.SBook) error {
	err := utils.WithMongoDBCollection(DBNAME, BOOKCOLLECTION, func(c *mgo.Collection) error{
		m := bson.M{"_id":book.Id}
		err := c.Update(m, book)
		if err != nil {
			utils.Logger.Println("updateBook error" + err.Error())
		}
		return err
	})

	return err
}

func DeleteBook(book *spider_dto.SBook) error {
	err := utils.WithMongoDBCollection(DBNAME, BOOKCOLLECTION, func(c *mgo.Collection) error{
		m := bson.M{"_id":book.Id}
		err := c.Remove(m)
		if err != nil {
			utils.Logger.Println("updateBook error" + err.Error())
		}
		return err
	})

	return err
}

func ListBookByName(book *spider_dto.SBook) (*spider_dto.SBook, error) {
	var res spider_dto.SBook
	err := utils.WithMongoDBCollection(DBNAME, BOOKCOLLECTION, func(c *mgo.Collection) error{
		m := bson.M{"name":book.Name}
		err := c.Find(m).One(&res)
		if err != nil {
			utils.Logger.Println("ListBookByName error" + err.Error())
		}
		return err
	})
	return &res, err
}

func ListBookById(req *spider_dto.SListCommon) (*spider_dto.SBook, error) {
	var res spider_dto.SBook
	err := utils.WithMongoDBCollection(DBNAME, BOOKCOLLECTION, func(c *mgo.Collection) error{
		m := bson.M{"_id":req.Id}
		err := c.Find(m).One(&res)
		if err != nil {
			utils.Logger.Println("ListBookByName error" + err.Error())
		}
		return err
	})
	return &res, err
}

func ListAllBook(req *spider_dto.SListCommon) (*[]spider_dto.SBook,error) {
	var res []spider_dto.SBook
	err := utils.WithMongoDBCollection(DBNAME, BOOKCOLLECTION, func(c *mgo.Collection) error{
		err := c.Find(bson.M{}).Select(bson.M{"chapter": 0}).Sort("_id").Skip(req.Offset).Limit(req.Limited).All(&res)
		if err != nil{
			utils.Logger.Println("ListBook error" + err.Error())
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return &res, err
}
func ListChapterById(req *spider_dto.SListCommon) (*spider_dto.SBook, error) {
	var res spider_dto.SBook
	err := utils.WithMongoDBCollection(DBNAME, BOOKCOLLECTION, func(c *mgo.Collection) error{
		m := bson.M{"_id":req.Id, "chapter._id":req.ChapterId}
		err := c.Find(m).Select(bson.M{"chapter.$": 1}).One(&res)
		if err != nil {
			utils.Logger.Println("ListBookByName error" + err.Error())
		}
		return err
	})
	return &res, err
}
/*
func InsertChapter(chapter *spider_dto.SChapter) error{

	err := utils.WithMongoDBCollection(DBNAME, chapter.BookName, func(c *mgo.Collection) error{
		err := c.Insert(chapter)

		if err != nil {
			utils.Logger.Println("InsertChapter error" + err.Error())
		}
		return err
	})

	return err
}

func UpdateChapter(chapter *spider_dto.SChapter) error {
	err := utils.WithMongoDBCollection(DBNAME, chapter.BookName, func(c *mgo.Collection) error{
		m := bson.M{"_id":chapter.Id}
		err := c.Update(m, chapter)
		if err != nil {
			utils.Logger.Println("UpdateChapter error" + err.Error())
		}
		return err
	})

	return err
}

func DeleteChapter(chapter *spider_dto.SChapter) error {
	err := utils.WithMongoDBCollection(DBNAME, chapter.BookName, func(c *mgo.Collection) error{
		m := bson.M{"_id":chapter.Id}
		err := c.Remove(m)
		if err != nil {
			utils.Logger.Println("DeleteChapter error" + err.Error())
		}
		return err
	})

	return err
}

func ListChapterByTitle(chapter *spider_dto.SChapter) (*spider_dto.SChapter, error) {
	var res spider_dto.SChapter
	err := utils.WithMongoDBCollection(DBNAME, chapter.BookName, func(c *mgo.Collection) error{
		m := bson.M{"title":chapter.Title}
		err := c.Find(m).One(&res)
		if err != nil {
			utils.Logger.Println("ListChapterByTitle error" + err.Error())
		}
		return err
	})
	return &res, err
}

*/
