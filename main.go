package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type thot struct {
	Id      int       `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

type thotsDb struct {
	db *sql.DB
}

func (t *thotsDb) tableExists() bool {
	_, err := t.db.Query("SELECT * FROM	thots")
	if err == nil {
		return true
	}

	return false
}

func (t *thotsDb) createTable() error {
	_, err := t.db.Exec(`CREATE TABLE "thots" ("id" INTEGER, "thot" TEXT NOT NULL, "created" DATETIME, PRIMARY KEY("id" AUTOINCREMENT))`)
	return err
}

func openDB() (thotsDb, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return thotsDb{}, err
	}
	db, err := sql.Open("sqlite3", filepath.Join(cwd, "thots.db"))

	t := thotsDb{db}
	if !t.tableExists() {
		err := t.createTable()
		if err != nil {
			return thotsDb{}, err
		}
	}

	return t, nil
}

type thotsService struct {
	thotsDb thotsDb
}

type ThotPayload struct {
	Thot string `json:"thot"`
}

func (ts thotsService) GetAllThots() ([]thot, error) {
	var thots []thot

	rows, err := ts.thotsDb.db.Query("select * from thots")
	if err != nil {
		return thots, fmt.Errorf("unable to fetch thots %v", err)
	}

	for rows.Next() {
		var thot thot
		err := rows.Scan(
			&thot.Id,
			&thot.Name,
			&thot.Created,
		)
		if err != nil {
			return thots, err
		}

		thots = append(thots, thot)
	}

	return thots, nil
}

func (ts thotsService) AddThot(thot string) error {
	result, err := ts.thotsDb.db.Exec("insert into thots (thot, created) values(?, ?)", thot, time.Now())
	if err != nil {
		return fmt.Errorf("error adding a new thot %v", err)
	}

	_, err = result.LastInsertId()
	return err
}

func (ts thotsService) UpdateThot(id string, thot string) error {
	result, err := ts.thotsDb.db.Exec("update thots set thot=? where id=?", thot, id)
	if err != nil {
		return fmt.Errorf("error updating thot for id %s %v", id, err)
	}

	_, err = result.LastInsertId()
	return err
}

func (ts thotsService) DeleteThot(id string) error {
	result, err := ts.thotsDb.db.Exec("delete from thots where id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting thot with %s - %v", id, err)
	}

	_, err = result.LastInsertId()
	return err
}

func NewThotsService(thotsDb thotsDb) thotsService {
	return thotsService{thotsDb: thotsDb}
}

func main() {
	thotsDb, err := openDB()

	thotsService := NewThotsService(thotsDb)

	if err != nil {
		log.Fatal("Could not connect to db", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		thots, err := thotsService.GetAllThots()
		if err != nil {
			fmt.Println("error fetching thots", err)
		}
		fmt.Println("thots", thots)
		io.WriteString(w, "hola amigos")
	})

	http.HandleFunc("/add-thot", func(w http.ResponseWriter, r *http.Request) {
		thotBody := ThotPayload{}

		err := json.NewDecoder(r.Body).Decode(&thotBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = thotsService.AddThot(thotBody.Thot)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/delete-thot/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		err := thotsService.DeleteThot(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/update-thot/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		var thotToUpdate ThotPayload

		err := json.NewDecoder(r.Body).Decode(&thotToUpdate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = thotsService.UpdateThot(id, thotToUpdate.Thot)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	if err := http.ListenAndServe(":4242", nil); err != nil {
		log.Fatal("Could not start server", err)
	}

	fmt.Println("server running on port localhost:4242")
}
