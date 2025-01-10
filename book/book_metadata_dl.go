package book

import (
	"amaryllis-api/model"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const api_url = "https://ndlsearch.ndl.go.jp/api/opensearch"

type APIBookData struct {
	XMLName xml.Name `xml:"rss"`
	Items   []struct {
		Title      string   `xml:"http://purl.org/dc/elements/1.1/ title"`
		Creator    []string `xml:"http://purl.org/dc/elements/1.1/ creator"`
		Identifier []struct {
			Type    string `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr"`
			Content string `xml:",chardata"`
		} `xml:"http://purl.org/dc/elements/1.1/ identifier"`
		Subject []struct {
			Type    string `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr"`
			Content string `xml:",chardata"`
		} `xml:"http://purl.org/dc/elements/1.1/ subject"`
		Volume    string   `xml:"http://ndl.go.jp/dcndl/terms/ volume"`
		PubYear   string   `xml:"http://purl.org/dc/elements/1.1/ date"`
		Publisher []string `xml:"http://purl.org/dc/elements/1.1/ publisher"`
	} `xml:"channel>item"`
}

func GetMetaData(isbn string) model.Book {
	time.Sleep(1 * time.Second)
	res, err := http.Get(fmt.Sprintf("%s?isbn=%s", api_url, isbn))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	res_byte, _ := io.ReadAll(res.Body)
	data := new(APIBookData)
	xml.Unmarshal(res_byte, data)
	var item_i int
	var ndc9 string
	var ndc10 string
	var ndlc string
	for i := 0; i < len(data.Items); i++ {
		for _, v := range data.Items[i].Identifier {
			if v.Type == "dcndl:ISBN" {
				if isbn == strings.Replace(v.Content, "-", "", -1) {
					item_i = i
					break
				}
			}
		}
		for _, v := range data.Items[i].Subject {
			if v.Type == "dcndl:NDC9" {
				ndc9 = v.Content
			} else if v.Type == "dcndl:NDC10" {
				ndc10 = v.Content
			} else if v.Type == "dcndl:NDLC" {
				ndlc = v.Content
			}
		}

	}
	book_data := new(model.Book)
	book_data.Title = data.Items[item_i].Title
	creator := data.Items[item_i].Creator[0]
	re := regexp.MustCompile("(, )|[0-9]|-|　| ")
	creator = re.ReplaceAllString(creator, "")
	book_data.Creator = creator
	book_data.Publisher = strings.Join(data.Items[item_i].Publisher, "")
	book_data.ISBN = isbn
	book_data.Volume = data.Items[item_i].Volume
	book_data.PubYear, err = strconv.Atoi(data.Items[item_i].PubYear)
	if err != nil {
		book_data.PubYear = 0
	}
	book_data.NDC9 = ndc9
	book_data.NDC10 = ndc10
	book_data.NDLC = ndlc
	return *book_data
}
