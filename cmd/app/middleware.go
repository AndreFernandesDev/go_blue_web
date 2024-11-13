package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *app) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *app) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *app) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := app.sessionManager.GetString(r.Context(), "authenticatedUserID")
		if len(id) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		uid, err := uuid.Parse(id)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		exists, err := app.DB.CheckUserExists(r.Context(), uid)
		if err != nil {
			app.serverError(w, err)
			return
		}

		if exists {
			app.infoLog.Println("Renew")
			app.sessionManager.RenewToken(r.Context())
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func (app *app) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

// func noSurf(next http.Handler) http.Handler {
// 	crsfHandler := nosurf.New(next)
// 	crsfHandler.SetBaseCookie(http.Cookie{
// 		HttpOnly: true,
// 		Path:     "/",
// 		Secure:   true,
// 	})
//
// 	return crsfHandler
// }
