package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/SerzhLimon/PaymentService/internal/models"
	"github.com/SerzhLimon/PaymentService/internal/transport"
)

type MockUsecase struct {
	mock.Mock
}

func (m *MockUsecase) WalletTransaction(req models.WalletTransaction) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockUsecase) GetBalance(walletID string) (models.GetBalanceResponse, error) {
	args := m.Called(walletID)
	return args.Get(0).(models.GetBalanceResponse), args.Error(1)
}

func setupRouter(s *transport.Server) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/api/v1/wallet", s.WalletTransaction)
	r.GET("/api/v1/wallets", s.GetBalance)
	return r
}

func TestWalletTransaction_Success(t *testing.T) {
	mockUsecase := new(MockUsecase)
	server := &transport.Server{Usecase: mockUsecase}
	router := setupRouter(server)

	requestBody := models.WalletTransaction{
		WalletID:  "7b7ad84a-cb3e-4734-8e80-98aef40122d2",
		Operation: "DEPOSIT",
		Amount:    100,
	}

	mockUsecase.On("WalletTransaction", requestBody).Return(nil)

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestWalletTransaction_BindError(t *testing.T) {
	mockUsecase := new(MockUsecase)
	server := &transport.Server{Usecase: mockUsecase}
	router := setupRouter(server)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid JSON format")
	mockUsecase.AssertNotCalled(t, "WalletTransaction")
}

func TestWalletTransaction_Failure(t *testing.T) {
	mockUsecase := new(MockUsecase)
	server := &transport.Server{Usecase: mockUsecase}
	router := setupRouter(server)

	requestBody := models.WalletTransaction{
		WalletID:  "243",
		Operation: "DEPOSIT",
		Amount:    100,
	}

	mockUsecase.On("WalletTransaction", requestBody).Return(errors.New("transaction failed"))

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "transaction failed")
	mockUsecase.AssertExpectations(t)
}

func Test_GetBalance_Success(t *testing.T) {
	mockUsecase := new(MockUsecase)
	server := &transport.Server{Usecase: mockUsecase}
	router := setupRouter(server)

	mockResult := models.GetBalanceResponse{Amount: 100.0}
	mockUsecase.On("GetBalance", "7b7ad84a-cb3e-4734-8e80-98aef40122d2").Return(mockResult, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/wallets?id=7b7ad84a-cb3e-4734-8e80-98aef40122d2", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"balance": 100.0}`, w.Body.String())
	mockUsecase.AssertExpectations(t)
}

func TestGetBalance_EmptyID(t *testing.T) {
	mockUsecase := new(MockUsecase)
	server := &transport.Server{Usecase: mockUsecase}
	router := setupRouter(server)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/wallets?id=", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "parameter 'id' is empty")
	mockUsecase.AssertNotCalled(t, "GetBalance")
}

func Test_GetBalance_Failure(t *testing.T) {
	mockUsecase := new(MockUsecase)
	server := &transport.Server{Usecase: mockUsecase}
	router := setupRouter(server)

	mockUsecase.On("GetBalance", "123").Return(models.GetBalanceResponse{}, errors.New("failed to get balance"))

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/wallets?id=123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "failed to get balance")
	mockUsecase.AssertExpectations(t)
}
