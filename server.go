package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func regionHandler(w *World) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		x1, _ := strconv.Atoi(r.URL.Query().Get("x1"))
		x2, _ := strconv.Atoi(r.URL.Query().Get("x2"))
		y1, _ := strconv.Atoi(r.URL.Query().Get("y1"))
		y2, _ := strconv.Atoi(r.URL.Query().Get("y2"))

		creatures := w.GetRegion(x1, y1, x2, y2)

		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(creatures)
	}
}
