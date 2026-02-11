package main

import (
	"net/http"
	"social/internal/pagination"
)

// getUserFeedHandler goDoc
//
// @Summary		Get User feed
// @Description	 Get User feed
// @Tags			Feed
// @Produce		json
// @Param			page_size	query		int	false	"Page Size"
// @Param			current_page	query		int	false	"Current page"
// @Param			sort	query		string false	"Sort"
// @Param			tags	query		 string	false "tags"
// @Param			search	query		string false	"search"
// @Success		200	{object}	[]repository.PostWithMetadata
// @Failure		400	{object}	error
// @Failure		404	{object}	error
// @Failure		500	{object}	error
// @Security		ApiKeyAuth
// @Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	p, err := pagination.ParsePaginationParams(r)

	if err := Validate.Struct(p); err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	f, err := pagination.ParseFilterParams(r)
	f.SortSafeList = []string{"id", "created_at", "-id", "-created_at"}
	if err := Validate.Struct(f); err != nil {
		app.badRequestErrorResponse(w, r, err)
	}

	feed, err := app.repository.Posts.GetUserFeed(ctx, 18, &p, f)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	p.Paginate()

	p.WriteHeaders(w)

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}
