package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
)

// SQLite we use it as an struct which contains a pointer to manage the db
type SQLite struct {
	db *sql.DB
}

// getConnection will open a new connection to our database and return it
func getConnection(dbName string) (*sql.DB, error) {

	ChceckDBStatus(dbName)

	// Create the data base connection and validate if there is any error
	db, err := sql.Open("sqlite3", "./" + dbName)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("DB connection successfully!")
	}

	return db, nil
}

// ChceckDBStatus validates if exists or not the db file
func ChceckDBStatus(dbName string) {
	if _, err := os.Stat("./" + dbName); err != nil {
		log.Println("DB doesnt exists yet, creating ./" + dbName + " ...")
		file, err := os.Create("./" + dbName)
		if err != nil {
			log.Fatal(err.Error())
		}
		err = file.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Println("DB ./" + dbName + " created")
	}
}

// MakeMigration execute init schema into the database
func MakeMigration(db *sql.DB) {
	cmd := exec.Command("sqlite3 ./${DATABASE_NAME} < ../../internal/data/init.sql")

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
	//b, err := ioutil.ReadFile("./internal/data/init.sql")
	//if err != nil {
	//	return err
	//}
	//
	//rows, err := db.Query(string(b))
	//if err != nil {
	//	return err
	//}
	//
	//return rows.Close()
}
