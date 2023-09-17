package repository

import (
	"cupuama-go/config"
	"cupuama-go/domain"
	"cupuama-go/logger"
	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	Get(id int64) (domain.Products, error)
	List() ([]domain.Products, error)
	Create(product *domain.Products) (int64, error)
	Delete(id int64) error
	Update(product *domain.Products) (domain.Products, error)
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

	stmt, err := pr.db.Prepare(
		"insert into products (name, unit) " +
			"values ($1, $2) returning id")

	if err != nil {
		logger.Log.Info("(CreateProduct:Prepare)" + err.Error())
		return -1, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(product.Name,
		product.Unit).Scan(&product.ID)

	if err != nil {
		logger.Log.Info("(CreateProduct:Exec)" + err.Error())
		return -1, err
	}

	return product.ID, nil

}

// Delete removes the product with the specified login from the database
func (pr *ProductRepositoryDB) Delete(id int64) error {

	_, err := pr.Get(id)
	if err != nil {
		logger.Log.Info("(DeleteProduct:Get) Product doesn't exist")
		return err
	}

	stmt, err := pr.db.Prepare(
		"delete from products where id = $1")

	if err != nil {
		logger.Log.Info("(DeleteProduct:Prepare)" + err.Error())
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		logger.Log.Info("(DeleteProduct:Physical:Exec)" + err.Error())

		stmt, err = pr.db.Prepare(
			"update products " +
				"set deleted = true, " +
				"date_updated = now() " +
				"where id = $1")

		if err != nil {
			logger.Log.Info("(DeteleProduct:Logic:Prepare)" + err.Error())
			return err
		}

		_, err = stmt.Exec(id)

		if err != nil {
			logger.Log.Info("(DeleteProduct:Logic:Exec)" + err.Error())
			return err
		}

	}

	return nil
}

// Update modify the data for the specified product
func (pr *ProductRepositoryDB) Update(product *domain.Products) (domain.Products, error) {

	_, err := pr.Get(product.ID)
	if err != nil {
		logger.Log.Info("(UpdateProduct:Get) Product doesn't exist")
		return domain.Products{}, err
	}

	stmt, err := pr.db.Prepare(
		"update products " +
			"set name = $1, " +
			"unit = $2, " +
			"date_updated = now() " +
			"where id = $3")

	if err != nil {
		logger.Log.Info("(UpdateProduct:Prepare)" + err.Error())
		return domain.Products{}, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.Unit, product.ID)

	if err != nil {
		logger.Log.Info("(UpdateProduct:Exec)" + err.Error())
		return domain.Products{}, err
	}

	return *product, nil

}
