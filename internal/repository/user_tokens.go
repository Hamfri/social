package repository

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"time"
)

const (
	ScopeActivation string = "activation"
)

type UserTokens interface {
	Create(context.Context, *sql.Tx, int64, time.Duration, string) (*UserToken, error)
	Delete(context.Context, *sql.Tx, int64, string) error
}

type UserToken struct {
	PlaintText string
	UserId     int64
	Token      []byte
	Expiry     time.Time
	Scope      string
}

type UserTokenRepository struct {
	DB *sql.DB
}

func generateToken(userId int64, ttl time.Duration, scope string) *UserToken {
	userToken := &UserToken{
		PlaintText: rand.Text(),
		UserId:     userId,
		Expiry:     time.Now().Add(ttl),
		Scope:      scope,
	}

	token := sha256.Sum256([]byte(userToken.PlaintText))
	userToken.Token = token[:]
	return userToken
}

func (r *UserTokenRepository) Create(ctx context.Context, tx *sql.Tx, userId int64, ttl time.Duration, scope string) (*UserToken, error) {
	token := generateToken(userId, ttl, scope)

	query := `
		INSERT INTO user_tokens (token, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)
	`

	args := []any{token.Token, token.UserId, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (r *UserTokenRepository) Delete(ctx context.Context, tx *sql.Tx, userId int64, scope string) error {
	query := `
		DELETE from user_tokens WHERE user_id = $1 AND scope=$2
	`
	args := []any{userId, scope}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()

	return err
}
