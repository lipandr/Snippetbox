package main

import (
	"errors"
	"fmt"
	"github.com/lipandr/Snippetbox/pkg/models"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := &templateData{Snippets: s}

	app.render(w, r, "home.page.tmpl", data)
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data := &templateData{Snippet: s}

	app.render(w, r, "show.page.tmpl", data)
}

// Add a new createSnippetForm handler, which for now returns a placeholder response.
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	formErrors := make(map[string]string)
	if strings.TrimSpace(title) == "" {
		formErrors["title"] = "Title can't be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		formErrors["title"] = "Title is too long (maximum is 100 characters)"
	}
	if strings.TrimSpace(content) == "" {
		formErrors["content"] = "This field cannot be blank"
	}

	if strings.TrimSpace(expires) == "" {
		formErrors["expires"] = "This field cannot be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		formErrors["expires"] = "This field is invalid"
	}

	if len(formErrors) > 0 {
		app.render(w, r, "create.page.tmpl", &templateData{
			FormErrors: formErrors,
			FormData:   r.PostForm,
		})
		return
	}

	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
