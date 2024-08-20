package http

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"order-svc/internal/dto"
	mockdb "order-svc/internal/repository/mock"
	"order-svc/internal/repository/sqlc"
	"order-svc/internal/usecase"
	"order-svc/pkg/producer"
	"testing"
)

func init() {
	orchestraProducer, _ = producer.NewKafkaProducer([]string{"localhost:29092"}, "orchestra-topic-test")
}

func TestBankRegisHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("RegisterBankAccount", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)                                 // Mock the store directly
		uc := usecase.NewBankRegistrationUsecase(store, orchestraProducer) // Pass the mocked store to the usecase
		handler := NewBankRegisHandler(uc)

		testCases := []struct {
			name           string
			setupMocks     func()
			setupRequest   func() *http.Request
			expectedStatus int
		}{
			{
				name: "successful bank registration",
				setupMocks: func() {
					expectedResponse := sqlc.BankAccountRegistration{
						ID:         1,
						CustomerID: "1",
						Username:   "testuser",
						Deposit:    0,
						Email:      "testuser@example.com",
						Status:     "complete",
					}
					store.EXPECT().CreateBankAccountRegistration(gomock.Any(), gomock.Any()).Return(expectedResponse, nil)
				},
				setupRequest: func() *http.Request {
					req := &dto.BankRegistrationRequest{
						Username: "testuser",
						Email:    "testuser@example.com",
						Deposit:  2000.0,
					}
					reqBody, _ := json.Marshal(req)
					httpReq, _ := http.NewRequest(http.MethodPost, "/bank", bytes.NewBuffer(reqBody))
					httpReq.Header.Set("Content-Type", "application/json")
					return httpReq
				},
				expectedStatus: http.StatusOK,
			},
			{
				name:       "invalid request body",
				setupMocks: func() {},
				setupRequest: func() *http.Request {
					httpReq, _ := http.NewRequest(http.MethodPost, "/bank", bytes.NewBufferString("invalid json"))
					httpReq.Header.Set("Content-Type", "application/json")
					return httpReq
				},
				expectedStatus: http.StatusBadRequest,
			},
			{
				name:       "validation error",
				setupMocks: func() {},
				setupRequest: func() *http.Request {
					req := &dto.BankRegistrationRequest{
						Username: "",
						Email:    "",
						Deposit:  500.0,
					}
					reqBody, _ := json.Marshal(req)
					httpReq, _ := http.NewRequest(http.MethodPost, "/bank", bytes.NewBuffer(reqBody))
					httpReq.Header.Set("Content-Type", "application/json")
					return httpReq
				},
				expectedStatus: http.StatusBadRequest,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tc.setupMocks()
				httpReq := tc.setupRequest()

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httpReq

				handler.RegisterBankAccount(c)

				assert.Equal(t, tc.expectedStatus, w.Code)
			})
		}
	})
}
