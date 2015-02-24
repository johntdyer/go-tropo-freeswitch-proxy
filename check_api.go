package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Status int

const (
	UNCHECKED Status = iota
	DOWN
	UP
)

type Site struct {
	url string
}

func (s Site) Status() (Status, error) {

	conn, err := net.DialTimeout("tcp", s.url, time.Duration(2*time.Second))

	api_status := UP

	if err, ok := err.(net.Error); ok && err.Timeout() {

		api_status = DOWN
	} else if err != nil {
		api_status = DOWN
	} else {

		fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
		_, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			api_status = DOWN
		}
	}
	return api_status, err
}
