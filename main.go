package main

import (
	"log"
	"net/http"
	"path/filepath"

	// "time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"

	"github.com/joho/godotenv"
	"github.com/plantexchange/app/handlers"

	// "github.com/plantexchange/app/models"
	"github.com/plantexchange/app/utils"
)

// Session store is initialized in utils/auth.go
func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found. Using environment variables.")
	}
}
func main() {
	// Initialize database
	utils.InitDB()
	defer utils.CloseDB()

	// Set up router
	r := mux.NewRouter()

	// Static files
	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// API Routes
	apiRouter := r.PathPrefix("/api").Subrouter()

	// Auth routes
	apiRouter.HandleFunc("/register", handlers.Register).Methods("POST")
	apiRouter.HandleFunc("/login", handlers.Login).Methods("POST")
	apiRouter.HandleFunc("/logout", handlers.Logout).Methods("POST")
	apiRouter.HandleFunc("/check-auth", handlers.CheckAuth).Methods("GET")

	// User routes
	apiRouter.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	apiRouter.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
	apiRouter.HandleFunc("/users/current", handlers.GetCurrentUser).Methods("GET")

	// Listing routes
	apiRouter.HandleFunc("/listings/search", handlers.SearchListings).Methods("GET")
	apiRouter.HandleFunc("/listings", handlers.GetListings).Methods("GET")
	apiRouter.HandleFunc("/listings", handlers.CreateListing).Methods("POST")
	apiRouter.HandleFunc("/listings/{id}", handlers.GetListing).Methods("GET")
	apiRouter.HandleFunc("/listings/{id}", handlers.UpdateListing).Methods("PUT")
	apiRouter.HandleFunc("/listings/{id}", handlers.DeleteListing).Methods("DELETE")

	// Message routes
	apiRouter.HandleFunc("/messages", handlers.GetMessages).Methods("GET")
	apiRouter.HandleFunc("/messages", handlers.SendMessage).Methods("POST")
	apiRouter.HandleFunc("/messages/{id}", handlers.GetMessage).Methods("GET")
	apiRouter.HandleFunc("/conversations", handlers.GetConversations).Methods("GET")
	apiRouter.HandleFunc("/conversations/{userId}", handlers.GetConversation).Methods("GET")

	// Favorites routes
	apiRouter.HandleFunc("/favorites", handlers.ToggleFavorite).Methods("POST")
	apiRouter.HandleFunc("/favorites", handlers.GetFavorites).Methods("GET")

	// HTML routes - serve appropriate templates
	r.HandleFunc("/", serveTemplate("index.html")).Methods("GET")
	r.HandleFunc("/login", serveTemplate("login.html")).Methods("GET")
	r.HandleFunc("/register", serveTemplate("register.html")).Methods("GET")
	r.HandleFunc("/profile", serveTemplate("profile.html")).Methods("GET")
	r.HandleFunc("/dashboard", serveTemplate("dashboard.html")).Methods("GET")
	r.HandleFunc("/messages", serveTemplate("messages.html")).Methods("GET")
	r.HandleFunc("/create-listing", serveTemplate("create-listing.html")).Methods("GET")
	r.HandleFunc("/listing/{id}", serveTemplate("listing.html")).Methods("GET")

	// Catch-all route to redirect to index
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	})

	// CORS setup
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Start server
	port := "8080"
	log.Printf("Starting server on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, corsMiddleware.Handler(r)))
}

// serveTemplate serves HTML templates
func serveTemplate(tmpl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join("templates", tmpl)
		http.ServeFile(w, r, path)
	}
}

