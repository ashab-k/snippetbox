package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic , app.logRequest , secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable , noSurf , app.authenticate)
	mux := pat.New()
	mux.Get("/" , dynamicMiddleware.ThenFunc(app.home));
	mux.Post("/snippet/create" , dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createSnippet))
	mux.Get("/snippet/create" , dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createSnippetForm))
	mux.Get("/snippet/:id" , dynamicMiddleware.ThenFunc(app.showSnippet))
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
    mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
    mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
    mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
    mux.Post("/user/logout", dynamicMiddleware.ThenFunc(app.logoutUser))

	 mux.Get("/ping", http.HandlerFunc(ping))

	fileserver := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/" , http.StripPrefix("/static" , fileserver))
	return standardMiddleware.Then(mux)
}