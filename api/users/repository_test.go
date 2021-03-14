package users

import (
	"log"
	"testing"

	"github.com/rfranzoia/cupuama-go/config"
	"github.com/rfranzoia/cupuama-go/utils"
)

var testUser Users

func init() {

	var app config.AppConfig

	// loads all queries into the application config cache
	qc, err := utils.CreateSQLCache("../../queries/*.sql")
	if err != nil {
		log.Fatal("cannot create queries cache")
	}

	app.SQLCache = qc
	app.UseCache = false

	NewUserService(&app)

	testUser = Users{
		Login:    utils.NewUUID(),
		Password: utils.GetMD5Hash("(Default12345.)"),
		Person: Person{
			FirstName:   "Jose",
			LastName:    "da Silva",
			DateOfBirth: "1975-03-21",
		},
	}
}

func TestCreate(t *testing.T) {
	err := testUser.Create(&testUser)
	if err != nil {
		t.Errorf("(TestCreate) cannot create user: %v", err)
	}
}

func TestGet(t *testing.T) {
	login := testUser.Login
	err := testUser.Create(&testUser)
	if err != nil {
		t.Errorf("(TestGet) cannot create user: %v", err)
	}

	_, err = testUser.Get(login)
	if err != nil {
		t.Errorf("(TestGet) error while searching for user: %v", err)
	}

}

func TestList(t *testing.T) {
	list, err := testUser.List()
	if err != nil {
		t.Errorf("(TestList) listing error: %v", err)

	} else if len(list) == 0 {
		t.Errorf("(TestList) list of users should not be empty")
	}
}

func TestDelete(t *testing.T) {
	login := testUser.Login
	err := testUser.Create(&testUser)
	if err != nil {
		t.Errorf("(TestDelete) cannot create user: %v", err)
	}

	err = testUser.Delete(login)
	if err != nil {
		t.Errorf("(TestDelete) error removing user: %v", err)
	}
}

func TestUpdate(t *testing.T) {

	login := testUser.Login
	err := testUser.Create(&testUser)
	if err != nil {
		t.Errorf("(TestUpdate) cannot create user: %v", err)
	}

	u, err := testUser.Get(login)
	if err != nil {
		t.Errorf("(TestUpdate) error while searching for user: %v", err)
	}

	u.Person.DateOfBirth = "1975-03-22"

	u, err = testUser.Update(&u)
	if err != nil {
		t.Errorf("(TestUpdate) error updating user: %v", err)
	}

	u, err = testUser.Get(login)
	if err != nil {
		t.Errorf("(TestUpdate) error while searching for user: %v", err)

	} else if u.Person.DateOfBirth == testUser.Person.DateOfBirth {
		t.Errorf("(TestUpdate) update user fail")
	}
}
