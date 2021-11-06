package main

import (
	"context"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"net/http"
	"nononsensecode.com/session-management/adapters/mem"
	redis2 "nononsensecode.com/session-management/adapters/redis"
	"nononsensecode.com/session-management/ports"
	"nononsensecode.com/session-management/server"
	"os"
	"strconv"
	"time"
)

var (
	client *redis.Client
	sessionManager *scs.SessionManager
)

func initRedisClient() {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	dbAsString := os.Getenv("REDIS_DB")
	db, err := strconv.Atoi(dbAsString)
	if err != nil {
		panic("please enter a valid db number")
	}
	client = redis.NewClient(&redis.Options{
		Addr: host + ":" + port,
		Password: password,
		DB: db,
	})
}

func initSessionManager(ctx context.Context, prefix string) {
	store := redis2.NewRedisStore(ctx, client, prefix)
	sessionManager = scs.New()
	sessionManager.Store = store
	sessionManager.IdleTimeout = 15 * time.Minute
	sessionManager.Lifetime = 1 * time.Hour
	sessionManager.Cookie.Name = "self_reading"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Persist = true
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	initRedisClient()
	initSessionManager(context.Background(), "user")
	server.RunHTTPServerOnAddr(":8080", func(router chi.Router) http.Handler {
		return ports.HandlerFromMux(
			ports.NewHttpServer(mem.UserRepo{}, sessionManager),
			router,
		)
	}, sessionManager)
}
