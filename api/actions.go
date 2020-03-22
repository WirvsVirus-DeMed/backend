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

// HandleBackendStateReq user is opening the frontend panel
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

// ProvideMedRessourceReq user wants to add a new ressource to the local DB
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

// ChangeMedRessourceReq user wants to edit or delete an ressource
func ChangeMedRessourceReq(msg []byte, p *model.Packet, database *sql.DB) []byte {
	// ChangeMedRessourceRequest
	var specpacket model.ChangeMedRessourceRequest
	err := json.Unmarshal(msg, &specpacket)

	if err != nil {
		log.Fatal(err)
	}

	if specpacket.Remove == true {
		fmt.Println("Remove")
		med, _ := db.GetViaID(database, specpacket.MedicineUUID)
		med[0].Delete(database)

		if err != nil {
			log.Fatal(err)
		}

	} else if !specpacket.Remove {
		me, _ := db.GetViaID(database, specpacket.MedicineUUID)
		med := me[0]
		medNew := &db.Medicine{med.UUID, specpacket.EditedMedicine.Title, specpacket.EditedMedicine.Description, specpacket.EditedMedicine.Owner, specpacket.EditedMedicine.Amount, specpacket.EditedMedicine.Pzn}
		medNew.Update(database)

	} else {
		log.Fatal("Not Compatible")
	}

	packet := &model.Packet{p.ID, idcounter, "ChangeMedRessourceResponse"}
	idcounter++
	res := &model.ChangeMedRessourceResponse{*packet}

	jrep, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}

	return jrep
}

// SearchMedRessourceReq User wants to request an search-task on all peer-clients -> IncommingMedRessourceResponse
func SearchMedRessourceReq(msg []byte, p *model.Packet, database *sql.DB) []byte {
	var specpacket model.SearchMedRessourceRequest
	err := json.Unmarshal(msg, &specpacket)

	if err != nil {
		log.Fatal(err)
	}

	//////////////////////////////////////////////////
	// Network stuff
	//
	//
	// Medicine should me in here
	meds := []*db.Medicine{}
	///////////////////////////////////////////////////

	packet := &model.Packet{p.ID, idcounter, "SearchMedRessourceResponse"}
	idcounter++
	res := &model.SearchMedRessourceResponse{meds, *packet}

	jrep, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}

	return jrep
}

// RequestMedRessourceRequest -> IncommingMedRessourceResponse
func RequestMedRessourceRequest(msg []byte, p *model.Packet, database *sql.DB) []byte {
	var specpacket model.RequestMedRessourceRequest
	err := json.Unmarshal(msg, &specpacket)

	if err != nil {
		log.Fatal(err)
	}

	/////////////////////////////////////////////////////
	// Do something with specpacket.MedicineUUID
	// Meds
	meds := db.GetViaID(database, specpacket.MedicineUUID)
	//
	//
	// Network stuff
	//
	//
	//
	// Fill those out with Information from the Network
	accepted := false
	addinfo := ""
	///////////////////////////////////////////////////

	packet := &model.Packet{p.ID, idcounter, "RequestMedRessourceResponse"}
	idcounter++
	res := &model.SearchMedRessourceResponse{accepted, addinfo, *packet}

	jrep, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}

	return jrep
}

// type IncommingMedRessourceRequest struct {
// 	Medicine db.Medicine `json:"ressource"`
// 	Packet
// }

// type IncommingMedRessourceResponse struct {
// 	Accepted bool   `json:"accepted"`
// 	AddInfos string `json:"additionalInformation"`
// 	Packet
// }

// IncommingMedRessourceRequest muss noch gemacht werden

// IncommingMedRessourceResponse
func IncommingMedRessourceResponse(msg []byte, p *model.Packet, database *sql.DB) []byte {

	return []byte{}
}
