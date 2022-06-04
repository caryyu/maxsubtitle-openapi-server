package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/caryyu/subtitle-open-server/internal/common"
	"github.com/caryyu/subtitle-open-server/internal/handler"
	"github.com/caryyu/subtitle-open-server/internal/resource"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
)

const (
	HTTP_PORT  = "HTTP_PORT"
	CACHE_PATH = "CACHE_PATH"
)

func main() {
	app := &common.App{
		Router:      mux.NewRouter(),
		Fetcher:     resource.NewA4kDotNet(),
		CacheClient: createCacheClient(),
	}

	initSetup(app)
	//http.Handle("/", app.Router)
	handler := configureRouters(app)
	port := ":" + os.Getenv(HTTP_PORT)
	log.Printf("listen on %s", port)
	log.Fatal(http.ListenAndServe(port, handler))
}

func configureRouters(app *common.App) http.Handler {
	app.Router.HandleFunc("/", handler.IndexHandler(app))
	app.Router.HandleFunc("/health", handler.HeathCheckHandler(app))
	app.Router.HandleFunc("/subtitle/{id}/download", handler.DownloadHandler(app)).Methods("GET")
	//app.Router.HandleFunc("/subtitle/search/{k}", handler.SearchHandler(app)).Methods("GET")
	//Equip the cache mechanism for handling fast response
	app.Router.Handle("/subtitle/search/{k}", app.CacheClient.Middleware(handler.SearchHandler(app))).Methods("GET")

	r := handlers.LoggingHandler(os.Stdout, app.Router)
	return r
}

func initSetup(app *common.App) {
	if port := os.Getenv(HTTP_PORT); len(port) == 0 {
		os.Setenv(HTTP_PORT, "3000")
	}
	if path := os.Getenv(CACHE_PATH); len(path) != 0 {
		app.Fetcher.CacheDir = path
	}
}

// https://github.com/victorspringer/http-cache
func createCacheClient() *cache.Client {
	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(10000000),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(10*time.Minute),
		cache.ClientWithRefreshKey("opn"),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return cacheClient
}
