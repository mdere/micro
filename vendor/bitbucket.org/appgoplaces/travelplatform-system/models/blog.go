package models

import (
	// "net/url"
	// "fmt"
	// orm "github.com/go-pg/pg/orm"
	"time"

	pg "github.com/go-pg/pg"

	blogdatapb "bitbucket.org/appgoplaces/service-protos/blog-data"
	"github.com/go-log/log"
	merr "github.com/micro/go-micro/errors"
	// "database/sql"
	// "gopkg.in/mgutz/dat.v1"
	// "gopkg.in/mgutz/dat.v1/sqlx-runner"
)

type status struct {
	statusID int    `sql:"blog_status_id"`
	status   string `sql:"status"`
}

var allowedStatuses = map[string]int{
	"active":   1,
	"disabled": 2,
	"updating": 3,
}

type BlogQuery struct {
	tableName       struct{} `sql:"crawler.blog"`
	BlogId          int64    `sql:"blog_id"`
	Url             string   `sql:"url"`
	BloggerDomainId int64    `sql:"blogger_domain_id"`
}

type BlogPlacesQuery struct {
	PlaceId          int64   `sql:"place_id"`
	Name             string  `sql:"name"`
	LocationId       int64   `sql:"location_id"`
	WorkPhone        int32   `sql:"work_phone"`
	FormattedAddress string  `sql:"formatted_address"`
	PlaceTypes       []int32 `sql:"place_types"`
	LocationAddress  string  `sql:"address"`
	City             string  `sql:"city"`
	State            string  `sql:"state"`
	CountryId        int64   `sql:"country_id"`
	Province         string  `sql:"province"`
	ZipCode          int32   `sql:"zip_code"`
	Longitude        float32 `sql:"longitude"`
	Latitude         float32 `sql:"latitude"`
}

type PlaceTagsQuery struct {
	PlaceTagId int64  `sql:"place_tag_id"`
	PlaceId    int64  `sql:"place_id"`
	TagId      int64  `sql:"tag_id"`
	TagName    string `sql:"tag_name"`
}

type BlogPlaceTagsQuery struct {
	BlogPlaceId int64  `sql:"blog_place_id"`
	PlaceTagId  int64  `sql:"place_tag_id"`
	BlogId      int64  `sql:"blog_id"`
	TagName     string `sql:"tag_name"`
}

type BlogTagsQuery struct {
	BlogTagId int64  `sql:"blog_tag_id"`
	BlogId    int64  `sql:"blog_id"`
	TagId     int64  `sql:"tag_id"`
	Name      string `sql:"tag_name"`
}

type Blog struct {
	TableName       struct{}  `sql:"crawler.blog"`
	BlogID          int64     `sql:"blog_id,pk"`
	BloggerDomainID int64     `sql:"blogger_domain_id"`
	CreateDate      time.Time `sql:"create_date"`
	URL             string    `sql:"url"`
}

func (db *Db) CreateBlogData(b *Blog) error {
	return db.Client.Insert(b)
}

func (db *Db) UpdateBlogData(b *Blog) error {
	return db.Client.Update(b)
}

func (db *Db) UpdateBlogPlaces(blogID int64, places []*blogdatapb.UpdatePlace) (bool, error) {
	var placesAdd []blogdatapb.BlogPlace
	var placesRemove []string
	for _, place := range places {
		if place.Action == "add" {
			placepb := blogdatapb.BlogPlace{
				BlogId:  blogID,
				VenueId: place.VenueId,
			}
			placesAdd = append(placesAdd, placepb)
		}
		if place.Action == "remove" {
			placesRemove = append(placesRemove, place.VenueId)
		}
	}
	if len(placesAdd) > 0 {
		addErr := db.Client.Insert(&placesAdd)
		if addErr != nil {
			return false, addErr
		}
	}
	if len(placesRemove) > 0 {
		removeIn := pg.In(placesRemove)
		_, removeErr := db.Client.Model((*blogdatapb.BlogPlace)(nil)).
			Where("blog_id = ? AND venue_id IN (?)", blogID, removeIn).Delete()
		if removeErr != nil {
			return false, removeErr
		}
	}
	return true, nil
}

func (db *Db) UpdateBlogTags(blogID int64, tags []*blogdatapb.UpdateTag) (bool, error) {
	var tagsAdd []blogdatapb.BlogTag
	var tagsRemove []int64
	for _, tag := range tags {
		if tag.Action == "add" {
			tagpb := blogdatapb.BlogTag{
				BlogId: blogID,
				TagId:  tag.TagId,
			}
			tagsAdd = append(tagsAdd, tagpb)
		}
		if tag.Action == "remove" {
			tagsRemove = append(tagsRemove, tag.TagId)
		}
	}
	if len(tagsAdd) > 0 {
		addErr := db.Client.Insert(&tagsAdd)
		if addErr != nil {
			return false, addErr
		}
	}
	if len(tagsRemove) > 0 {
		removeIn := pg.In(tagsRemove)
		_, removeErr := db.Client.Model((*blogdatapb.BlogTag)(nil)).
			Where("blog_id = ? AND tag_id IN (?)", blogID, removeIn).Delete()
		if removeErr != nil {
			return false, removeErr
		}
	}
	return true, nil
}

func (db *Db) GetBlogsResponse(domain_id int64) []*blogdatapb.Blog {
	var blogs []*blogdatapb.Blog
	err := db.Client.Model(&blogs).
		Column(
			"blog.id",
			"blog.url",
			"blog.blogger_domain_id",
			"blog.blog_status_id",
			"Places",
			"Tags",
			"BlogStatus").
		Where("blog.blogger_domain_id = ?", domain_id).
		Select()
	if err != nil {
		panic(err)
	}
	return blogs
}

func (db *Db) GetBlogByID(id int64) (BlogQuery, error) {
	var blogData BlogQuery
	log.Log(id)
	_, err := db.Client.Query(&blogData, `SELECT * FROM "crawler".blog WHERE blog_id = ?`, id)
	return blogData, err
}

func (db *Db) CheckBlogExists(url string) (bool, int64) {
	var blogData Blog
	log.Log(url)
	data, err := db.Client.Query(&blogData, `SELECT * FROM "crawler".blog WHERE url = ?`, url)

	if err != nil {
		panic(err)
	}

	return (data != nil && data.RowsReturned() > 0), blogData.BlogID
}

func (db *Db) UpdateBlogStatus(blogId int64, status string) error {
	if statusAllowed(status) {
		var blog blogdatapb.Blog
		_, err := db.Client.Model(&blog).
			Set("blog_status_id = ?", getStatusID(status)).
			Where("id = ?", blogId).
			Update()
		return err
	} else {
		return merr.Forbidden("model.blog", "status not allowed")
	}
}

func getStatusID(status string) int {
	return allowedStatuses[status]
}

func statusAllowed(status string) bool {
	if allowedStatuses[status] > 0 {
		return true
	}
	return false
}
