package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

// Corp represents corperation info
type Corp struct {
	Name            string
	Code            string
	Category        string
	Products        string
	ListingDate     string
	SettlementMonth string
	Representive    string
	Homepage        string
	Region          string
	UpdatedAt       time.Time
}

const (
	corpListURL = "https://kind.krx.co.kr/corpgeneral/corpList.do?method=download"
	corpJSON    = "./corpList.json"
)

var corpMap map[string]*Corp

func dailyUpdateCorps() {
	cron, err := newDailyCron(8, 0, 0, "Asia/Seoul")
	if err != nil {
		log.Println(err)
		return
	}

	for {
		<-cron.t.C
		log.Println(time.Now(), "- cron tick")
		updateCorpList()
		loadCorpMap()
	}
}

func transformToEucKr(b []byte) string {
	var bufs bytes.Buffer
	wr := transform.NewWriter(&bufs, korean.EUCKR.NewDecoder())
	defer wr.Close()

	wr.Write(b)
	return bufs.String()
}

func updateCorpList() {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(corpListURL)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		log.Println(err)
		return
	}

	now := time.Now().In(loc)
	corps := make(map[string]*Corp)

	tokenizer := html.NewTokenizer(resp.Body)
	tokenType := tokenizer.Next()
	// skip first <tr>
	for ; tokenType != html.ErrorToken; tokenType = tokenizer.Next() {
		if tokenType == html.StartTagToken && "tr" == tokenizer.Token().Data {
			break
		}
	}
	for tokenType = tokenizer.Next(); tokenType != html.ErrorToken; tokenType = tokenizer.Next() {
		if tokenType == html.StartTagToken && "tr" == tokenizer.Token().Data {
			parsed := make([]string, 0, 9)
			for {
				tokenType = tokenizer.Next()
				if tokenType == html.EndTagToken && "tr" == tokenizer.Token().Data {
					break
				}
				if tokenType == html.StartTagToken && "td" == tokenizer.Token().Data {
					tokenType = tokenizer.Next()
					if tokenType == html.TextToken {
						text := transformToEucKr(tokenizer.Text())
						text = strings.Trim(text, " \t\n")
						parsed = append(parsed, text)
					} else if tokenType == html.EndTagToken && "td" == tokenizer.Token().Data {
						parsed = append(parsed, "")
					}
				}
			}
			corp := new(Corp)
			corp.Name = parsed[0]
			corp.Code = parsed[1]
			corp.Category = parsed[2]
			corp.Products = parsed[3]
			corp.ListingDate = parsed[4]
			corp.SettlementMonth = parsed[5]
			corp.Representive = parsed[6]
			corp.Homepage = parsed[7]
			corp.Region = parsed[8]
			corp.UpdatedAt = now

			corps[corp.Name] = corp
		}
	}

	if fileExists(corpJSON) {
		os.Remove(corpJSON)
	}
	f, err := os.Create(corpJSON)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", " ")

	enc.Encode(corps)
}

func loadCorpMap() error {
	if !fileExists(corpJSON) {
		return errors.New("File not exists")
	}
	f, err := os.Open(corpJSON)
	if err != nil {
		return err
	}
	defer f.Close()
	byteValue, _ := ioutil.ReadAll(f)

	corpMap = nil
	corpMap = make(map[string]*Corp)
	err = json.Unmarshal(byteValue, &corpMap)

	return err
}
