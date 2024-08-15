package dto

type Product struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Stock int    `json:"stock"`
	Price int    `json:"price"`
}

type ProductResponse struct {
	Id       string  `json:"product_id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Amount   float64 `json:"amount"`
}

type ProductRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
