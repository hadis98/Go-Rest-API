package models

import (
	"fmt"
	"time"

	"github.com/hadis98/rest-api/db"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int64
}

// var events []Event = []Event{}

func (e *Event) Save() error {
	fmt.Println("[event save: ]", e.ID, e.UserID)
	//later add to database
	query := `
	INSERT INTO events (name,description,location,dateTime,user_id) 
	VALUES(?,?,?,?,?)`
	stmt, err := db.DB.Prepare(query) // when we use prepare, it is stored in memory so it will be reused in a highly efficient way
	if err != nil {
		return err
	}

	defer stmt.Close()
	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	//* we use Exec whenever we have a query that changes data in the database; like inserting data, updating data, ..

	if err != nil {
		return err
	}
	id, err := result.LastInsertId() //get the automatically generated id
	e.ID = id
	return err
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query) //* Query is used when we want to get back bunch of rows => so we use it to fetch data
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []Event
	for rows.Next() { //loop till we have a row to read
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID) //reads the content of the row we're currently reading
		fmt.Println("**before error**", err)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id = ?"
	row := db.DB.QueryRow(query, id) //*QueryRow is best choice when we know we only gonna get exactly one row as the result
	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (event Event) Update() error {
	query := `
		UPDATE events 
		SET name = ? , description = ?, location = ?, dateTime = ?
		WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.ID)
	return err
}

func (event Event) Delete() error {
	query := "DELETE FROM events WHERE id = ?"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.ID)
	return err

}

func (e Event) Register(userId int64) error {
	query := "INSERT INTO registration(event_id,user_id) VALUES(?,?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(e.ID, userId)
	return err

}

func (e Event) CancelRegistration(userId int64) error {

	query := "DELETE FROM registration WHERE event_id = ? AND user_id = ?"
	stmt, err := db.DB.Prepare(query)
	fmt.Println("[CancelRegistration]: ", stmt, "   ", err)
	if err != nil {
		return err
	}

	defer stmt.Close()
	result, err := stmt.Exec(e.ID, userId)
	fmt.Println("Result: ", result)
	return err
}
