package model

type Book struct {
	ISBN                 string `gorm:"primaryKey"`
	Title                string
	TitleTranscription   string
	Volume               string
	Creator              string
	CreatorTranscription string
	NDC9                 string
	NDC10                string
	NDLC                 string
	Publisher            string
	PubYear              int
	Image                string
}
