package api

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/WirvsVirus-DeMed/backend/db"
	"github.com/WirvsVirus-DeMed/backend/model"
)

// User is opening the frontend panel
func HandleBackendStateReq(msg []byte) []byte {
	var specpacket model.BackendStateRequest
	err := json.Unmarshal(msg, &specpacket)

	if err != nil {
		log.Fatal(err)
	}

	// meds, err := db.GetAll(database)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	med := &db.Medicine{"1", "21", "1", time.Now(), "1", 1, 1}
	meds := []db.Medicine{*med}

	packet := &model.Packet{1, -1, "BackendStateResponse"}
	res := &model.BackendStateResponse{meds, *packet}

	jrep, err := json.Marshal(res)
	fmt.Println("jrep")
	fmt.Println(string(jrep))
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(jrep))
	return jrep
}

// // User wants to add a new ressource to the local DB
// func ProvideMedRessourceReq(req model.ProvideMedRessourceRequest, database *sql.DB) model.ProvideMedRessourceResponse {

// }

// // User wants to request an search-task on all peer-clients
// func SearchMedRessourceReq(req model.SearchMedRessourceRequest, database *sql.DB) model.SearchMedRessourceResponse {

// }

// // User wants to edit or delete an ressource
// func ChangeMedRessourceReq(req model.ChangeMedRessourceRequest, database *sql.DB) model.ChangeMedRessourceResponse {

// }
