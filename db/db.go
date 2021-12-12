package db

import (
	"database/sql"
	"os"
	"time"

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
	if _, err := os.Stat("./data/db.sqlite"); os.IsNotExist(err) {
		file, err := os.Create("data/db.sqlite")
		if err != nil {
			panic(err)
		}
		file.Close()
	}

	var err error
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
			efficiency_factor FLOAT,
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
		SELECT id, front, back, efficiency_factor, repititions, next_repitition
		FROM cards
		WHERE next_repitition < datetime('now')
		ORDER BY next_repitition ASC
	`)
	if err != nil {
		panic(err)
	}

	// convert rows to Cards
	cards := []janki.Card{}
	for rows.Next() {
		var card janki.Card
		err := rows.Scan(&card.ID, &card.Front, &card.Back, &card.EfficiencyFactor, &card.Repititions, &card.NextRepitition)
		if err != nil {
			panic(err)
		}
		cards = append(cards, card)
	}

	return cards
}

// GetCard gets a card from the database
func GetCard(id string) (janki.Card, error) {
	var card janki.Card
	err := db.QueryRow(`
		SELECT id, front, back, efficiency_factor, repititions, next_repitition
		FROM cards
		WHERE id = ?
	`, id).Scan(&card.ID, &card.Front, &card.Back, &card.EfficiencyFactor, &card.Repititions, &card.NextRepitition)
	if err != nil {
		return card, err
	}

	return card, nil
}

// CreateCard creates a new card in the database
func CreateCard(front, back string) error {
	// insert the card into the database
	_, err := db.Exec(`
		INSERT INTO cards (front, back, efficiency_factor, repititions, next_repitition)
		VALUES (?, ?, ?, ?, ?)
	`, front, back, 2.5, 0, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// UpdateCard updates the card in the database
func UpdateCard(card janki.Card) error {

	// update the card in the database
	_, err := db.Exec(`
		UPDATE cards
		SET front = ?, back = ?, efficiency_factor = ?, repititions = ?, next_repitition = ?
		WHERE id = ?
	`, card.Front, card.Back, card.EfficiencyFactor, card.Repititions, card.NextRepitition, card.ID)
	if err != nil {
		return err
	}

	return nil
}
