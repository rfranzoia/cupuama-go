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

func (ois *OrderItemsStatus) ListOrderItemByOrderId(orderID int64) ([]OrderItems, error) {
	var items []OrderItems

	query := app.SQLCache["orders_orderItems_listByOrder.sql"]
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("(ListOrderItemByOrderId:Prepare)", err)
		return make([]OrderItems, 0), err
	}

	defer stmt.Close()

	rows, err := stmt.Query(orderID)
	if err != nil {
		log.Println("(ListOrderItemByOrderId:Query)", err)
		return make([]OrderItems, 0), err
	}

	defer rows.Close()

	for rows.Next() {
		var oi OrderItems

		if err := rows.Scan(&oi.ID, &oi.Product.ID, &oi.Product.Name, &oi.Fruit.ID, &oi.Fruit.Name, &oi.Quantity, &oi.UnitPrice); err != nil {
			log.Println("(ListOrderItemByOrderId:Scan)", err)
			return make([]OrderItems, 0), err
		}

		items = append(items, oi)
	}

	if err = rows.Err(); err != nil {
		log.Println("(ListOrder:Rows)", err)
		return make([]OrderItems, 0), err
	}

	return items, nil
}

// List retrieves a list of all non-deleted orders
func (ois *OrderItemsStatus) List() ([]OrderItemsStatus, error) {

	query := app.SQLCache["orders_list.sql"]
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("(ListOrder:Prepare)", err)
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Println("(ListOrder:Query)", err)
		return nil, err
	}

	defer rows.Close()

	var list []OrderItemsStatus

	for rows.Next() {
		var items []OrderItems
		var order OrderItemsStatus

		err := rows.Scan(&order.Order.ID, &order.Order.OrderDate, &order.Order.TotalPrice,
			&order.OrderStatus.ID, &order.OrderStatus.Status.Value, &order.OrderStatus.StatusChangeDate, &order.OrderStatus.Status.Description,
			&order.Order.Audit.Deleted, &order.Order.Audit.DateCreated, &order.Order.Audit.DateUpdated)
		if err != nil {
			log.Println("(ListOrder:Scan)", err)
			return nil, err
		}

		items, err = ois.ListOrderItemByOrderId(order.Order.ID)
		if err != nil {
			items = make([]OrderItems, 0)
		}

		order.OrderItems = items
		list = append(list, order)
	}

	if len(list) == 0 {
		log.Println("no order records found")
		err = errors.New("no records were found")
		return nil, err
	}

	if err = rows.Err(); err != nil {
		log.Println("(ListOrder:Rows)", err)
		return nil, err
	}

	return list, nil

}

