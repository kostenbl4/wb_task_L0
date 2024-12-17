package store

import (
	"database/sql"
)

type Items interface {
	GetByOrderUID(string) ([]Item, error)
	CreateMany(*sql.Tx, string, []Item) error
	//Delete(int) error
	//Update(*Item) error
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

type ItemStore struct {
	db *sql.DB
}

func (s *ItemStore) CreateMany(tx *sql.Tx, orderUID string, items []Item) error {

	for _, i := range items {
		query := `INSERT INTO items(order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		if _, err := tx.Exec(query, orderUID, i.ChrtID, i.TrackNumber, i.Price, i.Rid, i.Name, i.Sale, i.Size, i.TotalPrice, i.NmID, i.Brand, i.Status); err != nil {
			return err
		}
	}
	return nil
}

func (s *ItemStore) GetByOrderUID(orderUID string) ([]Item, error) {
	query := `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_uid = $1`

	rows, err := s.db.Query(query, orderUID)
	if err != nil {
		return nil, err
	}

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
