package products

type createProduct struct {
	Name string `json:"name"`
	PriceInCents int `json:"priceInCents"`
	Quantity  int `json:"quantity"`
}

type updateQuantity struct {
	Quantity int `json:"quantity"`
}