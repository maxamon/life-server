package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		fmt.Println("CheckOrigin")
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
	fmt.Println("Start wsHandler")
	return func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("New socket connection")
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
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Println("Msg", msg.Type, msg.Payload)

			switch msg.Type {
			case "get_region":
				creatures := w.GetRegion(msg.Payload.X1, msg.Payload.Y1, msg.Payload.X2, msg.Payload.Y2)
				err := conn.WriteJSON(creatures)
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}

			default:
				err := conn.WriteJSON(`{"error": "unknown msg type"`)
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
	}
}

func regionHandler(w *World) http.HandlerFunc {
	fmt.Println("Start regionHandler")
	return func(rw http.ResponseWriter, r *http.Request) {
		fmt.Print("New connection", r.URL.RawPath)

		x1, _ := strconv.Atoi(r.URL.Query().Get("x1"))
		x2, _ := strconv.Atoi(r.URL.Query().Get("x2"))
		y1, _ := strconv.Atoi(r.URL.Query().Get("y1"))
		y2, _ := strconv.Atoi(r.URL.Query().Get("y2"))

		creatures := w.GetRegion(x1, y1, x2, y2)

		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(creatures)
	}
}
