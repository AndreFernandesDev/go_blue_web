package main

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/AndreFernandesDev/boilerplate_web/internal/components"
	"github.com/AndreFernandesDev/boilerplate_web/internal/database"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/markbates/goth/gothic"
)

func (app *app) homeView(w http.ResponseWriter, r *http.Request) {
	v := app.generateViewData(r)
	render(w, r, 200, components.Home(v))
}

func (app *app) oauth(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	provider := params.ByName("provider")

	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	gothic.BeginAuthHandler(w, r)
}

func (app *app) oauthCallback(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	provider := params.ByName("provider")

	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		app.serverError(w, err)
		return
	}

	account, err := app.DB.GetAccount(r.Context(), user.UserID)
	if err != nil && err != sql.ErrNoRows {
		app.serverError(w, err)
		return
	}

	if err == sql.ErrNoRows {
		newUser, err := app.DB.SetUser(r.Context(), database.SetUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Username:  user.NickName,
			Firstname: user.FirstName,
			Lastname:  user.LastName,
			Email:     user.Email,
			Password:  "",
			AvatarUrl: user.AvatarURL,
		})
		if err != nil {
			app.serverError(w, err)
			return
		}

		newAccount, err := app.DB.SetAccount(r.Context(), database.SetAccountParams{
			ID:           uuid.New(),
			UserID:       newUser.ID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Provider:     user.Provider,
			ProviderID:   user.UserID,
			AccessToken:  user.AccessToken,
			RefreshToken: user.RefreshToken,
			ExpiresAt:    user.ExpiresAt,
		})

		if err != nil {
			app.serverError(w, err)
			return
		}

		app.sessionManager.Put(r.Context(), "authenticatedUserID", newAccount.UserID.String())
		app.sessionManager.Put(r.Context(), "flash", "You've been registered in successfully!")
	} else {
		app.infoLog.Println("Existing account:")
		app.infoLog.Println(account)
		app.sessionManager.Put(r.Context(), "authenticatedUserID", account.UserID.String())
		app.sessionManager.Put(r.Context(), "flash", "You've been logged in successfully!")
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)

}

func (app *app) userLogout(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
