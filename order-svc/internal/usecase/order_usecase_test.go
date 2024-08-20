package usecase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"order-svc/internal/dto"
	"order-svc/internal/dto/event"
	mockdb "order-svc/internal/repository/mock"
	"order-svc/internal/repository/sqlc"
	"order-svc/pkg/producer"
	"testing"
)

var orchestraProducer *producer.KafkaProducer

func init() {
	orchestraProducer, _ = producer.NewKafkaProducer([]string{"localhost:29092"}, "orchestra-topic-test")
}

func TestOrderUsecase_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	uc := NewOrderUsecase(store, orchestraProducer)

	ctx := context.Background()

	testCases := []struct {
		name        string
		setupMocks  func()
		input       *dto.OrderRequest
		expected    *sqlc.Order
		expectedErr error
	}{
		{
			name: "Successful order creation",
			setupMocks: func() {
				store.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(sqlc.Order{
					ID:         1,
					RefID:      gomock.Any().String(),
					CustomerID: "customer-001",
					Username:   "testuser",
					ProductID:  "product-001",
					Quantity:   2,
					Status:     dto.PROCESSING.String(),
				}, nil)
			},
			input: &dto.OrderRequest{
				CustomerID: "customer-001",
				Username:   "testuser",
				ProductID:  "product-001",
				Quantity:   2,
			},
			expected: &sqlc.Order{
				ID:         1,
				RefID:      gomock.Any().String(),
				CustomerID: "customer-001",
				Username:   "testuser",
				ProductID:  "product-001",
				Quantity:   2,
				Status:     dto.PROCESSING.String(),
			},
			expectedErr: nil,
		},
		{
			name: "Error creating order",
			setupMocks: func() {
				store.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(sqlc.Order{}, errors.New("database error"))
			},
			input: &dto.OrderRequest{
				CustomerID: "customer-002",
				Username:   "testuser2",
				ProductID:  "product-002",
				Quantity:   3,
			},
			expected:    nil,
			expectedErr: errors.New("database error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			result, err := uc.CreateOrder(ctx, tc.input)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.ID, result.ID)
				assert.Equal(t, tc.expected.CustomerID, result.CustomerID)
				assert.Equal(t, tc.expected.Username, result.Username)
				assert.Equal(t, tc.expected.ProductID, result.ProductID)
				assert.Equal(t, tc.expected.Quantity, result.Quantity)
				assert.Equal(t, tc.expected.Status, result.Status)
			}
		})
	}
}

func TestOrderUsecase_CancelOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	uc := NewOrderUsecase(store, orchestraProducer)

	ctx := context.Background()

	testCases := []struct {
		name        string
		setupMocks  func()
		input       *dto.OrderCancelRequest
		expected    *sqlc.Order
		expectedErr error
	}{
		{
			name: "Successful order cancellation",
			setupMocks: func() {
				store.EXPECT().FindOrderByID(gomock.Any(), int32(1)).Return(sqlc.Order{
					ID:         1,
					RefID:      "order-001",
					CustomerID: "customer-001",
					Username:   "testuser",
					ProductID:  "product-001",
					Quantity:   2,
					Status:     dto.COMPLETE.String(),
					Amount:     sql.NullFloat64{Float64: 100, Valid: true},
				}, nil)
				store.EXPECT().UpdateOrder(gomock.Any(), gomock.Any()).Return(sqlc.Order{
					ID:         1,
					RefID:      "order-001",
					CustomerID: "customer-001",
					Username:   "testuser",
					ProductID:  "product-001",
					Quantity:   2,
					Status:     dto.CANCEL_PROCESSING.String(),
					Amount:     sql.NullFloat64{Float64: 100, Valid: true},
				}, nil)
			},
			input: &dto.OrderCancelRequest{
				OrderID:  1,
				Username: "testuser",
			},
			expected: &sqlc.Order{
				ID:         1,
				RefID:      "order-001",
				CustomerID: "customer-001",
				Username:   "testuser",
				ProductID:  "product-001",
				Quantity:   2,
				Status:     dto.CANCEL_PROCESSING.String(),
				Amount:     sql.NullFloat64{Float64: 100, Valid: true},
			},
			expectedErr: nil,
		},
		{
			name: "Order not found",
			setupMocks: func() {
				store.EXPECT().FindOrderByID(gomock.Any(), int32(2)).Return(sqlc.Order{}, sql.ErrNoRows)
			},
			input: &dto.OrderCancelRequest{
				OrderID:  2,
				Username: "testuser",
			},
			expected:    nil,
			expectedErr: sql.ErrNoRows,
		},
		{
			name: "Cannot cancel non-complete order",
			setupMocks: func() {
				store.EXPECT().FindOrderByID(gomock.Any(), int32(3)).Return(sqlc.Order{
					ID:         3,
					RefID:      "order-003",
					CustomerID: "customer-003",
					Username:   "testuser",
					ProductID:  "product-003",
					Quantity:   2,
					Status:     dto.PROCESSING.String(),
				}, nil)
			},
			input: &dto.OrderCancelRequest{
				OrderID:  3,
				Username: "testuser",
			},
			expected:    nil,
			expectedErr: ErrCannotCancelOrder,
		},
		{
			name: "Unauthorized cancel attempt",
			setupMocks: func() {
				store.EXPECT().FindOrderByID(gomock.Any(), int32(4)).Return(sqlc.Order{
					ID:         4,
					RefID:      "order-004",
					CustomerID: "customer-004",
					Username:   "testuser",
					ProductID:  "product-004",
					Quantity:   2,
					Status:     dto.COMPLETE.String(),
				}, nil)
			},
			input: &dto.OrderCancelRequest{
				OrderID:  4,
				Username: "otheruser",
			},
			expected:    nil,
			expectedErr: ErrUnauthorizeCancelOrder,
		},
		{
			name: "Error updating order",
			setupMocks: func() {
				store.EXPECT().FindOrderByID(gomock.Any(), int32(5)).Return(sqlc.Order{
					ID:         5,
					RefID:      "order-005",
					CustomerID: "customer-005",
					Username:   "testuser",
					ProductID:  "product-005",
					Quantity:   2,
					Status:     dto.COMPLETE.String(),
					Amount:     sql.NullFloat64{Float64: 100, Valid: true},
				}, nil)
				store.EXPECT().UpdateOrder(gomock.Any(), gomock.Any()).Return(sqlc.Order{}, errors.New("update error"))
			},
			input: &dto.OrderCancelRequest{
				OrderID:  5,
				Username: "testuser",
			},
			expected:    nil,
			expectedErr: errors.New("update error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			result, err := uc.CancelOrder(ctx, tc.input)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.ID, result.ID)
				assert.Equal(t, tc.expected.RefID, result.RefID)
				assert.Equal(t, tc.expected.CustomerID, result.CustomerID)
				assert.Equal(t, tc.expected.Username, result.Username)
				assert.Equal(t, tc.expected.ProductID, result.ProductID)
				assert.Equal(t, tc.expected.Quantity, result.Quantity)
				assert.Equal(t, tc.expected.Status, result.Status)
				assert.Equal(t, tc.expected.Amount, result.Amount)
			}
		})
	}
}

