package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		log.Println("CheckOrigin")
		return true //for develop
	},
}

type RegionPayload struct {
	X1 int `json:"x1"`
	Y1 int `json:"y1"`
	X2 int `json:"x2"`
	Y2 int `json:"y2"`
}

type ClientMessage struct {
	Type    string        `json:"type"`
	Payload RegionPayload `json:"payload"`
}

func wsHandler(w *World) http.HandlerFunc {
	log.Println("Start wsHandler")
	return func(rw http.ResponseWriter, r *http.Request) {
		log.Println("New socket connection")
		conn, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		for {
			var msg ClientMessage
			err := conn.ReadJSON(&msg)
			if err != nil {
				return
			}

			log.Println("Msg", msg.Type, msg.Payload)

			switch msg.Type {
			case "get_region":
				creatures := w.GetRegion(msg.Payload.X1, msg.Payload.Y1, msg.Payload.X2, msg.Payload.Y2)
				err := conn.WriteJSON(creatures)
				if err != nil {
					return
				}

			default:
				conn.WriteJSON(map[string]string{"error": "unknown msg type"})
				return
			}
		}
	}
}

func regionHandler(w *World) http.HandlerFunc {
	log.Println("Start regionHandler")
	return func(rw http.ResponseWriter, r *http.Request) {
		log.Print("New connection", r.URL.RawPath)

		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			rw.WriteHeader(http.StatusOK)
			return
		}

		x1, err1 := strconv.Atoi(r.URL.Query().Get("x1"))
		x2, err2 := strconv.Atoi(r.URL.Query().Get("x2"))
		y1, err3 := strconv.Atoi(r.URL.Query().Get("y1"))
		y2, err4 := strconv.Atoi(r.URL.Query().Get("y2"))

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			http.Error(rw, "invalid query parameters", http.StatusBadRequest)
			return
		}

		if x1 < 0 || x2 >= w.Width || y1 < 0 || y2 >= w.Height || x1 > x2 || y1 > y2 {
			http.Error(rw, "coordinates out of bounds or invalid range", http.StatusBadRequest)
			return
		}

		creatures := w.GetRegion(x1, y1, x2, y2)

		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(creatures)
	}
}
