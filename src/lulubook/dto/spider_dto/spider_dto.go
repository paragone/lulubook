package spider_dto


type SBook struct{
	Name string      	`json:"name" bson:"name"`
	Image string 		`json:"image" bson:"image"`
	Url string 			`json:"url" bson:"url"`
	Chapters []SChapter `json:"chapters" bson:"chapters"`
}

type SChapter struct{
    Title string 	`json:"title" bson:"title"`
	Url string 		`json:"url" bson:"url"`
	Order int 		`json:"order" bson:"order"`
	Pre int 		`json:"pre" bson:"pre"`
	Next int 		`json:"next" bson:"next"`
	Content string 	`json:"content" bson:"content"`
}
