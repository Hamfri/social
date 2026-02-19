package main

import (
	"errors"
	"net/http"
	"social/internal/mailer"
	"social/internal/repository"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type loginCredentials struct {
	Email    string `json:"email" validate:"required,email,max=50"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}
type registerUserPayload struct {
	Username string `json:"username" validate:"required,max=50"`
	loginCredentials
}

// @Summary register
// @Description User account registration
// @Tags Accounts
// @Accept json
// @Produce json
// @Param payload body registerUserPayload true "registration Credentials"
// @Success 201 {object} repository.User
// @Failure 400 {object} ErrorResponse
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
		Role: repository.Role{
			Name: "user",
		},
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

	if err = app.jsonResponse(w, http.StatusCreated, user); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

type activateAccountPayload struct {
	Token string `json:"token" validate:"required,min=26,max=26" example:"7ND32BQSB2IYDR5CP2O42XME57"`
}

// @Summary activate account
// @Description User account activation
// @Tags Accounts
// @Accept json
// @Produce json
// @Param payload body activateAccountPayload true "Activate account payload "
// @Success 200 {object} repository.User
// @Failure 400 {object} ErrorResponse
// @Router /accounts/activate [put]
func (app *application) activateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var payload activateAccountPayload

	err := readJSON(w, r, &payload)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	if err = Validate.Struct(payload); err != nil {
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

	err = app.jsonResponse(w, http.StatusOK, user)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

// @Summary login
// @Description login
// @Tags Accounts
// @Accept json
// @Produce json
// @Param payload body loginCredentials true "Login credentials"
// @Success 200 {object} repository.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /accounts/login [post]
func (app *application) loginAccountHandler(w http.ResponseWriter, r *http.Request) {
	var payload loginCredentials

	err := readJSON(w, r, &payload)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	if err = Validate.Struct(payload); err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	ctx := r.Context()
	user, err := app.repository.Users.GetByEmail(ctx, payload.Email)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			app.unauthorizedResponse(w, r, err)
		default:
			app.internalServerErrorResponse(w, r, err)
		}
		return
	}

	if !user.Activated {
		app.badRequestErrorResponse(w, r, errors.New("please activate your account"))
		return
	}

	_, err = user.Password.Matches(payload.Password)
	if err != nil {
		app.unauthorizedResponse(w, r, err)
		return
	}

	claims := jwt.MapClaims{
		"sub": strconv.Itoa(int(user.ID)),
		"exp": time.Now().Add(app.config.auth.token.exp * time.Hour).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.aud,
	}

	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err = app.jsonResponse(w, http.StatusOK, token); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}

}

func (app *application) logoutAccountHandler(w http.ResponseWriter, r *http.Request) {}
