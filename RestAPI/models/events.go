package models

import (
	"fmt"
	"time"

	"example.com/APIs/DB"
)

type Event struct {
	Id          int64
	Title       string `binding:"required"`
	Description string `binding:"required"`
	Location    string `binding:"required"`
	DateTime    time.Time
	UserID      int64
}

var Events []Event

func (e *Event) Save() error {
	query := `INSERT INTO events(title, description, location, datetime, user_id) VALUES(?,?,?,?,?)`
	stmt, err := DB.DB.Prepare(query)
	if err != nil {
		fmt.Printf("Cannot process query. Error: %v", err)
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(e.Title, e.Description, e.Location, e.DateTime.Local().UTC().Nanosecond(), e.UserID)

	if err != nil {
		return err
	}
	insertedId, err := result.LastInsertId()
	e.Id = insertedId

	return err
}

func GetAllEvents() ([]Event, error) {
	query := `SELECT * FROM events`
	rows, err := DB.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.Id, &event.Title, &event.Description, &event.Location, &event.DateTime, &event.UserID)
		fmt.Printf("User id: %v", event.UserID)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, err
}

func GetEventById(eventID int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id = ?"
	row := DB.DB.QueryRow(query, eventID)
	var event Event

	err := row.Scan(&event.Id, &event.Title, &event.Description, &event.Location, &event.DateTime, &event.UserID)
	if err != nil {
		return nil, err
	}
	return &event, err
}

func (e *Event) UpdateEvent() error {
	query := `
	UPDATE events SET title = ?, description = ?, location = ?, dateTime = ?
	WHERE id = ?
	`
	stmt, err := DB.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(e.Title, e.Description, e.Location, e.DateTime.Local().UTC().Nanosecond(), e.Id)
	return err

}

func (e *Event) DeleteEvent() error {
	query := `
	DELETE FROM events WHERE id = ?
	`
	stmt, err := DB.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(e.Id)
	return err

}

func (e *Event) RegisterEvent(userID int64) error {
	query := `
	INSERT INTO registrations(event_id,user_id) VALUES(?, ?)
	`
	stmt, err := DB.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(e.Id, userID)

	return err
}

func (e *Event) CancelEvent(userID int64) error {
	query := `DELETE FROM registrations WHERE event_id = ? AND user_id = ?`
	stmt, err := DB.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(e.Id, userID)

	return err
}
