package repository

import (
	"cupuama-go/config"
	"cupuama-go/domain"
	"cupuama-go/logger"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
	"strings"
)

type OrderRepository interface {
	List() ([]domain.OrderItemsStatus, error)
	Get(orderID int64) (domain.OrderItemsStatus, error)
	Create(order *domain.OrderItemsStatus) (int64, error)
	CreateOrderItems(orderID int64, orderItems []domain.OrderItems, tx *sql.Tx) error
	CreateOrderStatus(orderID int64, os domain.OrderStatus, tx *sql.Tx) error
	UpdateOrder(orderID int64, oi []domain.OrderItems) error
	CancelOrder(orderID int64) error
	DeleteOrderItems(orderID int64, orderItems []domain.OrderItems) error
}

type OrderRepositoryDB struct {
	db  *sqlx.DB
	app *config.AppConfig
}

func NewOrderRepository(a *config.AppConfig) OrderRepository {
	return &OrderRepositoryDB{
		db:  a.DB,
		app: a,
	}
}

// List retrieves a list of all non-deleted orders
func (or *OrderRepositoryDB) List() ([]domain.OrderItemsStatus, error) {

	query := or.app.SQLCache["orders_list.sql"]
	stmt, err := or.db.Prepare(query)
	if err != nil {
		logger.Log.Info("(ListOrder:Prepare)" + err.Error())
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		logger.Log.Info("(ListOrder:Query)" + err.Error())
		return nil, err
	}

	defer rows.Close()

	var list []domain.OrderItemsStatus

	for rows.Next() {
		var items []domain.OrderItems
		var order domain.OrderItemsStatus

		err := rows.Scan(&order.Order.ID, &order.Order.OrderDate, &order.Order.TotalPrice,
			&order.OrderStatus.ID, &order.OrderStatus.Status.Value, &order.OrderStatus.StatusChangeDate, &order.OrderStatus.Status.Description,
			&order.Order.Audit.Deleted, &order.Order.Audit.DateCreated, &order.Order.Audit.DateUpdated)
		if err != nil {
			logger.Log.Info("(ListOrder:Scan)" + err.Error())
			return nil, err
		}

		items, err = or.ListOrderItemByOrderId(order.Order.ID)
		if err != nil {
			items = make([]domain.OrderItems, 0)
		}

		order.OrderItems = items
		list = append(list, order)
	}

	if len(list) == 0 {
		logger.Log.Info("no order records found")
		err = errors.New("no records were found")
		return nil, err
	}

	if err = rows.Err(); err != nil {
		logger.Log.Info("(ListOrder:Rows)" + err.Error())
		return nil, err
	}

	return list, nil

}

// Get retrieves an non-deleted order and its items
func (or *OrderRepositoryDB) Get(orderID int64) (domain.OrderItemsStatus, error) {

	query := or.app.SQLCache["orders_get.sql"]
	stmt, err := or.db.Prepare(query)
	if err != nil {
		logger.Log.Info("(ListOrder:Prepare)" + err.Error())
		return domain.OrderItemsStatus{}, err
	}

	defer stmt.Close()

	var order domain.OrderItemsStatus
	var items []domain.OrderItems

	err = stmt.QueryRow(&orderID).Scan(&order.Order.ID, &order.Order.OrderDate, &order.Order.TotalPrice,
		&order.OrderStatus.ID, &order.OrderStatus.Status.Value, &order.OrderStatus.StatusChangeDate, &order.OrderStatus.Status.Description,
		&order.Order.Audit.Deleted, &order.Order.Audit.DateCreated, &order.Order.Audit.DateUpdated)
	if err != nil {
		logger.Log.Info("(ListOrder:QueryRow)" + err.Error())
		return domain.OrderItemsStatus{}, err
	}

	items, err = or.ListOrderItemByOrderId(orderID)
	if err != nil {
		items = []domain.OrderItems{}
	}

	order.OrderItems = items

	if order.Order.ID == 0 {
		logger.Log.Info("no order record(s) found")
		err = errors.New("no order record(s) found")
		return domain.OrderItemsStatus{}, err
	}

	return order, nil
}

