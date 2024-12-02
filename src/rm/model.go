package rm

import (
	"time"
)

type RMDocument struct {
	Bookmarked     bool   `json:"Bookmarked"`
	ID             string `json:"ID"`
	ModifiedClient string `json:"ModifiedClient"`
	Parent         string `json:"Parent"`
	Type           string `json:"Type"`
	VisibleName    string `json:"VissibleName"`
}

type Document struct {
	ID           string    `json:"ID"`
	Name         string    `json:"Name"`
	ModifiedTime time.Time `json:"ModifiedTime"`
}
