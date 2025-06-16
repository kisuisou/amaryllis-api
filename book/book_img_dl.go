package book

import (
	"amaryllis-api/data_store"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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
	if err := godotenv.Load(); err != nil {
		log.Fatal("loading .env failed")
	}
	ctx := context.Background()
	last_fetch_str, err := data_store.Rdb.Get(ctx, "ndl_img_last_fetch").Result()
	last_fetch, _ := strconv.Atoi(last_fetch_str)
	if err == nil {
		for {
			now_unix_time := time.Now().Unix()
			if now_unix_time-int64(last_fetch) > 1 {
				break
			}
		}
	}
	res, err := http.Get(fmt.Sprintf("%s%s.jpg", ndl_url, isbn))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		app_id := os.Getenv("RAKUTEN_APP_ID")
		last_fetch_str, err := data_store.Rdb.Get(ctx, "rakuten_last_fetch").Result()
		last_fetch, _ := strconv.Atoi(last_fetch_str)
		if err == nil {
			for {
				now_unix_time := time.Now().Unix()
				if now_unix_time-int64(last_fetch) > 1 {
					break
				}
			}
		}
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
		err = data_store.Rdb.Set(ctx, "rakuten_last_fetch", strconv.Itoa(int(time.Now().Unix())), 0).Err()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = data_store.Rdb.Set(ctx, "ndl_img_last_fetch", strconv.Itoa(int(time.Now().Unix())), 0).Err()
		if err != nil {
			log.Fatal(err)
		}
	}
	file, _ := os.Create(fmt.Sprintf("./book_imgs/%s.jpg", isbn))
	defer file.Close()
	io.Copy(file, res.Body)
	return true
}
