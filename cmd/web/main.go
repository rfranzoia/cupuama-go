package main

import (
	"fmt"
	"log"

	"github.com/rfranzoia/cupuama-go/config"
	"github.com/rfranzoia/cupuama-go/database"
	"github.com/rfranzoia/cupuama-go/users"
	"github.com/rfranzoia/cupuama-go/utils"
)

func main() {
	defer database.GetConnection().Close()

	var app config.AppConfig

	// loads all queries into the application config cache
	qc, err := utils.CreateSQLCache()
	if err != nil {
		log.Fatal("cannot create queries cache")
	}

	app.SQLCache = qc
	app.UseCache = false

}

func userTesting() {
	user := users.Users{
		Login:    "Burgers",
		Password: "somthing22232",
		Person: users.Person{
			FirstName:   "Jose",
			LastName:    "da Silva",
			DateOfBirth: "1975-03-21",
		},
	}

	//users.Create(user)

	//users.Delete("UserLogin")
	user.Person.DateOfBirth = "1975-03-22"
	user.Login = "Burgers"
	users.Update(user)

	for _, u := range users.List() {
		fmt.Println("->", u)
	}

	//fmt.Println(users.Get("4372700044"))
}