// Create creates a new Order with Items and Status
func (or *OrderRepositoryDB) Create(order *domain.OrderItemsStatus) (int64, error) {

	ctx := context.Background()
	tx, err := or.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
	defer tx.Rollback()

	// creates the order
	insertQuery := or.app.SQLCache["orders_insert.sql"]
	stmt, err := tx.Prepare(insertQuery)
	if err != nil {
		logger.Log.Info("(CreateOrder:Prepare)" + err.Error())
		return -1, err

	}

	defer stmt.Close()

	// calculates the total price before inserting the order
	var totalPrice = 0.00
	for _, item := range order.OrderItems {
		totalPrice += item.UnitPrice * float64(item.Quantity)
	}
	var orderID int64

	err = stmt.QueryRow(&totalPrice).Scan(&orderID)
	if err != nil {
		logger.Log.Info("(CreateOrder:Exec)" + err.Error())
		_ = tx.Rollback()
		return -1, err
	}

	// creates all order items
	err = or.CreateOrderItems(orderID, order.OrderItems, tx)
	if err != nil {
		logger.Log.Info("(CreateOrderItems:Exec)" + err.Error())
		_ = tx.Rollback()
		return -1, err
	}

	// creates the first status: 0 - order-created
	os := domain.OrderStatus{
		Status: domain.OrderCreated,
	}

	if err = or.CreateOrderStatus(orderID, os, tx); err != nil {
		logger.Log.Info("(CreateOrderStatus:Exec)" + err.Error())
		_ = tx.Rollback()
		return -1, err
	}

	if err = tx.Commit(); err != nil {
		logger.Log.Info("(CreateOrder:Commit)" + err.Error())
		return -1, err
	}

	return orderID, nil
}

// CreateOrderItems insert a list of order items
func (or *OrderRepositoryDB) CreateOrderItems(orderID int64, orderItems []domain.OrderItems, tx *sql.Tx) error {

	var err error

	localCommit := false
	checkOrder := false

	if tx == nil {
		ctx := context.Background()
		tx, err = or.db.BeginTx(ctx, nil)
		if err != nil {
			logger.Log.Info("(CreateOrderItem:CreateTransaction)" + err.Error())
			return err
		}
		defer tx.Rollback()
		localCommit = true
		checkOrder = true
	}

	if checkOrder {
		orderExist := or.orderExists(orderID)
		if !orderExist {
			err := fmt.Errorf("order %d doesn't exist", orderID)
			logger.Log.Info("(CreateOrderItem:GetOrder)" + err.Error())
			_ = tx.Rollback()
			return err
		}
	}

	insertQuery := or.app.SQLCache["orders_orderItems_insert.sql"]
	for _, item := range orderItems {
		stmt, err := tx.Prepare(insertQuery)
		if err != nil {
			logger.Log.Info("(CreateOrderItem:Prepare)" + err.Error())
			_ = tx.Rollback()
			return err
		}

		err = stmt.QueryRow(&orderID, &item.Product.ID, &item.Fruit.ID, &item.Quantity, &item.UnitPrice).Scan(&item.ID)
		if err != nil {
			logger.Log.Info("(CreateOrderItem:Exec)" + err.Error())
			_ = tx.Rollback()
			return err
		}
	}

	if localCommit {
		err = tx.Commit()
		if err != nil {
			logger.Log.Info("(CreateOrderItem:Commit)" + err.Error())
			return err
		}
	}

	return nil
}

