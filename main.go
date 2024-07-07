package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"

	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pratikmane1299/thots/components"
	"github.com/pratikmane1299/thots/db"
	"github.com/pratikmane1299/thots/pages"
	"github.com/pratikmane1299/thots/services"
)

var static embed.FS

func main() {
	thotsDb, err := db.OpenDB()

	thotsService := services.NewThotsService(thotsDb)

	if err != nil {
		log.Fatal("Could not connect to db", err)
	}

	// http.Handle("/static/", http.FileServer(http.FS(static)))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServerFS(os.DirFS("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		thots, err := thotsService.GetAllThots()
		if err != nil {
			fmt.Println("error fetching thots", err)
		}

		pages.Index(thots).Render(r.Context(), w)
	})

	http.HandleFunc("/add-thot", func(w http.ResponseWriter, r *http.Request) {
		thot := r.PostFormValue("thot")

		err = thotsService.AddThot(thot)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		thots, err := thotsService.GetAllThots()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		components.ThotsList(thots).Render(r.Context(), w)
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

		var thotToUpdate services.ThotPayload

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
