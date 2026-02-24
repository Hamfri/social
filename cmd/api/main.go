package main

import (
	"expvar"
	"os"
	"runtime"
	"social/internal/auth"
	"social/internal/cache"
	"social/internal/db"
	"social/internal/env"
	"social/internal/mailer"
	"social/internal/ratelimiter"
	"social/internal/repository"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const version = "1.4.1"

//	@title			Social
//	@description	simple social network implementation
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	cfg := config{
		port:   env.GetString("APP_PORT", ":8080"),
		apiURL: env.GetString("APP_URL", "localhost:8080"),
		env:    env.GetString("APP_ENV", "development"),
		db: dbConfig{
			dsn:          env.GetString("DB_DSN", ""),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		smtp: smtp{
			host:     env.GetString("SMTP_HOST", ""),
			port:     env.GetInt("SMTP_PORT", 25),
			username: env.GetString("SMTP_USERNAME", ""),
			password: env.GetString("SMTP_PASSWORD", ""),
			sender:   env.GetString("SMTP_SENDER", ""),
		},
		auth: authConfig{
			basic: basicConfig{
				username: env.GetString("BASIC_AUTH_USERNAME", "admin"),
				password: env.GetString("BASIC_AUTH_PASSWORD", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("JWT_SECRET", ""),
				aud:    env.GetString("JWT_AUD", ""),
				iss:    env.GetString("JWT_ISS", ""),
				exp:    time.Duration(env.GetInt("JWT_EXP", 24)),
			},
		},
		redis: redisConfig{
			addr:    env.GetString("REDIS_ADDR", ""),
			pw:      env.GetString("REDIS_PASSWORD", ""),
			db:      env.GetInt("REDIS_DB", 1),
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
		ratelimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.GetInt("RATE_LIMITER_REQUESTS_COUNT", 20),
			TimeFrame:            time.Second * 5,
			Enabled:              env.GetBool("RATE_LIMITER_ENABLED", false),
		},
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(cfg.db.dsn, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("database connection pool established")

	// redis
	var rdb *redis.Client
	if cfg.redis.enabled {
		rdb = cache.NewRedisClient(cfg.redis.addr, cfg.redis.pw, cfg.redis.db)
		logger.Info("redis connection established")
	}

	mailer, err := mailer.NewMailtrap(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	JWTauthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.aud, cfg.auth.token.iss)

	ratelimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.ratelimiter.RequestsPerTimeFrame,
		cfg.ratelimiter.TimeFrame,
	)

	app := &application{
		config:        cfg,
		repository:    repository.New(db),
		redisCache:    cache.NewRedisStorage(rdb),
		logger:        logger,
		mailer:        mailer,
		wg:            &sync.WaitGroup{},
		authenticator: JWTauthenticator,
		rateLimiter:   ratelimiter,
	}

	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()

	// trigger release
	if err = app.run(mux); err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	}
}
