package users

import (
	"errors"
	"log"

	"github.com/rfranzoia/cupuama-go/database"
	"github.com/rfranzoia/cupuama-go/utils"
)

var db = database.GetConnection()

// Get retrieve an non-deleted user by login
func (*Users) Get(login string) (Users, error) {

	query := app.SQLCache["users_get_login.sql"]
	stmt, err := db.Prepare(query)

	if err != nil {
		log.Fatal("(GetUser:Prepare)", err)
		return Users{}, err
	}

	defer stmt.Close()
	var user Users

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
		log.Println("(GetUser:QueryRow:Scan)", err)
		return Users{}, err
	}

	return user, nil
}

//List retrieves all non deleted users
func (*Users) List() ([]Users, error) {

	query := app.SQLCache["users_list.sql"]
	stmt, err := db.Prepare(query)

	if err != nil {
		log.Println("(ListUser:Prepare)", err)
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Println("(ListUser:Query)", err)
		return nil, err
	}

	defer rows.Close()

	var list []Users

	for rows.Next() {
		var user Users

		err := rows.Scan(&user.Login,
			&user.Password,
			&user.Person.FirstName,
			&user.Person.LastName,
			&user.Person.DateOfBirth,
			&user.Audit.Deleted,
			&user.Audit.DateCreated,
			&user.Audit.DateUpdated)

		if err != nil {
			log.Println("(ListUser:Scan)", err)
			return nil, err
		}

		list = append(list, user)
	}

	if err = rows.Err(); err != nil {
		log.Println("(ListUser:Rows)", err)
		return nil, err
	}

	return list, nil
}

// Create inserts a new user into the database
func (*Users) Create(user *Users) error {

	stmt, err := db.Prepare(
		"insert into users (login, password, first_name, last_name, date_of_birth) " +
			"values ($1, $2, $3, $4, $5)")

	if err != nil {
		log.Println("(CreateUser:Prepare)", err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.Login, utils.GetMD5Hash(user.Password), user.Person.FirstName, user.Person.LastName, user.Person.DateOfBirth)

	if err != nil {
		log.Println("(CreateUser:Exec)", err)
		return err
	}

	return nil

}

// Delete removes the user with the specified login from the database
func (*Users) Delete(login string) error {

	_, err := model.Get(login)
	if err != nil {
		log.Println("(DeleteUser:Get)", "User doesn't exist")
		err = errors.New("user doesn't exists")
		return err
	}

	stmt, err := db.Prepare(
		"delete from users where login = $1")

	if err != nil {
		log.Println("(DeleteUser:Prepare)", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(login)

	if err != nil {
		log.Println("(DeleteUser:Physical:Exec)", err)

		stmt, err = db.Prepare(
			"update users " +
				"set deleted = true, " +
				"date_updated = now() " +
				"where login = $1")

		if err != nil {
			log.Println("(UpdateUser:Logical:Prepare)", err)
			return err
		}

		_, err = stmt.Exec(login)

		if err != nil {
			log.Println("(DeleteUser:Logical:Exec)", err)
			return err
		}
	}

	return nil

}

// Update modify the data (except password) for the specified user
func (*Users) Update(user *Users) (Users, error) {

	_, err := model.Get(user.Login)
	if err != nil {
		log.Println("(UpdateUser:Get)", "User doesn't exist")
		return Users{}, err
	}

	stmt, err := db.Prepare(
		"update users " +
			"set date_of_birth = $1, " +
			"first_name = $2, " +
			"last_name = $3, " +
			"date_updated = now() " +
			"where login = $4")

	if err != nil {
		log.Println("(UpdateUser:Prepare)", err)
		return Users{}, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.Person.DateOfBirth, user.Person.FirstName, user.Person.LastName, user.Login)

	if err != nil {
		log.Println("(UpdateUser:Exec)", err)
		return Users{}, err
	}

	return *user, nil
	// TODO: implement update password
}
