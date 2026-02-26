package main

import (
	"context"
	"errors"
	"net/http"
	"social/internal/repository"
)

type postCtx string

const postCtxKey = postCtx("posts")

type createPostPayload struct {
	Title   *string   `json:"title" validate:"required,max=100"`
	Content *string   `json:"content"  validate:"required,max=1000"`
	Tags    *[]string `json:"tags"`
}

// @Summary		Create a post
// @Description	create a post
// @Tags			Posts
// @Accept			json
// @Produce		json
// @Param payload body createPostPayload true "create post payload"
// @Success		201	{object}	repository.Post
// @Failure		400	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Security ApiKeyAuth
// @SecurityDefinitions.apiKey		ApiKeyAuth
// @in header
// @name Authorization
// @Router			/posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload createPostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	ctx := r.Context()

	user := app.getAuthUserContext(r)

	if err := Validate.Struct(payload); err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	post := &repository.Post{
		Title:   *payload.Title,
		Content: *payload.Content,
		Tags:    *payload.Tags,
		UserID:  user.ID,
	}

	if err := app.repository.Posts.Create(ctx, post); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

// @Summary		Get Posts
// @Description	Get post by ID
// @Tags			Posts
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Post ID"
// @Success		200	{object}	repository.Post
// @Failure		404	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Security ApiKeyAuth
// @SecurityDefinitions.apiKey		ApiKeyAuth
// @in header
// @name Authorization
// @Router			/posts/{id} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	posts := getPostFromCtx(r)

	ctx := r.Context()
	comments, err := app.repository.Comments.GetCommentsByPostID(ctx, posts.ID)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	posts.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, posts); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

type patchPostPayload struct {
	Title   *string   `json:"title" validate:"omitempty,max=100"`
	Content *string   `json:"content" validate:"omitempty,max=1000"`
	Tags    *[]string `json:"tags"`
}

// @Summary		Update Post
// @Description	Update post
// @Tags			Posts
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Post ID"
// @Param payload body patchPostPayload true "create post payload"
// @Success		200	{object}	repository.Posts
// @Failure		400	{object}	ErrorResponse
// @Failure		404	{object}	ErrorResponse
// @Failure		409	{object}    ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Security ApiKeyAuth
// @SecurityDefinitions.apiKey		ApiKeyAuth
// @in header
// @name Authorization
// @Router			/posts/{id} [patch]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)
	var payload patchPostPayload

	err := readJSON(w, r, &payload)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if payload.Tags != nil {
		post.Tags = *payload.Tags
	}

	if err = Validate.Struct(payload); err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	err = app.repository.Posts.UpdatePost(r.Context(), post)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEditConflict):
			app.conflictResponse(w, r, err)
		default:
			app.internalServerErrorResponse(w, r, err)
		}
		return
	}

	if err = app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

// @Summary		Delete Post
// @Description	Delete post
// @Tags			Posts
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Post ID"
// @Success		204
// @Failure		404	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Security ApiKeyAuth
// @SecurityDefinitions.apiKey		ApiKeyAuth
// @in header
// @name Authorization
// @Router			/posts/{id} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postId, err := app.readIntParam(r, "id")
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	if err = app.repository.Posts.DeletePost(r.Context(), postId); err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			app.notFoundErrorResponse(w, r, err)
		default:
			app.internalServerErrorResponse(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postId, err := app.readIntParam(r, "id")
		if err != nil {
			app.badRequestErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()

		post, err := app.repository.Posts.GetByID(ctx, postId)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrRecordNotFound):
				app.notFoundErrorResponse(w, r, err)
			default:
				app.internalServerErrorResponse(w, r, err)
			}
			return
		}

		// never mutate context
		// always create a new one
		ctx = context.WithValue(ctx, postCtxKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *repository.Post {
	post, _ := r.Context().Value(postCtxKey).(*repository.Post)
	return post
}
