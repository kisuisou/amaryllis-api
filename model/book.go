package model

type BookData struct {
	ISBN      string `gorm:"PrimaryKey"`
	Title     string
	Volume    string
	Creator   string
	NDC9      string
	NDC10     string
	NDLC      string
	Publisher string
	PubYear   int
}
