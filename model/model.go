package model

import "github.com/WirvsVirus-DeMed/backend/db"

type Packet struct {
	ID         int    `json:"requestId"`
	ResponseID int    `json:"responseId"`
	Type       string `json:"type"`
	// Data       map[string]interface{} `json:"data"`
}

type ProivideMedRessourceRequest struct {
	Packet   Packet
	Medicine db.Medicine `json:"ressource"`
}

type ProvideMedRessourceResponse struct {
	Packet  Packet
	Success bool `json:"success"`
}

type SearchMedRessourceRequest struct {
	Packet   Packet
	Keywords []string `json:"keywords"`
}

type SearchMedRessourceResponse struct {
	Packet    Packet
	Medicines []db.Medicine `json:"ressources"`
}

type RequestMedRessourceRequest struct {
	Packet       Packet
	MedicineUUID string `json:"ressourceUuid"`
}

type RequestMedRessourceResponse struct {
	Packet   Packet
	Accepted bool   `json:"accepted"`
	AddInfo  string `json:"additionalInformation"`
}

type BackendStateRequest struct {
	Packet Packet
}

type BackendStateResponse struct {
	Packet   Packet
	OwnItems []db.Medicine `json:"ownItems"`
}

type ChangeMedRessourceRequest struct {
	Packet         Packet
	MedicineUUID   string      `json:"ressourceUuid"`
	Remove         bool        `json:"remove"`
	EditedMedicine db.Medicine `json:"editedRessource"`
}

type ChangeMedRessourceResponse struct {
	Packet Packet
}

type IncommingMedRessourceRequest struct {
	Packet   Packet
	Medicine db.Medicine `json:"ressource"`
}

type IncommingMedRessourceResponse struct {
	Packet   Packet
	Accepted bool   `json:"accepted"`
	AddInfos string `json:"additionalInformation"`
}
