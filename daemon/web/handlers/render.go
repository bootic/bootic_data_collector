package handlers

import (
	"code.google.com/p/gorilla/mux"
	"fmt"
	htmlTemplate "html/template"
	"net/http"
	textTemplate "text/template"
)

func RenderHtml(name string) (handle func(http.ResponseWriter, *http.Request)) {

	name = fmt.Sprintf("webapp/html/%s", name)

	return func(res http.ResponseWriter, req *http.Request) {
		htmlTemplate.Must(htmlTemplate.ParseFiles(name)).Execute(res, req.Host)
	}
}

func RenderJavascripts(res http.ResponseWriter, req *http.Request) {

	script := fmt.Sprintf("./webapp/js/%s", mux.Vars(req)["script"])

	res.Header().Add("Content-Type", "application/javascript")

	textTemplate.Must(textTemplate.ParseFiles(script)).Execute(res, req.Host)
}
