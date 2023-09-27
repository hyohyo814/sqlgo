package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	var (
		User   = os.Getenv("DBUSER")
		Passwd = os.Getenv("DBPASS")
		DBName = "recordings"
	)

	var err error
	psqlInfo := fmt.Sprintf("user=%s password=%s dbname=%s", User, Passwd, DBName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")
}
