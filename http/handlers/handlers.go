package handlers

import (
	"code.google.com/p/gorilla/mux"
	"fmt"
	"net/http"
	//"twitter1/db"
)

func HelloWorld(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hello World")
}

func ShowTweets(res http.ResponseWriter, req *http.Request) {

	//user := new(db.User)

	fmt.Fprintf(res, "TODO: render tweets for '%s'\n", mux.Vars(req)["screenName"])
	// source, result, err := func () (source string, result string, err error) {

	//   if source, err = sourceUrl(mux.Vars(req)["source"]); err != nil {
	//     return
	//   }

	//   result, err = callSource(source)

	//   return
	// }()

	// if err == nil {
	//   fmt.Fprintf(res, "from source %s\n\n%s", source, result)
	// } else {
	//   fmt.Fprint(res, err)
	// }
}
