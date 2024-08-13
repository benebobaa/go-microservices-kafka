package dto

type Product struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Stock int    `json:"stock"`
	Price int    `json:"price"`
}

type ProductResponse struct {
	Product Product `json:"product"`
	Amount  float64 `json:"amount"`
}

type ProductRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
