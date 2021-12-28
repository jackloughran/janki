package main

import (
	"html/template"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/jackloughran/janki/db"
	"github.com/jackloughran/janki/janki"
	"github.com/jackloughran/janki/supermemo"
)

func main() {
	db.InitializeDB()

	http.HandleFunc("/", nextReviewHandler)
	http.HandleFunc("/assets/", assetsHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/review/", reviewHandler)
	http.HandleFunc("/edit", editHandler)

	go func() {
		time.Sleep(3 * time.Second)

		err := exec.Command("open", "http://localhost:1234").Start()

		if err != nil {
			panic(err)
		}
	}()

	// start a web server on port 1234
	http.ListenAndServe(":1234", nil)
}

func nextReviewHandler(w http.ResponseWriter, r *http.Request) {
	cardsToReview := db.CardsToReview()

	if len(cardsToReview) > 0 {
		http.Redirect(w, r, "/review/"+string(cardsToReview[0].ID), http.StatusFound)
		return
	}

	// render the review template and send it cardsToReview
	renderTemplate(w, "nocards", nil)
}

type reviewData struct {
	Card    janki.Card
	Flipped bool
}

// reviewHandler handles reviewing a single card
func reviewHandler(w http.ResponseWriter, r *http.Request) {
	// get the card ID from the URL
	cardID := r.URL.Path[len("/review/"):]

	// get the card from the database
	card, err := db.GetCard(cardID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if quality query param is passed, update the card
	if quality := r.URL.Query().Get("quality"); quality != "" {
		// convert quality to int
		qualityInt, err := strconv.Atoi(quality)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newCard := supermemo.AdjustCard(card, qualityInt)
		db.UpdateCard(newCard)

		// redirect to the next card
		http.Redirect(w, r, "/", http.StatusFound)
	}

	// flipped is true if the ?flipped=true query parameter is present
	flipped := r.URL.Query().Get("flipped") == "true"

	// render the review template and send it the card
	renderTemplate(w, "review", reviewData{card, flipped})
}

// createHandler renders the create template if it's a GET request, and creates a card if it's a POST request
func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// render the create template
		renderTemplate(w, "create", nil)
	} else {
		// create a new card
		err := db.CreateCard(r.FormValue("front"), r.FormValue("back"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// redirect to the create page to have another go
		http.Redirect(w, r, "/create", http.StatusFound)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		cards, err := db.GetAllCards()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		renderTemplate(w, "edit", cards)
	} else {
		card, err := db.GetCard(r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		card.Front = r.FormValue("front")
		card.Back = r.FormValue("back")

		err = db.UpdateCard(card)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/edit", http.StatusFound)
	}
}

// assetsHandler serves static assets from the "assets" directory
func assetsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

// renderTemplate renders a template with the given name and data
func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	// get the template
	tmpl, err := template.ParseFiles("templates/" + name + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// render the template
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
