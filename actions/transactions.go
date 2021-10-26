package actions

import (
	"fmt"
	"mespendyouspend/models"
	"net/http"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/x/responder"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Transaction)
// DB Table: Plural (transactions)
// Resource: Plural (Transactions)
// Path: Plural (/transactions)
// View Template Folder: Plural (/templates/transactions/)

// TransactionsResource is the resource for the Transaction model
type TransactionsResource struct {
	buffalo.Resource
}

const goalAmount = 350

// List gets all Transactions. This function is mapped to the path
// GET /transactions
func (v TransactionsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	transactions := models.Transactions{}

	q := tx.Where("spent_on >= ?", weekStart(time.Now()))

	// Retrieve all Transactions from the DB
	if err := q.EagerPreload("Spender").All(&transactions); err != nil {
		return err
	}

	totalAmount := 0
	for _, transaction := range transactions {
		totalAmount += transaction.Amount
	}
	c.Set("totalAmount", totalAmount)
	c.Set("amountRemaining", goalAmount-totalAmount)

	return responder.Wants("html", func(c buffalo.Context) error {
		c.Set("transactions", transactions)
		return c.Render(http.StatusOK, r.HTML("/transactions/index.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(transactions))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(transactions))
	}).Respond(c)
}

// Show gets the data for one Transaction. This function is mapped to
// the path GET /transactions/{transaction_id}
func (v TransactionsResource) Show(c buffalo.Context) error {
	return responder.Wants("html", func(c buffalo.Context) error {
		return c.Redirect(http.StatusSeeOther, "/transactions")
	}).Wants("json", func(c buffalo.Context) error {
		// Get the DB connection from the context
		tx, ok := c.Value("tx").(*pop.Connection)
		if !ok {
			return fmt.Errorf("no transaction found")
		}

		// Allocate an empty Transaction
		transaction := &models.Transaction{}

		// To find the Transaction the parameter transaction_id is used.
		if err := tx.EagerPreload("Spender").Find(transaction, c.Param("transaction_id")); err != nil {
			return c.Error(http.StatusNotFound, err)
		}
		return c.Render(200, r.JSON(transaction))
	}).Respond(c)
}

// New renders the form for creating a new Transaction.
// This function is mapped to the path GET /transactions/new
func (v TransactionsResource) New(c buffalo.Context) error {
	c.Set("transaction", &models.Transaction{SpentOn: time.Now()})

	return c.Render(http.StatusOK, r.HTML("/transactions/new.plush.html"))
}

// Create adds a Transaction to the DB. This function is mapped to the
// path POST /transactions
func (v TransactionsResource) Create(c buffalo.Context) error {
	transaction := &models.Transaction{}

	// Bind transaction to the html form elements
	if err := c.Bind(transaction); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	currentSpender := c.Value("currentSpender").(*models.Spender)
	transaction.SpenderID = currentSpender.ID

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(transaction)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the new.html template that the user can
			// correct the input.
			c.Set("transaction", transaction)

			return c.Render(http.StatusUnprocessableEntity, r.HTML("/transactions/new.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "transaction.created.success"))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/transactions")
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.JSON(transaction))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.XML(transaction))
	}).Respond(c)
}

// Edit renders a edit form for a Transaction. This function is
// mapped to the path GET /transactions/{transaction_id}/edit
func (v TransactionsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Transaction
	transaction := &models.Transaction{}

	if err := tx.Find(transaction, c.Param("transaction_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("transaction", transaction)
	return c.Render(http.StatusOK, r.HTML("/transactions/edit.plush.html"))
}

// Update changes a Transaction in the DB. This function is mapped to
// the path PUT /transactions/{transaction_id}
func (v TransactionsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Transaction
	transaction := &models.Transaction{}

	if err := tx.Find(transaction, c.Param("transaction_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind Transaction to the html form elements
	if err := c.Bind(transaction); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(transaction)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the edit.html template that the user can
			// correct the input.
			c.Set("transaction", transaction)

			return c.Render(http.StatusUnprocessableEntity, r.HTML("/transactions/edit.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "transaction.updated.success"))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/transactions/%v", transaction.ID)
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(transaction))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(transaction))
	}).Respond(c)
}

// Destroy deletes a Transaction from the DB. This function is mapped
// to the path DELETE /transactions/{transaction_id}
func (v TransactionsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Transaction
	transaction := &models.Transaction{}

	// To find the Transaction the parameter transaction_id is used.
	if err := tx.Find(transaction, c.Param("transaction_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(transaction); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a flash message
		c.Flash().Add("success", T.Translate(c, "transaction.destroyed.success"))

		// Redirect to the index page
		return c.Redirect(http.StatusSeeOther, "/transactions")
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(transaction))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(transaction))
	}).Respond(c)
}

func weekStart(in time.Time) time.Time {
	for {
		if in.Weekday() == time.Sunday {
			return in
		}
		in = in.AddDate(0, 0, -1)
	}
}
