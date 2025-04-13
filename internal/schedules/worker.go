package schedules

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jaykapade/mail-scheduler/internal/db"
	"github.com/jaykapade/mail-scheduler/internal/mailer"
)

func StartEmailWorker() {
	mailer.InitMailer()

	for {
		log.Println("running email worker")
		checkAndSendEmails()
		time.Sleep(time.Second * 10)
	}
}

func checkAndSendEmails() {
	now := time.Now()
	// Get all active schedules that are due to be sent
	query := `SELECT 
				id, user_id, greeting_id, start_date, frequency_type, frequency_interval, scheduled_time, is_active 
	FROM schedules 
				WHERE is_active = true 
		AND start_date <= $1
		AND (
					frequency_type = 'once' 
					OR frequency_type = 'daily'
				)
				AND (
					scheduled_time <= $1
				)
				AND (
					last_sent_at IS NULL 
					OR last_sent_at::date < current_date
		)
	`
	rows, err := db.Pool.Query(context.Background(), query, now)
	if err != nil {
		log.Println("Error querying schedules:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var schedule Schedule
		if err := rows.Scan(&schedule.ID, &schedule.UserID, &schedule.GreetingID, &schedule.StartDate, &schedule.FrequencyType, &schedule.FrequencyInterval, &schedule.ScheduledTime, &schedule.IsActive); err != nil {
			log.Println("Error scanning schedule:", err)
			continue
		}
		sendEmail(schedule)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over schedules:", err)
	}
}

func sendEmail(schedule Schedule) {
	log.Printf("📧 Sending email for schedule %s to user %s at %s", schedule.ID, schedule.UserID, schedule.StartDate)
	// TODO: use from Greeting itself
	fromEmail := os.Getenv("FROM_EMAIL")
	fromName := os.Getenv("FROM_NAME")

	// Get user email
	var userEmail string
	query := `SELECT email FROM users WHERE id = $1 LIMIT 1`
	err := db.Pool.QueryRow(context.Background(), query, schedule.UserID).Scan(&userEmail)
	if err != nil {
		log.Printf("Error getting user email: %v", err)
		updateScheduleError(schedule, fmt.Errorf("failed to get user email: %v", err))
		return
	}

	// Get greeting content
	var subject, body string
	query = `SELECT subject, body FROM greetings WHERE id = $1 LIMIT 1`
	err = db.Pool.QueryRow(context.Background(), query, schedule.GreetingID).Scan(&subject, &body)
	if err != nil {
		log.Printf("Error getting greeting content: %v", err)
		updateScheduleError(schedule, fmt.Errorf("failed to get greeting content: %v", err))
		return
	}

	// Send email
	err = mailer.SendEmail(fromEmail, fromName, userEmail, "Mailer Scheduler", subject, body)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		updateScheduleError(schedule, err)
		return
	}
	updateScheduleStatus(schedule)
}

func updateScheduleStatus(schedule Schedule) {
	// Update the schedule status in the database
	query := `UPDATE schedules SET last_sent_at = $1, latest_error = NULL WHERE id = $2`
	_, err := db.Pool.Exec(context.Background(), query, time.Now(), schedule.ID)

	if err != nil {
		log.Println("Error updating schedule status:", err)
	}
}

func updateScheduleError(schedule Schedule, err error) {
	query := `UPDATE schedules SET latest_error = $1 WHERE id = $2`
	_, err = db.Pool.Exec(context.Background(), query, err.Error(), schedule.ID)
	if err != nil {
		log.Printf("Error updating schedule: %v", err)
	}
}
