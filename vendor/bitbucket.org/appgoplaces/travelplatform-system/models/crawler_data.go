package models

import (
	"time"

	"github.com/go-log/log"
)

type CrawlerData struct {
	TableName     struct{}  `sql:"crawler.crawler_data"`
	CrawlerDataID int64     `sql:"crawler_data_id,pk"`
	BlogID        int64     `sql:"blog_id, fk"`
	Data          string    `sql:"data"`
	CreateDate    time.Time `sql:"create_date"`
	ModifyDate    time.Time `sql:"modify_date"`
	URL           string    `sql:"url"`
	Analyzed      bool      `sql:"analyzed"`
	LastCrawlDate time.Time `sql:"last_crawl_date"`
	WordFrequency string    `sql:"word_frequency"`
}

type CrawlerDataQuery struct {
	CrawlerDataID int64     `sql:"crawler_data_id,pk" json:"crawler_data_id"`
	BlogID        int64     `sql:"blog_id" json:"user_id"`
	Data          string    `sql:"data" json:"data"`
	CreateDate    time.Time `sql:"create_date" json:"create_date"`
	ModifyDate    time.Time `sql:"modify_date" json:"modify_date"`
	URL           string    `sql:"url" json:"url"`
	Analyzed      bool      `sql:"analyzed" json:"analyzed"`
	LastCrawlDate time.Time `sql:"last_crawl_date" json:"last_crawl_date"`
	WordFrequency string    `sql:"word_frequency" json:word_frequency"`
}

// type CrawlerData struct {
// 	BlogID        int64
// 	Data          string
// 	CreateDate    time.Time
// 	ModifyDate    time.Time
// 	URL           string
// 	Analyzed      bool
// 	LastCrawlDate time.Time
// }

// var (
// 	db = rdbms.Connect("crawler.models.crawler_data")
// )

func (db *Db) CreateCrawlerData(b *CrawlerData) {
	err := db.Client.Insert(b)

	if err != nil {
		panic(err)
	}
}

func (db *Db) UpdateCrawlerData(b *CrawlerData) {
	_, err := db.Client.Model(b).Where("blog_id = ?", b.BlogID).UpdateNotNull()

	if err != nil {
		panic(err)
	}

}

func (db *Db) GetCrawlerDataByID(id int64) (CrawlerDataQuery, error) {
	var crawlerData CrawlerDataQuery
	log.Log(id)
	_, err := db.Client.Query(&crawlerData, `SELECT * FROM "crawler".crawler_data WHERE crawler_data_id = ?`, id)
	return crawlerData, err
}

func (db *Db) CheckCrawlerDataExists(id int64) bool {
	var crawlerData CrawlerDataQuery
	log.Log(id)
	data, err := db.Client.Query(&crawlerData, `SELECT * FROM "crawler".crawler_data WHERE blog_id = ?`, id)

	if err != nil {
		panic(err)
	}
	return data.RowsReturned() > 0
}
