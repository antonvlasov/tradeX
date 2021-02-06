package database

import (
	"errors"
	"fmt"
	"log"
	"time"
)

//Get list of supported fields for sorting the order of events in reply
func supportedFields() []string {
	return []string{"date", "views", "clicks", "cost"}
}

type Event struct {
	Date   string  `json:"date"`
	Views  int     `json:"views"`
	Clicks int     `json:"clicks"`
	Cost   float64 `json:"cost"`
}
type Stat struct {
	Date   string  `json:"date"`
	Views  int     `json:"views"`
	Clicks int     `json:"clicks"`
	Cost   float64 `json:"cost"`
	Cpc    float64 `json:"cpc"`
	Cpm    float64 `json:"cpm"`
}

type Database int

//RPC wrapper for adding event
func (*Database) AddEvent(event Event, reply *int) error {
	err := Validate(event)
	if err != nil {
		return err
	}
	return addEvent(event)
}

//Validating that event contains correct data
func Validate(event Event) (err error) {
	err = verifyDate(event.Date)
	if err != nil {
		return
	}
	if event.Views < 0 || event.Clicks < 0 || event.Cost < 0 {
		return errors.New("views,clicks and cost should be positive")
	}
	return
}

//RPC wrapper for clearing the statistics
func (*Database) Clear(request int, reply *int) error {
	clearDB()
	return nil
}

type StatRequest struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Sortby string `json:"sortby"`
}

//RPC wrapper for selecting statistics
func (*Database) SelectStats(request StatRequest, reply *[]Stat) error {
	var err error
	*reply, err = selectStats(request.From, request.To, request.Sortby)
	if err != nil {
		return err
	}
	return nil
}

//Add event to the database
func addEvent(event Event) error {
	if db == nil {
		log.Fatal("db pointer is nil")
	}
	err := verifyDate(event.Date)
	if err != nil {
		fmt.Print("invalid date")
		return err
	}
	_, err = db.Exec("UPDATE Events SET views = views + $1, clicks = clicks + $2, cost = cost+ $3 WHERE date= $4",
		event.Views, event.Clicks, event.Cost, event.Date)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT OR IGNORE INTO Events VALUES($1,$2,$3,$4)", event.Date, event.Views, event.Clicks, event.Cost)
	if err != nil {
		return err
	}
	return nil
}

//Clear the database from all events
func clearDB() {
	if db == nil {
		log.Fatal("db pointer is nil")
	}
	db.Exec("DELETE FROM Events;")
}

//Select events for the period, ordered by field sortby
func selectStats(from, to, sortby string) (stats []Stat, err error) {
	if db == nil {
		log.Fatal("db pointer is nil")
	}
	if verifyDate(from) != nil || verifyDate(to) != nil {
		log.Print("invalid dates")
		return nil, errors.New("invalid dates")
	}
	sortby = getSortField(sortby)
	rows, err := db.Query("SELECT * FROM Events WHERE date BETWEEN $1 AND $2 ORDER BY "+sortby+" ASC", from, to, sortby)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var stat Stat
		rows.Scan(&stat.Date, &stat.Views, &stat.Clicks, &stat.Cost)
		countAverage(&stat)
		stats = append(stats, stat)
	}
	if stats == nil {
		stats = make([]Stat, 0)
	}
	return stats, nil
}

//Validating that sortfield is correct; setting it to default otherwise
func getSortField(sortby string) string {
	exists := false
	for _, field := range supportedFields() {
		if field == sortby {
			exists = true
			break
		}
	}
	if !exists {
		sortby = "date"
	}
	return sortby
}

//Ccount Cpc and Cpm
func countAverage(stat *Stat) {
	if stat.Clicks == 0 {
		stat.Cpc = -1
	} else {
		stat.Cpc = stat.Cost / float64(stat.Clicks)
	}
	if stat.Views == 0 {
		stat.Cpm = -1
	} else {
		stat.Cpm = stat.Cost / float64(stat.Views) * 1000
	}
}

//Verify that date is correct
func verifyDate(date string) error {
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return errors.New("date should be in format YYYY-MM-DD")
	}
	return nil
}
