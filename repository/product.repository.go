package repository

import (
	"cupuama-go/config"
	"cupuama-go/domain"
	"cupuama-go/logger"
	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	Create(product *domain.Products) (int64, error)
	Get(id int64) (domain.Products, error)
	List() ([]domain.Products, error)
	Update(product *domain.Products) (domain.Products, error)
	Delete(id int64) error
}

type ProductRepositoryDB struct {
	db  *sqlx.DB
	app *config.AppConfig
}

func NewProductRepository(a *config.AppConfig) ProductRepository {
	return &ProductRepositoryDB{
		db:  a.DB,
		app: a,
	}
}

// Get retrieve an non-deleted product by login
func (pr *ProductRepositoryDB) Get(id int64) (domain.Products, error) {

	query := pr.app.SQLCache["products_get_id.sql"]
	stmt, err := pr.db.Prepare(query)

	if err != nil {
		logger.Log.Fatal("(GetProduct:Prepare)" + err.Error())
		return domain.Products{}, err
	}

	defer stmt.Close()
	var product domain.Products

	err = stmt.QueryRow(id).
		Scan(&product.ID,
			&product.Name,
			&product.Unit,
			&product.Audit.Deleted,
			&product.Audit.DateCreated,
			&product.Audit.DateUpdated)

	if err != nil {
		logger.Log.Info("(GetProduct:QueryRow:Scan)" + err.Error())
		return domain.Products{}, err
	}

	return product, nil
}

// List retrieves all non deleted products
func (pr *ProductRepositoryDB) List() ([]domain.Products, error) {
	query := pr.app.SQLCache["products_list.sql"]
	stmt, err := pr.db.Prepare(query)

	if err != nil {
		logger.Log.Info("(ListProduct:Prepare)" + err.Error())
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		logger.Log.Info("(ListProduct:Query)" + err.Error())
		return nil, err
	}

	defer rows.Close()

	var list []domain.Products
	err = sqlx.StructScan(rows, &list)
	if err != nil {
		logger.Log.Info("(ListProduct:StructScan)" + err.Error())
		return nil, err
	}

	return list, nil
}

// Create inserts a new product into the database and returns the new ID
func (pr *ProductRepositoryDB) Create(product *domain.Products) (int64, error) {

	insertQuery := pr.app.SQLCache["products_insert.sql"]
	err := pr.db.QueryRow(insertQuery, product.Name, product.Unit).Scan(&product.ID)

	if err != nil {
		logger.Log.Info("(CreateProduct:Exec)" + err.Error())
		return -1, err
	}

	return product.ID, nil
}

// Update modify the data for the specified product
func (pr *ProductRepositoryDB) Update(product *domain.Products) (domain.Products, error) {

	_, err := pr.Get(product.ID)
	if err != nil {
		logger.Log.Info("(UpdateProduct:Get) Product doesn't exist")
		return domain.Products{}, err
	}

	updateQuery := pr.app.SQLCache["products_update.sql"]
	_, err = pr.db.Exec(updateQuery, product.Name, product.Unit, product.Audit.Deleted, product.ID)

	if err != nil {
		logger.Log.Info("(UpdateProduct:Exec)" + err.Error())
		return domain.Products{}, err
	}

	return *product, nil
}

// Delete removes the product with the specified login from the database
func (pr *ProductRepositoryDB) Delete(id int64) error {

	product, err := pr.Get(id)
	if err != nil {
		logger.Log.Info("(DeleteProduct:Get) Product doesn't exist")
		return err
	}

	deleteQuery := pr.app.SQLCache["products_delete.sql"]
	_, err = pr.db.Exec(deleteQuery, id)

	if err != nil {
		logger.Log.Info("(DeleteProduct:Physical:Exec) " + err.Error())

		product.Audit.Deleted = true
		_, err := pr.Update(&product)
		if err != nil {
			logger.Log.Info("(DeleteProduct:Logic:Exec) " + err.Error())
			return err
		}

	}

	return nil
}
