package db

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"lulubook/dto/spider_dto"
	"lulubook/utils"
)

func ListChapterById(req *spider_dto.SListCommon) (spider_dto.SChapter, error) {
	var res spider_dto.SChapter
	err := utils.WithMongoDBCollection(DBNAME, req.Id, func(c *mgo.Collection) error{
		m := bson.M{"_id":req.ChapterId}
		err := c.Find(m).One(&res)
		return err
	})
	return res, err
}
func ListBookChaptersById(req *spider_dto.SListCommon) ([]spider_dto.SChapter, error) {
	var res []spider_dto.SChapter
	err := utils.WithMongoDBCollection(DBNAME, req.Id, func(c *mgo.Collection) error{
		collation := &mgo.Collation{Locale: "zh",NumericOrdering:true}
		err := c.Find(bson.M{}).Collation(collation).Select(bson.M{"content":0}).Sort("_id").Skip(req.Offset).Limit(req.Limited).All(&res)
		return err
	})
	return res, err
}

func InsertChapter(chapter *spider_dto.SChapter) error{

	err := utils.WithMongoDBCollection(DBNAME, chapter.BookId, func(c *mgo.Collection) error{
		err := c.Insert(chapter)

		if err != nil {
			utils.Logger.Println("InsertChapter error" + err.Error())
		}
		return err
	})

	return err
}

func UpdateChapter(chapter *spider_dto.SChapter) error {
	err := utils.WithMongoDBCollection(DBNAME, chapter.BookId, func(c *mgo.Collection) error{
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
	err := utils.WithMongoDBCollection(DBNAME, chapter.BookId, func(c *mgo.Collection) error{
		m := bson.M{"_id":chapter.Id}
		err := c.Remove(m)
		if err != nil {
			utils.Logger.Println("DeleteChapter error" + err.Error())
		}
		return err
	})

	return err
}

func ListChapterByTitle(chapter *spider_dto.SChapter) (spider_dto.SChapter, error) {
	var res spider_dto.SChapter
	err := utils.WithMongoDBCollection(DBNAME, chapter.BookId, func(c *mgo.Collection) error{
		m := bson.M{"title":chapter.Title}
		err := c.Find(m).One(&res)
		return err
	})
	return res, err
}

