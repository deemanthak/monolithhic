package store

import (
	"github.com/google/uuid"

	coffeeco "github.com/deemanthak/monolithhic/internal"
)

type Store struct {
	ID              uuid.UUID
	Location        string
	ProductsForSale []coffeeco.Product
}
