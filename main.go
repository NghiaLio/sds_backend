package main

import (
	"log"
	"net/http"
	"os"

	"mobile-crud-backend/config"
	"mobile-crud-backend/handler"
	"mobile-crud-backend/repository"
	"mobile-crud-backend/service"
)

// enableCORS middleware handles Cross-Origin Resource Sharing.
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// requestLogger middleware prints incoming request details.
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[HTTP] %s %s - Remote: %s", r.Method, r.URL.String(), r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("Starting Mobile CRUD Backend with JWT Auth (SDS Mobile)...")

	// Database initialization
	dbPath := "tasks.db"
	db, err := config.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close()
	log.Printf("SQLite database loaded from '%s'\n", dbPath)

	// Auth Secret Key setup
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super-secret-key-mobile-crud-backend-2026"
		log.Println("[WARNING] JWT_SECRET environment variable is empty. Using fallback secret.")
	}

	// Layer initialization (Dependency Injection)
	// Repositories
	userRepo := repository.NewSQLiteUserRepository(db)
	categoryRepo := repository.NewSQLiteCategoryRepository(db)
	productRepo := repository.NewSQLiteProductRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, jwtSecret)
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo, categoryRepo)

	// Handlers and Middlewares
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	pingHandler := handler.NewPingHandler()
	authMiddleware := handler.NewAuthMiddleware(jwtSecret)

	// Router setup
	mux := http.NewServeMux()
	
	// Register Routes
	authHandler.RegisterRoutes(mux)
	categoryHandler.RegisterRoutes(mux, authMiddleware)
	productHandler.RegisterRoutes(mux, authMiddleware)
	pingHandler.RegisterRoutes(mux)

	// Apply Middlewares
	finalHandler := requestLogger(enableCORS(mux))

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	addr := ":" + port
	log.Printf("Server listening on http://localhost%s\n", addr)
	log.Println("Press Ctrl+C to stop the server.")

	if err := http.ListenAndServe(addr, finalHandler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
