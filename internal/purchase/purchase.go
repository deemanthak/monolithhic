package purchase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"

	coffeeco "github.com/deemanthak/monolithhic/internal"
	"github.com/deemanthak/monolithhic/internal/loyality"
	"github.com/deemanthak/monolithhic/internal/payment"
	"github.com/deemanthak/monolithhic/internal/store"
)

type Purchase struct {
	id                 uuid.UUID          `bson:"id"`
	Store              store.Store        `bson:"store"`
	ProductsToPurchase []coffeeco.Product `bson:"productsToPurchase"`
	total              money.Money        `bson:"total"`
	PaymentMeans       payment.Means      `bson:"paymentMeans"`
	timeOfPurchase     time.Time          `bson:"timeOfPurchase"`
	CardToken          *string            `bson:"cardToken"`
}

type CardChargeService interface {
	ChargeCard(ctx context.Context, amount money.Money, cardToken string) error
}
type Service struct {
	cardService  CardChargeService
	purchaseRepo Repository
}

func (p *Purchase) validateAndEnrich() error {
	if len(p.ProductsToPurchase) == 0 {
		return errors.New("purchase must consist of at least one product")
	}

	p.total = *money.New(0, "USD")

	for _, v := range p.ProductsToPurchase {
		newTotal, _ := p.total.Add(&v.BasePrice)
		p.total = *newTotal
	}

	if p.total.IsZero() {
		return errors.New("likely mistake; purchase should never be 0. Please validate")
	}

	p.id = uuid.New()
	p.timeOfPurchase = time.Now()

	return nil
}

func (s Service) CompletePurchase(ctx context.Context, purchase *Purchase, coffeeBuxCard *loyality.CoffeBux) error {
	if err := purchase.validateAndEnrich(); err != nil {
		return err
	}

	switch purchase.PaymentMeans {
	case payment.MEANS_CARD:
		if err := s.cardService.ChargeCard(ctx, purchase.total, *purchase.CardToken); err != nil {
			return errors.New("card charge failed, cancelling purchase")
		}
	case payment.MEANS_CASH:
	// TODO: For the reader to add :)
	case payment.MEANS_COFFEEBUX:
		if err := coffeeBuxCard.Pay(ctx, purchase.ProductsToPurchase); err != nil {
			return fmt.Errorf("failed to charge loyalty card: %w", err)
		}
	default:
		return errors.New("unknown payment type")
	}

	if err := s.purchaseRepo.Store(ctx, *purchase); err != nil {
		return errors.New("failed to store purchase")
	}

	if coffeeBuxCard != nil {
		coffeeBuxCard.AddStamp()
	}

	return nil
}
