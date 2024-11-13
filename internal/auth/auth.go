package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/twitterv2"
)

const (
	MaxAge = 86400 * 30
)

func NewAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env := os.Getenv("ENV")
	if env == "" {
		log.Fatal("ENV is not found in the environment")
	}

	url := os.Getenv("APP_URL")
	if env == "" {
		log.Fatal("APP_URL is not found in the environment")
	}

	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		log.Fatal("SESSION_SECRET is not found in the environment")
	}

	store := sessions.NewCookieStore([]byte(secret))

	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = env == "PROD"
	store.Options.SameSite = http.SameSiteLaxMode

	discordClientId := os.Getenv("DISCORD_CLIENT_ID")
	discordClientSecret := os.Getenv("DISCORD_CLIENT_SECRET")

	gothic.Store = store

	fmt.Println(os.Getenv("TWITTER_CLIENT_ID"))
	fmt.Println(os.Getenv("TWITTER_CLIENT_SECRET"))

	goth.UseProviders(
		discord.New(discordClientId, discordClientSecret, url+"/auth/callback/discord", discord.ScopeIdentify, discord.ScopeEmail),
		twitterv2.New(os.Getenv("TWITTER_CLIENT_ID"), os.Getenv("TWITTER_CLIENT_SECRET"), url+"/auth/callback/twitterv2"),
	)
}
