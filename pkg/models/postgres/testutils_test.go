package postgres

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"testing"
)

func newTestDB(t *testing.T) (*sql.DB, func()) {
	connStr := fmt.Sprintf("dbname=groupics_test user=groupics_test password=%s host=%s port=5432 connect_timeout=10 sslmode=disable", "pass", "localhost")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}

	script, err := ioutil.ReadFile("./testdata/db_test_02_setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	return db, func() {
		script, err := ioutil.ReadFile("./testdata/db_test_03_teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	}
}
