package internal

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

//go:embed templates
var files embed.FS

func ErrorMessage(w http.ResponseWriter, errorMsg string) {
	temp, err := template.ParseFS(files, "templates/error_msg.html")
	if err != nil {
		panic(err)
	}
	temp.Execute(w, errorMsg)
}

func InfoMessage(w http.ResponseWriter, msg string) {
	temp, err := template.ParseFS(files, "templates/info_msg.html")
	if err != nil {
		panic(err)
	}
	temp.Execute(w, msg)
}

type HostError struct {
	Error   error
	Message string
	Code    int
}

type HostErrorHandler func(http.ResponseWriter, *http.Request) *HostError

func (fn HostErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError
		log.Default().Print(e)
		http.Error(w, e.Message, e.Code)
	}
}
