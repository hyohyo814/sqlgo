package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}


func main() {
	// env
	var db *sql.DB
	var (
		User   = os.Getenv("DBUSER")
		Passwd = os.Getenv("DBPASS")
		DBName = "recordings"
	)

	var err error
	psqlInfo := fmt.Sprintf("user=%s password=%s dbname=%s", User, Passwd, DBName)

	// Open db conn
	db, err = sql.Open("postgres", psqlInfo)

	// Err handling
	if err != nil {
		log.Fatal("err opening connection", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal("DB not found!", pingErr)
	}

	fmt.Println("Connected!")
	data, err := albumsByArtist(db, "John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(data)
}

func albumsByArtist(db *sql.DB, name string) ([]Album, error){
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}
