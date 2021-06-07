package main

import (
	"database/sql"
	"os"
	"testing"
)

func createTempDBName(t *testing.T) (string, error) {
	rs, err := MakeRandomStr(10)
	if err != nil {
		return "", err
	}
	dbFileName := "/tmp/" + rs + ".db"
	t.Logf("DB file name: %s\n", dbFileName)
	return dbFileName, nil
}

func TestCreateDatabase(t *testing.T) {
	dbFileName, err := createTempDBName(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = createDatabase(dbFileName)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer os.Remove(dbFileName)

	if _, err := os.Stat(dbFileName); err != nil {
		t.Fatal(err)
		return
	}

	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		return
	}

	_, err = db.Query("select * from urls")
	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestGetAndAddUrlsRecord(t *testing.T) {
	dbFileName, err := createTempDBName(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	db, err := ConnectDB(dbFileName)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()

	// add new origin/tiny url
	origin := "https://example.com/hoge/var"
	tiny, err := db.AddTinyURL(origin)
	t.Logf("new tiny is \"%s\"\n", tiny)
	if err != nil {
		t.Fatal(err)
		return
	}

	// get above created tiny url
	result, err := db.GetOriginURL(tiny)
	t.Logf("result is \"%s\"\n", result)
	if err != nil {
		t.Fatal(err)
		return
	}
	if result != origin {
		t.Fatalf("real: %s  expected: %s\n", result, tiny)
		return
	}

	// add same origin/tiny url. Same tiny url is expected.
	sametiny, err := db.AddTinyURL(origin)
	t.Logf("same tiny is \"%s\"\n", sametiny)
	if err != nil {
		t.Fatal(err)
		return
	}
	if sametiny != tiny {
		t.Fatalf("real: %s  expected: %s\n", sametiny, tiny)
		return
	}
}
