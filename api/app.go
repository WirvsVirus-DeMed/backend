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

type Action func([]byte, *sql.DB) []byte

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Api(database *sql.DB) {
	// TODO Action hinzufügen
	actions := map[string]Action{
		"ProvideMedRessourceRequest":    ProvideMedRessourceReq,
		"SearchMedRessourceRequest":     HandleBackendStateReq,
		"SearchMedRessourceResponse":    HandleBackendStateReq,
		"RequestMedRessourceRequest":    HandleBackendStateReq,
		"RequestMedRessourceResponse":   HandleBackendStateReq,
		"BackendStateRequest":           HandleBackendStateReq, // Hier
		"ChangeMedRessourceRequest":     HandleBackendStateReq,
		"IncommingMedRessourceResponse": HandleBackendStateReq,
	}

	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
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
					jrep := value(msg, database)

					if err = conn.WriteMessage(msgType, jrep); err != nil {
						return
					}
				}
			}

			// Silly and not Clean but Golang has no fricking Generics
			// if packet.Type == "ProvideMedRessourceRequest" {
			// 	var specpacket model.ProvideMedRessourceRequest
			// 	err = json.Unmarshal(msg, &specpacket)

			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}
			// 	fmt.Println(specpacket)
			// } else if packet.Type == "SearchMedRessourceRequest" {
			// 	var specpacket model.SearchMedRessourceRequest
			// 	err = json.Unmarshal(msg, &specpacket)

			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}
			// 	fmt.Println(specpacket)
			// } else if packet.Type == "SearchMedRessourceResponse" {
			// 	var specpacket model.SearchMedRessourceResponse
			// 	err = json.Unmarshal(msg, &specpacket)

			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}
			// 	fmt.Println(specpacket)
			// } else if packet.Type == "RequestMedRessourceRequest" {
			// 	var specpacket model.RequestMedRessourceRequest
			// 	err = json.Unmarshal(msg, &specpacket)

			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}
			// 	fmt.Println(specpacket)
			// } else if packet.Type == "RequestMedRessourceResponse" {
			// 	var specpacket model.RequestMedRessourceResponse
			// 	err = json.Unmarshal(msg, &specpacket)

			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}
			// 	fmt.Println(specpacket)
			// } else if packet.Type == "BackendStateRequest" {
			// 	jrep := HandleBackendStateReq(msg)

			// 	if err = conn.WriteMessage(msgType, jrep); err != nil {
			// 		return
			// 	}
			// } else if packet.Type == "ChangeMedRessourceRequest" {
			// 	var specpacket model.ChangeMedRessourceRequest
			// 	err = json.Unmarshal(msg, &specpacket)

			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}
			// 	fmt.Println(specpacket)
			// } else if packet.Type == "IncommingMedRessourceResponse" {
			// 	var specpacket model.IncommingMedRessourceResponse
			// 	err = json.Unmarshal(msg, &specpacket)

			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}
			// 	fmt.Println(specpacket)
			// }

			// med := []db.Medicine{}
			// rep := &Packet{1, -1, "BackendStateResponse", }
			// jrep, err := json.Marshal(rep)
			// fmt.Println("jrep")
			// fmt.Println(string(jrep))
			// if err != nil {
			// 	log.Fatal(err)
			// }

			// Print the message to the console
			// fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			// Write message back to browser
			// if err = conn.WriteMessage(msgType, jrep); err != nil {
			// 	return
			// }
		}
	})

	http.Handle("/", http.FileServer(http.Dir("frontend/export")))
	http.ListenAndServe(":8080", nil)
}
