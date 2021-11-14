package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/caryyu/subtitle-open-server/internal/common"
	"github.com/caryyu/subtitle-open-server/internal/resource"
	"github.com/gorilla/mux"
)

// Index
func IndexHandler(app *common.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome"))
	}
}

// Check Application
func HeathCheckHandler(app *common.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ready"))
	}
}

// Search Subtitles
func SearchHandler(app *common.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var items []resource.Subtitle
		var err error

		if keyword := vars["k"]; len(keyword) != 0 {
			c := app.Fetcher
			if items, err = c.Search(keyword); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Oops, Something went wrong"))
				return
			}
		}

		var bytes []byte
		bytes, err = json.Marshal(items)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Oops, Something went wrong"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	}
}

// Download Subtitle
func DownloadHandler(app *common.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var file *resource.File
		var err error
		if id := vars["id"]; len(id) != 0 {
			c := app.Fetcher
			if file, err = c.GetFromCache(id); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Oops, Something went wrong"))
				return
			}
		}

		if file == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		w.Header().Set("Content-Disposition", "attachment; filename="+file.Name)
		w.WriteHeader(http.StatusOK)
		w.Write(file.Bytes)
	}
}
