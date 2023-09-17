package repository

import (
	"cupuama-go/config"
	"cupuama-go/domain"
	"cupuama-go/logger"
	"cupuama-go/utils"
	"database/sql"
	"errors"
	"fmt"
)

type UserRepository interface {
	GetByLogin(login string) (domain.Users, error)
	FindAll() ([]domain.Users, error)
	Create(user *domain.Users) error
	DeleteByLogin(login string) error
	UpdateByLogin(user *domain.Users) (domain.Users, error)
}

type UserRepositoryDB struct {
	db  *sql.DB
	app *config.AppConfig
}

func NewUserRepository(a *config.AppConfig) UserRepository {
	return &UserRepositoryDB{
		db:  a.DB,
		app: a,
	}
}

func (ur *UserRepositoryDB) GetByLogin(login string) (domain.Users, error) {

	query := ur.app.SQLCache["users_get_login.sql"]
	stmt, err := ur.db.Prepare(query)

	if err != nil {
		logger.Log.Fatal("(GetUser:Prepare)" + err.Error())
		return domain.Users{}, err
	}

	defer stmt.Close()
	var user domain.Users

	err = stmt.QueryRow(login).
		Scan(&user.Login,
			&user.Password,
			&user.Person.FirstName,
			&user.Person.LastName,
			&user.Person.DateOfBirth,
			&user.Audit.Deleted,
			&user.Audit.DateCreated,
			&user.Audit.DateUpdated)

	if err != nil {
		logger.Log.Info("(GetUser:QueryRow:Scan)" + err.Error())
		return domain.Users{}, err
	}

	return user, nil
}

// List retrieves all non deleted users
func (ur *UserRepositoryDB) FindAll() ([]domain.Users, error) {

	query := ur.app.SQLCache["users_list.sql"]
	stmt, err := ur.db.Prepare(query)

	if err != nil {
		logger.Log.Info("(ListUser:Prepare)" + err.Error())
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		logger.Log.Info("(ListUser:Query)" + err.Error())
		return nil, err
	}

	defer rows.Close()

	var list []domain.Users

	for rows.Next() {
		var user domain.Users

		err := rows.Scan(&user.Login,
			&user.Password,
			&user.Person.FirstName,
			&user.Person.LastName,
			&user.Person.DateOfBirth,
			&user.Audit.Deleted,
			&user.Audit.DateCreated,
			&user.Audit.DateUpdated)

		if err != nil {
			logger.Log.Info("(ListUser:Scan)" + err.Error())
			return nil, err
		}

		list = append(list, user)
	}

	if err = rows.Err(); err != nil {
		logger.Log.Info("(ListUser:Rows)" + err.Error())
		return nil, err
	}

	return list, nil
}

// Create inserts a new user into the database
func (ur *UserRepositoryDB) Create(user *domain.Users) error {

	stmt, err := ur.db.Prepare(
		"insert into users (login, password, first_name, last_name, date_of_birth) " +
			"values ($1, $2, $3, $4, $5)")

	if err != nil {
		logger.Log.Info("(CreateUser:Prepare)" + err.Error())
		return err
	}

	fmt.Println(user)

	defer stmt.Close()

	_, err = stmt.Exec(user.Login, utils.GetMD5Hash(user.Password), user.Person.FirstName, user.Person.LastName, user.Person.DateOfBirth)

	if err != nil {
		logger.Log.Info("(CreateUser:Exec)" + err.Error())
		return err
	}

	return nil

}

// Delete removes the user with the specified login from the database
func (ur *UserRepositoryDB) DeleteByLogin(login string) error {

	_, err := ur.GetByLogin(login)
	if err != nil {
		logger.Log.Info("(DeleteUser:Get) User doesn't exist")
		err = errors.New("user doesn't exists")
		return err
	}

	stmt, err := ur.db.Prepare(
		"delete from users where login = $1")

	if err != nil {
		logger.Log.Info("(DeleteUser:Prepare)" + err.Error())
	}

	defer stmt.Close()

	_, err = stmt.Exec(login)

	if err != nil {
		logger.Log.Info("(DeleteUser:Physical:Exec)" + err.Error())

		stmt, err = ur.db.Prepare(
			"update users " +
				"set deleted = true, " +
				"date_updated = now() " +
				"where login = $1")

		if err != nil {
			logger.Log.Info("(UpdateUser:Logical:Prepare)" + err.Error())
			return err
		}

		_, err = stmt.Exec(login)

		if err != nil {
			logger.Log.Info("(DeleteUser:Logical:Exec)" + err.Error())
			return err
		}
	}

	return nil

}

// Update modify the data (except password) for the specified user
func (ur *UserRepositoryDB) UpdateByLogin(user *domain.Users) (domain.Users, error) {

	_, err := ur.GetByLogin(user.Login)
	if err != nil {
		logger.Log.Info("(UpdateUser:Get) User doesn't exist")
		return domain.Users{}, err
	}

	stmt, err := ur.db.Prepare(
		"update users " +
			"set date_of_birth = $1, " +
			"first_name = $2, " +
			"last_name = $3, " +
			"date_updated = now() " +
			"where login = $4")

	if err != nil {
		logger.Log.Info("(UpdateUser:Prepare)" + err.Error())
		return domain.Users{}, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.Person.DateOfBirth, user.Person.FirstName, user.Person.LastName, user.Login)

	if err != nil {
		logger.Log.Info("(UpdateUser:Exec)" + err.Error())
		return domain.Users{}, err
	}

	return *user, nil
	// TODO: implement update password
}
