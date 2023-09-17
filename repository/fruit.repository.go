package repository

import (
	"cupuama-go/config"
	"cupuama-go/domain"
	"cupuama-go/errors"
	"cupuama-go/logger"
	"github.com/jmoiron/sqlx"
)

type FruitRepository interface {
	Create(fruit *domain.Fruits) (int64, *errors.AppError)
	Get(id int64) (domain.Fruits, *errors.AppError)
	List() ([]domain.Fruits, *errors.AppError)
	Update(fruit *domain.Fruits) (domain.Fruits, *errors.AppError)
	Delete(id int64) *errors.AppError
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
func (fr *FruitRepositoryDB) Get(id int64) (domain.Fruits, *errors.AppError) {

	query := fr.app.SQLCache["fruits_get_id.sql"]
	var fruit domain.Fruits

	err := fr.db.Get(&fruit, query, id)

	if err != nil {
		logger.Log.Error("(GetFruit)" + err.Error())
		return domain.Fruits{}, errors.NotFoundError("Fruit not found")
	}

	return fruit, nil
}

// List retrieves all non deleted fruits
func (fr *FruitRepositoryDB) List() ([]domain.Fruits, *errors.AppError) {
	var list []domain.Fruits
	query := fr.app.SQLCache["fruits_list.sql"]
	err := fr.db.Select(&list, query)
	if err != nil {
		logger.Log.Info("(ListFruit)" + err.Error())
		return []domain.Fruits{}, errors.UnexpectedError(err.Error())

	} else if len(list) == 0 {
		logger.Log.Info("(ListFruit:Result) no fruits found " + err.Error())
		return list, errors.NotFoundError("no fruits were found")
	}

	return list, nil
}

// Create inserts a new fruit into the database and returns the new ID
func (fr *FruitRepositoryDB) Create(fruit *domain.Fruits) (int64, *errors.AppError) {

	insertQuery := fr.app.SQLCache["fruits_insert.sql"]
	err := fr.db.QueryRow(insertQuery, fruit.Name, fruit.Harvest, fruit.Initials).Scan(&fruit.ID)

	if err != nil {
		logger.Log.Info("(CreateFruit:Exec)" + err.Error())
		return -1, errors.UnexpectedError("Error while creating a fruit")
	}

	return fruit.ID, nil

}

// Update modify the data for the specified fruit
func (fr *FruitRepositoryDB) Update(fruit *domain.Fruits) (domain.Fruits, *errors.AppError) {

	// check if the fruit ID is valid
	_, err := fr.Get(fruit.ID)
	if err != nil {
		logger.Log.Info("(UpdateFruit:Get) Fruit doesn't exist")
		return domain.Fruits{}, errors.NotFoundError("Fruit not found")
	}

	updateQuery := fr.app.SQLCache["fruits_update.sql"]
	_, err2 := fr.db.Exec(updateQuery, fruit.Name, fruit.Harvest, fruit.Initials, fruit.Audit.Deleted, fruit.ID)
	if err2 != nil {
		logger.Log.Info("(UpdateFruit:Exec) " + err2.Error())
		return domain.Fruits{}, errors.UnexpectedError("Error updating a fruit")
	}

	return *fruit, nil

}

// Delete removes the fruit with the specified login from the database
func (fr *FruitRepositoryDB) Delete(id int64) *errors.AppError {

	// check if the fruit ID is valid
	fruit, err := fr.Get(id)
	if err != nil {
		logger.Log.Info("(DeleteFruit:Get) Fruit doesn't exist")
		return err
	}

	deleteQuery := fr.app.SQLCache["fruits_delete.sql"]
	_, err2 := fr.db.Exec(deleteQuery, id)

	if err2 != nil {
		logger.Log.Info("(DeleteFruit:Physical:Exec) Fail " + err2.Error())

		fruit.Audit.Deleted = true
		_, err := fr.Update(&fruit)
		if err != nil {
			logger.Log.Info("(DeleteFruit:Logic:Exec) Fail " + err.Message)
			return err
		}

	}

	return nil
}
