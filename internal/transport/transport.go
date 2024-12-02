package transport

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/PaymentService/internal/models"
	"github.com/SerzhLimon/PaymentService/internal/repository"
	uc "github.com/SerzhLimon/PaymentService/internal/usecase"
)

type Server struct {
	Usecase uc.UseCase
}

func NewServer(database *sql.DB) *Server {
	pgClient := repository.NewPGRepository(database)
	uc := uc.NewUsecase(pgClient)

	return &Server{
		Usecase: uc,
	}
}

func (s *Server) WalletTransaction(c *gin.Context) {
	
	var request models.WalletTransaction
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.WithError(err).Error("error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format"})
		return
	}

	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugf("Parsed request: %s %s %d", request.WalletID, request.Operation, request.Amount)

	if err := s.Usecase.WalletTransaction(request); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "transaction failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "true"})
}

func (s *Server) GetBalance(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		err := errors.New("parametr 'id' is empty")
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter 'id' is empty"})
		return
	}

	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugf("Parsed request: %s", id)

	res, err := s.Usecase.GetBalance(id)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get balance"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) CreateWallet(c *gin.Context) {
	
	err := s.Usecase.CreateWallet()
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create wallet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "true"})
}
