package http

import (
	"code.google.com/p/gorilla/mux"
	"twitter1/http/handlers"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", handlers.HelloWorld)
	router.HandleFunc("/{screenName}", handlers.ShowTweets)
	return router
}
