package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *app) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// fileServer := http.FileServer(http.FS(ui.Files))
	router.ServeFiles("/assets/*filepath", http.Dir("assets"))

	dynamic := alice.New(app.sessionManager.LoadAndSave, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.homeView))
	router.Handler(http.MethodGet, "/home", dynamic.ThenFunc(app.homeView))
	router.Handler(http.MethodGet, "/auth/callback/:provider", dynamic.ThenFunc(app.oauthCallback))
	router.Handler(http.MethodGet, "/auth/login/:provider", dynamic.ThenFunc(app.oauth))

	// router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(h.UserSignupView))
	// router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(h.UserSignupPost))
	// router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(h.UserLoginView))
	// router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(h.UserLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/profile", protected.ThenFunc(app.homeView))
	//
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogout))
	//
	// router.Handler(http.MethodGet, "/clients", protected.ThenFunc(h.ClientsView))
	// router.Handler(http.MethodGet, "/clients/table", protected.ThenFunc(h.ClientsTable))
	// router.Handler(http.MethodGet, "/client/view/:id", protected.ThenFunc(h.ClientView))
	// router.Handler(http.MethodGet, "/client/edit/:id", protected.ThenFunc(h.ClientPutView))
	// router.Handler(http.MethodGet, "/client/create", protected.ThenFunc(h.ClientSetView))
	//
	// router.Handler(http.MethodPost, "/client", protected.ThenFunc(h.ClientPost))
	// router.Handler(http.MethodPut, "/client/:id", protected.ThenFunc(h.ClientPut))
	// router.Handler(http.MethodDelete, "/client/:id", protected.ThenFunc(h.ClientDelete))
	//
	// router.Handler(http.MethodGet, "/report/generate", protected.ThenFunc(h.ReportSet))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
