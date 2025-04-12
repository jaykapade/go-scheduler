package greetings

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jaykapade/mail-scheduler/internal/auth"
	"github.com/jaykapade/mail-scheduler/internal/db"
)

type Greeting struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"-"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreateGreetingHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := auth.GetUserID(r)
	userID, _ := uuid.Parse(userIDStr)

	var g Greeting
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	g.ID = uuid.New()
	g.UserID = userID
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()

	query := `INSERT INTO greetings (id, user_id, subject, body, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = db.Pool.Exec(context.Background(), query, g.ID, g.UserID, g.Subject, g.Body, g.CreatedAt, g.UpdatedAt)
	if err != nil {
		http.Error(w, "db insert failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(g)
}

func ListGreetingsHandler(w http.ResponseWriter, r *http.Request) {
	userId := auth.GetUserID(r)

	query := `SELECT id, subject, body, created_at, updated_at
			  FROM greetings
			  WHERE user_id = $1`
	rows, err := db.Pool.Query(context.Background(), query, userId)
	if err != nil {
		http.Error(w, "db query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var greetings []Greeting
	for rows.Next() {
		var g Greeting
		err = rows.Scan(&g.ID, &g.Subject, &g.Body, &g.CreatedAt, &g.UpdatedAt)
		if err != nil {
			http.Error(w, "db scan failed", http.StatusInternalServerError)
			return
		}
		greetings = append(greetings, g)
	}

	json.NewEncoder(w).Encode(greetings)
}

func UpdateGreetingHandler(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	greetingID := chi.URLParam(r, "id")

	var g Greeting
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	query := `UPDATE greetings SET subject=$1, body=$2, updated_at=NOW() where id=$3 AND user_id=$4`

	res, err := db.Pool.Exec(context.Background(), query, g.Subject, g.Body, greetingID, userID)
	if err != nil {
		http.Error(w, "db update failed", http.StatusInternalServerError)
		return
	}
	if res.RowsAffected() == 0 {
		http.Error(w, "greeting not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "updated successfully")
}

func DeleteGreetingHandler(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	greetingID := chi.URLParam(r, "id")

	query := `DELETE FROM greetings WHERE id=$1 AND user_id=$2`
	res, err := db.Pool.Exec(context.Background(), query, greetingID, userID)
	if err != nil {
		http.Error(w, "db delete failed", http.StatusInternalServerError)
		return
	}
	if res.RowsAffected() == 0 {
		http.Error(w, "greeting not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "deleted successfully")
}
