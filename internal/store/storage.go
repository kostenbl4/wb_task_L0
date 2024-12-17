package store

import "database/sql"

type Storage struct {
	Orders
	Items
	Payments
	Deliveries


}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Orders:   &OrderStore{db},
		Items:    &ItemStore{db},
		Payments: &PaymentStore{db},
		Deliveries: &DeliveryStore{db},
	}
}
