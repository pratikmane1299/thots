package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"
)

type Thot struct {
	Id      int       `json:"id"`
	Thot    string    `json:"thot"`
	Created time.Time `json:"created"`
}

type ThotsDb struct {
	Db *sql.DB
}

func (t *ThotsDb) tableExists() bool {
	_, err := t.Db.Query("SELECT * FROM	thots")
	if err == nil {
		return true
	}

	return false
}

func (t *ThotsDb) createTable() error {
	_, err := t.Db.Exec(`CREATE TABLE "thots" ("id" INTEGER, "thot" TEXT NOT NULL, "created" DATETIME, PRIMARY KEY("id" AUTOINCREMENT))`)
	return err
}

func OpenDB() (ThotsDb, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return ThotsDb{}, err
	}
	db, err := sql.Open("sqlite3", filepath.Join(cwd, "thots.db"))

	t := ThotsDb{db}
	if !t.tableExists() {
		err := t.createTable()
		if err != nil {
			return ThotsDb{}, err
		}
	}

	return t, nil
}
