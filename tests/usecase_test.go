package tests

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/SerzhLimon/PaymentService/internal/models"
	"github.com/SerzhLimon/PaymentService/internal/usecase"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) WalletTransactionDeposit(walletID uuid.UUID, amount int64) error {
	args := m.Called(walletID, amount)
	return args.Error(0)
}

func (m *MockRepository) WalletTransactionWithdraw(walletID uuid.UUID, amount int64) error {
	args := m.Called(walletID, amount)
	return args.Error(0)
}

func (m *MockRepository) GetBalance(walletID uuid.UUID) (models.GetBalanceResponse, error) {
	args := m.Called(walletID)
	return args.Get(0).(models.GetBalanceResponse), args.Error(1)
}

func (m *MockRepository) CreateWallet(walletID uuid.UUID) error {
	return nil
}

func TestWalletTransaction_Success_Deposit(t *testing.T) {
	mockRepo := new(MockRepository)
	usecase := usecase.NewUsecase(mockRepo)

	amount := int64(100)
	data := models.WalletTransaction{
		WalletID:  "7b7ad84a-cb3e-4734-8e80-98aef40122d2",
		Operation: "DEPOSIT",
		Amount:    amount,
	}

	mockRepo.On("WalletTransactionDeposit", mock.Anything, mock.Anything).Return(nil)

	err := usecase.WalletTransaction(data)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestWalletTransaction_Success_Withdraw(t *testing.T) {
	mockRepo := new(MockRepository)
	usecase := usecase.NewUsecase(mockRepo)

	amount := int64(50)
	data := models.WalletTransaction{
		WalletID:  "7b7ad84a-cb3e-4734-8e80-98aef40122d2",
		Operation: "WITHDRAW",
		Amount:    amount,
	}

	mockRepo.On("WalletTransactionWithdraw", mock.Anything, mock.Anything).Return(nil)

	err := usecase.WalletTransaction(data)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestWalletTransaction_InvalidUUID(t *testing.T) {
	mockRepo := new(MockRepository)
	usecase := usecase.NewUsecase(mockRepo)

	data := models.WalletTransaction{
		WalletID:  "invalid-uuid",
		Operation: "DEPOSIT",
		Amount:    100,
	}

	err := usecase.WalletTransaction(data)
	assert.Error(t, err)
}

func TestWalletTransaction_InvalidAmount(t *testing.T) {
	mockRepo := new(MockRepository)
	usecase := usecase.NewUsecase(mockRepo)

	data := models.WalletTransaction{
		WalletID:  "7b7ad84a-cb3e-4734-8e80-98aef40122d2",
		Operation: "DEPOSIT",
		Amount:    -100,
	}

	err := usecase.WalletTransaction(data)
	assert.Error(t, err)
}

func TestWalletTransaction_UnknownOperation(t *testing.T) {
	mockRepo := new(MockRepository)
	usecase := usecase.NewUsecase(mockRepo)

	data := models.WalletTransaction{
		WalletID:  "7b7ad84a-cb3e-4734-8e80-98aef40122d2",
		Operation: "TRANSFER",
		Amount:    100,
	}

	err := usecase.WalletTransaction(data)
	assert.Error(t, err)
}

func TestGetBalance_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	usecase := usecase.NewUsecase(mockRepo)

	walletID := "7b7ad84a-cb3e-4734-8e80-98aef40122d2"
	mockRepo.On("GetBalance", mock.Anything).Return(models.GetBalanceResponse{Amount: 100.0}, nil)

	result, err := usecase.GetBalance(walletID)
	assert.NoError(t, err)
	assert.Equal(t, 100.0, result.Amount)
	mockRepo.AssertExpectations(t)
}

func TestGetBalance_InvalidUUID(t *testing.T) {
	mockRepo := new(MockRepository)
	usecase := usecase.NewUsecase(mockRepo)

	// Test invalid UUID for GetBalance
	_, err := usecase.GetBalance("invalid-uuid")
	assert.Error(t, err)
}

func TestGetBalance_Failure(t *testing.T) {
	mockRepo := new(MockRepository)
	usecase := usecase.NewUsecase(mockRepo)

	// Test failure while retrieving balance
	walletID := "7b7ad84a-cb3e-4734-8e80-98aef40122d2"
	mockRepo.On("GetBalance", mock.Anything).Return(models.GetBalanceResponse{}, errors.New("failed to get balance"))

	_, err := usecase.GetBalance(walletID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get balance")
	mockRepo.AssertExpectations(t)
}
