package test

import (
	"database/sql"
	"fmt"
	"github.com/statistico/statistico-strategy/internal/trader/bootstrap"
	"testing"
)

func GetConnection(t *testing.T, tables []string) (*sql.DB, func()) {
	db := bootstrap.BuildConfig().Database

	dsn := "host=%s port=%s user=%s " + "password=%s dbname=%s sslmode=disable"

	psqlInfo := fmt.Sprintf(dsn, db.Host, db.Port, db.User, db.Password, db.Name)

	conn, err := sql.Open(db.Driver, psqlInfo)

	if err != nil {
		panic(err)
	}

	return conn, func() {
		for _, table := range tables {
			_, err := conn.Exec("delete from " + table)
			if err != nil {
				t.Fatalf("Failed to clear database. %s", err.Error())
			}
		}
	}
}
