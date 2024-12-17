package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
	}

	if err := respondJSON(w, http.StatusOK, data); err != nil {
		app.logger.Error("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
		errorJSON(w, http.StatusInternalServerError, err.Error())
	}
}