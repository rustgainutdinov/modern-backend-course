package infrastructure

import (
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"orderservice/pkg/orderservice/model"
	"time"
)

type Server struct {
	OrderService      model.OrderService
	OrderQueryService OrderQueryService
}

type MenuItemsList struct {
	MenuItems []model.MenuItem `json:"menuItems"`
}

type OrderQueryService interface {
	Orders() (*[]model.Order, error)
	Order(id string) (*model.FullOrderData, error)
}

func Router(srv *Server) http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/order/{ID}", srv.order).Methods(http.MethodGet)
	s.HandleFunc("/order/{ID}", srv.deleteOrder).Methods(http.MethodDelete)
	s.HandleFunc("/order/{ID}", srv.updateOrder).Methods(http.MethodPut)
	s.HandleFunc("/orders", srv.orders).Methods(http.MethodGet)
	s.HandleFunc("/order", srv.createOrder).Methods(http.MethodPost)
	return logMiddleware(r)
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method":    r.Method,
			"url":       r.URL,
			"remoteAdd": r.RemoteAddr,
			"userAgent": r.UserAgent(),
			"time":      time.Now(),
		}).Info("got a new request")
		h.ServeHTTP(w, r)
	})
}

func (s *Server) order(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["ID"]
	orderData, err := s.OrderQueryService.Order(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(orderData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err = io.WriteString(w, string(b)); err != nil {
		log.WithField("err", err).Error("write response error")
	}
}

func (s *Server) orders(w http.ResponseWriter, _ *http.Request) {
	orders, err := s.OrderQueryService.Orders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err = io.WriteString(w, string(b)); err != nil {
		log.WithField("err", err).Error("write response error")
	}
}

func (s *Server) createOrder(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	var createOrderReq MenuItemsList
	err = json.Unmarshal(b, &createOrderReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	orderUUID, err := s.OrderService.CreateOrder(createOrderReq.MenuItems)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err = io.WriteString(w, orderUUID.String()); err != nil {
		log.WithField("err", err).Error("write response error")
	}
}

func (s *Server) updateOrder(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	var menuItemsList MenuItemsList
	err = json.Unmarshal(b, &menuItemsList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	id := vars["ID"]
	err = s.OrderService.UpdateOrder(id, menuItemsList.MenuItems)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err = io.WriteString(w, "Ok"); err != nil {
		log.WithField("err", err).Error("write response error")
	}
}

func (s *Server) deleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["ID"]
	err := s.OrderService.DeleteOrder(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err = io.WriteString(w, "Ok"); err != nil {
		log.WithField("err", err).Error("write response error")
	}
}
