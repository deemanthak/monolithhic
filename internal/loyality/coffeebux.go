package loyality

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	coffeeco "github.com/deemanthak/monolithhic/internal"
	"github.com/deemanthak/monolithhic/internal/store"
)

type CoffeBux struct {
	ID                                   uuid.UUID
	store                                store.Store
	coffeeLover                          coffeeco.CoffeeLover
	FreeDrinksAvailable                  int
	RemainingDrinkPurchaseUntilFreeDrink int
}

func (c *CoffeBux) AddStamp() {
	if c.RemainingDrinkPurchaseUntilFreeDrink == 1 {
		c.RemainingDrinkPurchaseUntilFreeDrink = 10
		c.FreeDrinksAvailable += 1
	} else {
		c.RemainingDrinkPurchaseUntilFreeDrink--
	}
}

func (c *CoffeBux) Pay(ctx context.Context, products []coffeeco.Product) error {
	lp := len(products)
	if lp == 0 {
		return errors.New("nothing to buy")
	}
	if c.FreeDrinksAvailable < lp {
		return fmt.Errorf("not enough coffeeBux to cover entire purchase. Have %d, need %d", len(products),
			c.FreeDrinksAvailable)
	}

	c.FreeDrinksAvailable = c.FreeDrinksAvailable - lp
	return nil
}
