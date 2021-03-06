package spider_dto

import "time"

type SLibrary struct{
	Name string `json:"name" bson:"name"`
}

type SBook struct{
	Id    string  			`json:"_id" bson:"_id"`
	Name  string      	`json:"name" bson:"name"`
	Image string 		`json:"image" bson:"image"`
	Url   string 	    `json:"url" bson:"url"`
	ChapterNum int      `json:"chapternum" bson:"chapternum"`
	CreatedAt time.Time `json:"createdat" bson:"createdat"`
	UpdatedAt time.Time `json:"updatedat" bson:"updatedat"`
}

type SChapter struct{
	Id       string		`json:"_id" bson:"_id"`
	BookId   string     `json:"bookid" bson:"bookid"`
    Title    string 	`json:"title" bson:"title"`
	Url     string 		`json:"url" bson:"url"`
	Pre     int 		`json:"pre" bson:"pre"`
	Next    int 		`json:"next" bson:"next"`
	Content string 		`json:"content" bson:"content"`
	CreatedAt time.Time `json:"createdat" bson:"createdat"`
	UpdatedAt time.Time `json:"updatedat" bson:"updatedat"`
}

type SListCommon struct{
	Id 	  string
	ChapterId 	  string
	Offset    int            `form:"offset,default=0" json:"offset,default=0"`
	Limited   int            `form:"limited,default=20" json:"limited,default=20"`
	Order     string         `form:"order,default=asc" json:"order,default=asc"`
}

type SpiderRequest struct{
	Action     string               `json:"action"`
	Name       string               `json:"name"`
	Url        string               `json:"url"`
}


type DbRequest struct{
	Action     string               `json:"action"`
}