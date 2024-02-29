package handler

import (
	"redis-server-go/resp"
	"sync"
)

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

// hashmap inside hashmap
var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
	"HSET": hset,
	"HGET": hget,
}

func ping(args []resp.Value) resp.Value {
	// only one ping string from cli
	if len(args) == 0 {
		return resp.Value{Typ: "string", Str: "PONG"}
	}

	// more than one ping string, return back the org bulk after ping word
	return resp.Value{Typ: "string", Str: args[0].Bulk}
}

func set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "wrong number of arguments for SET"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	// to prevent multi thread go from modifying the SET
	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return resp.Value{Typ: "string", Str: "OK"}
}

func get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "wrong number of arguments for GET"}
	}

	key := args[0].Bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	if !ok {
		return resp.Value{Typ: "string", Str: "key not found"}
	}
	SETsMu.RUnlock()

	return resp.Value{Typ: "string", Str: value}
}

func hset(args []resp.Value) resp.Value {
	if len(args) != 3 {
		return resp.Value{Typ: "error", Str: "wrong number of arguments for HSET"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return resp.Value{Typ: "string", Str: "OK"}
}

func hget(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "wrong number of arguments for HGET"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	HSETsMu.Lock()
	value, ok := HSETs[hash][key]
	HSETsMu.Unlock()

	if !ok {
		return resp.Value{Typ: "error", Str: "cannot find the key for HGET"}
	}

	return resp.Value{Typ: "string", Str: value}
}