// CreateOrderStatus creates a new Order Status for an order
func (or *OrderRepositoryDB) CreateOrderStatus(orderID int64, os domain.OrderStatus, tx *sql.Tx) error {

	var err error

	localCommit := false
	checkOrder := false

	if tx == nil {
		ctx := context.Background()
		tx, err = or.db.BeginTx(ctx, nil)
		if err != nil {
			logger.Log.Info("(CreateOrderStatus:CreateTransaction)" + err.Error())
			return err
		}
		defer tx.Rollback()
		localCommit = true
		checkOrder = true
	}

	if checkOrder {
		if orderExist := or.orderExists(orderID); !orderExist {
			err := fmt.Errorf("order %d doesn't exist", orderID)
			logger.Log.Info("(CreateOrderStatus:GetOrder)" + err.Error())
			tx.Rollback()
			return err
		}
	}

	if os.Status.Value < 0 {
		err := fmt.Errorf("cannot create negative status")
		logger.Log.Info("(CreateOrderStatus:checkNegative)" + err.Error())
		tx.Rollback()
		return err

	}

	query := or.app.SQLCache["orders_list_max_status.sql"]
	stmt, err := tx.Prepare(query)
	if err != nil {
		logger.Log.Info("(CreateOrderStatus:ListMax:Prepare)" + err.Error())
		tx.Rollback()
		return err
	}

	var latestStatus int64
	err = stmt.QueryRow(&orderID).Scan(&orderID, &latestStatus)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			logger.Log.Info("(CreateOrderStatus:ListMax:Exec)" + err.Error())
			tx.Rollback()
			return err
		} else {
			latestStatus = -1
		}
	}

	if localCommit {
		// prevents the creation of a status that's not valid
		if latestStatus > os.Status.Value {
			err = errors.New("cannot set order to previous status")
			logger.Log.Info("(CreateOrderStatus:validationPrevious)" + err.Error())
			tx.Rollback()
			return err

		} else if os.Status.Value != 9 && os.Status.Value != (latestStatus+1) {
			err = fmt.Errorf("status order is not correct: got %d and should be %d", os.Status.Value, (latestStatus + 1))
			logger.Log.Info("(CreateOrderStatus:validationNext)" + err.Error())
			tx.Rollback()
			return err

		} else if os.Status.Value == 9 && latestStatus >= 4 {
			err = fmt.Errorf("cannot cancel order %d after status ´%s´", orderID, domain.OrderStatusMap[4].Description)
			logger.Log.Info("(CreateOrderStatus:validationCancel)" + err.Error())
			tx.Rollback()
			return err
		}
	}

	query = or.app.SQLCache["orders_orderStatus_insert.sql"]
	stmt, err = tx.Prepare(query)

	if err != nil {
		logger.Log.Info("(CreateOrderStatus:Prepare)" + err.Error())
		tx.Rollback()
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(orderID, os.Status.Value, os.Status.Description).Scan(&os.ID)

	if err != nil {
		logger.Log.Info("(CreateOrderStatus:Exec)" + err.Error())
		tx.Rollback()
		return err
	}

	if localCommit {
		query = or.app.SQLCache["orders_orderItems_sumByOrder.sql"]
		stmt, err = tx.Prepare(query)
		if err != nil {
			logger.Log.Info("(CreateOrdertatus:sumItemsPrice:prepare)" + err.Error())
			tx.Rollback()
			return err
		}
		var totalPrice float64
		err = stmt.QueryRow(&orderID).Scan(&totalPrice)
		if err != nil {
			logger.Log.Info("(CreateOrdertatus:sumItemsPrice:QueryScan)" + err.Error())
			tx.Rollback()
			return err
		}
		if err = or.updateOrderAudit(orderID, totalPrice, tx); err != nil {
			logger.Log.Info("(CreateOrdertatus:updateDateUpdated)" + err.Error())
			tx.Rollback()
			return err
		}

		err = tx.Commit()
		if err != nil {
			logger.Log.Info("(CreateOrderStatus:Commit)" + err.Error())
			return err
		}
	}

	return nil

}

// ListOrderItemByOrderId retrives a slice with the OrderItems of an Order
func (or *OrderRepositoryDB) ListOrderItemByOrderId(orderID int64) ([]domain.OrderItems, error) {
	var items []domain.OrderItems

	query := or.app.SQLCache["orders_orderItems_listByOrder.sql"]
	stmt, err := or.db.Prepare(query)
	if err != nil {
		logger.Log.Info("(ListOrderItemByOrderId:Prepare)" + err.Error())
		return make([]domain.OrderItems, 0), err
	}

	defer stmt.Close()

	rows, err := stmt.Query(orderID)
	if err != nil {
		logger.Log.Info("(ListOrderItemByOrderId:Query)" + err.Error())
		return make([]domain.OrderItems, 0), err
	}

	defer rows.Close()

	for rows.Next() {
		var oi domain.OrderItems

		if err := rows.Scan(&oi.ID, &oi.Product.ID, &oi.Product.Name, &oi.Fruit.ID, &oi.Fruit.Name, &oi.Quantity, &oi.UnitPrice); err != nil {
			logger.Log.Info("(ListOrderItemByOrderId:Scan)" + err.Error())
			return make([]domain.OrderItems, 0), err
		}

		items = append(items, oi)
	}

	if err = rows.Err(); err != nil {
		logger.Log.Info("(ListOrder:Rows)" + err.Error())
		return make([]domain.OrderItems, 0), err
	}

	return items, nil
}

