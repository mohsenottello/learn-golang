// Entry point. Wiring only: load config, open dependencies, build the router, serve.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/mohsenottello/golang/internal/config"
	"github.com/mohsenottello/golang/internal/httpapi"
	"github.com/mohsenottello/golang/internal/mail"
	"github.com/mohsenottello/golang/internal/store"
	"github.com/mohsenottello/golang/internal/translate"
	"github.com/mohsenottello/golang/web"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	st, err := store.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("store: ", err)
	}
	defer st.Close()

	if err := st.Migrate(); err != nil {
		log.Fatal("migrate: ", err)
	}

	if cfg.AdminEmail != "" && cfg.AdminPassword != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
		created, err := st.EnsureOwner(ctx, strings.ToLower(cfg.AdminEmail), envOr("ADMIN_NAME", "Owner"), string(hash))
		if err != nil {
			log.Fatal("owner seed: ", err)
		}
		if created {
			log.Println("seeded owner account", cfg.AdminEmail)
		}
	}

	srv := &http.Server{
		Addr: cfg.Addr,
		Handler: httpapi.NewRouter(httpapi.Deps{
			Store:      st,
			Mailer:     mail.New(cfg.SMTP),
			Translator: translate.New(cfg.DeepLKey),
			BaseURL:    cfg.BaseURL,
			Static:     web.Dist,
		}),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Println("listening on", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	log.Println("shutting down")
	srv.Shutdown(shutdownCtx)
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
