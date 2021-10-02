package actions

import (
	"fmt"
	"mespendyouspend/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
)

// AuthenticateUser authenticates users before handling requests
func AuthenticateUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		spender, err := IsAuthentic(c)
		if err != nil {
			c.Logger().Info("Sign in failed:", err)
			c.Flash().Add("danger", "Please sign in to continue.")
			return c.Redirect(302, "rootPath()")
		}

		c.Set("currentSpender", spender)
		return next(c)
	}
}

// IsAuthentic determines if a user is logged into the current context
func IsAuthentic(c buffalo.Context) (*models.Spender, error) {
	token, err := c.Cookies().Get(tokenCookieName)
	if err != nil {
		err = fmt.Errorf("failed to get token cookie: %w", err)
		return nil, err
	}

	tx := c.Value("tx").(*pop.Connection)
	spenderToken := models.SpenderToken{}
	err = tx.Where("expires_at > CURRENT_TIMESTAMP").
		Where("token = ?", token).
		EagerPreload("Spender").
		First(&spenderToken)

	return &spenderToken.Spender, err
}
