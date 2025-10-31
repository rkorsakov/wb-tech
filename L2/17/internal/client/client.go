package client

import (
	ui2 "17/internal/ui"
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

type Client struct {
	timeout time.Duration
	conn    *net.TCPConn
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewClient(timeout time.Duration) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{timeout: timeout, ctx: ctx, cancel: cancel}
}

func (c *Client) Connect(host, port string) error {
	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return fmt.Errorf("error resolving host addr: %w", err)
	}
	conn, err := net.DialTimeout("tcp", addr.String(), c.timeout)
	if err != nil {
		return fmt.Errorf("error connecting to host: %w", err)
	}
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		conn.Close()
		return fmt.Errorf("expected TCP connection, got %T", conn)
	}
	c.conn = tcpConn
	return nil
}

func (c *Client) send(data []byte) error {
	_, err := c.conn.Write(data)
	if err != nil {
		return fmt.Errorf("error sending data: %w", err)
	}
	return nil
}

func (c *Client) StartInputHandler() {
	ui := ui2.NewUI()

	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				data, err := ui.Read()
				if err == io.EOF {
					c.Close()
					return
				}
				if err != nil {
					fmt.Printf("Error reading input: %v\n", err)
					continue
				}

				err = c.send(data)
				if err != nil {
					fmt.Printf("Error sending data: %v\n", err)
					return
				}
			}
		}
	}()
}

func (c *Client) StartConnectionHandler() {
	buffer := make([]byte, 1024)
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			n, err := c.conn.Read(buffer)
			if err != nil {
				fmt.Printf("Error receiving data: %v\n", err)
				return
			}
			fmt.Printf("Received: %s", string(buffer[:n]))
		}
	}
}

func (c *Client) Close() {
	c.cancel()
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) Done() <-chan struct{} {
	return c.ctx.Done()
}
