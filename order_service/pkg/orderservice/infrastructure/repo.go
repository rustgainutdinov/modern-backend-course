package infrastructure

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"orderservice/pkg/orderservice/model"
)

type repo struct {
	db *sqlx.DB
}

func (s *repo) Create(order model.FullOrderData) (*uuid.UUID, error) {
	u := uuid.New()
	q := `INSERT INTO "order" (id_order, created_at, cost) VALUES ($1, now(), $2)`
	_, err := s.db.Exec(q, u, order.Cost)
	if err != nil {
		return nil, err
	}
	return &u, s.AddOrderMenuItems(order.MenuItems, u.String())
}

func (s *repo) Delete(orderID string) error {
	q := `DELETE FROM "order" WHERE id_order = $1`
	_, err := s.db.Exec(q, orderID)
	return err
}

func (s *repo) AddOrderMenuItems(items []model.MenuItem, orderID string) error {
	q := `INSERT INTO "menu_item" (id_menu_item, id_order, quantity) VALUES ($1, $2, $3)`
	for _, item := range items {
		_, err := s.db.Exec(q, item.Id, orderID, item.Quantity)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewRepo(db *sqlx.DB) model.OrderRepo {
	return &repo{db}
}
