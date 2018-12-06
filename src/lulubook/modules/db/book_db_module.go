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

func ListBookById(book *spider_dto.SBook) (*spider_dto.SBook, error) {
	var res spider_dto.SBook
	err := utils.WithMongoDBCollection(DBNAME, BOOKCOLLECTION, func(c *mgo.Collection) error{
		m := bson.M{"_id":book.Id}
		err := c.Find(m).One(&res)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &res, err
}



func ListAllBook(req *spider_dto.SListCommon) ([]spider_dto.SBook,error) {
	var res []spider_dto.SBook
	err := utils.WithMongoDBCollection(DBNAME, BOOKCOLLECTION, func(c *mgo.Collection) error{
		collation := &mgo.Collation{Locale: "zh",NumericOrdering:true}
		err := c.Find(bson.M{}).Collation(collation).Sort("_id").Skip(req.Offset).Limit(req.Limited).All(&res)
		return err
	})
	if err != nil {
		return nil, err
	}
	return res, err
}



