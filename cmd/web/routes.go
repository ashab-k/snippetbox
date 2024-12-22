package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic , app.logRequest , secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable)
	mux := pat.New()
	mux.Get("/" , dynamicMiddleware.ThenFunc(app.home));
	mux.Post("/snippet/create" , dynamicMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/create" , dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Get("/snippet/:id" , dynamicMiddleware.ThenFunc(app.showSnippet))

	fileserver := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/" , http.StripPrefix("/static" , fileserver))
	return standardMiddleware.Then(mux)
}