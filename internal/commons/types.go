package commons

type HandlerEnvironment struct {
	StorageUrl string
}

type CartId string

type Address struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Street    string `json:"street"`
	HouseNr   string `json:"house_nr"`
	City      string `json:"city"`
	State     string `json:"state"`
	Country   string `json:"country"`
	Zip       string `json:"zip"`
}

type CheckoutRequest struct {
	CartId  CartId  `json:"cart_id"`
	Address Address `json:"address"`
}

type CheckoutResponse struct {
	Shipping []string `json:"shipping"`
	Payment  []string `json:"payment"`
}

type ShippingMethod string

type PaymentIntent struct {
	Id     string `json:"id"`
	Method string `json:"method"`
}

type PaymentRequest struct {
	CartId         CartId         `json:"cart_id"`
	Address        Address        `json:"address"`
	ShippingMethod ShippingMethod `json:"shipping_method"`
	PaymentIntent  PaymentIntent  `json:"payment_intent"`
}

type CartItem struct {
	ItemId string  `json:"itemId"`
	Amount int     `json:"amount"`
	Price  float64 `json:"price"`
}

type OrderRequest struct {
	Address        Address        `json:"address"`
	ShippingMethod ShippingMethod `json:"shipping_method"`
}
