package book

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const book_img_url = "https://ndlsearch.ndl.go.jp/thumbnail/"

func GetBookImg(isbn string) {
	res, err := http.Get(fmt.Sprintf("%s%s.jpg", book_img_url, isbn))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	file, _ := os.Create(fmt.Sprintf("./book_imgs/%s.jpg", isbn))
	defer file.Close()
	io.Copy(file, res.Body)
}
