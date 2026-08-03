package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/tojdor/nasrinShop/internal/auth"
	"github.com/tojdor/nasrinShop/internal/categories"
	"github.com/tojdor/nasrinShop/internal/materials"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("файл .env не найден, использую переменные окружения")
	}

	databaseURL := os.Getenv("DATABASE_URL")

	if err := runMigrations(databaseURL); err != nil {
		log.Fatalf("ошибка применения миграций: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	mux := http.NewServeMux()

	categoriesHandler := categories.NewHandler(categories.NewService(categories.NewStorage(pool)))
	materialsHandler := materials.NewHandler(materials.NewService(materials.NewStorage(pool)))

	categoriesHandler.RegisterRoutes(mux, auth.RequireAdmin)
	materialsHandler.RegisterRoutes(mux, auth.RequireAdmin)

	mux.HandleFunc("GET /admin/ping", auth.RequireAdmin(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/", fs)

	log.Println("сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func runMigrations(databaseURL string) error {
	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	log.Println("миграции применены")
	return nil
}