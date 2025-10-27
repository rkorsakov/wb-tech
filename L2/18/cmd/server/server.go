package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	hndlr "18/internal/handler"
	"18/internal/service"
)

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	}
}

func main() {
	port := flag.String("port", ":8080", "server port")
	flag.Parse()

	log.Println("Server starting on port:", *port)

	calendar := service.NewCalendarService()
	handler := hndlr.NewHandler(calendar)

	http.HandleFunc("/create_event", loggingMiddleware(handler.CreateEvent))
	http.HandleFunc("/update_event", loggingMiddleware(handler.UpdateEvent))
	http.HandleFunc("/delete_event", loggingMiddleware(handler.DeleteEvent))
	http.HandleFunc("/events_for_day", loggingMiddleware(handler.GetDayEvents))
	http.HandleFunc("/events_for_week", loggingMiddleware(handler.GetWeekEvents))
	http.HandleFunc("/events_for_month", loggingMiddleware(handler.GetMonthEvents))

	err := http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}