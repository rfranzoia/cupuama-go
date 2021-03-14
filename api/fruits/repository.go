package fruits

import (
	"errors"
	"log"

	"github.com/rfranzoia/cupuama-go/database"
)

var db = database.GetConnection()

// Get retrieve an non-deleted fruit by login
func (*Fruits) Get(id int64) (Fruits, error) {

	query := app.SQLCache["fruits_get_id.sql"]
	stmt, err := db.Prepare(query)

	if err != nil {
		log.Fatal("(GetFruit:Prepare)", err)
		return Fruits{}, err
	}

	defer stmt.Close()
	var fruit Fruits

	err = stmt.QueryRow(id).
		Scan(&fruit.ID,
			&fruit.Name,
			&fruit.Harvest,
			&fruit.Initials,
			&fruit.Audit.Deleted,
			&fruit.Audit.DateCreated,
			&fruit.Audit.DateUpdated)

	if err != nil {
		log.Println("(GetFruit:QueryRow:Scan)", err)
		return Fruits{}, err
	}

	return fruit, nil
}

//List retrieves all non deleted fruits
func (*Fruits) List() ([]Fruits, error) {

	query := app.SQLCache["fruits_list.sql"]
	stmt, err := db.Prepare(query)

	if err != nil {
		log.Println("(ListFruit:Prepare)", err)
		empty := []Fruits{}
		return empty, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Println("(ListFruit:Query)", err)
		empty := []Fruits{}
		return empty, err
	}

	defer rows.Close()

	var list []Fruits

	for rows.Next() {
		var fruit Fruits

		err := rows.Scan(&fruit.ID,
			&fruit.Name,
			&fruit.Harvest,
			&fruit.Initials,
			&fruit.Audit.Deleted,
			&fruit.Audit.DateCreated,
			&fruit.Audit.DateUpdated)

		if err != nil {
			log.Println("(ListFruit:Scan)", err)
			empty := []Fruits{}
			return empty, err
		}

		list = append(list, fruit)
	}

	err = rows.Err()
	if err != nil {
		log.Println("(ListFruit:Rows)", err)
		empty := []Fruits{}
		return empty, err
	} else if len(list) == 0 {
		err = errors.New("No fruits were found")
		log.Println("(ListFruit:Result)", err)
		empty := []Fruits{}
		return empty, err
	}

	return list, nil
}

// Create inserts a new fruit into the database and returns the new ID
func (*Fruits) Create(fruit *Fruits) (int64, error) {

	stmt, err := db.Prepare(
		"insert into fruits (name, harvest, initials, date_created) " +
			"values ($1, $2, $3, now()) returning id")

	if err != nil {
		log.Println("(CreateFruit:Prepare)", err)
		return -1, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(fruit.Name,
		fruit.Harvest,
		fruit.Initials).Scan(&fruit.ID)

	if err != nil {
		log.Println("(CreateFruit:Exec)", err)
		return -1, err
	}

	return fruit.ID, nil

}

// Delete removes the fruit with the specified login from the database
func (*Fruits) Delete(id int64) error {

	_, err := model.Get(id)
	if err != nil {
		log.Println("(DeleteFruit:Get)", "Fruit doesn't exist")
		return err
	}

	stmt, err := db.Prepare(
		"delete from fruits where id = $1")

	if err != nil {
		log.Println("(DeleteFruit:Prepare)", err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		log.Println("(DeleteFruit:Physical:Exec)", err)

		stmt, err = db.Prepare(
			"update fruits " +
				"set deleted = true, " +
				"date_updated = now() " +
				"where id = $1")

		if err != nil {
			log.Println("(DeteleFruit:Logic:Prepare)", err)
			return err
		}

		_, err = stmt.Exec(id)

		if err != nil {
			log.Println("(DeleteFruit:Logic:Exec)", err)
			return err
		}

	}

	return nil
}

// Update modify the data for the specified fruit
func (*Fruits) Update(fruit *Fruits) (Fruits, error) {

	_, err := model.Get(fruit.ID)
	if err != nil {
		log.Println("(UpdateFruit:Get)", "Fruit doesn't exist")
		return Fruits{}, err
	}

	stmt, err := db.Prepare(
		"update fruits " +
			"set name = $1, " +
			"harvest = $2, " +
			"initials = $3, " +
			"date_updated = now() " +
			"where id = $4")

	if err != nil {
		log.Println("(UpdateFruit:Prepare)", err)
		return Fruits{}, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(fruit.Name, fruit.Harvest, fruit.Initials, fruit.ID)

	if err != nil {
		log.Println("(UpdateFruit:Exec)", err)
		return Fruits{}, err
	}

	return *fruit, nil

}
