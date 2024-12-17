package main

import (
	"fmt"

	"github.com/kostenbl4/wb_task_L0/internal/store"
)

func (app *application) GetConsumerFunc() func([]byte) error {

	return func(msg []byte) error {
		var in store.Order

		if err := bytesToJSON(msg, &in); err != nil {
			app.logger.Error("failed to parse json: " + err.Error())
			return err
		}
		
		tx, err := app.store.Orders.BeginTx()

		if err != nil {
			app.logger.Error("failed to start transaction: " + err.Error())
			return err
		}

		defer tx.Rollback()

		if err := app.store.Orders.Create(tx, &in); err != nil {
			app.logger.Error("failed to create order: " + err.Error())
			return err
		}

		if err := app.store.Deliveries.Create(tx, in.OrderUID, in.Delivery); err != nil {
			app.logger.Error(fmt.Sprintf("failed to create delivery related to order_uid=%v: %v", in.OrderUID, err))
			return err
		}

		if err := app.store.Payments.Create(tx, in.OrderUID, in.Payment); err != nil {
			app.logger.Error(fmt.Sprintf("failed to create payment related to order_uid=%v: %v", in.OrderUID, err))
			return err
		}

		if err := app.store.Items.CreateMany(tx, in.OrderUID, in.Items); err != nil {
			app.logger.Error(fmt.Sprintf("failed to create items related to order_uid=%v: %v", in.OrderUID, err))
			return err
		}

		if err := tx.Commit(); err != nil {
			app.logger.Error("failed to commit transaction: " + err.Error())
			return err
		}

		app.logger.Debug(fmt.Sprintf("order with order_uid=%v created", in.OrderUID))
		app.cache.Set(in.OrderUID, in)
		
		return nil
	}
}