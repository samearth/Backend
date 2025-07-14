package main

import (
	profile2 "github.com/MentorsPath/Backend/internal/api/profile"
	auth2 "github.com/MentorsPath/Backend/internal/user"
	"log"
	"net/http"
	"time"

	"github.com/MentorsPath/Backend/config"
	"github.com/MentorsPath/Backend/database"
	"github.com/MentorsPath/Backend/database/repositories"
	"github.com/MentorsPath/Backend/internal/api/handlers"
	"github.com/MentorsPath/Backend/internal/auth"
	"github.com/MentorsPath/Backend/internal/middleware"
	"github.com/MentorsPath/Backend/pkg/jwt"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	// Load config from config path (should include DBURL, JWT_SECRET, REFRESH_SECRET, PORT)
	if err := config.Init("./config"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to MySQL DB (RDS) using the DSN in DBURL
	db, err := database.InitDB(viper.GetString("DB_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate all models
	//if err := db.AutoMigrate(
	//	&models.User{},
	//	&models.Profile{},
	//	&models.MentorProfile{},
	//	&models.MenteeProfile{},
	//	&models.Skill{},
	//	&models.UserSkill{},
	//); err != nil {
	//	log.Fatalf("Failed to auto-migrate models: %v", err)
	//}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	profileRepo := repositories.NewProfileRepository(db)

	// Initialize JWT generators
	jwtGen := jwt.NewGenerator(viper.GetString("JWT_SECRET"))
	refreshGen := jwt.NewGenerator(viper.GetString("REFRESH_SECRET"))

	// Initialize services
	authService := auth.NewAuthService(userRepo, profileRepo, jwtGen, refreshGen)
	profileService := profile2.NewProfileService(profileRepo)

	userService := auth2.NewUserService(userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	profileHandler := handlers.NewProfileHandler(profileService, userService)

	// Initialize router
	r := mux.NewRouter()

	sessionHandler := handlers.NewSessionHandler(userService)

	// Public routes
	r.HandleFunc("/auth/signup", authHandler.Signup).Methods("POST")
	r.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	r.HandleFunc("/auth/refresh", authHandler.Refresh).Methods("POST")

	// Protected routesz
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware(viper.GetString("JWT_SECRET")))

	api.HandleFunc("/session", sessionHandler.GetCurrentUser).Methods("GET")

	// Profile routes (Mentor and Mentee CRUD)
	api.HandleFunc("/mentor", profileHandler.GetMentorProfile).Methods("GET")
	api.HandleFunc("/mentor", profileHandler.UpdateMentorProfile).Methods("PUT")

	api.HandleFunc("/mentee", profileHandler.GetMenteeProfile).Methods("GET")
	api.HandleFunc("/mentee", profileHandler.UpdateMenteeProfile).Methods("PUT")

	api.HandleFunc("/profile", profileHandler.GetProfile).Methods("GET")
	api.HandleFunc("/profile", profileHandler.UpdateProfile).Methods("PUT")

	// HTTP Server
	port := viper.GetString("PORT")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on port %s...", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
