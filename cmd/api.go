package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	repo "github.com/pick-cee/go-ecommerce-api/internal/adapters/postgresql/sqlc"
	"github.com/pick-cee/go-ecommerce-api/internal/orders"
	"github.com/pick-cee/go-ecommerce-api/internal/products"
	"github.com/pick-cee/go-ecommerce-api/internal/users"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	repository := repo.New(app.db)

	// Middlewares
	r.Use(middleware.RequestID) // important for Rate limiting
	r.Use(middleware.RealIP) // important for rate limiting, tracking and analytics
	r.Use(middleware.Logger) 
	r.Use(middleware.Recoverer) // recover from crashes

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("All good!"))
	})

	productsService := products.NewService(repository, app.db)
	productsHandler := products.NewHandler(productsService)

	r.Get("/products", productsHandler.ListProducts)
	r.With(users.AuthMiddleware).Post("/products", productsHandler.CreateProduct)
	r.With(users.AuthMiddleware).Patch("/products/{id}", productsHandler.UpdateProductQuantity)
	r.With(users.AuthMiddleware).Delete("/products/{id}", productsHandler.DeleteProduct)

	ordersService := orders.NewService(repository, app.db)
	ordersHandler := orders.NewHandler(ordersService)
	r.With(users.AuthMiddleware).Post("/orders", ordersHandler.PlaceOrder)

	usersService := users.NewService(repository, app.db)
	usersHandler := users.NewHandler(usersService)
	r.Post("/users/signup", usersHandler.CreateUser)
	r.Post("/users/signin", usersHandler.LoginUser)

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr: app.config.addrr,
		Handler: h,
		WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second * 10,
		IdleTimeout: time.Minute,
	}

	log.Printf("Server has started ar addr %s", app.config.addrr)

	return srv.ListenAndServe()
}

type application struct {
	config config
	db *pgx.Conn
}

type config struct{
	addrr string
	db dbConfig
}

type dbConfig struct {
	dsn string
}