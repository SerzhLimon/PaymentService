package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/SerzhLimon/PaymentService/internal/models"
)

type Repository interface {
	WalletTransactionDeposit(id uuid.UUID, amount int64) error
	WalletTransactionWithdraw(id uuid.UUID, amount int64) error
	GetBalance(id uuid.UUID) (models.GetBalanceResponse, error)
	CreateWallet(id uuid.UUID) error
}

type pgRepo struct {
	db *sql.DB
}

func NewPGRepository(db *sql.DB) Repository {
	return &pgRepo{db: db}
}

func (r *pgRepo) WalletTransactionDeposit(id uuid.UUID, amount int64) error {
	txOptions := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}

	tx, err := r.db.BeginTx(context.Background(), txOptions)
	if err != nil {
		err := errors.Errorf("pgRepo.WalletTransactionDeposit %v", err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	result, err := tx.Exec(queryWalletTransactionDeposit, id, amount)
	if err != nil {
		err := errors.Errorf("pgRepo.WalletTransactionDeposit %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		err := errors.Errorf("pgRepo.WalletTransactionDeposit %v", err)
		return err
	}
	if rowsAffected < 1 {
		err := errors.Errorf("pgRepo.WalletTransactionDeposit: no rows affected")
		return err
	}

	if err := tx.Commit(); err != nil {
		err := errors.Errorf("pgRepo.WalletTransactionDeposit %v", err)
		return err
	}

	return nil
}

func (r *pgRepo) WalletTransactionWithdraw(id uuid.UUID, amount int64) error {
	txOptions := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}

	tx, err := r.db.BeginTx(context.Background(), txOptions)
	if err != nil {
		err := errors.Errorf("pgRepo.WalletTransactionWithdraw %v", err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	result, err := tx.Exec(queryWalletTransactionWithdraw, id, amount)
	if err != nil {
		err := errors.Errorf("pgRepo.WalletTransactionWithdraw %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		err := errors.Errorf("pgRepo.WalletTransactionWithdraw %v", err)
		return err
	}
	if rowsAffected < 1 {
		err := errors.Errorf("pgRepo.WalletTransactionWithdraw: no rows affected")
		return err
	}

	if err := tx.Commit(); err != nil {
		err := errors.Errorf("pgRepo.WalletTransactionWithdraw %v", err)
		return err
	}

	return nil
}

func (r *pgRepo) GetBalance(id uuid.UUID) (models.GetBalanceResponse, error) {
	var res models.GetBalanceResponse
	if err := r.db.QueryRow(queryGetBalance, id).Scan(&res.Amount); err != nil {
		err := errors.Errorf("pgRepo.GetBalance %v", err)
		return res, err
	}
	return res, nil
}

func (r *pgRepo) CreateWallet(id uuid.UUID) error {
	_, err := r.db.Exec(queryCreateWallet, id)
	if err != nil {
		err := errors.Errorf("pgRepo.CreateWallet %v", err)
		return err
	}
	return  nil
}
