package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jaykapade/mail-scheduler/internal/auth"
	"github.com/jaykapade/mail-scheduler/internal/db"
	"github.com/jaykapade/mail-scheduler/internal/greetings"
	"github.com/jaykapade/mail-scheduler/internal/schedules"
)

func main() {
	err := db.Init()
	if err != nil {
		log.Fatalf("❌ DB Init failed: %v", err)
	}

	r := chi.NewRouter()

	r.Post("/register", auth.RegisterHandler)
	r.Post("/login", auth.LoginHandler)

	r.Route("/greetings", func(r chi.Router) {
		r.Get("/", auth.JWTMiddleware(greetings.ListGreetingsHandler))
		r.Post("/", auth.JWTMiddleware(greetings.CreateGreetingHandler))
		r.Put("/{id}", auth.JWTMiddleware(greetings.UpdateGreetingHandler))
		r.Delete("/{id}", auth.JWTMiddleware(greetings.DeleteGreetingHandler))
	})

	r.Route("/schedules", func(r chi.Router) {
		r.Get("/", auth.JWTMiddleware(schedules.ListSchedulesHandler))
		r.Post("/", auth.JWTMiddleware(schedules.CreateScheduleHandler))
	})

	// Start the email worker in a new goroutine
	go schedules.StartEmailWorker()

	log.Println("🚀 Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
