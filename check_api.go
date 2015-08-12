package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

// Status of health check
type Status int

const (
	// UNCHECKED -  if unchecked
	UNCHECKED Status = iota
	// DOWN - is down
	DOWN
	// UP - is up
	UP
)

// Site url
type Site struct {
	url string
}

// Status - Check status of end point
func (s Site) Status() (Status, error) {

	conn, err := net.DialTimeout("tcp", s.url, time.Duration(2*time.Second))

	apiStatus := UP

	if err, ok := err.(net.Error); ok && err.Timeout() {

		apiStatus = DOWN
	} else if err != nil {
		apiStatus = DOWN
	} else {

		fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
		_, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			apiStatus = DOWN
		}
	}
	return apiStatus, err
}
