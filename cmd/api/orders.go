package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kostenbl4/wb_task_L0/internal/store"
)

func (app *application) GetOrderById(w http.ResponseWriter, r *http.Request) {

	orderUID := chi.URLParam(r, "id")

	if orderUID == "" {
		errorJSON(w, http.StatusBadRequest, "order_uid is empty")
		return
	}

	cacheOrder, found := app.cache.Get(orderUID)
	if found {
		app.logger.Debug(fmt.Sprintf("order with order_uid=%v found in cache", orderUID))
		writeJSON(w, http.StatusOK, cacheOrder)
		return
	}

	app.logger.Debug(fmt.Sprintf("order with order_uid=%v not found in cache", orderUID))

	order, err := app.store.Orders.GetByUID(orderUID)

	if err != nil {
		app.logger.Error(fmt.Sprintf("failed to get order by order_uid=%v: %v", orderUID, err))
		if err == sql.ErrNoRows {
			errorJSON(w, http.StatusNotFound, "order not found")
			return
		}
		errorJSON(w, http.StatusBadRequest, "failed to get order: "+err.Error())
		return
	}

	delivery, err := app.store.Deliveries.GetByOrderUID(orderUID)

	if err != nil {
		app.logger.Error(fmt.Sprintf("failed to get delivery by order_uid=%v: %v", orderUID, err))
		if err == sql.ErrNoRows {
			order.Delivery = store.Delivery{}
		} else {
			errorJSON(w, http.StatusBadRequest, "failed to get delivery: "+err.Error())
			return
		}
	} else {
		order.Delivery = delivery
	}

	payment, err := app.store.Payments.GetByOrderUID(orderUID)

	if err != nil {
		app.logger.Error(fmt.Sprintf("failed to get payment by order_uid=%v: %v", orderUID, err))
		if err == sql.ErrNoRows {
			order.Payment = store.Payment{}
		} else {
			errorJSON(w, http.StatusBadRequest, "failed to get payment: "+err.Error())
			return
		}
	} else {
		order.Payment = payment
	}

	items, err := app.store.Items.GetByOrderUID(orderUID)

	if err != nil {
		app.logger.Error(fmt.Sprintf("failed to get items by order_uid=%v: %v", orderUID, err))
		if err == sql.ErrNoRows {
			order.Items = []store.Item{}
		} else {
			errorJSON(w, http.StatusBadRequest, "failed to get items: "+err.Error())
			return
		}
	} else {
		order.Items = items
	}

	app.cache.Set(orderUID, *order)

	if err := writeJSON(w, http.StatusOK, order); err != nil {
		app.logger.Error("failed to respond json: " + err.Error())
		errorJSON(w, http.StatusBadRequest, "failed to respond json: "+err.Error())
		return
	}

}

func (app *application) CreateOrder(w http.ResponseWriter, r *http.Request) {

	var in store.Order

	if err := readJSON(r, &in); err != nil {
		app.logger.Error("failed to read json: " + err.Error())
		errorJSON(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}

	tx, err := app.store.Orders.BeginTx()

	if err != nil {
		app.logger.Error("failed to start transaction: " + err.Error())
		errorJSON(w, http.StatusInternalServerError, "failed to start transaction: "+err.Error())
		return
	}

	defer tx.Rollback()

	if err := app.store.Orders.Create(tx, &in); err != nil {
		app.logger.Error("failed to create order: " + err.Error())
		errorJSON(w, http.StatusBadRequest, "failed to create order: "+err.Error())
		return
	}

	if err := app.store.Deliveries.Create(tx, in.OrderUID, in.Delivery); err != nil {
		app.logger.Error(fmt.Sprintf("failed to create delivery related to order_uid=%v: %v", in.OrderUID, err))
		return
	}

	if err := app.store.Payments.Create(tx, in.OrderUID, in.Payment); err != nil {
		app.logger.Error(fmt.Sprintf("failed to create payment related to order_uid=%v: %v", in.OrderUID, err))
		return
	}

	if err := app.store.Items.CreateMany(tx, in.OrderUID, in.Items); err != nil {
		app.logger.Error(fmt.Sprintf("failed to create items related to order_uid=%v: %v", in.OrderUID, err))
		return
	}

	if err := tx.Commit(); err != nil {
		app.logger.Error("failed to commit transaction: " + err.Error())
		errorJSON(w, http.StatusInternalServerError, "failed to commit transaction: "+err.Error())
		return
	}

	app.cache.Set(in.OrderUID, in)

	if err := writeJSON(w, http.StatusCreated, "order created"); err != nil {
		app.logger.Error("failed to respond json: " + err.Error())
		errorJSON(w, http.StatusBadRequest, "failed to respond json: "+err.Error())
		return
	}

}

func (app *application) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	orderUID := chi.URLParam(r, "id")

	if err := app.store.Orders.Delete(orderUID); err != nil {
		app.logger.Error("failed to delete order: " + err.Error())
		errorJSON(w, http.StatusBadRequest, "failed to delete order: "+err.Error())
		return
	}

	app.cache.Delete(orderUID)

	if err := writeJSON(w, http.StatusOK, "order deleted"); err != nil {
		app.logger.Error("failed to respond json: " + err.Error())
	}
}
