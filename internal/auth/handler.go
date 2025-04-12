package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jaykapade/mail-scheduler/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}

	query := `INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3)`
	uid := uuid.New()

	_, err = db.Pool.Exec(context.Background(), query, uid, req.Email, string(hashedPwd))
	if err != nil {
		http.Error(w, "email already exists", http.StatusConflict)
		return
	}

	token, err := GenerateJWT(uid.String())
	if err != nil {
		http.Error(w, "error generating token", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	var user User
	query := `SELECT id, password_hash FROM users WHERE email=$1 LIMIT 1`
	row := db.Pool.QueryRow(context.Background(), query, req.Email)
	err = row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "wrong password", http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(user.ID.String())
	if err != nil {
		http.Error(w, "error generating token", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
