package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/qiangxue/go-rest-api/internal/auth"
	"github.com/qiangxue/go-rest-api/internal/business"
	"github.com/qiangxue/go-rest-api/internal/config"
	"github.com/qiangxue/go-rest-api/internal/user"
	"github.com/qiangxue/go-rest-api/pkg/dbcontext"
	"github.com/qiangxue/go-rest-api/pkg/log"
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
	r.HandleFunc("/public1", publicEndpoint1).Methods("GET")

	// Protected Endpoints
	r.Handle("/protected1", auth.AuthenticateMiddleware(http.HandlerFunc(protectedEndpoint1), cfg.JWTSigningKey)).Methods("GET")
	r.Handle("/protected2", auth.AuthenticateMiddleware(http.HandlerFunc(protectedEndpoint2), cfg.JWTSigningKey)).Methods("GET")

	r.HandleFunc("/public2", publicEndpoint2).Methods("GET")

	business.RegisterBusinessHandlers(r,
		business.NewService(business.NewRepository(db, logger), user.NewRepository(db, logger), logger),
		logger,
		cfg.JWTSigningKey)

	business.RegisterHandlers(r,
		business.NewService(business.NewRepository(db, logger), user.NewRepository(db, logger), logger),
		logger,
		cfg.JWTSigningKey)

	//businessCategory.RegisterHandlers(r,
	//	businessCategory.NewService(businessCategory.NewRepository(db, logger), logger),
	//	logger)
	//
	auth.RegisterHandlers(r,
		auth.NewService(cfg.JWTSigningKey, cfg.JWTExpiration, logger, user.NewRepository(db, logger)),
		logger)

	// jwt stuffs
	//jwtMiddleware := auth.JwtMiddleware(r, cfg.JWTSigningKey)
	//authMiddleware := auth.AuthenticateMiddleware(r, cfg.JWTSigningKey)
	//authMW := AuthMiddleware(r)
	//r.PathPrefix("/api/v1/businesses").Subrouter().Use(jwtMiddleware)

	//r.Handle("/api/v1/blocked", AuthMiddleware(http.HandlerFunc(blockedHandler))).Methods("GET")

	//business.RegisterHandlers(r, jwtMiddleware,
	//	business.NewService(business.NewRepository(db, logger), user.NewRepository(db, logger), logger),
	//	logger)
	return r
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Auth middleware hit!")
		next.ServeHTTP(w, r)
	})
}

func blockedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("you're authorized"))
}
func regularHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a regular handler"))
}

func AuthMiddleware2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for a Bearer token in the Authorization header
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// In a real application, you would validate the token here.
		// For simplicity, we'll assume any token is valid.

		next.ServeHTTP(w, r)
	})
}

func publicEndpoint1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is a public endpoint 1")
}

func publicEndpoint2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is a public endpoint 2")
}

func protectedEndpoint1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is a protected endpoint 1")
}

func protectedEndpoint2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is a protected endpoint 2")
}

/*
// buildHandler sets up the HTTP routing and builds an HTTP handler.
func buildHandler(logger log.Logger, db *dbcontext.DB, cfg *config.Config) http.Handler {
	router := routing.New()

	router.Use(
		accesslog.Handler(logger),
		errors.Handler(logger),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.AllowAll),
	)

	healthcheck.RegisterHandlers(router, Version)

	rg := router.Group("/v1")

	authHandler := auth.Handler(cfg.JWTSigningKey)

	album.RegisterHandlers(rg.Group(""),
		album.NewService(album.NewRepository(db, logger), logger),
		authHandler, logger,
	)
	businessCategory.RegisterHandlers(rg.Group(""),
		businessCategory.NewService(businessCategory.NewRepository(db, logger), logger),
		authHandler, logger,
	)

	auth.RegisterHandlers(rg.Group(""),
		auth.NewService(cfg.JWTSigningKey, cfg.JWTExpiration, logger),
		logger,
	)

	return router
}
*/

/*
// logDBQuery returns a logging function that can be used to log SQL queries.
func logDBQuery(logger log.Logger) dbx.QueryLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB query successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB query error: %v", err)
		}
	}
}

// logDBExec returns a logging function that can be used to log SQL executions.
func logDBExec(logger log.Logger) dbx.ExecLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB execution successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB execution error: %v", err)
		}
	}
}
*/
