#!/bin/sh
export VERSION=1.5.3
docker run -p 8080:8080  --rm --name=social -it --network social_social-net \
    -e DB_DSN="postgresql://postgres:jZWctXPbDEunT98bbTA6lD8JKRxxP@db.jvqportyhdfLkfqqqppf.supabase.co:5432/social?sslmode=require" \
    -e APP_PORT=":8080" \
    -e APP_URL="localhost:8080" \
    -e APP_ENV="development" \
    -e DB_MAX_OPEN_CONNS=30 \
    -e DB_MAX_IDLE_CONNS=30 \
    -e DB_MAX_IDLE_TIME="15m" \
    -e SMTP_HOST="sandbox.smtp.mailtrap.io" \
    -e SMTP_PORT=25 \
    -e SMTP_USERNAME="" \
    -e SMTP_PASSWORD="" \
    -e SMTP_SENDER="Hatari <hatari@hadaa.com>" \
    -e BASIC_AUTH_USERNAME="admin" \
    -e BASIC_AUTH_PASSWORD="admin" \
    -e JWT_AUD="localhost" \
    -e JWT_SECRET="" \
    -e JWT_ISS="localhost" \
    -e JWT_EXP=72 \
    -e REDIS_ADDR="localhost:6379" \
    -e REDIS_PASSWORD="" \
    -e REDIS_DB=0 \
    -e REDIS_ENABLED=false \
    -e RATE_LIMITER_REQUESTS_COUNT=20 \
    -e RATE_LIMITER_ENABLED=true \
    -e CORS_ALLOWED_ORIGIN="localhost:4200" \
    api:${VERSION}