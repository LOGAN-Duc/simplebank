package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/techschool/simplebank/util"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("could not load config", err)
	}
	testDB, err = sql.Open(config.DbDriver, config.DbSource)

	if err != nil {
		log.Fatal("could not connect to db", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
