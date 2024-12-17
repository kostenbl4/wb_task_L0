package main

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(r *http.Request, data any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func respondJSON(w http.ResponseWriter, status int, data any) error {
	type msg struct {
		Data any `json:"data"`
	}

	return writeJSON(w, status, &msg{Data: data})
}

func errorJSON(w http.ResponseWriter, status int, errorMsg string) error {

	type msg struct {
		Error string `json:"error"`
	}

	return respondJSON(w, status, &msg{Error: errorMsg})
}

func bytesToJSON(b []byte, data any) error {
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	return nil
}
