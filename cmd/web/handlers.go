package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ashab-k/snippetbox/pkg/forms"
	"github.com/ashab-k/snippetbox/pkg/models"
)

func (app *application)home(w http.ResponseWriter, r *http.Request) {
	s , err := app.snippets.Latest()
	if err != nil {
		app.serverError(w , err)
	}


	app.render(w , r , "home.page.tmpl" , &templateData{ 
		Snippets: s ,
	})
}

func (app *application)showSnippet(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Query().Get(":id"))
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

	app.render(w ,r , "show.page.tmpl" , &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
    app.render(w,r , "create.page.tmpl" , &templateData{
		Form: forms.New(nil),
	})
}



func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

  err := r.ParseForm()

  if err != nil {
	 app.clientError(w, http.StatusBadRequest)
	 return
	
  }

  form := forms.New(r.PostForm)
  form.Required("title" , "content" , "expires")
  form.MaxLength("Title" ,100)
  form.PermittedValues("expires" ,"365" , "7" , "1")

	if  !form.Valid() {
        app.render(w, r, "create.page.tmpl", &templateData{Form: form})
        return
    }

    // Because the form data (with type url.Values) has been anonymously embedded
    // in the form.Form struct, we can use the Get() method to retrieve
    // the validated value for a particular form field.
    id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
    if err != nil {
        app.serverError(w, err)
        return
    }

	app.session.Put(r , "flash" , "Snippet Added Successfully")
    http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
        Form: forms.New(nil),
    })
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	 err := r.ParseForm()
	 if err != nil {
		app.clientError(w , http.StatusBadRequest)
	 }
	form := forms.New(r.PostForm)
	form.Required("name", "email" , "password")
	form.MatchesPattern("email" , forms.EmailRX)
	form.MinLength("password" , 10)
	
	if !form.Valid(){
		app.render(w , r, "signup.page.tmpl" , &templateData{
			Form: form,
		})
	}
	err = app.users.Insert(form.Get("name") , form.Get("email") , form.Get("password"))

	if err == models.ErrDuplicateEmail{
		form.Errors.Add("email" , "email already in use")
		app.render(w , r,  "signup.page.tmpl" , &templateData{
			Form: form,
		})
		return
	}else if err != nil {
		app.serverError(w , err)
		return 
	}
	
	app.session.Put(r , "flash" , "Your Signup was successful. Please Log in")
	http.Redirect(w , r ,"/user/login" , http.StatusSeeOther)
	

    fmt.Fprintln(w, "Create a new user...")
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w,r , "login.page.tmpl" , &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
    id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
    if err == models.ErrInvalidCredentials {
        form.Errors.Add("generic", "Email or Password is incorrect")
        app.render(w, r, "login.page.tmpl", &templateData{Form: form})
        return
    } else if err != nil {
        app.serverError(w, err)
        return
    }

    app.session.Put(r, "userID", id)

    http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)

}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r , "userID")
	app.session.Put(r , "flash" , "User logged out successfully")
	http.Redirect(w , r ,"/user/login" , http.StatusSeeOther)
}

func ping(w http.ResponseWriter , r *http.Request){
	w.Write([]byte("OK"))
}