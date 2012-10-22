package web

import (
	"code.google.com/p/gorilla/mux"
	"datagram.io/daemon/web/handlers"
)

func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/", handlers.RenderHtml("main.html"))
	router.HandleFunc("/js/{script:.*}", handlers.RenderJavascripts)

	return router
}
