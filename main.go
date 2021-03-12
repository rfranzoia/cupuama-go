package main

import (
	"fmt"
	"log"

	"github.com/rfranzoia/cupuama-go/config"
	"github.com/rfranzoia/cupuama-go/database"
	"github.com/rfranzoia/cupuama-go/orders"
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

	orders.NewOrderAPI(&app)

	// order := orders.OrderItemsStatus{
	// 	Order: orders.Orders{
	// 		TotalPrice: 150.0,
	// 	},
	// 	OrderItems: []orders.OrderItems{
	// 		orders.OrderItems{
	// 			Product: products.Products{
	// 				ID: 3,
	// 			},
	// 			Fruit: fruits.Fruits{
	// 				ID: 5,
	// 			},
	// 			Quantity:  10,
	// 			UnitPrice: 7.5,
	// 		},
	// 	},
	// }

	// orders.Create(order)

	// os := orders.OrderStatus{
	// 	Order: orders.Orders{
	// 		ID: 8,
	// 	},
	// 	Status:            1,
	// 	StatusDescription: "order-confirmed",
	// }

	// orders.CreateOrderStatus(os)

	// orders := orders.List()
	// for _, o := range orders {
	// 	fmt.Println("order:", o.Order, o.OrderStatus)
	// 	for _, i := range o.OrderItems {
	// 		fmt.Println("item:", i)
	// 	}
	// 	fmt.Println("---")
	// }

	order := orders.Get(4)
	fmt.Println("order:", order.Order, order.OrderStatus)
	for _, i := range order.OrderItems {
		fmt.Println("item:", i)
	}
	fmt.Println("---")

	// f := fruits.Fruits{
	// 	Name:     "Pera",
	// 	Harvest:  "Nunca",
	// 	Initials: "PERA",
	// }

	//fruits.Create(f)
	//fruits.Update(f)
	//fruits.Delete(6)

	// for _, f := range fruits.List() {
	// 	fmt.Println("->", f)
	// }
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
