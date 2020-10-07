package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

var (
	router = mux.NewRouter()
)

// StartRouter read the routes and start the server
func StartRouter() {
	Routes()
	http.ListenAndServe(":9090", router)
}
