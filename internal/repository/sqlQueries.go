package repository

const (
	queryWalletTransactionDeposit = `
		UPDATE wallets
		SET balance = balance + $2, updated_at = now()
		WHERE id = $1
	`

	queryWalletTransactionWithdraw = `
		UPDATE wallets
		SET balance = balance - $2, updated_at = now()
		WHERE id = $1
	`

	queryGetBalance = `
		SELECT balance
		FROM wallets
		WHERE id = $1
	`

	queryCreateWallet = `
		INSERT INTO wallets (id, balance, created_at)
		VALUES ($1, 0, NOW())
	`
)