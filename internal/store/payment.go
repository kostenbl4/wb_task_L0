package store

import (
	"database/sql"
)

type Payments interface {
	GetByOrderUID(string) (Payment, error)
	Create(*sql.Tx, string, Payment) error
	//Delete(int) error
	//Update(*Payment) error
}

type Payment struct {
	Transaction  string  `json:"transaction"`
	RequestID    *string `json:"request_id"`
	Currency     string  `json:"currency"`
	Provider     string  `json:"provider"`
	Amount       int     `json:"amount"`
	PaymentDt    int64   `json:"payment_dt"`
	Bank         string  `json:"bank"`
	DeliveryCost int     `json:"delivery_cost"`
	GoodsTotal   int     `json:"goods_total"`
	CustomFee    int     `json:"custom_fee"`
}

type PaymentStore struct {
	db *sql.DB
}

func (s *PaymentStore) Create(tx *sql.Tx, orderUID string, payment Payment) error {
	query := `INSERT INTO payment(order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	if _, err := tx.Exec(query, orderUID, payment.Transaction, payment.RequestID, payment.Currency, payment.Provider, payment.Amount, payment.PaymentDt, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee); err != nil {
		return err
	}
	return nil
}

func (s *PaymentStore) GetByOrderUID(orderUID string) (Payment, error) {
	query := `SELECT transaction, request_id, currency, provider, amount, delivery_cost, goods_total, custom_fee FROM payment WHERE order_uid = $1`
	var payment Payment
	if err := s.db.QueryRow(query, orderUID).Scan(&payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider, &payment.Amount, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee); err != nil {
		return Payment{}, err
	}
	return payment, nil
}
