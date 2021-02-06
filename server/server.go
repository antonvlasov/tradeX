package server

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/antonvlasov/tradeX/database"
)

type Listener int
type Reply struct {
	Data string
}

//Start the RPC Server. Function blocks.
func Serve(dbpath string) {
	database.Connect(dbpath)
	defer database.Close()
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:1778")
	if err != nil {
		log.Fatal(err)
	}
	inbound, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	handler := new(database.Database)
	rpc.Register(handler)
	for {
		conn, err := inbound.Accept()
		if err != nil {
			continue
		}
		jsonrpc.ServeConn(conn)
	}
}
