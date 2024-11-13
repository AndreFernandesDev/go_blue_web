package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/AndreFernandesDev/boilerplate_web/internal/types"
	"github.com/a-h/templ"
)

func (app *app) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *app) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *app) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func render(w http.ResponseWriter, r *http.Request, status int, c templ.Component) error {
	w.WriteHeader(status)
	return c.Render(r.Context(), w)
}

func (app *app) generateViewData(r *http.Request) *types.ViewData {
	return &types.ViewData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		// CSRFToken:       nosurf.Token(r),
	}
}

// func (app *app) decodePostForm(r *http.Request, dst any) error {
// 	err := r.ParseForm()
// 	if err != nil {
// 		return err
// 	}
//
// 	err = app.formDecoder.Decode(dst, r.PostForm)
// 	if err != nil {
// 		var invalidDecoderError *form.InvalidDecoderError
//
// 		if errors.As(err, &invalidDecoderError) {
// 			panic(err)
// 		}
//
// 		return err
// 	}
//
// 	return nil
// }

func (app *app) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)

	if !ok {
		return false
	}

	return isAuthenticated
}
