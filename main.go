package main

import (
	"net/http"
	"text/template"

	"github.com/jackloughran/janki/db"
)

func main() {
	db.InitializeDB()

	http.HandleFunc("/", reviewHandler)
	http.HandleFunc("/assets/", assetsHandler)

	// start a web server on port 1234
	http.ListenAndServe(":1234", nil)
}

func reviewHandler(w http.ResponseWriter, r *http.Request) {
	cardsToReview := db.CardsToReview()

	// render the review template and send it cardsToReview
	renderTemplate(w, "review", cardsToReview)
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
