package main

import (
	"context"
	"errors"
	"net/http"
	"social/internal/repository"
)

type userCtx string

const userCtxKey = userCtx("user")

// @Summary		Get user
// @Description	get user by ID
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"User ID"
// @Success		200	{object}	repository.User
// @Failure		400	{object}	error
// @Failure		404	{object}	error
// @Failure		500	{object}	error
// @Security		ApiKeyAuth
// @Router			/users/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

type AuthenticatedUser struct {
	UserID int64 `json:"user_id"`
}

// @Summary		follow user
// @Description	Follow a user
// @Tags			Users
// @Produce		json
// @Param			id	path	int	true	"User ID"
// @Success		204
// @Failure		400	string		error	"you are already following that user"
// @Failure		409	string	error	"self follow not allowed"
// @Failure		404	string	error "record not found"
// @Failure		500	{object}	error
// @Security		ApiKeyAuth
// @Router			/users/{id}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser := getUserFromCtx(r)

	// Todo // use authenticated user id
	var authenticatedUser AuthenticatedUser

	if err := readJSON(w, r, &authenticatedUser); err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	userFollow := repository.UserFollow{
		FollowedID: followedUser.ID,
		FollowerID: authenticatedUser.UserID,
	}

	err := app.repository.UserFollows.Follow(r.Context(), &userFollow)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrAlreadyFollowing):
			app.badRequestErrorResponse(w, r, err)
		case errors.Is(err, repository.ErrNoSelfFollow):
			app.conflictResponse(w, r, err)
		default:
			app.internalServerErrorResponse(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}

}

// @Summary		Unfollow user
// @Description	Unfollow a user
// @Tags			Users
// @Produce		json
// @Param			id	path	int	true	"User ID"
// @Success		204
// @Failure		404	{object}	error
// @Failure		500	{object}	error
// @Security		ApiKeyAuth
// @Router			/users/{id}/unfollow [delete]
func (app *application) unFollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser := getUserFromCtx(r)

	// Todo // use authenticated user id
	var authenticatedUser AuthenticatedUser

	if err := readJSON(w, r, &authenticatedUser); err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	userFollow := repository.UserFollow{
		FollowedID: followedUser.ID,
		FollowerID: authenticatedUser.UserID,
	}

	ctx := r.Context()

	err := app.repository.UserFollows.Unfollow(ctx, userFollow)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			app.notFoundErrorResponse(w, r, err)
		default:
			app.internalServerErrorResponse(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

func (app *application) usersContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := app.readIntParam(r, "id")
		if err != nil {
			app.badRequestErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.repository.Users.GetByID(ctx, userId)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrRecordNotFound):
				app.notFoundErrorResponse(w, r, err)
			default:
				app.internalServerErrorResponse(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, userCtxKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *repository.User {
	user, _ := r.Context().Value(userCtxKey).(*repository.User)
	return user
}
