package goddis

import (
	"net"
	"time"
)

type Goddis struct {
	Port    int
	State   state
	Clients map[net.Addr]Client
}

type state struct {
	StringMap map[string]string
	HashMap   map[string]map[string]string
	Expiry    map[string]time.Time
}

func NewGoddis() *Goddis {
	goddis := new(Goddis)
	goddis.State.HashMap = make(map[string]map[string]string)
	goddis.State.StringMap = make(map[string]string)
	goddis.State.Expiry = make(map[string]time.Time)
	goddis.Clients = make(map[net.Addr]Client)
	return goddis
}
