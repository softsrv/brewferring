package provider

import (
	"context"

	"github.com/terminaldotshop/terminal-sdk-go"
	"github.com/terminaldotshop/terminal-sdk-go/option"
)

type OAuthConfig struct {
	Issuer        string   `json:"issuer"`
	AuthEndpoint  string   `json:"authorization_endpoint"`
	TokenEndpoint string   `json:"token_endpoint"`
	JWKSEndpoint  string   `json:"jwks_uri"`
	ResponseTypes []string `json:"response_types_supported"`
	ClientID      string
	ClientSecret  string
	RedirectURI   string
}

// The following models are from terminal
type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
}

type Order struct {
	ID          string
	Total       float64
	TrackingURL string
	ProductID   string
}

type Address struct {
	ID     string
	Name   string
	Street string
}

type Card struct {
	ID       string
	Brand    string
	ExpMonth int64
	ExpYear  int64
}

type Profile struct {
	ID    string
	Email string
}

type Provider interface {
	GetProfile(context.Context) (Profile, error)
	ListProducts(context.Context) ([]Product, error)
	ListCards(context.Context) ([]Card, error)
	ListAddresses(context.Context) ([]Address, error)
	ListOrders(context.Context) ([]Order, error)
	CreateOrder(ctx context.Context, productID string, cardID string, addressID string) (orderID string, err error)
}

type TerminalProvider struct {
	client *terminal.Client
}

func NewProvider(accessToken string) Provider {
	return &TerminalProvider{
		client: terminal.NewClient(option.WithBearerToken(accessToken)),
	}
}

func (t TerminalProvider) GetProfile(ctx context.Context) (Profile, error) {
	response, err := t.client.Profile.Me(ctx)
	if err != nil {
		return Profile{}, err
	}
	profile := Profile{
		ID:    response.Data.User.ID,
		Email: response.Data.User.Email,
	}
	return profile, err
}
func (t TerminalProvider) ListProducts(ctx context.Context) ([]Product, error) {

	response, err := t.client.Product.List(ctx)
	if err != nil {
		return nil, err
	}

	var products []Product
	for _, p := range response.Data {
		// since cron is already a subscription service, we dont' support it
		if p.Name == "cron" {
			continue
		}
		products = append(products, Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			// TODO: not sure about this. what are the variants??
			Price: float64(p.Variants[0].Price) / 100, // Convert cents to dollars
		})
	}

	return products, nil

}

func (t TerminalProvider) ListCards(ctx context.Context) ([]Card, error) {
	response, err := t.client.Card.List(ctx)
	if err != nil {
		return nil, err
	}

	var cards []Card
	for _, c := range response.Data {

		cards = append(cards, Card{
			ID:       c.ID,
			Brand:    c.Brand,
			ExpMonth: c.Expiration.Month,
			ExpYear:  c.Expiration.Year,
		})
	}

	return cards, nil
}
func (t TerminalProvider) ListAddresses(ctx context.Context) ([]Address, error) {
	response, err := t.client.Address.List(ctx)
	if err != nil {
		return nil, err
	}

	var addresses []Address
	for _, a := range response.Data {

		addresses = append(addresses, Address{
			ID:     a.ID,
			Name:   a.Name,
			Street: a.Street1,
		})
	}

	return addresses, nil
}
func (t TerminalProvider) ListOrders(ctx context.Context) ([]Order, error) {
	response, err := t.client.Order.List(ctx)
	if err != nil {
		return nil, err
	}

	var orders []Order
	for _, o := range response.Data {

		orders = append(orders, Order{
			ID:          o.ID,
			Total:       float64(o.Amount.Subtotal + o.Amount.Shipping),
			TrackingURL: o.Tracking.URL,
			ProductID:   o.Items[0].ProductVariantID,
		})
	}

	return orders, nil
}
func (t TerminalProvider) CreateOrder(ctx context.Context, productID string, cardID string, addressID string) (orderID string, err error) {

	variant := make(map[string]int64)
	variant[productID] = 1

	order, err := t.client.Order.New(ctx, terminal.OrderNewParams{
		AddressID: terminal.String(addressID),
		CardID:    terminal.String(cardID),
		Variants:  terminal.F(variant),
	})
	if err != nil {
		return "", err
	}

	return order.Data, err
}
