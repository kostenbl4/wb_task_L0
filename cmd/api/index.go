package main

import (
	"html/template"
	"net/http"
)


func (app *application) HandleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		app.logger.Error("failed to parse html template: " + err.Error())
		errorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		app.logger.Error("failed to execute html template: " + err.Error())
		errorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
}