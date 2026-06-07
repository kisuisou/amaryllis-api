package model

type Book struct {
	ID                   uint `gorm:"primaryKey;autoIncrement"`
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
	MetaDataStatus       string
	ImageStatus          string
}

type BookIdentifier struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	BookID uint   `gorm:"index;not null"`
	Type   string `gorm:"uniqueIndex:idx_book_identifiers_type_value;not null"`
	Value  string `gorm:"uniqueIndex:idx_book_identifiers_type_value;not null"`
}
