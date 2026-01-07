package payment

import (
	"errors"
	"fmt"
	"log"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/checkout/session"
)

type PaymentClient interface {
	CreatePayment(amount float64, userId uint, orderId uint) (*stripe.CheckoutSession, error)
	GetPaymentStatus(pId string) (*stripe.CheckoutSession, error)
}

type payment struct {
	stripeSecretKey string
	succeccUrl      string
	cancelUrl       string
}

func (p payment) CreatePayment(amount float64, userId uint, orderId uint) (*stripe.CheckoutSession, error) {
	stripe.Key = p.stripeSecretKey
	amountInCents := amount * 100

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					UnitAmount: stripe.Int64(int64(amountInCents)),
					Currency:   stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Electronics"),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(p.succeccUrl),
		CancelURL:  stripe.String(p.cancelUrl),
	}

	params.AddMetadata("order_id", fmt.Sprintf("%d", orderId))
	params.AddMetadata("user_id", fmt.Sprintf("%d", userId))

	session, err := session.New(params)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		return nil, errors.New("payment create session failed")
	}

	return session, nil
}

func (p payment) GetPaymentStatus(pId string) (*stripe.CheckoutSession, error) {
	stripe.Key = p.stripeSecretKey
	session, err := session.Get(pId, nil)
	if err != nil {
		log.Printf("Error getting session: %v", err)
		return nil, errors.New("payment get session failed")
	}

	return session, nil
}

func NewPaymentClient(stripeSecretKey, succeccUrl, cancelUrl string) PaymentClient {
	return &payment{
		stripeSecretKey: stripeSecretKey,
		succeccUrl:      succeccUrl,
		cancelUrl:       cancelUrl,
	}
}
