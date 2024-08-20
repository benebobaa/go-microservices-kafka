package handler

import (
	"fmt"
	"sync"
)

func NewUserHandler() *UserHandler {
	db := make(map[string]User)

	user1 := User{
		ID:            "USER-001",
		Username:      "beneboba",
		AccountBankID: "AC-001",
		Email:         "bene@gmail.com",
	}

	user2 := User{
		ID:            "USER-002",
		Username:      "archlinux",
		AccountBankID: "AC-002",
		Email:         "bene@beneboba.me",
	}

	user3 := User{
		ID:            "USER-003",
		Username:      "astrovim",
		AccountBankID: "AC-003",
		Email:         "admin@beneboba.me",
	}

	db[user1.ID] = user1
	db[user2.ID] = user2
	db[user3.ID] = user3

	return &UserHandler{
		db:    db,
		mutex: &sync.RWMutex{},
	}
}

func NewProductHandler() *ProductHandler {
	db := make(map[string]Product)

	product1 := Product{
		ID:    "P-001",
		Name:  "Product 1",
		Stock: 10,
		Price: 1000,
	}

	product2 := Product{
		ID:    "P-002",
		Name:  "Product 2",
		Stock: 10,
		Price: 2000,
	}

	product3 := Product{
		ID:    "P-003",
		Name:  "Product 3",
		Stock: 3,
		Price: 3000,
	}

	db[product1.ID] = product1
	db[product2.ID] = product2
	db[product3.ID] = product3

	return &ProductHandler{
		db:    db,
		mutex: &sync.RWMutex{},
	}
}

func NewPaymentHandler() *PaymentHandler {

	dbt := make(map[string]Transaction)
	dbb := make(map[string]Balance)

	h := &PaymentHandler{
		dbT:   dbt,
		dbB:   dbb,
		mutex: &sync.RWMutex{},
	}

	if err := h.LoadData(); err != nil {
		fmt.Println("Error loading data:", err)
	}

	return h
}
