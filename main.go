package main

import (
	"fmt"
	"net"
	"redis-server-go/aof"
	"redis-server-go/handler"
	"redis-server-go/resp"

	"strings"
)

func main() {
	fmt.Println("Listening on port :6379")

	// Create a new server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	Aof, err := aof.NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer Aof.Close()

	// read from file when restarting
	Aof.Read(func(value resp.Value) {
		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		handler, ok := handler.Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})

	// Listen for connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		Resp := resp.NewResp(conn)
		value, err := Resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.Typ != "array" || len(value.Array) == 0 {
			fmt.Println("Invalid request")
			continue
		}

		// take the first argu which should be either set, get, hset, hget
		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		writer := resp.NewWriter(conn)

		// check the command
		handler, ok := handler.Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write((resp.Value{Typ: "string", Str: ""}))
			continue
		}

		// write to file
		if command == "SET" || command == "HSET" {
			Aof.Write(value)
		}

		// check the rest of the args
		result := handler(args)
		fmt.Println("result: \n", result)
		writer.Write(result)

	}
}
