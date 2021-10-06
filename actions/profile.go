package actions

import (
	"fmt"
	"mespendyouspend/models"
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// showProfileHandler renders the current spender's profile
func showProfileHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("profile/show"))
}

// updateProfileHandler updates the current spender's profile
func updateProfileHandler(c buffalo.Context) error {
	currentSpender := c.Value("currentSpender").(*models.Spender)

	// Bind
	currentSpender.Name = c.Request().FormValue("Name")
	// Update
	err := models.DB.Update(currentSpender)
	if err != nil {
		err = fmt.Errorf("failed to update profile: %w", err)
		return err
	}

	return c.Redirect(http.StatusSeeOther, "profilePath()")
}
