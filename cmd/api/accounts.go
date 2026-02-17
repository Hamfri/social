package main

import (
	"errors"
	"net/http"
	"social/internal/mailer"
	"social/internal/repository"
)

type registerUserPayload struct {
	Username string `json:"username" validate:"required,max=50"`
	Email    string `json:"email" validate:"required,email,max=50"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// @Summary User account registration
// @Description User account registration
// @Tags Accounts
// @Accept json
// @Produce json
// @Param payload body registerUserPayload true "User Credentials"
// @Success 201 {object} repository.User "User registered"
// @Failure 400 {object} error
// @Failure 400 {object} error
// @Router /accounts/register [post]
func (app *application) registerAccountHandler(w http.ResponseWriter, r *http.Request) {
	var payload registerUserPayload

	err := readJSON(w, r, &payload)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	ctx := r.Context()

	user := repository.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	err = user.Password.Set(payload.Password)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	token, err := app.repository.Users.CreateAndInvite(ctx, &user)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEmailTaken):
			app.badRequestErrorResponse(w, r, err)
		case errors.Is(err, repository.ErrUsernameTaken):
			app.badRequestErrorResponse(w, r, err)
		default:
			app.internalServerErrorResponse(w, r, err)
		}
		return
	}

	app.backgroundTaskRunner(func() {
		data := map[string]any{
			"username": user.Username,
			"token":    *token,
		}

		recipient := user.Username + "<" + user.Email + ">"
		err = app.mailer.Send(recipient, mailer.UserActivationTemplate, data)
		if err != nil {
			app.logger.Error(err.Error())
		}
	})

	if err = writeJSON(w, http.StatusCreated, user); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

type activateAccountPayload struct {
	Token string `json:"token" validate:"required,min=26,max=26"`
}

func (app *application) activateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var payload activateAccountPayload

	err := readJSON(w, r, &payload)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	ctx := r.Context()

	user, err := app.repository.Activate(ctx, repository.ScopeActivation, payload.Token)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			app.badRequestErrorResponse(w, r, err)
		case errors.Is(err, repository.ErrUsernameTaken):
			app.badRequestErrorResponse(w, r, err)
		case errors.Is(err, repository.ErrEmailTaken):
			app.badRequestErrorResponse(w, r, err)
		default:
			app.internalServerErrorResponse(w, r, err)

		}
		return
	}

	err = writeJSON(w, http.StatusOK, user)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

func (app *application) loginAccountHandler(w http.ResponseWriter, r *http.Request) {}

func (app *application) logoutAccountHandler(w http.ResponseWriter, r *http.Request) {}
