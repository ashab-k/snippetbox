package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/ashab-k/snippetbox/pkg/models"
)

func (app *application)home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        app.notFound(w)
        return
    }

	files := []string{"./ui/html/home.page.tmpl" , "./ui/html/base.layout.tmpl" , "./ui/html/footer.partial.tmpl"}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w , err)
		return
	}

	err = ts.Execute(w , nil)

	if err != nil {
		app.serverError(w ,err)
	}


    w.Write([]byte("Hello from Snippetbox"))
}

func (app *application)showSnippet(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id < 1 {
       app.notFound(w)
        return
    }

	s, err := app.snippets.Get(id)
    if err == models.ErrNoRecord {
        app.notFound(w)
        return
    } else if err != nil {
        app.serverError(w, err)
        return
    }

    // Write the snippet data as a plain-text HTTP response body.
    fmt.Fprintf(w, "%v", s)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        w.Header().Set("Allow", "POST")
        app.clientError(w ,http.StatusMethodNotAllowed)
        return
    }
	title := "O snail"
    content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
    expires := "7"

    // Pass the data to the SnippetModel.Insert() method, receiving the
    // ID of the new record back.
    id, err := app.snippets.Insert(title, content, expires)
    if err != nil {
        app.serverError(w, err)
        return
    }

    // Redirect the user to the relevant page for the snippet.
    http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
