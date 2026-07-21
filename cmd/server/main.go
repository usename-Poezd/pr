package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"poll-app/internal/handler"
	"poll-app/internal/repository/postgres"
	"poll-app/internal/service"
	"strings"
	"time"
)

func main() {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://" + env("DB_USER", "postgres") + ":" + env("DB_PASSWORD", "postgres") + "@" + env("DB_HOST", "db") + ":" + env("DB_PORT", "5432") + "/" + env("DB_NAME", "polls") + "?sslmode=" + env("DB_SSLMODE", "disable")
	}
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		slog.Error("database", "error", err)
		return
	}
	defer db.Close()
	for i := 0; i < 30; i++ {
		if err = db.Ping(ctx); err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		slog.Error("database unavailable", "error", err)
		return
	}
	b, err := os.ReadFile(filepath.Join("migrations", "001_init.sql"))
	if err != nil {
		panic(err)
	}
	if _, err = db.Exec(ctx, string(b)); err != nil {
		panic(err)
	}
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(200) })
	r.Mount("/api", handler.New(service.NewPollService(postgres.NewRepository(db))).Routes())
	r.Handle("/*", static("frontend/dist"))
	addr := ":" + env("PORT", "8080")
	slog.Info("listening", "addr", addr)
	if err = http.ListenAndServe(addr, r); err != nil {
		slog.Error("server", "error", err)
	}
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func static(root string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api" || strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		files := http.Dir(root)
		p := strings.TrimPrefix(filepath.Clean(r.URL.Path), "/")
		if p == "." {
			p = ""
		}
		if f, e := files.Open(p); e == nil {
			f.Close()
			http.FileServer(files).ServeHTTP(w, r)
			return
		} else if !errors.Is(e, fs.ErrNotExist) {
			http.NotFound(w, r)
			return
		}
		r.URL.Path = "/"
		http.FileServer(files).ServeHTTP(w, r)
	})
}
