package products

import (
	"log"
	"testing"

	"github.com/rfranzoia/cupuama-go/config"
	"github.com/rfranzoia/cupuama-go/utils"
)

var testProduct Products

func init() {

	var app config.AppConfig

	// loads all queries into the application config cache
	qc, err := utils.CreateSQLCache("../../queries/*.sql")
	if err != nil {
		log.Fatal("cannot create queries cache")
	}

	app.SQLCache = qc
	app.UseCache = false

	NewProductService(&app)

	rndName := utils.NewUUID()
	rndName = utils.Substring(rndName, 0, 10)

	testProduct = Products{
		Name: rndName,
		Unit: "UN",
	}
}

func TestCreate(t *testing.T) {
	_, err := testProduct.Create(&testProduct)
	if err != nil {
		t.Errorf("(TestCreate) cannot create product: %v", err)
	}
}

func TestGet(t *testing.T) {
	id, err := testProduct.Create(&testProduct)
	if err != nil {
		t.Errorf("(TestGet) cannot create product: %v", err)
	}

	_, err = testProduct.Get(id)
	if err != nil {
		t.Errorf("(TestGet) error while searching for product: %v", err)
	}

}

func TestList(t *testing.T) {
	list, err := testProduct.List()
	if err != nil {
		t.Errorf("(TestList) listing error: %v", err)

	} else if len(list) == 0 {
		t.Errorf("(TestList) list of products should not be empty")
	}
}

func TestDelete(t *testing.T) {
	id, err := testProduct.Create(&testProduct)
	if err != nil {
		t.Errorf("(TestDelete) cannot create product: %v", err)
	}

	err = testProduct.Delete(id)
	if err != nil {
		t.Errorf("(TestDelete) error removing product: %v", err)
	}
}

func TestUpdate(t *testing.T) {

	id, err := testProduct.Create(&testProduct)
	if err != nil {
		t.Errorf("(TestUpdate) cannot create product: %v", err)
	}

	f, err := testProduct.Get(id)
	if err != nil {
		t.Errorf("(TestUpdate) error searching for product: %v", err)
	}

	f.Unit = "PKG"

	f, err = testProduct.Update(&f)
	if err != nil {
		t.Errorf("(TestUpdate) error updating product: %v", err)
	}

	f, err = testProduct.Get(id)
	if err != nil {
		t.Errorf("(TestUpdate) error searching for product: %v", err)

	} else if f.Unit == testProduct.Unit {
		t.Errorf("(TestUpdate) update product fail")
	}
}