// DeleteOrderItems remove Items from an order
func (or *OrderRepositoryDB) DeleteOrderItems(orderID int64, orderItems []domain.OrderItems) error {

	if len(orderItems) == 0 {
		logger.Log.Info("(DeleteOrderItems:NoItemsToDelete)")
		return errors.New("no order item was informed to be deleted")
	}

	var err error

	ctx := context.Background()
	tx, err := or.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.Info("(DeleteOrderItems:CreateTransaction)" + err.Error())
		return err
	}
	defer tx.Rollback()

	if !or.orderExists(orderID) {
		err = errors.New("order doesnt exists")
		logger.Log.Info("(DeleteOrderItems:OrderDoesntExists)" + err.Error())
		return err
	}

	// delete the selected items
	for _, item := range orderItems {
		// check if the item belongs to the provided order
		if belongsTo := or.doesItemBelongsToOrder(orderID, item); !belongsTo {
			err = fmt.Errorf("order item %d doesn't belong to order %d", item.ID, orderID)
			logger.Log.Info("(DeleteOrderItems:deleteItem:Prepare)" + err.Error())
			return err
		}

		// actually delete de valid item
		query := or.app.SQLCache["orders_orderItems_delete.sql"]
		stmt, err := tx.Prepare(query)
		if err != nil {
			logger.Log.Info("(DeleteOrderItems:deleteItem:Prepare)" + err.Error())
			return err
		}
		_, err = stmt.Exec(&item.ID)
		if err != nil {
			logger.Log.Info("(DeleteOrderItems:deleteItem:Exec)" + err.Error())
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		logger.Log.Info("(DeleteOrderItems:Commit)" + err.Error())
		return err
	}

	return nil
}

// UpdateOrder exchange items from an order with new ones
func (or *OrderRepositoryDB) UpdateOrder(orderID int64, oi []domain.OrderItems) error {

	var err error

	ctx := context.Background()
	tx, err := or.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.Info("(UpdateOrder:CreateTransaction)" + err.Error())
		return err
	}
	defer tx.Rollback()

	if exists := or.orderExists(orderID) && or.orderHasStatus(orderID, domain.OrderCreated); !exists {
		err := fmt.Errorf("order %d doesn't exist", orderID)
		logger.Log.Info("(UpdateOrder:GetOrder)" + err.Error())
		tx.Rollback()
		return err
	}

	// delete all items from an order
	query := or.app.SQLCache["orders_orderItems_deleteAll.sql"]
	stmt, err := or.db.Prepare(query)
	if err != nil {
		logger.Log.Info("(UpdateOrder:deleteAll:Prepare)" + err.Error())
		return err
	}
	_, err = stmt.Exec(&orderID)
	if err != nil {
		logger.Log.Info("(UpdateOrder:deleteAll:exec)" + err.Error())
		return err
	}

	// insert the new provided items into the order
	err = or.CreateOrderItems(orderID, oi, tx)
	if err != nil {
		logger.Log.Info("(UpdateOrder:CreateOrderItems)" + err.Error())
		tx.Rollback()
		return err
	}

	// recalculate TotalPrice
	var totalPrice = 0.00
	for _, item := range oi {
		totalPrice += (item.UnitPrice * float64(item.Quantity))
	}

	if err = or.updateOrderAudit(orderID, totalPrice, tx); err != nil {
		tx.Rollback()
		logger.Log.Info("(UpdateOrder:UpdateOrderAudit)" + err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		logger.Log.Info("(UpdateOrder:Commit)" + err.Error())
		return err
	}

	return nil
}

