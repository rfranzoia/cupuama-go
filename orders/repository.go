package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/rfranzoia/cupuama-go/database"
)

var db = database.GetConnection()

// List retrieves a list of all non-deleted orders
func (*OrderItemsStatus) List(orderID int64) ([]OrderItemsStatus, error) {

	orderList := app.SQLCache["orders_list.sql"]
	stmt, err := db.Prepare(orderList)

	if err != nil {
		log.Println("(ListOrder:Prepare)", err)
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(orderID, orderID)
	if err != nil {
		log.Println("(ListOrder:Query)", err)
		return nil, err
	}

	defer rows.Close()

	var list []OrderItemsStatus
	var items []OrderItems

	var ois OrderItemsStatus
	var currentOIS OrderItemsStatus

	for rows.Next() {
		var oi OrderItems

		err := rows.Scan(&currentOIS.Order.ID, &currentOIS.Order.OrderDate, &currentOIS.Order.TotalPrice,
			&currentOIS.OrderStatus.ID, &currentOIS.OrderStatus.Status, &currentOIS.OrderStatus.StatusChangeDate, &currentOIS.OrderStatus.StatusDescription,
			&oi.ID, &oi.Product.ID, &oi.Product.Name, &oi.Fruit.ID, &oi.Fruit.Name, &oi.Quantity, &oi.UnitPrice)

		if err != nil {
			log.Println("(ListOrder:Scan)", err)
			return nil, err
		}

		oi.Order = currentOIS.Order
		items = append(items, oi)

		// in other words, if is the first record
		if ois.Order.ID == 0 {
			ois.Order.ID = currentOIS.Order.ID
			ois.Order.OrderDate = currentOIS.Order.OrderDate
			ois.Order.TotalPrice = currentOIS.Order.TotalPrice
			ois.OrderStatus = currentOIS.OrderStatus
			ois.OrderStatus.Order = currentOIS.Order

		} else if ois.Order.ID != currentOIS.Order.ID {
			ois.OrderItems = items
			list = append(list, ois)

			items = []OrderItems{}
			ois = OrderItemsStatus{}

			ois.Order.ID = currentOIS.Order.ID
			ois.Order.OrderDate = currentOIS.Order.OrderDate
			ois.Order.TotalPrice = currentOIS.Order.TotalPrice
			ois.OrderStatus = currentOIS.OrderStatus
			ois.OrderStatus.Order = currentOIS.Order

			currentOIS = OrderItemsStatus{}
		}

	}

	ois.OrderItems = items

	list = append(list, ois)

	err = rows.Err()
	if err != nil {
		log.Println("(ListOrder:Rows)", err)
		return nil, err
	}

	return list, nil

}

// Create creates a new Order with Items and Status
func (*OrderItemsStatus) Create(ois OrderItemsStatus) (OrderItemsStatus, error) {

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// creates the order
	insertQuery := app.SQLCache["orders_insert.sql"]
	stmt, err := tx.Prepare(insertQuery)
	defer stmt.Close()

	if err != nil {
		log.Println("(CreateOrder:Prepare)", err)
		return OrderItemsStatus{}, err
	}

	err = stmt.QueryRow(&ois.Order.TotalPrice).Scan(&ois.Order.ID)
	if err != nil {
		log.Println("(CreateOrder:Exec)", err)
		tx.Rollback()
		return OrderItemsStatus{}, err
	}

	// creates all order items
	err = ois.CreateOrderItems(ois.Order.ID, ois.OrderItems, tx)
	if err != nil {
		log.Println("(CreateOrderItems:Exec)", err)
		tx.Rollback()
		return OrderItemsStatus{}, err
	}

	// creates the first status: 0 - order-created
	os := OrderStatus{
		Order: Orders{
			ID: ois.Order.ID,
		},
		Status:            0,
		StatusDescription: "order-created",
	}

	err = ois.CreateOrderStatus(os, tx)
	if err != nil {
		log.Println("(CreateOrderStatus:Exec)", err)
		tx.Rollback()
		return OrderItemsStatus{}, err
	}

	err = tx.Commit()
	if err != nil {
		log.Println("(CreateOrder:Commit)", err)
		return OrderItemsStatus{}, err
	}

	return ois, nil
}

// CreateOrderItems insert a list of order items
func (*OrderItemsStatus) CreateOrderItems(orderID int64, orderItems []OrderItems, tx *sql.Tx) error {

	var err error

	localCommit := false
	checkOrder := false

	if tx == nil {
		ctx := context.Background()
		tx, err = db.BeginTx(ctx, nil)
		if err != nil {
			log.Println("(CreateOrderItem:CreateTransaction)", err)
			return err
		}
		localCommit = true
		checkOrder = true
	}

	if checkOrder {
		orderExist := orderExists(orderID)
		if !orderExist {
			err := fmt.Errorf("order %d doesn't exist", orderID)
			log.Println("(CreateOrderItem:GetOrder)", err)
			tx.Rollback()
			return err
		}
	}

	insertQuery := app.SQLCache["orders_orderItem_insert.sql"]
	for _, item := range orderItems {
		stmt, err := tx.Prepare(insertQuery)
		if err != nil {
			log.Println("(CreateOrderItem:Prepare)", err)
			tx.Rollback()
			return err
		}

		err = stmt.QueryRow(&orderID, &item.Product.ID, &item.Fruit.ID, &item.Quantity, &item.UnitPrice).Scan(&item.ID)
		if err != nil {
			log.Println("(CreateOrderItem:Exec)", err)
			tx.Rollback()
			return err
		}
	}

	if localCommit {
		err = tx.Commit()
		if err != nil {
			log.Println("(CreateOrderItem:Commit)", err)
			return err
		}
	}

	return nil
}

// CreateOrderStatus creates a new Order Status for an order
func (*OrderItemsStatus) CreateOrderStatus(os OrderStatus, tx *sql.Tx) error {

	var err error

	localCommit := false
	checkOrder := false

	if tx == nil {
		ctx := context.Background()
		tx, err = db.BeginTx(ctx, nil)
		if err != nil {
			log.Println("(CreateOrderStatus:CreateTransaction)", err)
			return err
		}
		localCommit = true
		checkOrder = true
	}

	if checkOrder {
		orderExist := orderExists(os.Order.ID)
		if !orderExist {
			err := fmt.Errorf("order %d doesn't exist", os.Order.ID)
			log.Println("(CreateOrderStatus:GetOrder)", err)
			tx.Rollback()
			return err
		}
	}

	if os.Status > 0 {
		query := app.SQLCache["orders_list_max_status.sql"]
		stmt, err := tx.Prepare(query)
		if err != nil {
			log.Println("(CreateOrderStatus:ListMax:Prepare)", err)
			tx.Rollback()
			return err
		}

		var latestStatus int64
		err = stmt.QueryRow(&os.Order.ID).Scan(&os.Order.ID, &latestStatus)
		if err != nil {
			log.Println("(CreateOrderStatus:ListMax:Exec)", err)
			tx.Rollback()
			return err
		}

		// prevents the creation of a status that's not valid
		if latestStatus > os.Status || os.Status < 0 {
			err = errors.New("cannot set order to previous status")
			log.Println("(CreateOrderStatus:validation)", err)
			tx.Rollback()
			return err
		}
	}

	query := app.SQLCache["orders_orderStatus_insert.sql"]
	stmt, err := tx.Prepare(query)

	if err != nil {
		log.Println("(CreateOrderStatus:Prepare)", err)
		tx.Rollback()
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(os.Order.ID, os.Status, os.StatusDescription).Scan(&os.ID)

	if err != nil {
		log.Println("(CreateOrderStatus:Exec)", err)
		tx.Rollback()
		return err
	}

	if localCommit {
		err = tx.Commit()
		if err != nil {
			log.Println("(CreateOrderStatus:Commit)", err)
			return err
		}
	}

	return nil

}

// DeleteOrderItems remove Items from an order
func (*OrderItemsStatus) DeleteOrderItems(orderID int64, oi []OrderItems) error {

	var err error

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("(CreateOrderStatus:CreateTransaction)", err)
		return err
	}

	orderExist := orderExists(orderID)
	if !orderExist {
		err := fmt.Errorf("order %d doesn't exist", orderID)
		log.Println("(CreateOrderStatus:GetOrder)", err)
		tx.Rollback()
		return err
	}

	// delete the selected items
	for _, item := range oi {
		fmt.Println("deleting", item)
	}

	err = tx.Commit()
	if err != nil {
		log.Println("(CreateOrderStatus:Commit)", err)
		return err
	}

	return nil
}

// UpdateOrder exchange items from an order with new ones
func (ois *OrderItemsStatus) UpdateOrder(orderID int64, oi []OrderItems) error {

	var err error

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("(UpdateOrder:CreateTransaction)", err)
		return err
	}

	orderExist := orderExists(orderID)
	if !orderExist {
		err := fmt.Errorf("order %d doesn't exist", orderID)
		log.Println("(UpdateOrder:GetOrder)", err)
		tx.Rollback()
		return err
	}

	// delete all items from an order
	// insert the new provided items into the order
	err = ois.CreateOrderItems(orderID, oi, tx)
	if err != nil {
		log.Println("(UpdateOrder:GetOrder)", err)
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println("(UpdateOrder:Commit)", err)
		return err
	}

	return nil
}

// orderExists checks if an order exists
func orderExists(orderID int64) bool {
	query := app.SQLCache["orders_get_id.sql"]
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("(OrderExists:Prepare)", err)
		return false
	}
	var order Orders
	err = stmt.QueryRow(&orderID).Scan(&order.ID)
	if err != nil {
		log.Println("(OrderExists:Exec)", err)
		return false
	}
	return order.ID != 0
}
