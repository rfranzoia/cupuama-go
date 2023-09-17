package repository

import (
	"cupuama-go/config"
	"cupuama-go/domain"
	"cupuama-go/logger"
	"database/sql"
	"errors"
)

type FruitRepository interface {
	Get(id int64) (domain.Fruits, error)
	List() ([]domain.Fruits, error)
	Create(fruit *domain.Fruits) (int64, error)
	Delete(id int64) error
	Update(fruit *domain.Fruits) (domain.Fruits, error)
}
type FruitRepositoryDB struct {
	db  *sql.DB
	app *config.AppConfig
}

func NewFruitRepository(a *config.AppConfig) FruitRepository {
	return &FruitRepositoryDB{
		db:  a.DB,
		app: a,
	}
}

// Get retrieve an non-deleted fruit by login
func (fr *FruitRepositoryDB) Get(id int64) (domain.Fruits, error) {

	query := fr.app.SQLCache["fruits_get_id.sql"]
	stmt, err := fr.db.Prepare(query)

	if err != nil {
		logger.Log.Fatal("(GetFruit:Prepare)" + err.Error())
		return domain.Fruits{}, err
	}

	defer stmt.Close()
	var fruit domain.Fruits

	err = stmt.QueryRow(id).
		Scan(&fruit.ID,
			&fruit.Name,
			&fruit.Harvest,
			&fruit.Initials,
			&fruit.Audit.Deleted,
			&fruit.Audit.DateCreated,
			&fruit.Audit.DateUpdated)

	if err != nil {
		logger.Log.Info("(GetFruit:QueryRow:Scan)")
		return domain.Fruits{}, err
	}

	return fruit, nil
}

// List retrieves all non deleted fruits
func (fr *FruitRepositoryDB) List() ([]domain.Fruits, error) {

	query := fr.app.SQLCache["fruits_list.sql"]
	stmt, err := fr.db.Prepare(query)

	if err != nil {
		logger.Log.Info("(ListFruit:Prepare)" + err.Error())
		empty := []domain.Fruits{}
		return empty, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		logger.Log.Info("(ListFruit:Query)" + err.Error())
		empty := []domain.Fruits{}
		return empty, err
	}

	defer rows.Close()

	var list []domain.Fruits

	for rows.Next() {
		var fruit domain.Fruits

		err := rows.Scan(&fruit.ID,
			&fruit.Name,
			&fruit.Harvest,
			&fruit.Initials,
			&fruit.Audit.Deleted,
			&fruit.Audit.DateCreated,
			&fruit.Audit.DateUpdated)

		if err != nil {
			logger.Log.Info("(ListFruit:Scan)" + err.Error())
			empty := []domain.Fruits{}
			return empty, err
		}

		list = append(list, fruit)
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Info("(ListFruit:Rows)" + err.Error())
		empty := []domain.Fruits{}
		return empty, err
	} else if len(list) == 0 {
		err = errors.New("no fruits were found")
		logger.Log.Info("(ListFruit:Result)" + err.Error())
		empty := []domain.Fruits{}
		return empty, err
	}

	return list, nil
}

// Create inserts a new fruit into the database and returns the new ID
func (fr *FruitRepositoryDB) Create(fruit *domain.Fruits) (int64, error) {

	stmt, err := fr.db.Prepare(
		"insert into fruits (name, harvest, initials, deleted, date_created) " +
			"values ($1, $2, $3, false, now()) returning id")

	if err != nil {
		logger.Log.Info("(CreateFruit:Prepare)" + err.Error())
		return -1, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(fruit.Name,
		fruit.Harvest,
		fruit.Initials).Scan(&fruit.ID)

	if err != nil {
		logger.Log.Info("(CreateFruit:Exec)" + err.Error())
		return -1, err
	}

	return fruit.ID, nil

}

// Delete removes the fruit with the specified login from the database
func (fr *FruitRepositoryDB) Delete(id int64) error {

	_, err := fr.Get(id)
	if err != nil {
		logger.Log.Info("(DeleteFruit:Get) Fruit doesn't exist")
		return err
	}

	stmt, err := fr.db.Prepare(
		"delete from fruits where id = $1")

	if err != nil {
		logger.Log.Info("(DeleteFruit:Prepare)" + err.Error())
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		logger.Log.Info("(DeleteFruit:Physical:Exec)" + err.Error())

		stmt, err = fr.db.Prepare(
			"update fruits " +
				"set deleted = true, " +
				"date_updated = now() " +
				"where id = $1")

		if err != nil {
			logger.Log.Info("(DeteleFruit:Logic:Prepare)" + err.Error())
			return err
		}

		_, err = stmt.Exec(id)

		if err != nil {
			logger.Log.Info("(DeleteFruit:Logic:Exec)" + err.Error())
			return err
		}

	}

	return nil
}

// Update modify the data for the specified fruit
func (fr *FruitRepositoryDB) Update(fruit *domain.Fruits) (domain.Fruits, error) {

	_, err := fr.Get(fruit.ID)
	if err != nil {
		logger.Log.Info("(UpdateFruit:Get) Fruit doesn't exist")
		return domain.Fruits{}, err
	}

	stmt, err := fr.db.Prepare(
		"update fruits " +
			"set name = $1, " +
			"harvest = $2, " +
			"initials = $3, " +
			"date_updated = now() " +
			"where id = $4")

	if err != nil {
		logger.Log.Info("(UpdateFruit:Prepare)" + err.Error())
		return domain.Fruits{}, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(fruit.Name, fruit.Harvest, fruit.Initials, fruit.ID)

	if err != nil {
		logger.Log.Info("(UpdateFruit:Exec)" + err.Error())
		return domain.Fruits{}, err
	}

	return *fruit, nil

}
