package store

import (
	"database/sql"
	"time"
)

type Orders interface {
	GetByUID(string) (*Order, error)
	Create(*sql.Tx, *Order) error
	Delete(string) error
	BeginTx() (*sql.Tx, error)
}

type Order struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature *string   `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

type OrderStore struct {
	db *sql.DB
}

func (s *OrderStore) GetByUID(uid string) (*Order, error) {
	query := `SELECT * FROM orders WHERE order_uid = $1`
	var order Order
	err := s.db.QueryRow(query, uid).Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *OrderStore) Create(tx *sql.Tx, order *Order) error {

	query := `INSERT INTO orders(order_uid, track_number, entry, locale, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) 
				VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	if _, err := tx.Exec(query, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard); err != nil {
		return err
	}

	return nil
}

func (s *OrderStore) Delete(uid string) error {
	query := `DELETE FROM orders WHERE order_uid = $1`
	if _, err := s.db.Exec(query, uid); err != nil {
		return err
	}
	return nil
}

func (s *OrderStore) BeginTx() (*sql.Tx, error) {
	return s.db.Begin()
}
