package rdbms

// https://github.com/go-pg/pg

import (
    "fmt"
    "time"
    "github.com/go-pg/pg"
    "github.com/micro/go-log"
    "bitbucket.org/appgoplaces/travelplatform-system/lib/conf"
    
    "database/sql"

    _ "github.com/lib/pq"
    "gopkg.in/mgutz/dat.v1"
    "gopkg.in/mgutz/dat.v1/sqlx-runner"
)

var (
    env = conf.GetConf()
)
// Connect connects to a MongoDB instance
func Connect(serviceName string) *pg.DB {
    log.Log("Connecting " + serviceName + " to the RDBMS DataBase...")
    db := pg.Connect(&pg.Options{
        Addr: env.Rdbms.Address,
        User: env.Rdbms.User,
        Password: env.Rdbms.Password,
        Database: env.Rdbms.Database,
    })
    db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
    	query, err := event.FormattedQuery()
    	if err != nil {
    		panic(err)
    	}

    	log.Logf("%s %s", time.Since(event.StartTime), query)
    })
    log.Log("Sucessfully connected " + serviceName + " to RDBMS Database")
    return db
}

func ConnectNew(serviceName string) *runner.DB {
    log.Log("Connecting " + serviceName + " to new client for RDBMS DataBase...")
    // create a normal database connection through database/sql
    log.Log(env.Rdbms)
    option := fmt.Sprintf(
        "dbname=%s user=%s password=%s host=%s sslmode=disable",
        env.Rdbms.Database,
        env.Rdbms.User,
        env.Rdbms.Password,
        env.Rdbms.Host)
    db, err := sql.Open("postgres", option)
    if err != nil {
        panic(err)
    }

    // ensures the database can be pinged with an exponential backoff (15 min)
    runner.MustPing(db)

    // set to reasonable values for production
    db.SetMaxIdleConns(4)
    db.SetMaxOpenConns(16)

    // set this to enable interpolation
    dat.EnableInterpolation = true

    // set to check things like sessions closing.
    // Should be disabled in production/release builds.
    dat.Strict = false

    // Log any query over 10ms as warnings. (optional)
    runner.LogQueriesThreshold = 10 * time.Millisecond

    return runner.NewDB(db, "postgres")
}