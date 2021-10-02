package actions

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"mespendyouspend/models"
	"net/http"
	"os"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

const (
	tokenLength        = 35
	tokenCookieName    = "mespend-token"
	tokenValidDuration = 30 * 24 * time.Hour
)

func init() {
	gothic.Store = App().SessionStore

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/google/callback")),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/github/callback")),
	)
}

func AuthCallback(c buffalo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}
	// Do something with the user, maybe register them/sign them in
	spender := models.Spender{}
	tx := c.Value("tx").(*pop.Connection)
	err = tx.Where("email = ?", user.Email).First(&spender)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			spender, err = createSpender(tx, user)
			if err != nil {
				err = fmt.Errorf("failed to create spender: %w", err)
				return err
			}
			return c.Redirect(307, "transactionsPath()")
		} else {
			err = fmt.Errorf("failed to check for existing spender: %w", err)
			return err
		}
	}

	if user.Name != "" && user.Name != spender.Name {
		spender.Name = user.Name
		err = tx.Update(&spender)
		if err != nil {
			err = fmt.Errorf("failed to update spender: %w", err)
			return err
		}
	}

	token := GenerateSecureToken(tokenLength)
	c.Logger().Debug(len(token))
	err = tx.Create(&models.SpenderToken{
		SpenderID: spender.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(tokenValidDuration),
	})
	if err != nil {
		err = fmt.Errorf("failed to save auth token: %w", err)
		return err
	}

	ck := http.Cookie{
		Name:    tokenCookieName,
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(tokenValidDuration),
	}
	http.SetCookie(c.Response(), &ck)

	return c.Redirect(306, "rootPath()")
}

func createSpender(tx *pop.Connection, user goth.User) (models.Spender, error) {
	spender := models.Spender{
		Email: user.Email,
		Name:  user.Name,
	}
	err := tx.Create(&spender)
	return spender, err
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
