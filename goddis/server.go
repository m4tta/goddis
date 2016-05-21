package goddis

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
)

type Command struct {
	command string
	args    []string
	client  Client
}

// All I do is read data off connections
func (g *Goddis) handleConnection(c net.Conn, cmdChan chan<- Command, addChan chan<- Client, rmChan chan<- Client) {
	client := Client{c}
	addChan <- client
	defer func() {
		rmChan <- client
	}()

	r := bufio.NewReader(c)
	for {
		buf := make([]byte, 1000000)
		n, err := r.Read(buf)
		if err != nil {
			break
		}
		data := buf[0:n]
		cmd, _ := g.parseCommand(data)
		cmd.client = client
		cmdChan <- *cmd
	}
}

// parse received data into a command
func (g *Goddis) parseCommand(data []byte) (*Command, error) {
	command := new(Command)
	if !strings.HasPrefix(string(data), "*") {
		return nil, errors.New("Protocol error")
	}
	args := strings.Split(string(data), "\r\n")
	var cmds []string
	for _, s := range args {
		if !strings.HasPrefix(s, "$") && !strings.HasPrefix(s, "*") && len(s) > 0 {
			cmds = append(cmds, s)
		}
	}
	if len(cmds) != 0 {
	command.command = cmds[0]
	command.args = cmds[1:]
	}
	return command, nil
}

func (g *Goddis) validCommand(cmd Command) bool {
	t := reflect.TypeOf(g)
	m, exists := t.MethodByName(cmd.command)
	if exists {
		return true
	}
	fmt.Println(exists)
	fmt.Println(m)
	// might allow me to do generic command processor
	// atleast get it to validate arg count maybe
	// think of other things
	// im sure there is a better way to do command processing
	return false
}

// TODO: Better command validation
func (g *Goddis) processCommand(cmd Command) {
	switch strings.ToUpper(cmd.command) {
	case "PING":
		cmd.client.Pong()
	case "ECHO":
		if len(cmd.args) > 0 {
			cmd.client.BulkString(cmd.args[0])
		} else {
			cmd.client.Error(IncorrectArgs(cmd))
		}
	case "EXISTS":
		if len(cmd.args) > 1 {
			cmd.client.SendBool(g.Exists(cmd.args[0]))
		} else {
			cmd.client.Error(IncorrectArgs(cmd));
		}
	case "DEL":
		cmd.client.SendInt(g.Del(cmd.args[0:]...))
	case "EXPIRE":
		if len(cmd.args) > 1 {
			cmd.client.SendBool(g.Expire(cmd.args[0], cmd.args[1]))
		} else {
			cmd.client.Error(IncorrectArgs(cmd));
		}
	case "SET":
		// TODO: Validate syntax for a SET command with optional params
		if len(cmd.args) > 1 {
			if ok := g.Set(cmd.args[0], cmd.args[1], cmd.args[2:]...); ok {
				cmd.client.Ok()
			} else {
				cmd.client.SyntaxError()
			}
		} else {
			cmd.client.Error(IncorrectArgs(cmd))
		}
	case "INCR":
		if len(cmd.args) == 1 {
			i, err := g.Incr(cmd.args[0])
			if err != nil {
				cmd.client.Error(err.Error())
			} else {
				cmd.client.SendInt(i)
			}
		} else {
			cmd.client.Error(IncorrectArgs(cmd));
		}
	case "INCRBY":
		if len(cmd.args) == 2 {
			i, err := g.IncrBy(cmd.args[0], cmd.args[1])
			if err != nil {
				cmd.client.Error(err.Error())
			} else {
				cmd.client.SendInt(i)
			}
		} else {
			cmd.client.Error(IncorrectArgs(cmd));
		}
	case "GET":
		if len(cmd.args) == 1 {
			if value, ok := g.Get(cmd.args[0]); ok {
				cmd.client.BulkString(value)
			} else {
				cmd.client.SendNull()
			}
		} else {
			cmd.client.Error(IncorrectArgs(cmd));
		}
	case "MGET":
		if len(cmd.args) >= 1 {
			vals := g.MGet(cmd.args[0:]...)
			cmd.client.SendArray(vals...)
		} else {
			cmd.client.Error(IncorrectArgs(cmd))
		}
	case "HSET":
		if len(cmd.args) == 3 {
			if ok := g.HSet(cmd.args[0], cmd.args[1], cmd.args[2]); ok {
				cmd.client.SendBool(true)
			} else {
				cmd.client.SendBool(false)
			}
		} else {
			cmd.client.Error(IncorrectArgs(cmd));
		}
	case "HGET":
		if len(cmd.args) == 2 {
			val, ok := g.HGet(cmd.args[0], cmd.args[1])
			if ok {
				cmd.client.BulkString(val)
			} else {
				cmd.client.SendNull()
			}
		} else {
			cmd.client.Error(IncorrectArgs(cmd))
		}
	case "HMGET":
		if len(cmd.args) >= 2 {
			vals := g.HMGet(cmd.args[0], cmd.args[1:]...)
			cmd.client.SendArray(vals...)
		} else {
			cmd.client.Error(IncorrectArgs(cmd))
		}
	case "HINCRBY":
		if len(cmd.args) == 3 {
			i, err := g.HIncrBy(cmd.args[0], cmd.args[1], cmd.args[2])
			if err != nil {
				cmd.client.Error(err.Error())
			} else {
				cmd.client.SendInt(i)
			}
		} else {
			cmd.client.Error(IncorrectArgs(cmd))
		}
	default:
		cmd.client.Error(UnknownCmd(cmd))
	}
}

func IncorrectArgs(cmd Command) string {
	return "wrong number of arguments for '" + cmd.command + "' command"
}

func UnknownCmd(cmd Command) string {
	if cmd.command != "" {
		return "unknown command '" + cmd.command + "'"
	} else {
		return "No command"
	}
}

func (g *Goddis) marshalChannels(cmdChan <-chan Command, addChan <-chan Client, rmChan <-chan Client) {
	for {
		select {
		case cmd, ok := <-cmdChan:
			if ok {
				g.processCommand(cmd)
			}
		case client, ok := <-addChan:
			if ok {
				g.Clients[client.conn.RemoteAddr()] = client
			}
		case client, ok := <-rmChan:
			if ok {
				client.conn.Close()
				delete(g.Clients, client.conn.RemoteAddr())
			}
		}
	}
}

func (g *Goddis) Listen(port int) {
	cmdChan := make(chan Command, 50)
	addChan := make(chan Client, 50)
	rmChan := make(chan Client, 50)

	// Start listening server
	serv, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Listening on port", port)
	defer serv.Close()

	// Channel processer
	go g.marshalChannels(cmdChan, addChan, rmChan)

	// listener loop
	for {
		conn, err := serv.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		// Handle new connections in own goroutine
		go g.handleConnection(conn, cmdChan, addChan, rmChan)
	}
}
