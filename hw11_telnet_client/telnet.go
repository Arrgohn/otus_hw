package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Telnet struct {
	adr     string
	timeout time.Duration
	ctx     context.Context
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
	cl      io.Closer
}

func (t *Telnet) Connect() error {
	conn, err := net.DialTimeout("tcp", t.adr, t.timeout)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	t.conn = conn

	return nil
}

func (t *Telnet) Close() error {
	err := t.conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func (t *Telnet) Send() error {
	_, err := io.Copy(t.conn, t.in)
	if err != nil {
		return fmt.Errorf("error occurred while sending: %w", err)
	}

	fmt.Fprintln(os.Stderr, "...EOF")

	return nil
}

func (t *Telnet) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	if err != nil {
		return fmt.Errorf("error occurred while receiving: %w", err)
	}
	fmt.Fprintln(os.Stderr, "...Connection was closed by peer")

	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Telnet{adr: address, timeout: timeout, in: in, out: out}
}
