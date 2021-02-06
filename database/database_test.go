package database

import (
	"testing"
)

func TestDatabase(t *testing.T) {
	Connect("../resources/test/")
	clearDB()

	event1 := Event{"2021-02-04", 27000, 9000, 3000}
	event2 := Event{"2021-02-04", 13000, 946, 5690}
	event3 := Event{"2021-02-03", 2, 0, 3000}
	event4 := Event{"2021-02-05", 645, 238, 346}
	events := []Event{event1, event2, event3, event4}

	for i := range events {
		err := addEvent(events[i])
		if err != nil {
			t.Errorf("Error occured: %v", err)
		}
	}

	date1 := "2021-02-03"
	date2 := "2021-02-04"
	date3 := "2021-02-05"
	date4 := "2021-02-06"

	results_full := []Stat{
		{"2021-02-03", 2, 0, 3000, -1, 1.5e+06},
		{"2021-02-04", 40000, 9946, 8690, 0.8737180776191433, 217.25},
		{"2021-02-05", 645, 238, 346, 1.453781512605042, 536.4341085271318},
	}
	results_part := []Stat{
		{"2021-02-03", 2, 0, 3000, -1, 1.5e+06},
		{"2021-02-04", 40000, 9946, 8690, 0.8737180776191433, 217.25},
	}
	results_by_cost := []Stat{
		{"2021-02-05", 645, 238, 346, 1.453781512605042, 536.4341085271318},
		{"2021-02-03", 2, 0, 3000, -1, 1.5e+06},
		{"2021-02-04", 40000, 9946, 8690, 0.8737180776191433, 217.25},
	}
	sortby := "date"
	stats, err := selectStats(date1, date2, sortby)
	if err != nil {
		t.Errorf("Error occured: %v", err)
	}
	if len(stats) != len(results_part) {
		t.Errorf("wrong stats")
	}
	for i := range stats {
		if stats[i] != results_part[i] {
			t.Errorf("wrong stats, expected %v, got %v", results_part[i], stats[i])
		}
	}
	sortby = "gibberish"
	stats, err = selectStats(date1, date3, sortby)
	if err != nil {
		t.Errorf("Error occured: %v", err)
	}
	if len(stats) != len(results_full) {
		t.Errorf("wrong stats")
	}
	for i := range stats {
		if stats[i] != results_full[i] {
			t.Errorf("wrong stats, expected %v, got %v", results_full[i], stats[i])
		}
	}
	sortby = "cost"
	stats, err = selectStats(date1, date4, sortby)
	if err != nil {
		t.Errorf("Error occured: %v", err)
	}
	if len(stats) != len(results_by_cost) {
		t.Errorf("wrong stats")
	}
	for i := range stats {
		if stats[i] != results_by_cost[i] {
			t.Errorf("wrong stats, expected %v, got %v", results_by_cost[i], stats[i])
		}
	}
	stats, err = selectStats(date4, date1, sortby)
	if err != nil {
		t.Errorf("Error occured: %v", err)
	}
	if len(stats) != 0 {
		t.Errorf("wrong stats")
	}

	clearDB()

	stats, err = selectStats(date1, date4, sortby)
	if err != nil {
		t.Errorf("Error occured: %v", err)
	}
	if len(stats) != 0 {
		t.Errorf("wrong stats")
	}

	Close()
}

func TestWrappers(t *testing.T) {
	Connect("../resources/test/")
	clearDB()

	handler := new(Database)
	var mock *int

	event1 := Event{"2021-02-04", 27000, 9000, 3000}
	event2 := Event{"2021-02-04", 13000, 946, 5690}
	event3 := Event{"2021-02-03", 2, 0, 3000}
	event4 := Event{"2021-02-05", 645, 238, 346}
	events := []Event{event1, event2, event3, event4}

	for i := range events {
		err := handler.AddEvent(events[i], mock)
		if err != nil {
			t.Errorf("Error occured: %v", err)
		}
	}

	date1 := "2021-02-03"
	date2 := "2021-02-04"
	date3 := "2021-02-05"
	date4 := "2021-02-06"

	results_full := []Stat{
		{"2021-02-03", 2, 0, 3000, -1, 1.5e+06},
		{"2021-02-04", 40000, 9946, 8690, 0.8737180776191433, 217.25},
		{"2021-02-05", 645, 238, 346, 1.453781512605042, 536.4341085271318},
	}
	results_part := []Stat{
		{"2021-02-03", 2, 0, 3000, -1, 1.5e+06},
		{"2021-02-04", 40000, 9946, 8690, 0.8737180776191433, 217.25},
	}
	results_by_cost := []Stat{
		{"2021-02-05", 645, 238, 346, 1.453781512605042, 536.4341085271318},
		{"2021-02-03", 2, 0, 3000, -1, 1.5e+06},
		{"2021-02-04", 40000, 9946, 8690, 0.8737180776191433, 217.25},
	}
	sortby := "date"
	var stats []Stat
	err := handler.SelectStats(StatRequest{date1, date2, sortby}, &stats)
	if err != nil {
		t.Errorf("Error occured: %v", err)
	}
	if len(stats) != len(results_part) {
		t.Errorf("wrong stats")
	}
	for i := range stats {
		if stats[i] != results_part[i] {
			t.Errorf("wrong stats, expected %v, got %v", results_part[i], stats[i])
		}
	}
	sortby = "gibberish"
	err = handler.SelectStats(StatRequest{date1, date3, sortby}, &stats)
	if err != nil {
		t.Errorf("Error occured: %v", err)
	}
	if len(stats) != len(results_full) {
		t.Errorf("wrong stats")
	}
	for i := range stats {
		if stats[i] != results_full[i] {
			t.Errorf("wrong stats, expected %v, got %v", results_full[i], stats[i])
		}
	}
	sortby = "cost"

	err = handler.SelectStats(StatRequest{date1, date4, sortby}, &stats)
	if err != nil {
		t.Errorf("Error occured: %v", err)
	}
	if len(stats) != len(results_by_cost) {
		t.Errorf("wrong stats")
	}
	for i := range stats {
		if stats[i] != results_by_cost[i] {
			t.Errorf("wrong stats, expected %v, got %v", results_by_cost[i], stats[i])
		}
	}

	err = handler.SelectStats(StatRequest{date4, date1, sortby}, &stats)
	if err != nil {
		t.Errorf("Error occured: %v", err)
	}
	if len(stats) != 0 {
		t.Errorf("wrong stats")
	}

	handler.Clear(0, mock)

	err = handler.SelectStats(StatRequest{date1, date4, sortby}, &stats)
	if err != nil {
		t.Errorf("Error occured: %v", err)
	}
	if len(stats) != 0 {
		t.Errorf("wrong stats")
	}

	Close()
}

func TestValidation(t *testing.T) {
	stringEvent1 := Event{"2021-02-04", 27000, 9000, 3000}
	stringEvent2 := Event{"2021-02-04", 0, 0, 30}
	stringEvent3 := Event{"2021-02-04", -20, 9000, 3000}
	stringEvent4 := Event{"2021.02.04", 10, 9000, 3000}

	input := []Event{stringEvent1, stringEvent2, stringEvent3, stringEvent4}

	validated := []bool{true, true, false, false}
	for i := range input {
		err := Validate(input[i])
		if (err == nil) != validated[i] {
			t.Errorf("incorrect validation on item %v: got %v", input[i], err == nil)
		}
	}
}
