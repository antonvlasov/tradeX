package server

import (
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
	"testing"

	"github.com/antonvlasov/tradeX/database"
)

func TestServer(t *testing.T) {
	go Serve("../resources/test/")
	var client *rpc.Client
	var err error
	for client, err = jsonrpc.Dial("tcp", "localhost:1778"); err != nil; client, err = jsonrpc.Dial("tcp", "localhost:1778") {
		fmt.Println("connecting...")
	}
	var mock *int
	err = client.Call("Database.Clear", 0, &mock)
	if err != nil {
		t.Errorf("Error occured: %v", err)
	}

	event1 := database.Event{"2021-02-04", 27000, 9000, 3000}
	event2 := database.Event{"2021-02-04", 13000, 946, 5690}
	event3 := database.Event{"2021-02-03", 2, 0, 3000}
	event4 := database.Event{"2021-02-05", 645, 238, 346}
	events := []database.Event{event1, event2, event3, event4}

	for i := range events {
		err = client.Call("Database.AddEvent", events[i], mock)
		if err != nil {
			t.Errorf("Error occured: %v", err)
		}
	}

	date1 := "2021-02-03"
	date2 := "2021-02-04"
	date3 := "2021-02-05"
	date4 := "2021-02-06"

	results_full := []database.Stat{
		{"2021-02-03", 2, 0, 3000, -1, 1.5e+06},
		{"2021-02-04", 40000, 9946, 8690, 0.8737180776191433, 217.25},
		{"2021-02-05", 645, 238, 346, 1.453781512605042, 536.4341085271318},
	}
	results_part := []database.Stat{
		{"2021-02-03", 2, 0, 3000, -1, 1.5e+06},
		{"2021-02-04", 40000, 9946, 8690, 0.8737180776191433, 217.25},
	}
	results_by_cost := []database.Stat{
		{"2021-02-05", 645, 238, 346, 1.453781512605042, 536.4341085271318},
		{"2021-02-03", 2, 0, 3000, -1, 1.5e+06},
		{"2021-02-04", 40000, 9946, 8690, 0.8737180776191433, 217.25},
	}
	sortby := "date"
	var stats []database.Stat
	err = client.Call("Database.SelectStats", database.StatRequest{date1, date2, sortby}, &stats)
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

	err = client.Call("Database.SelectStats", database.StatRequest{date1, date3, sortby}, &stats)
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

	err = client.Call("Database.SelectStats", database.StatRequest{date1, date4, sortby}, &stats)
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

	err = client.Call("Database.SelectStats", database.StatRequest{date4, date1, sortby}, &stats)
	if err != nil {
		t.Errorf("Error occured: %v", err)
	}
	if len(stats) != 0 {
		t.Errorf("wrong stats")
	}

	client.Call("Database.Clear", 0, mock)

	err = client.Call("Database.SelectStats", database.StatRequest{date1, date4, sortby}, &stats)
	if err != nil {
		t.Errorf("Error occured: %v", err)
	}
	if len(stats) != 0 {
		t.Errorf("wrong stats")
	}
}

// func TestItems(t *testing.T) {
// 	arith := new(Arith)

// 	rpc.Register(arith)
// 	rpc.HandleHTTP()

// 	l, e := net.Listen("tcp", ":1234")
// 	if e != nil {
// 		log.Fatal("Dial error:", e)
// 	}
// 	conn, e := l.Accept()
// 	if e != nil {
// 		log.Fatal("Accept error:", e)
// 	}
// 	fmt.Println("accepted")
// 	//go http.Serve(l, nil)
// 	jsonrpc.ServeConn(conn)

// 	serverAddress := "127.0.0.1"

// 	client, err := jsonrpc.Dial("tcp", serverAddress+":1234")
// 	if err != nil {
// 		log.Fatal("dialing:", err)
// 	}
// 	// Synchronous call
// 	args := Args{7, 8}
// 	var reply int
// 	err = client.Call("Arith.Multiply", args, &reply)
// 	if err != nil {
// 		log.Fatal("arith error:", err)
// 	}
// 	fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)
// }
