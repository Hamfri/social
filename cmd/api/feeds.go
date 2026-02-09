package main

import (
	"net/http"
	"social/internal/pagination"
)

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
