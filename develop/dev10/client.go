package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"
)

type Client struct {
	address    string
	timeout    time.Duration
	reader     io.ReadCloser
	writer     io.Writer
	connection net.Conn
	readerScan *bufio.Scanner
	connScan   *bufio.Scanner
}

func NewTCPClient(address string, timeout time.Duration, reader io.ReadCloser, writer io.Writer) *Client {
	return &Client{
		address: address,
		timeout: timeout,
		reader:  reader,
		writer:  writer,
	}
}
func (t *Client) Connect() (err error) {

	t.connection, err = net.DialTimeout("tcp", t.address, t.timeout)
	t.connScan = bufio.NewScanner(t.connection)
	t.readerScan = bufio.NewScanner(t.reader)

	return
}

func (t *Client) Close() (err error) {
	if t.connection != nil {
		err = t.connection.Close()
	}
	return
}

func (t *Client) Send() (err error) {
	if t.connection == nil {
		return
	}
	if !t.readerScan.Scan() {
		return errors.New("no data")
	}
	_, err = t.connection.Write(append(t.readerScan.Bytes(), '\n'))
	return
}

func (t *Client) Receive() (err error) {
	if t.connection == nil {
		return
	}
	if !t.connScan.Scan() {
		return errors.New("connection closed")
	}
	_, err = t.writer.Write(append(t.connScan.Bytes(), '\n'))
	return
}
