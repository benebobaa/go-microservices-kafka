package usecase

import (
	"context"
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

var bankRegisProducer *producer.KafkaProducer

func init() {
	bankRegisProducer, _ = producer.NewKafkaProducer([]string{"localhost:29092"}, "bank-regis-topic-test")
}

func TestBankRegistrationUsecase_RegisterBankAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	uc := NewBankRegistrationUsecase(store, bankRegisProducer)

	ctx := context.Background()

	testCases := []struct {
		name        string
		setupMocks  func()
		input       *dto.BankRegistrationRequest
		expected    *sqlc.BankAccountRegistration
		expectedErr error
	}{
		{
			name: "Successful bank account registration",
			setupMocks: func() {
				store.EXPECT().CreateBankAccountRegistration(gomock.Any(), gomock.Any()).Return(sqlc.BankAccountRegistration{
					CustomerID: "customer-001",
					Username:   "testuser",
					Email:      "test@example.com",
					Status:     dto.PROCESSING.String(),
					Deposit:    1000,
				}, nil)
			},
			input: &dto.BankRegistrationRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Deposit:  1000,
			},
			expected: &sqlc.BankAccountRegistration{
				CustomerID: "customer-001",
				Username:   "testuser",
				Email:      "test@example.com",
				Status:     dto.PROCESSING.String(),
				Deposit:    1000,
			},
			expectedErr: nil,
		},
		{
			name: "Error creating bank account registration",
			setupMocks: func() {
				store.EXPECT().CreateBankAccountRegistration(gomock.Any(), gomock.Any()).Return(sqlc.BankAccountRegistration{}, errors.New("database error"))
			},
			input: &dto.BankRegistrationRequest{
				Username: "testuser2",
				Email:    "test2@example.com",
				Deposit:  500,
			},
			expected:    nil,
			expectedErr: errors.New("database error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			result, err := uc.RegisterBankAccount(ctx, tc.input)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.CustomerID, result.CustomerID)
				assert.Equal(t, tc.expected.Username, result.Username)
				assert.Equal(t, tc.expected.Email, result.Email)
				assert.Equal(t, tc.expected.Status, result.Status)
				assert.Equal(t, tc.expected.Deposit, result.Deposit)
			}
		})
	}
}

func TestBankRegistrationUsecase_UpdateBankRegistrationMessaging(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	uc := NewBankRegistrationUsecase(store, bankRegisProducer)

	ctx := context.Background()

	testCases := []struct {
		name        string
		setupMocks  func()
		input       event.GlobalEvent[dto.BankRegistrationUpdate, any]
		expectedErr error
	}{
		{
			name: "Successful update",
			setupMocks: func() {
				store.EXPECT().UpdateBankAccountRegistration(gomock.Any(), gomock.Any()).Return(sqlc.BankAccountRegistration{
					CustomerID: "customer-001",
					Username:   "testuser",
					Email:      "test@example.com",
					Status:     "COMPLETED",
				}, nil)
			},
			input: event.GlobalEvent[dto.BankRegistrationUpdate, any]{
				EventType: event.BANK_ACCOUNT_REGISTRATION.String(),
				Payload: event.BasePayload[dto.BankRegistrationUpdate, any]{
					Request: dto.BankRegistrationUpdate{
						CustomerID: "customer-001",
						Username:   "testuser",
						Email:      "test@example.com",
						Status:     "COMPLETED",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Error updating bank registration",
			setupMocks: func() {
				store.EXPECT().UpdateBankAccountRegistration(gomock.Any(), gomock.Any()).Return(sqlc.BankAccountRegistration{}, errors.New("update failed"))
			},
			input: event.GlobalEvent[dto.BankRegistrationUpdate, any]{
				EventType: event.BANK_ACCOUNT_REGISTRATION.String(),
				Payload: event.BasePayload[dto.BankRegistrationUpdate, any]{
					Request: dto.BankRegistrationUpdate{
						CustomerID: "customer-002",
						Username:   "testuser2",
						Email:      "test2@example.com",
						Status:     "PROCESSING",
					},
				},
			},
			expectedErr: errors.New("update failed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			err := uc.UpdateBankRegistrationMessaging(ctx, tc.input)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
