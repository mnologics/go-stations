package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {

	// register routes
	mux := http.NewServeMux()

	// Health Check エンドポイントの登録
	_healthz := handler.NewHealthzHandler()
	mux.Handle("/healthz", _healthz)

	return mux
}
