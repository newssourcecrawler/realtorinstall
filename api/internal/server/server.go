package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/newssourcecrawler/realtorinstall/api/internal/server/middleware"
	"github.com/newssourcecrawler/realtorinstall/api/internal/server/routes"
)

func Start() {
	router := gin.Default()
	router.Use(cors.Default())

	// Build DBs + repos + services + handlers in one place (factor as needed)
	dbs := openAllDBs()         // create data/ dir, open each .db, ApplyMigrations
	repos := initRepos(dbs)     // userRepo, propRepo, ...
	svcs := initServices(repos) // authSvc, propSvc, ...
	hs := initHandlers(svcs)    // authH, propH, ...

	// Public route
	router.POST("/login", hs.Auth.Login)

	// JWT + permission middleware
	router.Use(middleware.Auth(svcs.Auth, repos.User))
	routes.Register(router, hs)

	// Listen + TLS + graceful shutdown
	srv := &http.Server{
		Addr:      ":8443",
		Handler:   router,
		TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
	}
	go func() {
		fmt.Println("Listening on https://localhost:8443")
		srv.ListenAndServeTLS("certs/server.crt", "certs/server.key")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
