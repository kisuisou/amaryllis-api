package model

type Book struct {
	ISBN      string `gorm:"primaryKey"`
	Title     string
	Volume    string
	Creator   string
	NDC9      string
	NDC10     string
	NDLC      string
	Publisher string
	PubYear   int
}