// CancelOrder change status of an order to Canceled
func (or *OrderRepositoryDB) CancelOrder(orderID int64) error {

	ctx := context.Background()
	tx, err := or.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.Info("(UpdateOrder:CreateTransaction)" + err.Error())
		return err
	}
	defer tx.Rollback()

	if !or.orderExists(orderID) {
		err := errors.New("order doesn't exists")
		logger.Log.Info("(CancelOrder:OrderNotFound)")
		return err
	}

	order, err := or.Get(orderID)
	if err != nil {
		logger.Log.Info("(CancelOrder:GetOrder)" + err.Error())
		return err
	}

	if order.OrderStatus.Status.Equals(domain.OrderCanceled) {
		err := errors.New("order is already canceled")
		logger.Log.Info("(CancelOrder:OrderCanceledAlready)")
		return err

	} else if order.OrderStatus.Status.Value > domain.OrderReadyForDelivery.Value {
		err := errors.New("cannot cancel orders that haven been dispatched or delivered")
		logger.Log.Info("(CancelOrder:OrderDispatchedOrDelivered)")
		return err

	} else if order.OrderStatus.Status.Value < domain.OrderInPreparation.Value {
		// delete all items from an order
		query := or.app.SQLCache["orders_orderItems_deleteAll.sql"]
		stmt, err := tx.Prepare(query)
		if err != nil {
			logger.Log.Info("(CancelOrder:deleteAll:Prepare)" + err.Error())
			return err
		}

		_, err = stmt.Exec(&orderID)
		if err != nil {
			logger.Log.Info("(CancelOrder:deleteAll:exec)" + err.Error())
			return err
		}
	}

	os := domain.OrderStatus{
		Status: domain.OrderCanceled,
	}

	err = or.CreateOrderStatus(orderID, os, tx)
	if err != nil {
		logger.Log.Info("(CancelOrder:CreateCancelStatus)" + err.Error())
		tx.Rollback()
		return err
	}

	if err = or.updateOrderAudit(orderID, 0.00, tx); err != nil {
		logger.Log.Info("(CancelOrder:UpdateDateUpdated)" + err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		logger.Log.Info("(CancelOrder:Commit)" + err.Error())
		return err
	}

	return nil
}

func (or *OrderRepositoryDB) updateOrderAudit(orderID int64, totalPrice float64, tx *sql.Tx) error {
	var err error

	localCommit := false

	// uses the transaction from calling method if one has it
	if tx == nil {
		ctx := context.Background()
		tx, err = or.db.BeginTx(ctx, nil)
		if err != nil {
			logger.Log.Info("(updateOrderAudit:CreateTransaction)" + err.Error())
			return err
		}
		defer tx.Rollback()
		localCommit = true
	}

	// update the date_upated on order
	query := or.app.SQLCache["orders_updateTotalPrice.sql"]
	stmt, err := tx.Prepare(query)
	if err != nil {
		logger.Log.Info("(updateOrderAudit:Prepare)" + err.Error())
		return err
	}

	_, err = stmt.Exec(&totalPrice, &orderID)
	if err != nil {
		logger.Log.Info("(updateOrderAudit:exec)" + err.Error())
		if localCommit {
			tx.Rollback()
		}
		return err
	}

	// finish everything up
	if localCommit {
		err = tx.Commit()
		if err != nil {
			logger.Log.Info("(updateOrderAudit:Commit)" + err.Error())
			return err
		}
	}

	return nil
}

func (or *OrderRepositoryDB) orderHasStatus(orderID int64, status domain.OrderStatusType) bool {
	query := or.app.SQLCache["orders_list_max_status.sql"]
	stmt, err := or.db.Prepare(query)
	if err != nil {
		logger.Log.Info("(orderHasStatus:Prepare)" + err.Error())
		return false
	}

	var latestStatus int64
	err = stmt.QueryRow(&orderID).Scan(&orderID, &latestStatus)
	if err != nil {
		logger.Log.Info("(orderHasStatus:Exec)" + err.Error())
		return false
	}

	return (latestStatus == status.Value)
}

// orderExists checks if an order exists
func (or *OrderRepositoryDB) orderExists(orderID int64) bool {
	query := or.app.SQLCache["orders_get_id.sql"]
	stmt, err := or.db.Prepare(query)
	if err != nil {
		logger.Log.Info("(OrderExists:Prepare)" + err.Error())
		return false
	}
	var order domain.Orders
	err = stmt.QueryRow(&orderID).Scan(&order.ID)
	if err != nil {
		logger.Log.Info("(OrderExists:Exec)" + err.Error())
		return false
	}
	return order.ID != 0
}

func (or *OrderRepositoryDB) doesItemBelongsToOrder(orderID int64, orderItem domain.OrderItems) bool {
	query := or.app.SQLCache["orders_orderItems_selectByIdAndOrderId.sql"]
	stmt, err := or.db.Prepare(query)
	if err != nil {
		logger.Log.Info("(doesItemBelongsToOrder:Prepare)" + err.Error())
		return false
	}

	if err = stmt.QueryRow(&orderID, &orderItem.ID).Scan(&orderItem.ID); err != nil {
		logger.Log.Info("(doesItemBelongsToOrder:QueryRow)" + err.Error())
		return false
	}

	return true
}
