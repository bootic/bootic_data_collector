package web

import (
	"code.google.com/p/gorilla/mux"
	"text/template"
	"net/http"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", root)
	return router
}

func root (res http.ResponseWriter, req *http.Request) {
	template.Must(template.ParseFiles("views/home.html")).Execute(res, req.Host)
}
