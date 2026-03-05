package main

import (
	"context"
	"fmt"
	"log"
	"saloon-backend/config"
	"saloon-backend/database"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	ctx := context.Background()

	fmt.Println("--- SALONS ---")
	rows, _ := db.Query(ctx, "SELECT id, name, opening_time::text, closing_time::text FROM salons")
	for rows.Next() {
		var id, name, open, close string
		rows.Scan(&id, &name, &open, &close)
		fmt.Printf("ID: %s, Name: %s, Open: %s, Close: %s\n", id, name, open, close)
	}

	fmt.Println("\n--- STAFF WORKING HOURS ---")
	rows, _ = db.Query(ctx, `
		SELECT s.name, h.day_of_week, h.start_time::text, h.end_time::text, h.is_off, s.id
		FROM staff_working_hours h 
		JOIN staff s ON s.id = h.staff_id
	`)
	for rows.Next() {
		var name, start, end, id string
		var day int
		var off bool
		rows.Scan(&name, &day, &start, &end, &off, &id)
		fmt.Printf("Staff: %s, Day: %d, Start: %s, End: %s, Off: %v, ID: %s\n", name, day, start, end, off, id)
	}
}
