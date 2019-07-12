package database

import (
	"os"
	"io/ioutil"
	"path/filepath"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var persistentDb *sqlx.DB
var testing = false

func GetDb() (*sqlx.DB) {
	return persistentDb
}

func Initialize() error {
	var err error
	persistentDb, err = openDb()

	if err != nil {
		return err
	}

	return runMigrations()
}

func SetTestingEnvironment() {
	testing = true
}

func Close() {
	persistentDb.Close()
}

func openDb() (*sqlx.DB, error) {
	filename := getFilename()

	if testing {
		os.Remove(filename)
	}

	return sqlx.Connect("sqlite3", filename)
}

func runMigrations() error {
	pwd, _ := os.Getwd()
	filename := filepath.Join(pwd, "database", "schema.sql")
	file, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	_, err = persistentDb.Exec(string(file))

	return err
}

func getFilename() string {
	pwd, _ := os.Getwd()

	if testing {
		return filepath.Join(os.TempDir(), "test.db")
	} else {
		return filepath.Join(pwd, "database", "sqlite", "gofigure.db")
	}
}