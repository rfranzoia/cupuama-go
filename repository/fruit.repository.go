package repository

import (
	"cupuama-go/config"
	"cupuama-go/domain"
	"cupuama-go/logger"
	"errors"
	"github.com/jmoiron/sqlx"
)

type FruitRepository interface {
	Create(fruit *domain.Fruits) (int64, error)
	Get(id int64) (domain.Fruits, error)
	List() ([]domain.Fruits, error)
	Update(fruit *domain.Fruits) (domain.Fruits, error)
	Delete(id int64) error
}
type FruitRepositoryDB struct {
	db  *sqlx.DB
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
	var fruit domain.Fruits
	err := fr.db.Get(&fruit, query, id)

	if err != nil {
		logger.Log.Fatal("(GetFruit)" + err.Error())
		return domain.Fruits{}, err
	}

	return fruit, nil
}

// List retrieves all non deleted fruits
func (fr *FruitRepositoryDB) List() ([]domain.Fruits, error) {
	var list []domain.Fruits
	query := fr.app.SQLCache["fruits_list.sql"]
	err := fr.db.Select(&list, query)
	if err != nil {
		logger.Log.Info("(ListFruit)" + err.Error())
		return []domain.Fruits{}, err

	} else if len(list) == 0 {
		err = errors.New("no fruits were found")
		logger.Log.Info("(ListFruit:Result)" + err.Error())
		return list, err
	}

	return list, nil
}

// Create inserts a new fruit into the database and returns the new ID
func (fr *FruitRepositoryDB) Create(fruit *domain.Fruits) (int64, error) {

	insertQuery := fr.app.SQLCache["fruits_insert.sql"]
	stmt, err := fr.db.Prepare(insertQuery)

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

// Update modify the data for the specified fruit
func (fr *FruitRepositoryDB) Update(fruit *domain.Fruits) (domain.Fruits, error) {

	// check if the fruit ID is valid
	_, err := fr.Get(fruit.ID)
	if err != nil {
		logger.Log.Info("(UpdateFruit:Get) Fruit doesn't exist")
		return domain.Fruits{}, err
	}

	updateQuery := fr.app.SQLCache["fruits_update.sql"]
	_, err = fr.db.Exec(updateQuery, fruit.Name, fruit.Harvest, fruit.Initials, fruit.Audit.Deleted, fruit.ID)

	if err != nil {
		logger.Log.Info("(UpdateFruit:Exec) " + err.Error())
		return domain.Fruits{}, err
	}

	return *fruit, nil

}

// Delete removes the fruit with the specified login from the database
func (fr *FruitRepositoryDB) Delete(id int64) error {

	// check if the fruit ID is valid
	fruit, err := fr.Get(id)
	if err != nil {
		logger.Log.Info("(DeleteFruit:Get) Fruit doesn't exist")
		return err
	}

	deleteQuery := fr.app.SQLCache["fruits_delete.sql"]
	_, err = fr.db.Exec(deleteQuery, id)

	if err != nil {
		logger.Log.Info("(DeleteFruit:Physical:Exec) Fail " + err.Error())

		fruit.Audit.Deleted = true
		_, err := fr.Update(&fruit)
		if err != nil {
			logger.Log.Info("(DeleteFruit:Logic:Exec) Fail " + err.Error())
			return err
		}

	}

	return nil
}
