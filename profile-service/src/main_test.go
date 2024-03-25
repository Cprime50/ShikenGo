package src

import (
	"log"
	"os"
	"testing"

	"github.com/Cprime50/user/db"
)

func TestMain(m *testing.M) {

	log.Println("Running tests...")
	Db, err := db.ConnectTest()
	if err != nil {
		log.Fatal(err)
	}
	err = db.Migrate(Db)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Db.Close()

	os.Exit(m.Run())
}
