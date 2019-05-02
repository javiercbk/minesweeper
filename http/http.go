package http

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/javiercbk/minesweeper/auth"
	"github.com/javiercbk/minesweeper/game"
	"github.com/javiercbk/minesweeper/http/response"
	"github.com/javiercbk/minesweeper/http/security"
	"github.com/javiercbk/minesweeper/player"

	"gopkg.in/go-playground/validator.v9"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	gommonLog "github.com/labstack/gommon/log"

	// imports the postgres sql driver
	_ "github.com/lib/pq"
)

// Config contains all the configurations to initialize an http server
type Config struct {
	Address   string
	JWTSecret string
}

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func httpErrorHandlerFactory(logger *log.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			if code == 404 {
				response.NewErrorResponse(c, code, fmt.Sprintf("resource %s was not found", c.Request().URL))
				return
			}
		}
		logger.Printf("Error in server %v", err)
		response.NewErrorResponse(c, code, http.StatusText(code))
	}
}

// Serve http connections
func Serve(cnf Config, logger *log.Logger, db *sql.DB) error {
	router := echo.New()
	router.HTTPErrorHandler = httpErrorHandlerFactory(logger)
	router.Validator = &customValidator{validator: validator.New()}
	router.Logger.SetLevel(gommonLog.INFO)
	router.Use(middleware.Recover())
	router.Use(middleware.Secure())
	router.Use(middleware.BodyLimit("1M"))
	router.Use(middleware.Gzip())
	initRoutes(router, cnf.JWTSecret, logger, db)
	srv := newServer(router, cnf.Address)
	go func() {
		// serve connections
		logger.Printf("Listening http connections on %s", cnf.Address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			router.Logger.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can't be catched, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	router.Logger.Printf("Shutdown Server ...\n")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		router.Logger.Fatal("Server Shutdown:", err)
	}
	<-ctx.Done()
	router.Logger.Printf("timeout of 5 seconds.\n")
	router.Logger.Printf("Server exiting\n")
	return nil
}

func initRoutes(router *echo.Echo, jwtSecret string, logger *log.Logger, db *sql.DB) {
	jwtMiddleware := security.JWTMiddlewareFactory(jwtSecret)
	apiRouter := router.Group("/api")
	{
		authRouter := apiRouter.Group("/auth")
		authHandler := auth.NewHandler(logger, db)
		authHandler.Routes(authRouter, jwtSecret)
	}
	{
		gamesRouter := apiRouter.Group("/games")
		gamesRouter.Use(jwtMiddleware)
		gameHandler := game.NewHandler(logger, db)
		gameHandler.Routes(gamesRouter)
	}
	{
		playerRouter := apiRouter.Group("/players")
		playerHandler := player.NewHandler(logger, db)
		playerHandler.Routes(playerRouter, jwtMiddleware)
	}
}

func newServer(handler http.Handler, address string) *http.Server {
	// see https://blog.cloudflare.com/exposing-go-on-the-internet/
	tlsConfig := &tls.Config{
		// Causes servers to use Go's default ciphersuite preferences,
		// which are tuned to avoid attacks. Does nothing on clients.
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519, // Go 1.8 only
		},

		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			// Best disabled, as they don't provide Forward Secrecy,
			// but might be necessary for some clients
			// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	return &http.Server{
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig:    tlsConfig,
		Handler:      handler,
	}
}
