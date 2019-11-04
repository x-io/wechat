package api

import (
	"sync"
)

//API API Manager
type API struct {
	accessTokenLock sync.RWMutex
	ticketLock      sync.RWMutex
}

//New API
func New() *API {
	return &API{}
}
