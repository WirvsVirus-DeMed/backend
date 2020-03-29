package api // api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/WirvsVirus-DeMed/backend/model"
	"github.com/gorilla/websocket"
)

type Action func([]byte, *model.Packet, *sql.DB) []byte

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Api(database *sql.DB) {
	// TODO Action hinzufügen
	actions := map[string]Action{
		"ProvideMedRessourceRequest":    ProvideMedRessourceReq,
		"SearchMedRessourceRequest":     SearchMedRessourceReq,
		"RequestMedRessourceRequest":    HandleBackendStateReq,
		"BackendStateRequest":           HandleBackendStateReq, // Hier
		"ChangeMedRessourceRequest":     ChangeMedRessourceReq,
		"IncommingMedRessourceResponse": HandleBackendStateReq, // we ask
	}

	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		// For now allow all connections
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			var packet model.Packet
			err = json.Unmarshal(msg, &packet)
			fmt.Println(packet)

			if err != nil {
				log.Fatal(err)
			}

			for key, value := range actions {
				if key == packet.Type {
					jrep := value(msg, &packet, database)

					if err = conn.WriteMessage(msgType, jrep); err != nil {
						return
					}
				}
			}
		}
	})

	http.Handle("/", http.FileServer(http.Dir("frontend/export")))
	http.ListenAndServe(":8080", nil)
}
