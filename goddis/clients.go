package goddis

import (
	"net"
	"strconv"
)

type Client struct {
	conn net.Conn
}

func (c *Client) Write(s string) {
	c.conn.Write([]byte(s))
}

func (c *Client) Ok() {
	c.SString("OK")
}

func (c *Client) Pong() {
	c.SString("PONG")
}

func (c *Client) Error(s string) {
	c.Write("-ERR " + s + "\r\n")
}

func (c *Client) ErrorType(s string) {
	c.Write("-WRONGTYPE " + s + "\r\n")
}

func (c *Client) BulkString(s string) {
	if len(s) > 0 {
		c.Write("$" + strconv.Itoa(len(s)) + "\r\n" + s + "\r\n")
	} else {
		c.SendNull()
	}
}

func (c *Client) SendArray(values ...string) {
	var message string
	length := len(values)
	message = "*" + strconv.Itoa(length) + "\r\n"
	for _, v := range values {
		if len(v) == 0 {
			message += "$-1\r\n"
		} else {
			message += "$" + strconv.Itoa(len(v)) + "\r\n" + v + "\r\n"
		}
	}
	c.Write(message)
}

func (c *Client) SString(s string) {
	c.Write("+" + s + "\r\n")
}

func (c *Client) SendInt(i int) {
	c.Write(":" + strconv.Itoa(i) + "\r\n")
}

func (c *Client) SendBool(b bool) {
	c.Write(":" + strconv.Itoa(Btoi(b)) + "\r\n")
}

func (c *Client) SendNull() {
	c.Write("$-1\r\n")
}

func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
