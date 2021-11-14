package common

import (
	"github.com/caryyu/subtitle-open-server/internal/resource"
	"github.com/gorilla/mux"
)

type App struct {
	Router  *mux.Router
	Fetcher *resource.A4kDotNet
}
