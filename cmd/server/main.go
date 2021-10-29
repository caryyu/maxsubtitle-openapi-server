package main

import (
	"log"
	"net/http"
	"os"

	"github.com/caryyu/subtitle-open-server/internal/common"
	"github.com/caryyu/subtitle-open-server/internal/handler"
	"github.com/gorilla/mux"
)

const (
	HTTP_PORT = "HTTP_PORT"
)

func main() {
	initSetup()
	app := &common.App{
		Router: mux.NewRouter(),
	}
	http.Handle("/", app.Router)
	handler := configureRouters(app)
	port := ":" + os.Getenv(HTTP_PORT)
	log.Printf("listen on %s", port)
	log.Fatal(http.ListenAndServe(port, handler))
}

func configureRouters(app *common.App) http.Handler {
	app.Router.HandleFunc("/", handler.IndexHandler(app))
	app.Router.HandleFunc("/health", handler.HeathCheckHandler(app))
	app.Router.HandleFunc("/subtitle/search/{k}", handler.SearchHandler(app)).Methods("GET")
	app.Router.HandleFunc("/subtitle/{id}/download", handler.DownloadHandler(app)).Methods("GET")
	return app.Router
}

func initSetup() {
	if port := os.Getenv(HTTP_PORT); len(port) == 0 {
		os.Setenv(HTTP_PORT, "3000")
	}
}
