package models

import (
	rdbms "bitbucket.org/appgoplaces/travelplatform-system/db/rdbms"
	pg "github.com/go-pg/pg"
)

func init() {
	// Register many to many model so ORM can better recognize m2m relation.
	// This should be done before dependant models are used.
	// orm.RegisterTable((*BlogPlace)(nil))
	// orm.RegisterTable((*BlogTag)(nil))
}

type Db struct {
	Client *pg.DB
	// NewClient *runner.DB
}

func NewDb(service string) *Db {
	db := Db{
		Client: rdbms.Connect(service),
	}
	return &db
}

// I want to be able to inject interface types and make it ready to be called