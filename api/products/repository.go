package products

import (
	"log"

	"github.com/rfranzoia/cupuama-go/database"
)

var db = database.GetConnection()

// Get retrieve an non-deleted product by login
func (*Products) Get(id int64) (Products, error) {

	query := app.SQLCache["products_get_id.sql"]
	stmt, err := db.Prepare(query)

	if err != nil {
		log.Fatal("(GetProduct:Prepare)", err)
		return Products{}, err
	}

	defer stmt.Close()
	var product Products

	err = stmt.QueryRow(id).
		Scan(&product.ID,
			&product.Name,
			&product.Unit,
			&product.Audit.Deleted,
			&product.Audit.DateCreated,
			&product.Audit.DateUpdated)

	if err != nil {
		log.Println("(GetProduct:QueryRow:Scan)", err)
		return Products{}, err
	}

	return product, nil
}

//List retrieves all non deleted products
func (*Products) List() ([]Products, error) {

	query := app.SQLCache["products_list.sql"]
	stmt, err := db.Prepare(query)

	if err != nil {
		log.Println("(ListProduct:Prepare)", err)
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Println("(ListProduct:Query)", err)
		return nil, err
	}

	defer rows.Close()

	var list []Products

	for rows.Next() {
		var product Products

		err := rows.Scan(&product.ID,
			&product.Name,
			&product.Unit,
			&product.Audit.Deleted,
			&product.Audit.DateCreated,
			&product.Audit.DateUpdated)

		if err != nil {
			log.Println("(ListProduct:Scan)", err)
			return nil, err
		}

		list = append(list, product)
	}

	err = rows.Err()
	if err != nil {
		log.Println("(ListProduct:Rows)", err)
		return nil, err
	}

	return list, nil
}

// Create inserts a new product into the database and returns the new ID
func (*Products) Create(product *Products) (int64, error) {

	stmt, err := db.Prepare(
		"insert into products (name, unit) " +
			"values ($1, $2) returning id")

	if err != nil {
		log.Println("(CreateProduct:Prepare)", err)
		return -1, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(product.Name,
		product.Unit).Scan(&product.ID)

	if err != nil {
		log.Println("(CreateProduct:Exec)", err)
		return -1, err
	}

	return product.ID, nil

}

// Delete removes the product with the specified login from the database
func (*Products) Delete(id int64) error {

	_, err := model.Get(id)
	if err != nil {
		log.Println("(DeleteProduct:Get)", "Product doesn't exist")
		return err
	}

	stmt, err := db.Prepare(
		"delete from products where id = $1")

	if err != nil {
		log.Println("(DeleteProduct:Prepare)", err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		log.Println("(DeleteProduct:Physical:Exec)", err)

		stmt, err = db.Prepare(
			"update products " +
				"set deleted = true, " +
				"date_updated = now() " +
				"where id = $1")

		if err != nil {
			log.Println("(DeteleProduct:Logic:Prepare)", err)
			return err
		}

		_, err = stmt.Exec(id)

		if err != nil {
			log.Println("(DeleteProduct:Logic:Exec)", err)
			return err
		}

	}

	return nil
}

// Update modify the data for the specified product
func (*Products) Update(product *Products) (Products, error) {

	_, err := model.Get(product.ID)
	if err != nil {
		log.Println("(UpdateProduct:Get)", "Product doesn't exist")
		return Products{}, err
	}

	stmt, err := db.Prepare(
		"update products " +
			"set name = $1, " +
			"unit = $2, " +
			"date_updated = now() " +
			"where id = $3")

	if err != nil {
		log.Println("(UpdateProduct:Prepare)", err)
		return Products{}, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.Unit, product.ID)

	if err != nil {
		log.Println("(UpdateProduct:Exec)", err)
		return Products{}, err
	}

	return *product, nil

}
