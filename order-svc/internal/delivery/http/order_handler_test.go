package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"order-svc/internal/dto"
	"order-svc/internal/middleware"
	mockdb "order-svc/internal/repository/mock"
	"order-svc/internal/repository/sqlc"
	"order-svc/internal/usecase"
	"order-svc/pkg"
	"order-svc/pkg/producer"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var orchestraProducer *producer.KafkaProducer

func init() {
	orchestraProducer, _ = producer.NewKafkaProducer([]string{"localhost:29092"}, "orchestra-topic-test")
}

func TestOrderHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("CreateOrder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)
		uc := usecase.NewOrderUsecase(store, orchestraProducer)
		handler := NewOrderHandler(uc)

		testCases := []struct {
			name           string
			setupMocks     func()
			setupRequest   func() (*http.Request, *pkg.UserInfo)
			expectedStatus int
		}{
			{
				name: "successful order creation",
				setupMocks: func() {
					expectedOrder := &sqlc.Order{
						ID:         1,
						RefID:      "test-ref",
						CustomerID: "user-id",
						Username:   "testuser",
						ProductID:  "p-001",
						Quantity:   2,
						Status:     dto.PROCESSING.String(),
					}
					store.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(*expectedOrder, nil)
				},
				setupRequest: func() (*http.Request, *pkg.UserInfo) {
					req := &dto.OrderRequest{
						ProductID: "p-001",
						Quantity:  2,
					}
					reqBody, _ := json.Marshal(req)
					httpReq, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(reqBody))
					httpReq.Header.Set("Content-Type", "application/json")
					userInfo := &pkg.UserInfo{
						ID:       "user-id",
						Username: "testuser",
					}
					return httpReq, userInfo
				},
				expectedStatus: http.StatusOK,
			},
			{
				name:       "invalid request body",
				setupMocks: func() {},
				setupRequest: func() (*http.Request, *pkg.UserInfo) {
					httpReq, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString("invalid json"))
					httpReq.Header.Set("Content-Type", "application/json")
					userInfo := &pkg.UserInfo{
						ID:       "user-id",
						Username: "testuser",
					}
					return httpReq, userInfo
				},
				expectedStatus: http.StatusBadRequest,
			},
			{
				name:       "validation error",
				setupMocks: func() {},
				setupRequest: func() (*http.Request, *pkg.UserInfo) {
					req := &dto.OrderRequest{
						ProductID: "",
						Quantity:  0,
					}
					reqBody, _ := json.Marshal(req)
					httpReq, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(reqBody))
					httpReq.Header.Set("Content-Type", "application/json")
					userInfo := &pkg.UserInfo{
						ID:       "user-id",
						Username: "testuser",
					}
					return httpReq, userInfo
				},
				expectedStatus: http.StatusBadRequest,
			},
			{
				name: "usecase error",
				setupMocks: func() {
					store.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(sqlc.Order{}, errors.New("database error"))
				},
				setupRequest: func() (*http.Request, *pkg.UserInfo) {
					req := &dto.OrderRequest{
						ProductID: "p-001",
						Quantity:  2,
					}
					reqBody, _ := json.Marshal(req)
					httpReq, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(reqBody))
					httpReq.Header.Set("Content-Type", "application/json")
					userInfo := &pkg.UserInfo{
						ID:       "user-id",
						Username: "testuser",
					}
					return httpReq, userInfo
				},
				expectedStatus: http.StatusInternalServerError,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tc.setupMocks()
				httpReq, userInfo := tc.setupRequest()

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httpReq
				c.Set(middleware.ClaimsKey, userInfo)

				handler.CreateOrder(c)

				assert.Equal(t, tc.expectedStatus, w.Code)
			})
		}
	})

	t.Run("CancelOrder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)
		uc := usecase.NewOrderUsecase(store, orchestraProducer)
		handler := NewOrderHandler(uc)

		testCases := []struct {
			name           string
			setupMocks     func()
			setupRequest   func() (*http.Request, *pkg.UserInfo)
			expectedStatus int
		}{
			{
				name: "successful order cancellation",
				setupMocks: func() {
					// Mock FindOrderByID
					existingOrder := sqlc.Order{
						ID:         1,
						RefID:      "test-ref",
						CustomerID: "user-id",
						Username:   "testuser",
						ProductID:  "p-001",
						Quantity:   2,
						Status:     dto.COMPLETE.String(),
						Amount:     sql.NullFloat64{Float64: 100, Valid: true},
					}
					store.EXPECT().FindOrderByID(gomock.Any(), int32(1)).Return(existingOrder, nil)

					// Mock UpdateOrder
					updatedOrder := sqlc.Order{
						ID:         1,
						RefID:      "test-ref",
						CustomerID: "user-id",
						Username:   "testuser",
						ProductID:  "p-001",
						Quantity:   2,
						Status:     dto.CANCEL_PROCESSING.String(),
						Amount:     sql.NullFloat64{Float64: 100, Valid: true},
					}
					store.EXPECT().UpdateOrder(gomock.Any(), gomock.Any()).Return(updatedOrder, nil)

				},
				setupRequest: func() (*http.Request, *pkg.UserInfo) {
					httpReq, _ := http.NewRequest(http.MethodPost, "/orders/cancel?order_id=1", nil)
					userInfo := &pkg.UserInfo{
						ID:       "user-id",
						Username: "testuser",
					}
					return httpReq, userInfo
				},
				expectedStatus: http.StatusOK,
			},
			{
				name: "order not found",
				setupMocks: func() {
					store.EXPECT().FindOrderByID(gomock.Any(), int32(1)).Return(sqlc.Order{}, sql.ErrNoRows)
				},
				setupRequest: func() (*http.Request, *pkg.UserInfo) {
					httpReq, _ := http.NewRequest(http.MethodPost, "/orders/cancel?order_id=1", nil)
					userInfo := &pkg.UserInfo{
						ID:       "user-id",
						Username: "testuser",
					}
					return httpReq, userInfo
				},
				expectedStatus: http.StatusInternalServerError,
			},
			{
				name: "cannot cancel completed order",
				setupMocks: func() {
					existingOrder := sqlc.Order{
						ID:         1,
						RefID:      "test-ref",
						CustomerID: "user-id",
						Username:   "testuser",
						ProductID:  "p-001",
						Quantity:   2,
						Status:     dto.PROCESSING.String(),
					}
					store.EXPECT().FindOrderByID(gomock.Any(), int32(1)).Return(existingOrder, nil)
				},
				setupRequest: func() (*http.Request, *pkg.UserInfo) {
					httpReq, _ := http.NewRequest(http.MethodPost, "/orders/cancel?order_id=1", nil)
					userInfo := &pkg.UserInfo{
						ID:       "user-id",
						Username: "testuser",
					}
					return httpReq, userInfo
				},
				expectedStatus: http.StatusBadRequest,
			},
			{
				name: "unauthorized cancel",
				setupMocks: func() {
					existingOrder := sqlc.Order{
						ID:         1,
						RefID:      "test-ref",
						CustomerID: "user-id",
						Username:   "otheruser",
						ProductID:  "p-001",
						Quantity:   2,
						Status:     dto.COMPLETE.String(),
					}
					store.EXPECT().FindOrderByID(gomock.Any(), int32(1)).Return(existingOrder, nil)
				},
				setupRequest: func() (*http.Request, *pkg.UserInfo) {
					httpReq, _ := http.NewRequest(http.MethodPost, "/orders/cancel?order_id=1", nil)
					userInfo := &pkg.UserInfo{
						ID:       "user-id",
						Username: "testuser",
					}
					return httpReq, userInfo
				},
				expectedStatus: http.StatusUnauthorized,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tc.setupMocks()
				httpReq, userInfo := tc.setupRequest()

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httpReq
				c.Set(middleware.ClaimsKey, userInfo)

				handler.CancelOrder(c)

				assert.Equal(t, tc.expectedStatus, w.Code)
			})
		}
	})

	t.Run("FindAll", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)
		uc := usecase.NewOrderUsecase(store, orchestraProducer)
		handler := NewOrderHandler(uc)

		testCases := []struct {
			name           string
			setupMocks     func()
			setupRequest   func() (*http.Request, *pkg.UserInfo)
			expectedStatus int
		}{
			{
				name: "successful find all orders",
				setupMocks: func() {
					orders := []sqlc.Order{
						{
							ID:         1,
							RefID:      "test-ref-1",
							CustomerID: "user-id",
							Username:   "testuser",
							ProductID:  "p-001",
							Quantity:   2,
							Status:     dto.PROCESSING.String(),
						},
						{
							ID:         2,
							RefID:      "test-ref-2",
							CustomerID: "user-id",
							Username:   "testuser",
							ProductID:  "p-002",
							Quantity:   1,
							Status:     dto.COMPLETE.String(),
						},
					}
					store.EXPECT().FindOrdersByUsername(gomock.Any(), "testuser").Return(orders, nil)
				},
				setupRequest: func() (*http.Request, *pkg.UserInfo) {
					httpReq, _ := http.NewRequest(http.MethodGet, "/orders", nil)
					userInfo := &pkg.UserInfo{
						ID:       "user-id",
						Username: "testuser",
					}
					return httpReq, userInfo
				},
				expectedStatus: http.StatusOK,
			},
			{
				name:       "missing username",
				setupMocks: func() {},
				setupRequest: func() (*http.Request, *pkg.UserInfo) {
					httpReq, _ := http.NewRequest(http.MethodGet, "/orders", nil)
					userInfo := &pkg.UserInfo{
						ID: "user-id",
					}
					return httpReq, userInfo
				},
				expectedStatus: http.StatusBadRequest,
			},
			{
				name: "internal server error",
				setupMocks: func() {
					store.EXPECT().FindOrdersByUsername(gomock.Any(), "testuser").Return(nil, errors.New("database error"))
				},
				setupRequest: func() (*http.Request, *pkg.UserInfo) {
					httpReq, _ := http.NewRequest(http.MethodGet, "/orders", nil)
					userInfo := &pkg.UserInfo{
						ID:       "user-id",
						Username: "testuser",
					}
					return httpReq, userInfo
				},
				expectedStatus: http.StatusInternalServerError,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tc.setupMocks()
				httpReq, userInfo := tc.setupRequest()

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httpReq
				c.Set(middleware.ClaimsKey, userInfo)

				handler.FindAll(c)

				assert.Equal(t, tc.expectedStatus, w.Code)
			})
		}
	})
}
