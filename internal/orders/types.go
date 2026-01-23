package orders

type orderItem struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

type createOrderRequest struct {
	CustomerID int `json:"customerId"`
	Items 		[]orderItem `json:"items"`
}