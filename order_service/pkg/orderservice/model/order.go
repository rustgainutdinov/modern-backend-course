package model

import (
	"errors"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type OrderRepo interface {
	Create(order FullOrderData) (*uuid.UUID, error)
	AddOrderMenuItems(items []MenuItem, orderID string) error
	Delete(orderID string) error
}

type OrderService interface {
	CreateOrder(menuItems []MenuItem) (*uuid.UUID, error)
	DeleteOrder(orderID string) error
	UpdateOrder(orderID string, menuItems []MenuItem) error
}

type FullOrderData struct {
	ID               string        `json:"id"`
	OrderAtTimestamp time.Duration `json:"orderAtTimestamp"`
	Cost             int           `json:"cost"`
	MenuItems        []MenuItem    `json:"menuItems"`
}

type Order struct {
	Id        string     `json:"id"`
	MenuItems []MenuItem `json:"menuItems"`
}

type MenuItem struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type orderService struct {
	repo OrderRepo
}

func (s *orderService) CreateOrder(menuItems []MenuItem) (*uuid.UUID, error) {
	if len(menuItems) == 0 {
		return nil, errors.New("count of menu items must be more than 0")
	}
	cost := rand.Intn(50) + 50
	fullOrderData := FullOrderData{Cost: cost, MenuItems: menuItems}
	return s.repo.Create(fullOrderData)
}

func (s *orderService) DeleteOrder(orderID string) error {
	return s.repo.Delete(orderID)
}

func (s *orderService) UpdateOrder(orderID string, menuItems []MenuItem) error {
	return s.repo.AddOrderMenuItems(menuItems, orderID)
}

func NewOrderService(repo OrderRepo) OrderService {
	return &orderService{repo}
}