// Get retrieves an order
func (ois *OrderItemsStatus) Get(orderID int64) (OrderItemsStatus, error) {

	query := app.SQLCache["orders_get.sql"]
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("(ListOrder:Prepare)", err)
		return OrderItemsStatus{}, err
	}

	defer stmt.Close()

	var order OrderItemsStatus
	var items []OrderItems

	err = stmt.QueryRow(&orderID).Scan(&order.Order.ID, &order.Order.OrderDate, &order.Order.TotalPrice,
		&order.OrderStatus.ID, &order.OrderStatus.Status.Value, &order.OrderStatus.StatusChangeDate, &order.OrderStatus.Status.Description,
		&order.Order.Audit.Deleted, &order.Order.Audit.DateCreated, &order.Order.Audit.DateUpdated)
	if err != nil {
		log.Println("(ListOrder:QueryRow)", err)
		return OrderItemsStatus{}, err
	}

	items, err = ois.ListOrderItemByOrderId(orderID)
	if err != nil {
		items = []OrderItems{}
	}

	order.OrderItems = items

	if order.Order.ID == 0 {
		log.Println("no order record(s) found")
		err = errors.New("no order record(s) found")
		return OrderItemsStatus{}, err
	}

	return order, nil
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
		Status: OrderCreated,
	}

	if err = ois.CreateOrderStatus(ois.Order.ID, os, tx); err != nil {
		log.Println("(CreateOrderStatus:Exec)", err)
		tx.Rollback()
		return OrderItemsStatus{}, err
	}

	if err = tx.Commit(); err != nil {
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
func (*OrderItemsStatus) CreateOrderStatus(orderID int64, os OrderStatus, tx *sql.Tx) error {

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
		orderExist := orderExists(orderID)
		if !orderExist {
			err := fmt.Errorf("order %d doesn't exist", orderID)
			log.Println("(CreateOrderStatus:GetOrder)", err)
			tx.Rollback()
			return err
		}
	}

	if os.Status.Value < 0 {
		err := fmt.Errorf("cannot create negative status")
		log.Println("(CreateOrderStatus:checkNegative)", err)
		tx.Rollback()
		return err

	} else if os.Status.Value > 0 {
		query := app.SQLCache["orders_list_max_status.sql"]
		stmt, err := tx.Prepare(query)
		if err != nil {
			log.Println("(CreateOrderStatus:ListMax:Prepare)", err)
			tx.Rollback()
			return err
		}

		var latestStatus int64
		err = stmt.QueryRow(&orderID).Scan(&orderID, &latestStatus)
		if err != nil {
			log.Println("(CreateOrderStatus:ListMax:Exec)", err)
			tx.Rollback()
			return err
		}

		// prevents the creation of a status that's not valid
		if latestStatus > os.Status.Value {
			err = errors.New("cannot set order to previous status")
			log.Println("(CreateOrderStatus:validationPrevious)", err)
			tx.Rollback()
			return err

		} else if os.Status.Value != 9 && os.Status.Value != (latestStatus+1) {
			err = errors.New(fmt.Sprintf("status order is not correct: got %d and should be %d", os.Status.Value, (latestStatus + 1)))
			log.Println("(CreateOrderStatus:validationNext)", err)
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

	err = stmt.QueryRow(orderID, os.Status.Value, os.Status.Description).Scan(&os.ID)

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
func (ois *OrderItemsStatus) DeleteOrderItems(orderID int64, oi ...OrderItems) error {

	if len(oi) == 0 {
		log.Println("(DeleteOrderItems:NoItemsToDelete)", oi)
		return nil
	}

	var err error

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("(DeleteOrderItems:CreateTransaction)", err)
		return err
	}

	order, err := ois.Get(orderID)
	if err != nil {
		err := fmt.Errorf("order %d doesn't exist", orderID)
		log.Println("(DeleteOrderItems:GetOrder)", err)
		return err
	}

	// delete the selected items if the item exists in the order
	for _, item := range oi {
		for _, orderItem := range order.OrderItems {
			if orderItem.Product.ID == item.Product.ID && orderItem.Fruit.ID == item.Fruit.ID {
				query := app.SQLCache["orders_orderItems_deleteOne.sql"]
				stmt, err := tx.Prepare(query)
				if err != nil {
					log.Println("(DeleteOrderItems:deleteItem:Prepare)", err)
					return err
				}
				_, err = stmt.Exec(&orderID, &item.Product.ID, &item.Fruit.ID)
				if err != nil {
					log.Println("(DeleteOrderItems:deleteItem:Exec)", err)
					return err
				}
			}
		}
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

	orderIsValid := orderExists(orderID) && orderHasStatus(orderID, OrderCreated)
	if !orderIsValid {
		err := fmt.Errorf("order %d doesn't exist", orderID)
		log.Println("(UpdateOrder:GetOrder)", err)
		tx.Rollback()
		return err
	}

	// delete all items from an order
	query := app.SQLCache["orders_orderItems_deleteAll.sql"]
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("(UpdateOrder:deleteAll:Prepare)", err)
		return err
	}

	_, err = stmt.Exec(&orderID)
	if err != nil {
		log.Println("(UpdateOrder:deleteAll:exec)", err)
		return err
	}

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

// CancelOrder change status of an order to Canceled
func (ois *OrderItemsStatus) CancelOrder(orderID int64) error {

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("(UpdateOrder:CreateTransaction)", err)
		return err
	}

	if !orderExists(orderID) {
		err := errors.New("order doesn't exists")
		log.Println("(CancelOrder:OrderNotFount)")
		return err
	}

	order, err := ois.Get(orderID)
	if err != nil {
		log.Println("(CancelOrder:GetOrder)", err)
		return err
	}

	if order.OrderStatus.Status.Value == OrderCanceled.Value {
		err := errors.New("order is already canceled")
		log.Println("(CancelOrder:OrderCanceledAlready)")
		return err

	} else if order.OrderStatus.Status.Value > OrderReadyForDelivery.Value {
		err := errors.New("cannot cancel orders that haven been dispatched or delivered")
		log.Println("(CancelOrder:OrderDispatchedOrDelivered)")
		return err

	} else if order.OrderStatus.Status.Value < OrderInPreparation.Value {
		// delete all items from an order
		query := app.SQLCache["orders_orderItems_deleteAll.sql"]
		stmt, err := tx.Prepare(query)
		if err != nil {
			log.Println("(CancelOrder:deleteAll:Prepare)", err)
			return err
		}

		_, err = stmt.Exec(&orderID)
		if err != nil {
			log.Println("(CancelOrder:deleteAll:exec)", err)
			return err
		}
	}

	os := OrderStatus{
		Status: OrderCanceled,
	}

	err = ois.CreateOrderStatus(orderID, os, tx)
	if err != nil {
		log.Println("(CancelOrder:CreateCancelStatus)", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println("(CancelOrder:Commit)", err)
		return err
	}

	return nil
}

func orderHasStatus(orderID int64, status OrderStatusType) bool {
	query := app.SQLCache["orders_list_max_status.sql"]
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("(orderHasStatus:Prepare)", err)
		return false
	}

	var latestStatus int64
	err = stmt.QueryRow(&orderID).Scan(&orderID, &latestStatus)
	if err != nil {
		log.Println("(orderHasStatus:Exec)", err)
		return false
	}

	return (latestStatus == status.Value)
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
