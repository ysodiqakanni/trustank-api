package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/ysodiqakanni/trustank-api/internal/auth"
	"github.com/ysodiqakanni/trustank-api/internal/business"
	"github.com/ysodiqakanni/trustank-api/internal/businessCategory"
	"github.com/ysodiqakanni/trustank-api/internal/config"
	"github.com/ysodiqakanni/trustank-api/internal/user"
	"github.com/ysodiqakanni/trustank-api/pkg/dbcontext"
	"github.com/ysodiqakanni/trustank-api/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Version indicates the current version of the application.
var Version = "1.0.0"

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

func main() {
	flag.Parse()
	// create root logger tagged with server version
	logger := log.New().With(nil, "version", Version)

	// load application configurations
	cfg, err := config.Load(*flagConfig, logger)
	if err != nil {
		logger.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}
	// connect to the mongo database
	escapedPassword := url.QueryEscape(cfg.DbPassword)
	connStr := fmt.Sprintf(cfg.DbConnectionString, escapedPassword)
	db, err := NewMongoDB(connStr, cfg.DbName)

	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}

	defer db.Client().Disconnect(context.Background())

	// build HTTP server
	address := fmt.Sprintf(":%v", cfg.ServerPort)
	hs := &http.Server{
		Addr:    address,
		Handler: buildHandler(logger, dbcontext.New(db), cfg),
	}

	// start the HTTP server with graceful shutdown
	// go routing.GracefulShutdown(hs, 10*time.Second, logger.Infof)
	logger.Infof("server %v is running at %v", Version, address)
	if err := hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(err)
		os.Exit(-1)
	}
}

func NewMongoDB(connStr, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Ping successful!")
	// mongoClient = client

	return client.Database(dbName), nil
}

func buildHandler(logger log.Logger, db *dbcontext.DB, cfg *config.Config) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/api/healthcheck", HealthCheckHandler).Methods("GET")

	business.RegisterBusinessHandlers(r,
		business.NewService(business.NewRepository(db, logger), user.NewRepository(db, logger), logger),
		logger,
		cfg.JWTSigningKey)

	business.RegisterHandlers(r,
		business.NewService(business.NewRepository(db, logger), user.NewRepository(db, logger), logger),
		logger,
		cfg.JWTSigningKey)

	businessCategory.RegisterHandlers(r,
		businessCategory.NewService(businessCategory.NewRepository(db, logger), logger),
		logger,
		cfg.JWTSigningKey)

	auth.RegisterHandlers(r,
		auth.NewService(cfg.JWTSigningKey, cfg.JWTExpiration, logger, user.NewRepository(db, logger)),
		logger)

	return r
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Everything is dope from this side :)")
}
