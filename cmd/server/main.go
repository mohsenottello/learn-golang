// Entry point. Wiring only: load config, open dependencies, build the router, serve.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mohsenottello/golang/internal/config"
	"github.com/mohsenottello/golang/internal/httpapi"
	"github.com/mohsenottello/golang/internal/mail"
	"github.com/mohsenottello/golang/internal/store"
	"github.com/mohsenottello/golang/web"
)

func main() {
	cfg := config.Load()

	st, err := store.Open(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal("store: ", err)
	}
	defer st.Close()

	if err := st.Migrate(); err != nil {
		log.Fatal("migrate: ", err)
	}

	srv := &http.Server{
		Addr: cfg.Addr,
		Handler: httpapi.NewRouter(httpapi.Deps{
			Store:   st,
			Mailer:  mail.New(cfg.SMTP),
			BaseURL: cfg.BaseURL,
			Static:  web.Dist,
		}),
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Graceful shutdown: finish in-flight requests on SIGINT/SIGTERM.
	go func() {
		log.Println("listening on", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Println("shutting down")
	srv.Shutdown(ctx)
}