func TestOrderUsecase_UpdateOrderMessaging(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	uc := NewOrderUsecase(store, orchestraProducer)

	ctx := context.Background()

	testCases := []struct {
		name        string
		setupMocks  func()
		input       event.GlobalEvent[dto.OrderUpdateRequest, any]
		expectedErr error
	}{
		{
			name: "Successful update - ORDER_PROCESS",
			setupMocks: func() {
				store.EXPECT().FindOrderByRefID(gomock.Any(), "ref-001").Return(sqlc.Order{
					RefID:    "ref-001",
					Amount:   sql.NullFloat64{Float64: 100, Valid: true},
					Quantity: 2,
				}, nil)
				store.EXPECT().UpdateOrder(gomock.Any(), gomock.Any()).Return(sqlc.Order{
					RefID:    "ref-001",
					Amount:   sql.NullFloat64{Float64: 150, Valid: true},
					Quantity: 3,
					Status:   "PROCESSING",
				}, nil)
			},
			input: event.GlobalEvent[dto.OrderUpdateRequest, any]{
				EventType: event.ORDER_PROCESS.String(),
				Payload: event.BasePayload[dto.OrderUpdateRequest, any]{
					Request: dto.OrderUpdateRequest{
						RefID:    "ref-001",
						Amount:   150,
						Quantity: 3,
						Status:   "PROCESSING",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Successful update - non-ORDER_PROCESS",
			setupMocks: func() {
				store.EXPECT().FindOrderByRefID(gomock.Any(), "ref-002").Return(sqlc.Order{
					RefID:    "ref-002",
					Amount:   sql.NullFloat64{Float64: 200, Valid: true},
					Quantity: 4,
				}, nil)
				store.EXPECT().UpdateOrder(gomock.Any(), gomock.Any()).Return(sqlc.Order{
					RefID:    "ref-002",
					Amount:   sql.NullFloat64{Float64: 200, Valid: true},
					Quantity: 4,
					Status:   "COMPLETED",
				}, nil)
			},
			input: event.GlobalEvent[dto.OrderUpdateRequest, any]{
				EventType: "OTHER_EVENT",
				Payload: event.BasePayload[dto.OrderUpdateRequest, any]{
					Request: dto.OrderUpdateRequest{
						RefID:  "ref-002",
						Status: "COMPLETED",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Error finding order",
			setupMocks: func() {
				store.EXPECT().FindOrderByRefID(gomock.Any(), "ref-003").Return(sqlc.Order{}, errors.New("order not found"))
			},
			input: event.GlobalEvent[dto.OrderUpdateRequest, any]{
				EventType: event.ORDER_PROCESS.String(),
				Payload: event.BasePayload[dto.OrderUpdateRequest, any]{
					Request: dto.OrderUpdateRequest{
						RefID: "ref-003",
					},
				},
			},
			expectedErr: errors.New("order not found"),
		},
		{
			name: "Error updating order",
			setupMocks: func() {
				store.EXPECT().FindOrderByRefID(gomock.Any(), "ref-004").Return(sqlc.Order{
					RefID:    "ref-004",
					Amount:   sql.NullFloat64{Float64: 300, Valid: true},
					Quantity: 5,
				}, nil)
				store.EXPECT().UpdateOrder(gomock.Any(), gomock.Any()).Return(sqlc.Order{}, errors.New("update failed"))
			},
			input: event.GlobalEvent[dto.OrderUpdateRequest, any]{
				EventType: event.ORDER_PROCESS.String(),
				Payload: event.BasePayload[dto.OrderUpdateRequest, any]{
					Request: dto.OrderUpdateRequest{
						RefID:    "ref-004",
						Amount:   350,
						Quantity: 6,
						Status:   "PROCESSING",
					},
				},
			},
			expectedErr: errors.New("update failed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			err := uc.UpdateOrderMessaging(ctx, tc.input)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
