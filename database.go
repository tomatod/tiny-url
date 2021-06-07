package main

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type DB struct {
	*sql.DB
}

const SQL_CREATE_URLS = `
	create table urls (
		tiny text not null primary key, 
		origin text not null unique
	);
`

func ConnectDB(dbFileName string) (*DB, error) {
	// If specified database file is not found, new database file is created.
	if _, err := os.Stat(dbFileName); err != nil {
		Warnf("Specified database file \"%s\" was not found. So new database (.db file) will be created.\n", dbFileName)
		if createDatabase(dbFileName) != nil {
			Errorf("DatabaseError: Creating database \"%s\" was failed.\n", dbFileName)
			return nil, err
		}
	}

	_db, err := sql.Open("sqlite3", dbFileName)
	var db DB
	db.DB = _db
	return &db, err
}

func createDatabase(fileName string) error {
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(SQL_CREATE_URLS)
	return err
}

func (db *DB) GetOriginURL(tiny string) (string, error) {
	rows, err := db.Query("SELECT origin From urls where tiny = $1", tiny)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		return "", errors.New("DatabaseError: Specified tiny path \"" + tiny + "\" was not found.")
	}

	var origin string
	if err = rows.Scan(&origin); err != nil {
		Warnf("Select urls table query result couldn't be read. Error: \"%v\"\n", err)
		return "", err
	}

	return origin, nil
}

func (db *DB) GetTinyURL(origin string) (string, error) {
	rows, err := db.Query("SELECT tiny FROM urls WHERE origin = $1", origin)
	if err != nil {
		Warnf("Select query of urls table is failed.")
		return "", err
	}
	defer func() {
		rows.Close()
		Debugf("Row is closed.")
	}()
	if rows.Next() {
		var tiny string
		if err = rows.Scan(&tiny); err != nil {
			return "", err
		}
		return tiny, nil
	}
	return "", nil
}

func (db *DB) AddTinyURL(origin string) (string, error) {
	tiny, err := db.GetTinyURL(origin)
	if err != nil {
		Warnf("GetTinyURL() is failed.")
		return "", err
	}
	if tiny != "" {
		return tiny, nil
	}

	tx, err := db.Begin()
	if err != nil {
		return "", err
	}
	defer func() {
		if p := recover(); p != nil {
			Errorf("DatabaseError: Panic occur.")
			tx.Rollback()
		}
		if err != nil {
			Warnf("Transaction is rollbacked.")
			tx.Rollback()
			return
		}
	}()

	tiny, err = MakeRandomStr(10)
	if err != nil {
		Warnf("Making random string by MakeRandomStr(). This is unexpected error. Error: \"%v\"\n", err)
		return "", err
	}

	insert, err := tx.Prepare("INSERT INTO urls VALUES(?, ?)")
	if err != nil {
		Warnf("Insert query of urls table is failed.")
		tiny, _ = db.GetTinyURL(origin)
		if tiny != "" {
			return tiny, nil
		}
		return "", err
	}
	defer func() {
		insert.Close()
		Debugf("Insert statement is closed.")
	}()

	if _, err = insert.Exec(tiny, origin); err != nil {
		Warnf("Faild to add new record to urls in execute query. Error: %v \n", err)
		return "", err
	}
	if err = tx.Commit(); err != nil {
		Warnf("Faild to add new record to urls in commit result. Error: %v \n", err)
		return "", err
	}

	Infof("New URL is added. origin:'%s' tiny:'%s'\n", origin, tiny)
	return tiny, nil
}
