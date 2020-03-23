package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/shohrukh56/mux/pkg/mux"
	"github.com/shohrukh56/sales/cmd/sales/app"
	"github.com/shohrukh56/sales/pkg/core/sales"
	"github.com/shohrukh56/DI/pkg/di"
	"github.com/shohrukh56/jwt/pkg/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
	"net"
	"net/http"
	"os"
)

var (
	host = flag.String("host", "0.0.0.0", "Server host")
	port = flag.String("port", "9999", "Server port")
	dsn  = flag.String("dsn", "postgres://user:pass@localhost:5431/sales", "Postgres DSN")
)
//-host 0.0.0.0 -port 9999 -dsn postgres://user:pass@localhost:5431/sales
const (
	envHost = "HOST"
	envPort = "PORT"
	envDSN  = "DATABASE_URL"
)

type DSN string

func main() {
	flag.Parse()
	serverHost := checkENV(envHost, *host)
	serverPort := checkENV(envPort, *port)
	serverDsn := checkENV(envDSN, *dsn)
	addr := net.JoinHostPort(serverHost, serverPort)
	secret := jwt.Secret("secret")
	start(addr, serverDsn, secret)
}
func checkENV(env string, loc string) string {
	str, ok := os.LookupEnv(env)
	if !ok {
		return loc
	}
	return str
}
func start(addr string, dsn string,  secret jwt.Secret) {
	container := di.NewContainer()
	container.Provide(
		app.NewServer,
		mux.NewExactMux,
		sales.NewService,
		func() DSN { return DSN(dsn) },
		func() jwt.Secret { return secret },
		func(dsn DSN) *pgxpool.Pool {
			pool, err := pgxpool.Connect(context.Background(), string(dsn))
			if err != nil {
				panic(fmt.Errorf("can't create pool: %w", err))
			}
			return pool
		},
	)

	container.Start()

	var appServer *app.Server
	container.Component(&appServer)
	panic(http.ListenAndServe(addr, appServer))
}
