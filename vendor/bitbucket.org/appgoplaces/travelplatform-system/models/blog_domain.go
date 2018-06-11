package models

import (
	"errors"

	blogDatapb "bitbucket.org/appgoplaces/service-protos/blog-data"
)

type BloggerDomain struct {
	TableName       struct{} `sql:"crawler.blogger_domain"`
	BloggerDomainId int64    `sql:"blogger_domain_id,pk"`
	UserId          int64    `sql:"user_id"`
	Url             string   `sql:"url"`
	Enabled         bool     `sql:"enabled"`
}

type BloggerDomainQuery struct {
	BloggerDomainId int64  `sql:"blogger_domain_id,pk" json:"blogger_domain_id"`
	UserId          int64  `sql:"user_id" json:"user_id"`
	Url             string `sql:"url" json:"url"`
	Enabled         bool   `sql:"enabled" json:"enabled"`
}

// type UserInfo struct {
// 	TableName struct{} `sql:"users"`
// 	ID        int64    "sql:id, pk"
// 	FirstName string   "sql:firstname"
// 	LastName  string   "sql:lastname"
// 	Email     string   "sql:email"
// }

type UserInfoQuery struct {
	ID        int64  `sql:"id,pk" json:"id"`
	FirstName string `sql:"firstname" json:"firstname"`
	LastName  string `sql:"lastname" json:"lastname"`
	Email     bool   `sql:"email" json:"email"`
}

type GetBlogDomain struct {
	Id        int64  `sql:"blogger_domain_id"`
	Url       string `sql:"url"`
	UserId    int64  `sql:"user_id"`
	Enabled   bool   `sql:"enabled"`
	BlogCount int32  `sql:blog_count`
}

// AddBlogDomain method takes in url and checks if it exists in DB,
// If not it will add the domain, else it will return a false
func (db *Db) AddBlogDomain(url string) (bool, error) {
	var blogDomain blogDatapb.BlogDomain
	err := db.Client.Model(&blogDomain).Where("url = ?", url).Select()
	if err != nil {
		if err.Error() != "pg: no rows in result set" {
			return false, err
		}
	}
	if blogDomain.BloggerDomainId > 0 {
		err = errors.New("blog domain already exists")
		return true, err
	}
	blogDomain.Url = url
	err = db.Client.Insert(&blogDomain)
	if err != nil {
		return false, err
	}
	return true, nil
}

// UpdateBlogDomain method takes in blogdomain proto object and updates it
func (db *Db) UpdateBlogDomain(b *blogDatapb.BlogDomain) (bool, error) {
	err := db.Client.Update(b)
	if err != nil {
		if err.Error() == "pg: no rows in result set" {
			return true, err
		}
		return false, err
	}
	return err == nil, err
}

// GetBlogDomain method takes blogdomain proto object with BloggerDomainId
// Set in struct and updates whole BlogDomain Struct pointer
func (db *Db) GetBlogDomain(b *blogDatapb.BlogDomain) (bool, error) {
	err := db.Client.Select(b)
	if err != nil {
		if err.Error() == "pg: no rows in result set" {
			return true, err
		}
		return false, err
	}
	return err == nil, err
}

// DeleteBlogDomain method takes blogdomain proto object and Deletes if
// object's BloggerDomainID exists in DB
func (db *Db) DeleteBlogDomain(b *blogDatapb.BlogDomain) (bool, error) {
	err := db.Client.Delete(b)
	if err != nil {
		if err.Error() == "pg: no rows in result set" {
			return true, err
		}
		return false, err
	}
	return err == nil, err
}

func (db *Db) GetBlogDomains() ([]*blogDatapb.GetBlogDomain, error) {
	var blogDomains []GetBlogDomain
	pbBlogDomains := make([]*blogDatapb.GetBlogDomain, 0, len(blogDomains))
	_, err := db.Client.Query(&blogDomains,
		`SELECT blog_domain.blogger_domain_id, blog_domain.url, blog_domain.user_id, blog_domain.enabled, COUNT(blog.id) AS blog_count
		FROM "crawler".blogger_domain AS blog_domain
		LEFT JOIN "crawler".blog AS blog ON blog_domain.blogger_domain_id = blog.blogger_domain_id
		GROUP BY blog_domain.blogger_domain_id;`)
	if err != nil {
		return pbBlogDomains, err
	}
	for _, blogDomain := range blogDomains {
		blogDatapb := blogDatapb.GetBlogDomain(blogDomain)
		pbBlogDomains = append(pbBlogDomains, &blogDatapb)
	}
	return pbBlogDomains, nil
}

func (db *Db) CreateBloggerDomain(b *BloggerDomain) error {
	return db.Client.Insert(b)
}

func (db *Db) UpdateBloggerDomain(b *BloggerDomain) error {
	return db.Client.Update(b)
}

func (db *Db) GetAdminSystemAccountUserID() int64 {
	var user UserInfoQuery

	_, err := db.Client.Query(&user, `SELECT * FROM "public".users WHERE firstname LIKE ?`, "admin_service_account")
	if err != nil {
		panic(err)
	}

	return user.ID
}

func (db *Db) GetBloggerDomainsByUserID(id int64) ([]BloggerDomainQuery, error) {
	var bloggerDomains []BloggerDomainQuery
	_, err := db.Client.Query(&bloggerDomains, `SELECT * FROM "crawler".blogger_domain WHERE user_id = ?`, id)
	return bloggerDomains, err
}

func (db *Db) GetBloggerDomainBlogByURL(url string) BloggerDomainQuery {
	var blog BloggerDomainQuery
	_, err := db.Client.Query(&blog, `SELECT * FROM "crawler".blogger_domain WHERE url LIKE ?`, url)

	if err != nil {
		panic(err)
	}
	return blog
}

func (db *Db) CheckBloggerDomainExists(url string) (bool, int64) {
	var bloggerDomainData BloggerDomain
	data, err := db.Client.Query(&bloggerDomainData, `SELECT * FROM "crawler".blogger_domain WHERE url = ?`, url)

	if err != nil {
		panic(err)
	}

	return (data != nil && data.RowsReturned() > 0), bloggerDomainData.BloggerDomainId
}

// Create ToProto methods here
