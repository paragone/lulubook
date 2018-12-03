package spider_dto

import "time"

type SLibrary struct{
	Name string `json:"name" bson:"name"`
}

type SBook struct{
	Name  string      	`json:"name" bson:"name"`
	Image string 		`json:"image" bson:"image"`
	Url   []string 	    `json:"url" bson:"url"`
	Chapters []SChapter `json:"chapter" bson:"chapter"`
}

type SChapter struct{
	BookName string     `json:"bookname" bson:"bookname"`
    Title    string 	`json:"title" bson:"title"`
	Url     []string 	`json:"url" bson:"url"`
	Order   int 		`json:"order" bson:"order"`
	Pre     int 		`json:"pre" bson:"pre"`
	Next    int 		`json:"next" bson:"next"`
	Content string 		`json:"content" bson:"content"`
	CreatedAt time.Time `json:"createdat" bson:"createdat"`
	UpdatedAt time.Time `json:"updatedat" bson:"updatedat"`
}

type SListCommon struct{
	CollectionName 	  string
	Offset    int            `form:"offset,default=0" json:"offset,default=0"`
	Limited   int            `form:"limited,default=0" json:"limited,default=20"`
	Order     string         `form:"order,default=asc" json:"order,default=asc"`
}

type SpiderRequest struct{
	Action     string               `json:"action"`
	Name       string               `json:"name"`
	Url        string               `json:"url"`
}

type SpiderResponse struct{
	ErrorCode int             `json:"error_code"`
	ErrorDesc string          `json:"error_desc"`
}