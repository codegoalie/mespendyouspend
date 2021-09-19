package actions

import (
	"database/sql"
	"errors"
	"fmt"
	"mespendyouspend/models"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
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
			err = createSpender(tx, user)
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

	// TODO: Start session?

	return c.Redirect(306, "rootPath()")
}

func createSpender(tx *pop.Connection, user goth.User) error {
	return tx.Create(models.Spender{
		Email: user.Email,
		Name:  user.Name,
	})
}
