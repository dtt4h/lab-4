package main

import (
	"context"
	_ "embed"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"lab-4/backend/internal/config"
	"lab-4/backend/internal/database"
	"lab-4/backend/internal/handler"
	"lab-4/backend/internal/mailer"
	"lab-4/backend/internal/middleware"
	"lab-4/backend/internal/repository"
	"lab-4/backend/internal/service"
)

func main() {
	configureLogger()
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is required")
	}
	pool, err := database.NewPool(context.Background(), config.DatabaseConfig{URL: os.Getenv("DATABASE_URL")})
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	if err := database.ApplyMigrations(context.Background(), pool, envOr("MIGRATIONS_DIR", "/migrations")); err != nil {
		log.Fatal(err)
	}
	bootstrapAdmin(context.Background(), pool, os.Getenv("BOOTSTRAP_ADMIN_EMAIL"))
	users := service.NewUsersService(repository.NewUsersRepository(pool))
	incidents := service.NewIncidentsService(repository.NewIncidentsRepository(pool))
	employees := service.NewEmployeesService(repository.NewEmployeesRepository(pool))
	hotels := service.NewHotelsService(repository.NewHotelsRepository(pool))
	locations := service.NewLocationsService(repository.NewLocationsRepository(pool))
	categories := service.NewIncidentCategoriesService(repository.NewIncidentCategoriesRepository(pool))
	bookings := service.NewBookingsService(repository.NewBookingsRepository(pool))
	rooms := service.NewRoomsService(repository.NewRoomsRepository(pool))
	authHandler := handler.NewAuthHandler(users, secret, 15*time.Minute)
	incidentsHandler := handler.NewIncidentsHandler(incidents)
	employeesHandler := handler.NewEmployeesHandler(employees)
	catalogsHandler := handler.NewCatalogsHandler(hotels, locations, categories)
	bookingMailer := mailer.NewSMTP(mailer.SMTPConfig{Host: os.Getenv("SMTP_HOST"), Port: envInt("SMTP_PORT", 587), Username: os.Getenv("SMTP_USERNAME"), Password: os.Getenv("SMTP_PASSWORD"), From: os.Getenv("SMTP_FROM")})
	bookingsHandler := handler.NewBookingsHandler(bookings, rooms, users, bookingMailer)
	adminHandler := handler.NewAdminHandler(users)
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger)
	r.Use(middleware.CORS(envOr("FRONTEND_ORIGIN", "http://localhost:5173")))
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	r.Get("/openapi.yaml", serveOpenAPI)
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})
	r.Get("/swagger/", serveSwaggerUI)
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/public/hotels", catalogsHandler.Hotels)
		r.Get("/public/rooms", bookingsHandler.Rooms)
		r.Post("/public/bookings", bookingsHandler.PublicCreate)
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.Refresh)
		r.Post("/auth/logout", authHandler.Logout)
		r.With(middleware.RequireAuth(secret)).Get("/auth/me", authHandler.Me)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAuth(secret))
			r.With(middleware.RequireRoles(users, "security_manager", "hotel_manager", "super_admin")).Get("/hotels", catalogsHandler.Hotels)
			r.With(middleware.RequireRoles(users, "security_manager", "hotel_manager", "super_admin")).Get("/categories", catalogsHandler.Categories)
			r.With(middleware.RequireRoles(users, "security_manager", "hotel_manager", "super_admin")).Get("/locations", catalogsHandler.Locations)
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRoles(users, "security_manager", "super_admin"))
				r.Get("/incidents", incidentsHandler.List)
				r.Post("/incidents", incidentsHandler.Create)
				r.Get("/incidents/{id}", incidentsHandler.GetByID)
				r.Put("/incidents/{id}", incidentsHandler.Update)
				r.Delete("/incidents/{id}", incidentsHandler.Delete)
			})
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRoles(users, "hotel_manager", "super_admin"))
				r.Get("/employees", employeesHandler.List)
				r.Post("/employees", employeesHandler.Create)
				r.Get("/employees/{id}", employeesHandler.GetByID)
				r.Put("/employees/{id}", employeesHandler.Update)
				r.Delete("/employees/{id}", employeesHandler.Delete)
				r.Post("/hotels", catalogsHandler.CreateHotel)
				r.Delete("/hotels/{id}", catalogsHandler.DeleteHotel)
				r.Post("/categories", catalogsHandler.CreateCategory)
				r.Delete("/categories/{id}", catalogsHandler.DeleteCategory)
				r.Post("/locations", catalogsHandler.CreateLocation)
				r.Delete("/locations/{id}", catalogsHandler.DeleteLocation)
			})
			r.With(middleware.RequirePermission(users, "room.read")).Get("/rooms", bookingsHandler.Rooms)
			r.With(middleware.RequirePermission(users, "booking.create")).Post("/bookings", bookingsHandler.Create)
			r.With(middleware.RequirePermission(users, "booking.read.own")).Get("/bookings/my", bookingsHandler.My)
			r.With(middleware.RequirePermission(users, "booking.pay")).Post("/bookings/{id}/pay", bookingsHandler.Pay)
			r.With(middleware.RequirePermission(users, "booking.read")).Get("/admin/bookings", bookingsHandler.AdminList)
			r.With(middleware.RequirePermission(users, "booking.update")).Patch("/admin/bookings/{id}/status", bookingsHandler.AdminStatus)
			r.With(middleware.RequirePermission(users, "roles.manage")).Get("/admin/users", adminHandler.ListUsers)
			r.With(middleware.RequirePermission(users, "roles.manage")).Patch("/admin/users/{id}/role", adminHandler.SetRole)
			r.With(middleware.RequireRoles(users, "room_manager", "super_admin")).Get("/admin/room-types", bookingsHandler.RoomTypes)
			r.With(middleware.RequireRoles(users, "room_manager", "super_admin")).Post("/admin/room-types", bookingsHandler.CreateRoomType)
			r.With(middleware.RequireRoles(users, "room_manager", "super_admin")).Post("/admin/rooms", bookingsHandler.CreateRoom)
		})
	})
	address := os.Getenv("HTTP_ADDR")
	if address == "" {
		address = ":8080"
	}
	server := &http.Server{
		Addr:              address,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	errs := make(chan error, 1)
	go func() { errs <- server.ListenAndServe() }()
	log.Printf("listening on %s", address)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	select {
	case err := <-errs:
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	case sig := <-stop:
		log.Printf("shutting down after %s", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("graceful shutdown failed: %v", err)
		}
	}
}

//go:embed openapi.yaml
var openAPISpec []byte

func serveOpenAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	_, _ = w.Write(openAPISpec)
}

func serveSwaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = io.WriteString(w, `<!doctype html><html><head><meta charset="utf-8"><title>Hotel Ops API</title><link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"></head><body><div id="swagger-ui"></div><script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script><script>SwaggerUIBundle({url:"/openapi.yaml",dom_id:"#swagger-ui",persistAuthorization:true})</script></body></html>`)
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil || value < 1 || value > 65535 {
		return fallback
	}
	return value
}

func configureLogger() {
	path := os.Getenv("LOG_FILE")
	if path == "" {
		return
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("cannot open log file: %v", err)
		return
	}
	log.SetOutput(io.MultiWriter(os.Stdout, file))
}

func bootstrapAdmin(ctx context.Context, pool *pgxpool.Pool, email string) {
	if email == "" {
		return
	}
	if _, err := pool.Exec(ctx, `UPDATE users SET role='super_admin' WHERE email=$1`, email); err != nil {
		log.Printf("cannot set bootstrap admin: %v", err)
	}
}
