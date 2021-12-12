package db

import (
	"database/sql"
	"os"

	// sqlite
	_ "github.com/mattn/go-sqlite3"

	"github.com/jackloughran/janki/janki"
)

var db *sql.DB

func InitializeDB() {
	// create the data directory if it doesn't exist
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		os.Mkdir("./data", 0755)
	}

	// create the data/db.sqlite file if it doesn't exist
	file, err := os.Create("data/db.sqlite")
	if err != nil {
		panic(err)
	}
	file.Close()

	// open db at data/db.sqlite
	db, err = sql.Open("sqlite3", "data/db.sqlite")
	if err != nil {
		panic(err)
	}

	// create the cards table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cards (
			id INTEGER PRIMARY KEY,
			front TEXT,
			back TEXT,
			repititions INTEGER,
			next_repitition DATETIME
		)
	`)
	if err != nil {
		panic(err)
	}
}

func CardsToReview() []janki.Card {
	// get all cards where next_repitition is before now
	rows, err := db.Query(`
		SELECT id, front, back, repititions, next_repitition
		FROM cards
		WHERE next_repitition < datetime('now')
	`)
	if err != nil {
		panic(err)
	}

	// convert rows to Cards
	cards := []janki.Card{}
	for rows.Next() {
		var card janki.Card
		err := rows.Scan(&card.ID, &card.Front, &card.Back, &card.Repititions, &card.NextRepitition)
		if err != nil {
			panic(err)
		}
		cards = append(cards, card)
	}

	return cards
}
