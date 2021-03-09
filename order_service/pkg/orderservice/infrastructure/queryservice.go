package infrastructure

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"orderservice/pkg/orderservice/model"
	"time"
)

type queryService struct {
	db *sqlx.DB
}

func (s *queryService) Orders() (*[]model.Order, error) {
	q := `SELECT id_order FROM "order"`
	var orderIDs []uuid.UUID
	err := s.db.Select(&orderIDs, q)
	if err != nil {
		return nil, err
	}
	orders := make([]model.Order, 0)
	for _, orderID := range orderIDs {
		orderData, err := s.Order(orderID.String())
		if err != nil {
			return nil, err
		}
		orders = append(orders, model.Order{Id: orderData.ID, MenuItems: orderData.MenuItems})
	}
	return &orders, nil
}

func (s *queryService) Order(id string) (*model.FullOrderData, error) {
	q := `SELECT * FROM "order" WHERE id_order = $1`
	var orders []sqlxOrder
	err := s.db.Select(&orders, q, id)
	if err != nil {
		return nil, err
	}
	order := orders[0]
	menuItems, err := s.OrderMenuItems(order.ID)
	if err != nil {
		return nil, err
	}
	return &model.FullOrderData{ID: order.ID, Cost: order.Cost, OrderAtTimestamp: time.Duration(order.CreatedAt.Unix()), MenuItems: *menuItems}, err
}

func (s *queryService) OrderMenuItems(id string) (*[]model.MenuItem, error) {
	q := `SELECT id_menu_item, quantity FROM "menu_item" WHERE id_order = $1`
	var result []sqlxMenuItem
	err := s.db.Select(&result, q, id)
	if err != nil {
		return nil, err
	}
	menuItems := make([]model.MenuItem, 0)
	for _, resultItem := range result {
		menuItems = append(menuItems, model.MenuItem{Id: resultItem.ID, Quantity: resultItem.Quantity})
	}
	return &menuItems, nil
}

func NewQueryService(db *sqlx.DB) OrderQueryService {
	return &queryService{db}
}

type sqlxOrder struct {
	ID        string    `db:"id_order"`
	CreatedAt time.Time `db:"created_at"`
	Cost      int       `db:"cost"`
}

type sqlxMenuItem struct {
	ID       string `db:"id_menu_item"`
	Quantity int    `db:"quantity"`
}
