package main

import (
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"

	"github.com/joho/godotenv"
	"github.com/plantexchange/app/handlers"
	"github.com/plantexchange/app/models"
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

	// Create some sample data
	createSampleData()

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
	apiRouter.HandleFunc("/listings", handlers.GetListings).Methods("GET")
	apiRouter.HandleFunc("/listings", handlers.CreateListing).Methods("POST")
	apiRouter.HandleFunc("/listings/{id}", handlers.GetListing).Methods("GET")
	apiRouter.HandleFunc("/listings/{id}", handlers.UpdateListing).Methods("PUT")
	apiRouter.HandleFunc("/listings/{id}", handlers.DeleteListing).Methods("DELETE")
	apiRouter.HandleFunc("/listings/search", handlers.SearchListings).Methods("GET")

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

// createSampleData creates initial data for development
func createSampleData() {
	// Check if we already have users in the database
	users := utils.GetUsers()
	if len(users) > 0 {
		log.Println("Sample data already exists, skipping creation")
		return
	}

	log.Println("Creating sample data...")

	// Create sample users
	user1 := models.User{
		Email:     "Raj@example.com",
		Username:  "Raj_gardener",
		Password:  utils.HashPassword("password123"),
		Name:      "Raj Sharma",
		Location:  "Delhi",
		Bio:       "Avid gardener with a focus on native plants and vegetables.",
		CreatedAt: time.Now(),
	}

	user2 := models.User{
		Email:     "badal@example.com",
		Username:  "badal_plants",
		Password:  utils.HashPassword("password123"),
		Name:      "Badal Singh",
		Location:  "Mumbai",
		Bio:       "Succulent collector and indoor plant enthusiast.",
		CreatedAt: time.Now(),
	}

	// Save users and get their IDs
	user1ID := utils.SaveUser(user1)
	user2ID := utils.SaveUser(user2)

	if user1ID == "" || user2ID == "" {
		log.Println("Failed to create sample users")
		return
	}

	log.Printf("Created users with IDs: %s, %s", user1ID, user2ID)

	// Create sample listings
	listing1 := models.Listing{
		UserID:      user1ID,
		Title:       "Monstera Deliciosa Cuttings",
		Description: "Healthy cuttings from my 3-year-old monstera plant. Well rooted and ready for potting.",
		Type:        "cutting",
		PlantType:   "indoor",
		Price:       15.00,
		TradeFor:    "Pothos varieties, philodendrons",
		Location:    "Delhi",
		Images:      []string{"https://images.unsplash.com/photo-1466781783364-36c955e42a7f"},
		CreatedAt:   time.Now(),
		Status:      "available",
	}

	listing2 := models.Listing{
		UserID:      user1ID,
		Title:       "Heirloom Tomato Seeds",
		Description: "Seeds from my prize-winning Cherokee Purple tomatoes. Great for warm climates.",
		Type:        "seed",
		PlantType:   "vegetable",
		Price:       5.00,
		TradeFor:    "Any herb seeds",
		Location:    "Delhi",
		Images:      []string{"https://images.unsplash.com/photo-1501004318641-b39e6451bec6"},
		CreatedAt:   time.Now(),
		Status:      "available",
	}

	listing3 := models.Listing{
		UserID:      user2ID,
		Title:       "Echeveria Succulent Collection",
		Description: "Set of 5 different echeveria varieties. All are 2 years old and in excellent health.",
		Type:        "plant",
		PlantType:   "succulent",
		Price:       25.00,
		TradeFor:    "Other rare succulents",
		Location:    "Hyderabad",
		Images:      []string{"https://images.unsplash.com/photo-1492282442770-077ea21f0f81"},
		CreatedAt:   time.Now(),
		Status:      "available",
	}

	// Save listings and get their IDs
	listing1ID := utils.SaveListing(listing1)
	listing2ID := utils.SaveListing(listing2)
	listing3ID := utils.SaveListing(listing3)

	if listing1ID == "" || listing2ID == "" || listing3ID == "" {
		log.Println("Failed to create sample listings")
		return
	}

	log.Printf("Created listings with IDs: %s, %s, %s", listing1ID, listing2ID, listing3ID)

	// Create sample messages
	message1 := models.Message{
		FromID:    user2ID,
		ToID:      user1ID,
		ListingID: listing1ID,
		Content:   "Hi Raj, I'm interested in your monstera cuttings. Would you consider trading for some of my rare succulents?",
		CreatedAt: time.Now(),
	}

	message2 := models.Message{
		FromID:    user1ID,
		ToID:      user2ID,
		ListingID: listing1ID,
		Content:   "Hi Badal, I'd definitely be interested in trading. What types of succulents do you have available?",
		CreatedAt: time.Now().Add(time.Hour),
	}

	utils.SaveMessage(message1)
	utils.SaveMessage(message2)

	// Create sample favorites
	utils.AddFavorite(user2ID, listing1ID) // Bob favorites Alice's monstera listing

	log.Println("Sample data created successfully")
}
