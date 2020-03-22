package model

import "github.com/WirvsVirus-DeMed/backend/db"

type Packet struct {
	ID         int    `json:"requestId"`
	ResponseID int    `json:"responseId"`
	Type       string `json:"type"`
	// Data       map[string]interface{} `json:"data"`
}

type ProvideMedRessourceRequest struct {
	Medicine *db.Medicine `json:"ressource"`
	Packet
}

type ProvideMedRessourceResponse struct {
	Success bool `json:"success"`
	Packet
}

type SearchMedRessourceRequest struct {
	Keywords []string `json:"keywords"`
	Packet
}

type SearchMedRessourceResponse struct {
	Medicines []*db.Medicine `json:"ressources"`
	Packet
}

type RequestMedRessourceRequest struct {
	MedicineUUID string `json:"ressourceUuid"`
	Packet
}

type RequestMedRessourceResponse struct {
	Accepted bool   `json:"accepted"`
	AddInfo  string `json:"additionalInformation"`
	Packet
}

type BackendStateRequest struct {
	Packet
}

type BackendStateResponse struct {
	OwnItems []*db.Medicine `json:"ownItems"`
	Packet
}

type ChangeMedRessourceRequest struct {
	MedicineUUID   string       `json:"ressourceUuid"`
	Remove         bool         `json:"remove"`
	EditedMedicine *db.Medicine `json:"editedRessource"`
	Packet
}

type ChangeMedRessourceResponse struct {
	Packet
}

type IncommingMedRessourceRequest struct {
	Medicine db.Medicine `json:"ressource"`
	Packet
}

type IncommingMedRessourceResponse struct {
	Accepted bool   `json:"accepted"`
	AddInfos string `json:"additionalInformation"`
	Packet
}
