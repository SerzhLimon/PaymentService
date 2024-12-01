package usecase

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/SerzhLimon/PaymentService/internal/models"
	"github.com/SerzhLimon/PaymentService/internal/repository"
)

type operation int

const (
	unknown  operation = 0
	deposit  operation = 1
	withdraw operation = 2
)

type Usecase struct {
	pgPepo repository.Repository
}

type UseCase interface {
	WalletTransaction(models.WalletTransaction) error
	GetBalance(id string) (models.GetBalanceResponse, error)
}

func NewUsecase(pgPepo repository.Repository) UseCase {
	return &Usecase{pgPepo: pgPepo}
}

func (u *Usecase) WalletTransaction(data models.WalletTransaction) error {
	var err error
	id, err := u.parsedUUID(data.WalletID)
	if err != nil {
		err = errors.Errorf("usecase.WalletTransaction %v", err)
		return err
	}

	if err = u.parsedAmount(data.Amount); err != nil {
		err = errors.Errorf("usecase.WalletTransaction %v", err)
		return err
	}

	operation := u.parsedOperation(data.Operation)
	switch operation {
	case deposit:
		return u.pgPepo.WalletTransactionDeposit(id, data.Amount)
	case withdraw:
		return u.pgPepo.WalletTransactionWithdraw(id, data.Amount)
	default:
		err = errors.New("usecase.WalletTransaction: unknown transaction")
	}

	return err
}

func (u *Usecase) GetBalance(walletID string) (models.GetBalanceResponse, error) {
	id, err := u.parsedUUID(walletID)
	if err != nil {
		err = errors.Errorf("usecase.GetBalance %v", err)
		return models.GetBalanceResponse{}, err
	}
	
	return u.pgPepo.GetBalance(id)
}

func (u *Usecase) parsedUUID(data string) (uuid.UUID, error) {
	id, err := uuid.Parse(data)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (u *Usecase) parsedOperation(data string) operation {
	var res operation
	switch data {
	case "DEPOSIT":
		res = deposit
	case "WITHDRAW":
		res = withdraw
	default:
		res = unknown
	}
	return res
}

func (u *Usecase) parsedAmount(data int64) error {
	if data < 0 {
		err := errors.New("amount must be > 0")
		return err
	}
	return nil
}
