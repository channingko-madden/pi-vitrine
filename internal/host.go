package internal

import (
	"html/template"
	"log"
	"net/http"
)

func ErrorMessage(w http.ResponseWriter, errorMsg string) {
	temp, err := template.ParseFiles("internal/templates/error_msg.html")
	if err == nil {
		temp.Execute(w, errorMsg)
	} else {
		log.Default().Print(err)
	}

}
