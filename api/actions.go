package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/WirvsVirus-DeMed/backend/db"
	"github.com/WirvsVirus-DeMed/backend/model"
)

// TODO: Error handling
// TODO: Request und Response id ändern
var idcounter = 1

// User is opening the frontend panel
func HandleBackendStateReq(msg []byte, p *model.Packet, database *sql.DB) []byte {
	var specpacket model.BackendStateRequest
	err := json.Unmarshal(msg, &specpacket)

	if err != nil {
		log.Fatal(err)
	}

	meds, err := db.GetAll(database)
	if err != nil {
		log.Fatal(err)
	}

	packet := &model.Packet{p.ID, idcounter, "BackendStateResponse"}
	idcounter++
	res := &model.BackendStateResponse{meds, *packet}

	jrep, err := json.Marshal(res)
	fmt.Println("jrep")
	fmt.Println(string(jrep))
	if err != nil {
		log.Fatal(err)
	}

	return jrep
}

// ProvideMedRessourceReq User wants to add a new ressource to the local DB
func ProvideMedRessourceReq(msg []byte, p *model.Packet, database *sql.DB) []byte {
	var specpacket model.ProvideMedRessourceRequest
	err := json.Unmarshal(msg, &specpacket)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(specpacket.Medicine)

	specpacket.Medicine.Add(database)
	if err != nil {
		log.Fatal(err)
	}

	packet := &model.Packet{p.ID, idcounter, "ProvideMedRessourceRequest"}
	idcounter++
	res := &model.ProvideMedRessourceResponse{true, *packet}

	jrep, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}

	return jrep
}

// ChangeMedRessourceReq User wants to edit or delete an ressource
// func ChangeMedRessourceReq(msg []byte, database *sql.DB) []byte {

// }

// SearchMedRessourceReq User wants to request an search-task on all peer-clients
// func SearchMedRessourceReq(msg []byte, p *model.Packet, database *sql.DB) []byte {

// }
