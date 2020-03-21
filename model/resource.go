package model

import (
	"time"
)

type Medicine struct {
	UUID       string    `json:"uuid"`
	Title      string    `json:"title"`
	Desciption string    `json:"desciption"`
	CreatedAt  time.Time `json:"createdAt"`
	Owner      string    `json:"owner"`
}

type Packet struct {
	UUID string
	Type string
	data []byte
}
