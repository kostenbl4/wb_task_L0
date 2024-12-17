package store

import (
	"database/sql"
)

type Deliveries interface {
	GetByOrderUID(string) (Delivery, error)
	Create(*sql.Tx, string, Delivery) error
	//Delete(int) error
	//Update(*Delivery) error
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type DeliveryStore struct {
	db *sql.DB
}

func (s *DeliveryStore) Create(tx *sql.Tx, orderUID string, delivery Delivery) error {
	query := `INSERT INTO delivery(order_uid, name, phone, zip, city, address, region, email) VALUES($1, $2, $3, $4, $5, $6, $7, $8)`
	if _, err := tx.Exec(query, orderUID, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email); err != nil {
		return err
	}
	return nil
}

func (s *DeliveryStore) GetByOrderUID(orderUID string) (Delivery, error) {
	query := `SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = $1`
	var delivery Delivery
	err := s.db.QueryRow(query, orderUID).Scan(&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email)
	if err != nil {
		return Delivery{}, err
	}
	return delivery, nil
}
