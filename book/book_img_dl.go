package book

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const ndl_url = "https://ndlsearch.ndl.go.jp/thumbnail/"
const rakuten_url = "https://app.rakuten.co.jp/services/api/BooksBook/Search/20170404?format=json"

type rakuten_api_res struct {
	Items []struct {
		Item struct {
			LargeImageUrl string `json:"largeImageUrl"`
		}
	}
}

func GetBookImg(isbn string) bool {
	time.Sleep(1 * time.Second)
	if err := godotenv.Load(); err != nil {
		log.Fatal("loading .env failed")
	}
	res, err := http.Get(fmt.Sprintf("%s%s.jpg", ndl_url, isbn))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		app_id := os.Getenv("RAKUTEN_APP_ID")
		res, err = http.Get(fmt.Sprintf("%s&isbn=%s&applicationId=%s", rakuten_url, isbn, app_id))
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		res_text, _ := io.ReadAll(res.Body)
		res_json := new(rakuten_api_res)
		json.Unmarshal(res_text, res_json)
		var img_url string
		if len(res_json.Items) == 0 {
			return false
		}
		img_url = res_json.Items[0].Item.LargeImageUrl
		if img_url == "" {
			return false
		}
		res, err = http.Get(img_url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
	}
	file, _ := os.Create(fmt.Sprintf("./book_imgs/%s.jpg", isbn))
	defer file.Close()
	io.Copy(file, res.Body)
	return true
}
