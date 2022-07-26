package types

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
)

type StorageAsk struct {
	StorageAsk *storagemarket.StorageAsk `json:"storage_ask"`
	Address    address.Address `json:"address"`
}

