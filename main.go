package main

import (
	"net/http"
	"text/template"

	"github.com/jackloughran/janki/db"
	"github.com/jackloughran/janki/janki"
)

func main() {
	db.InitializeDB()

	http.HandleFunc("/", reviewHandler)
	http.HandleFunc("/assets/", assetsHandler)
	http.HandleFunc("/create", createHandler)

	// start a web server on port 1234
	http.ListenAndServe(":1234", nil)
}

type reviewData struct {
	Cards   []janki.Card
	Flipped bool
}

func reviewHandler(w http.ResponseWriter, r *http.Request) {
	cardsToReview := db.CardsToReview()

	// flipped is true if the ?flipped=true query parameter is present
	flipped := r.URL.Query().Get("flipped") == "true"

	// render the review template and send it cardsToReview
	renderTemplate(w, "review", reviewData{cardsToReview, flipped})
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
