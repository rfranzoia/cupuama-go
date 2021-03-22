package fruits

import (
	"log"
	"strings"
	"testing"

	"github.com/rfranzoia/cupuama-go/config"
	"github.com/rfranzoia/cupuama-go/utils"
)

var testFruit Fruits

func init() {

	var app config.AppConfig

	// loads all queries into the application config cache
	qc, err := utils.CreateSQLCache("../../queries/*.sql")
	if err != nil {
		log.Fatal("cannot create queries cache")
	}

	app.SQLCache = qc
	app.UseCache = false

	NewFruitService(&app)

	rndName := utils.NewUUID()
	rndName = utils.Substring(rndName, 0, 10)
	rndInit := strings.ToUpper(utils.Substring(rndName, 0, 4))

	testFruit = Fruits{
		Name:     rndName,
		Initials: rndInit,
		Harvest:  "All Year",
	}
}

func TestCreate(t *testing.T) {
	_, err := testFruit.Create(&testFruit)
	if err != nil {
		t.Errorf("(TestCreate) cannot create fruit: %v", err)
	}
}

func TestGet(t *testing.T) {
	id, err := testFruit.Create(&testFruit)
	if err != nil {
		t.Errorf("(TestGet) cannot create fruit: %v", err)
	}

	_, err = testFruit.Get(id)
	if err != nil {
		t.Errorf("(TestGet) error while searching for fruit: %v", err)
	}

}

func TestList(t *testing.T) {
	list, err := testFruit.List()
	if err != nil {
		t.Errorf("(TestList) listing error: %v", err)

	} else if len(list) == 0 {
		t.Errorf("(TestList) list of fruits should not be empty")
	}
}

func TestDelete(t *testing.T) {
	id, err := testFruit.Create(&testFruit)
	if err != nil {
		t.Errorf("(TestDelete) cannot create fruit: %v", err)
	}

	err = testFruit.Delete(id)
	if err != nil {
		t.Errorf("(TestDelete) error removing fruit: %v", err)
	}
}

func TestUpdate(t *testing.T) {

	id, err := testFruit.Create(&testFruit)
	if err != nil {
		t.Errorf("(TestUpdate) cannot create fruit: %v", err)
	}

	f, err := testFruit.Get(id)
	if err != nil {
		t.Errorf("(TestUpdate) error searching for fruit: %v", err)
	}

	f.Harvest = "January"

	f, err = testFruit.Update(&f)
	if err != nil {
		t.Errorf("(TestUpdate) error updating fruit: %v", err)
	}

	f, err = testFruit.Get(id)
	if err != nil {
		t.Errorf("(TestUpdate) error searching for fruit: %v", err)

	} else if f.Harvest == testFruit.Harvest {
		t.Errorf("(TestUpdate) update fruit fail")
	}
}
