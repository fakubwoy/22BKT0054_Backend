package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/fakubwoy/go-file-share/internal/auth"
	"github.com/fakubwoy/go-file-share/internal/config"
	"github.com/fakubwoy/go-file-share/internal/models"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func RegisterHandler(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error decoding request: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		passwordHash, err := auth.HashPassword(req.Password)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		user := &models.User{
			Email:        req.Email,
			PasswordHash: passwordHash,
		}

		if err := user.Create(db); err != nil {
			log.Printf("Database error creating user: %v", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		token, err := auth.GenerateJWTToken(user.ID, cfg)
		if err != nil {
			log.Printf("Error generating token: %v", err)
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		response := AuthResponse{Token: token}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func LoginHandler(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Email == "" || req.Password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		user, err := models.GetUserByEmail(db, req.Email)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := auth.GenerateJWTToken(user.ID, cfg)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		response := AuthResponse{Token: token}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
