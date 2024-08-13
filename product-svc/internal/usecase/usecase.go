package usecase

import (
	"log"
	"product-svc/http_client"
	"product-svc/internal/dto"
)

type Usecase struct {
	userClient *http_client.UserClient
}

func NewUsecase(userClient *http_client.UserClient) *Usecase {
	return &Usecase{userClient: userClient}
}

// func (u *Usecase) ValidateUser(request *dto.UserValidateRequest) (*dto.UserValidateResponse, error) {
//
// 	var response dto.UserValidateResponse
//
// 	err := u.userClient.Call("POST", request, &response)
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	response.Status = "success"
// 	response.Message = "User validated successfully"
//
// 	return &response, nil
// }

func (u *Usecase) ReserveProduct() (*dto.Product, error) {

	log.Println("Reserving product")

	return &dto.Product{
		ID:    1,
		Name:  "Product 1",
		Price: 100,
	}, nil
}
