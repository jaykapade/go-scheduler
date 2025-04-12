package schedules

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jaykapade/mail-scheduler/internal/auth"
	"github.com/jaykapade/mail-scheduler/internal/db"
)

type Schedule struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"-"`
	GreetingID        uuid.UUID `json:"greeting_id"`
	StartDate         time.Time `json:"start_date"`
	FrequencyType     string    `json:"frequency_type"`
	FrequencyInterval int       `json:"frequency_interval"`
	TimeOfDay         string    `json:"time_of_day"`
	IsActive          bool      `json:"is_active"`
	LastSentAt        time.Time `json:"last_sent_at"`
	LatestError       *string   `json:"latest_error"`
}

func CreateScheduleHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := auth.GetUserID(r)
	userId, _ := uuid.Parse(userIdStr)

	var schedule Schedule
	err := json.NewDecoder(r.Body).Decode(&schedule)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate that the start date is in the future
	if schedule.StartDate.Before(time.Now()) {
		http.Error(w, "start date must be in the future", http.StatusBadRequest)
		return
	}

	// Parse the time_of_day as a proper time object
	timeOfDayParsed, err := time.Parse("15:04:05", schedule.TimeOfDay)
	if err != nil {
		http.Error(w, "invalid time format for time_of_day", http.StatusBadRequest)
		return
	}

	var greetingID uuid.UUID
	query := `SELECT id FROM greetings WHERE id = $1 AND user_id = $2`
	err = db.Pool.QueryRow(context.Background(), query, schedule.GreetingID, userId).Scan(&greetingID)
	if err != nil {
		http.Error(w, "greeting not found", http.StatusNotFound)
		return
	}

	schedule.ID = uuid.New()
	schedule.UserID = userId

	// Insert the schedule into the database
	query = `INSERT INTO SCHEDULES (id, user_id, greeting_id, start_date, frequency_type, frequency_interval, time_of_day, is_active, last_sent_at, latest_error)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = db.Pool.Exec(context.Background(), query, schedule.ID, schedule.UserID, schedule.GreetingID, schedule.StartDate, schedule.FrequencyType, schedule.FrequencyInterval, timeOfDayParsed, schedule.IsActive, schedule.LastSentAt, schedule.LatestError)
	if err != nil {
		http.Error(w, "failed to create schedule", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(schedule)
}

func ListSchedulesHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := auth.GetUserID(r)
	userId, _ := uuid.Parse(userIdStr)

	query := `SELECT id, greeting_id, start_date, frequency_type, frequency_interval, time_of_day, is_active, last_sent_at, latest_error FROM schedules WHERE user_id = $1`
	rows, err := db.Pool.Query(context.Background(), query, userId)
	if err != nil {
		http.Error(w, "failed to fetch schedules", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var schedules []Schedule
	for rows.Next() {
		var schedule Schedule
		err = rows.Scan(&schedule.ID, &schedule.GreetingID, &schedule.StartDate, &schedule.FrequencyType, &schedule.FrequencyInterval, &schedule.TimeOfDay, &schedule.IsActive, &schedule.LastSentAt, &schedule.LatestError)
		if err != nil {
			http.Error(w, "failed to scan schedule", http.StatusInternalServerError)
			return
		}
		schedules = append(schedules, schedule)
	}

	json.NewEncoder(w).Encode(schedules)
}
