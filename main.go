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
	// Connection feedback
	fmt.Println("Connected!")

	// query multiple rows
	albums, err := albumsByArtist(db, "John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	// query single row
	alb, err := albumByID(db, 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)

	// db insert
	albID, err := addAlbum(db, Album{
		Title: "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price: 49.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", albID)
}

func albumsByArtist(db *sql.DB, name string) ([]Album, error){
	var albums []Album
	query := `SELECT * FROM album WHERE artist=$1`

	rows, err := db.Query(query, name)
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

func albumByID(db *sql.DB, id int64) (Album, error) {
	var alb Album
	query := `SELECT * FROM album WHERE id = $1`

	row := db.QueryRow(query, id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumByID %d: %v", id, err)
		}
		return alb, fmt.Errorf("albumsByID %d: %v", id, err)
	}
	return alb, nil
}

func addAlbum(db *sql.DB, alb Album) (int64, error) {
	insert := `INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id`

	var id int64
	err := db.QueryRow(insert, alb.Title, alb.Artist, alb.Price).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}

