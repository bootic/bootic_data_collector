package web

import (
	"code.google.com/p/gorilla/mux"
	"text/template"
	"net/http"
	//"fmt"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", tmpl("views/home.html"))
	router.HandleFunc("/realtime_ftw", tmpl("views/realtime.html"))
	router.HandleFunc("/app.js", javascript("js/app.js"))
	return router
}

func tmpl(name string) (handle func (http.ResponseWriter, *http.Request)) {
	return func(res http.ResponseWriter, req *http.Request) {
		template.Must(template.ParseFiles(name)).Execute(res, req.Host)
	}
}

func javascript(name string) (handle func (http.ResponseWriter, *http.Request)) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/javascript")
		template.Must(template.ParseFiles(name)).Execute(res, req.Host)
	}
}