package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailTaken    = errors.New("email is already taken")
	ErrUsernameTaken = errors.New("username already taken")
)

type Users interface {
	GetByID(context.Context, int64) (*User, error)
	Create(context.Context, *sql.Tx, *User) error
	Update(context.Context, *sql.Tx, *User) error
	CreateAndInvite(context.Context, *User) (*string, error)
	GetUserByToken(context.Context, *sql.Tx, string, string) (*User, error)
	Activate(context.Context, string, string) (*User, error)
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type UserRepository struct {
	DB *sql.DB
	// Avoid doing this in large codebases
	// use services instead to avoid creating spaghetti code
	// Dependency injection
	UserTokens
}

func (r *UserRepository) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users
		(username, email, password)
		VALUES($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	args := []any{user.Username, user.Email, user.Password.hash}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	var pqErr *pq.Error

	if err != nil {
		switch {
		case errors.As(err, &pqErr) && pqErr.Constraint == "users_email_key":
			return ErrEmailTaken
		case errors.As(err, &pqErr) && pqErr.Constraint == "users_username_key":
			return ErrUsernameTaken
		default:
			return err
		}
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		UPDATE users SET username = $1, email = $2, password = $3, activated = $4 WHERE id = $5
	`

	args := []any{user.Username, user.Email, user.Password.hash, user.Activated, user.ID}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()

	var pqErr *pq.Error

	if err != nil {
		switch {
		case errors.As(err, &pqErr) && pqErr.Constraint == "users_email_key":
			return ErrEmailTaken
		case errors.As(err, &pqErr) && pqErr.Constraint == "users_username_key":
			return ErrUsernameTaken
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, ID int64) (*User, error) {
	query := `
		SELECT id, username, email, password FROM users
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)

	defer cancel()

	var post User

	err := r.DB.QueryRowContext(ctx, query, ID).Scan(
		&post.ID,
		&post.Username,
		&post.Email,
		&post.Password,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (r *UserRepository) GetUserByToken(ctx context.Context, tx *sql.Tx, scope, plainTextToken string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(plainTextToken))

	query := `
		SELECT u.id, u.username, u.email, u.password, u.activated
		FROM users u
		INNER JOIN user_tokens t ON t.user_id = u.id
		WHERE t.scope = $1 
		AND t.token = $2
		AND t.expiry > $3
	`

	args := []any{scope, tokenHash[:], time.Now()}
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User

	err := tx.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (r *UserRepository) CreateAndInvite(ctx context.Context, user *User) (*string, error) {
	return WithTxAndResult(ctx, r.DB, func(tx *sql.Tx) (*string, error) {
		if err := r.Create(ctx, tx, user); err != nil {
			return nil, err
		}

		token, err := r.UserTokens.Create(ctx, tx, user.ID, 24*time.Hour, ScopeActivation)
		if err != nil {
			return nil, err
		}

		return &token.PlaintText, nil
	})
}

func (r *UserRepository) Activate(ctx context.Context, scope, plainTextToken string) (*User, error) {
	return WithTxAndResult(ctx, r.DB, func(tx *sql.Tx) (*User, error) {
		user, err := r.GetUserByToken(ctx, tx, ScopeActivation, plainTextToken)
		if err != nil {
			return nil, err
		}

		user.Activated = true
		err = r.Update(ctx, tx, user)
		if err != nil {
			return nil, err
		}

		err = r.UserTokens.Delete(ctx, tx, user.ID, scope)
		if err != nil {
			return nil, err
		}

		return user, nil
	})
}

type password struct {
	plaintText *string
	hash       []byte
}

func (p *password) Set(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.plaintText = &password
	p.hash = hash

	return nil
}

func (p *password) Matches(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, err // errors.New("password or email")
		default:
			return false, err
		}
	}

	return true, nil
}
