package main

import (
	"encoding/json"
	"net/http"
)


func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}