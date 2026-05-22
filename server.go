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

type WorldBoundsPayload struct {
	MinX int `json:"min_x"`
	MinY int `json:"min_y"`
	MaxX int `json:"max_x"`
	MaxY int `json:"max_y"`
}

type ClientMessage struct {
	Type    string        `json:"type"`
	Payload RegionPayload `json:"payload"`
}

type ServerMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type RegionResponse struct {
	Creatures []*Creature   `json:"creatures"`
	Terrain   [][]byte      `json:"terrain"`
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

		// Send world bounds on connection
		boundsMsg := ServerMessage{
			Type: "world_bounds",
			Payload: WorldBoundsPayload{
				MinX: 0,
				MinY: 0,
				MaxX: w.Width,
				MaxY: w.Height,
			},
		}
		if err := conn.WriteJSON(boundsMsg); err != nil {
			log.Printf("Failed to send world bounds: %v", err)
			return
		}

		for {
			var msg ClientMessage
			err := conn.ReadJSON(&msg)
			if err != nil {
				return
			}

			log.Println("Msg", msg.Type, msg.Payload)

			switch msg.Type {
			case "get_region":
				creatures, terrain := w.GetRegionWithTerrain(msg.Payload.X1, msg.Payload.Y1, msg.Payload.X2, msg.Payload.Y2)
				response := RegionResponse{
					Creatures: creatures,
					Terrain:   terrain,
				}
				err := conn.WriteJSON(response)
				if err != nil {
					return
				}

			case "get_stats":
				stats := w.GetStats()
				response := ServerMessage{
					Type:    "stats",
					Payload: stats,
				}
				err := conn.WriteJSON(response)
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

		creatures, terrain := w.GetRegionWithTerrain(x1, y1, x2, y2)

		response := RegionResponse{
			Creatures: creatures,
			Terrain:   terrain,
		}

		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(response)
	}
}

func statsHandler(w *World) http.HandlerFunc {
	log.Println("Start statsHandler")
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			rw.WriteHeader(http.StatusOK)
			return
		}

		stats := w.GetStats()

		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(stats)
	}
}
