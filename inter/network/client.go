package network

import (
	"bufio"
	"crypto/tls"
	"encoding/hex"

	"fmt"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/Synaxis/bfheroesFesl/storage/level"

	"github.com/sirupsen/logrus"
)

type Clients struct {
	mu        *sync.Mutex
	connected map[ClientKey]*Client
}

func newClients() *Clients {
	return &Clients{
		connected: make(map[ClientKey]*Client, 500),
		mu:        new(sync.Mutex),
	}
}

func (cls *Clients) Add(cl *Client) {
	cls.mu.Lock()
	cls.connected[cl.Key()] = cl
	cls.mu.Unlock()
}

func (cls *Clients) Remove(cl *Client) {
	cls.mu.Lock()
	delete(cls.connected, cl.Key())
	cls.mu.Unlock()
}

type ClientKey struct {
	name, addr string
}

func (ck *ClientKey) String() string {
	return fmt.Sprintf("%s_%s", ck.name, ck.addr)
}

type Client struct {
	name       string
	conn       net.Conn
	recvBuffer []byte
	eventChan  chan ClientEvent
	receiver   chan ClientEvent
	sender     chan codec.Packet
	IsActive   bool
	reader     *bufio.Reader
	HashState  *level.State
	IpAddr     net.Addr
	State      ClientState

	Options ClientOptions
}

type ClientOptions struct {
	FESL bool
}

func newClientTCP(name string, conn net.Conn, fesl bool) *Client {
	return &Client{
		name:      name,
		conn:      conn,
		IpAddr:    conn.RemoteAddr(),
		eventChan: make(chan ClientEvent, 500),
		reader:    bufio.NewReader(conn),
		IsActive:  true,
		Options: ClientOptions{
			FESL: fesl,
		},
	}
}

func newClientTLS(name string, conn *tls.Conn) *Client {
	return &Client{
		name:      name,
		conn:      conn,
		IpAddr:    conn.RemoteAddr(),
		IsActive:  true,
		eventChan: make(chan ClientEvent, 500),
		Options: ClientOptions{
			FESL: true, 
		},
	}
}

func (client *Client) handleRequestTLS() {
	client.IsActive = true
	buf := make([]byte, 8096) // buffer
	tempBuf := []byte{}

	for client.IsActive {
		n, err := client.readBuf(buf)
		if err != nil {
			return
		}

		if tempBuf != nil {
			tempBuf = append(tempBuf, buf[:n]...)
			tempBuf = client.readFESLTLS(buf[:n])
		} else {
			tempBuf = client.readFESLTLS(buf[:n])
		}
		buf = make([]byte, 8096) // new fresh buffer
	}
}

func (client *Client) handleRequest() {
	client.IsActive = true
	buf := make([]byte, 8096) // buffer
	tempBuf := []byte{}

	for client.IsActive {
		n, err := client.readBuf(buf)
		if err != nil {
			return
		}

		if client.Options.FESL {
			if tempBuf != nil {
				tempBuf = append(tempBuf, buf[:n]...)
				tempBuf = client.readFESL(buf[:n])
			} else {
				tempBuf = client.readFESL(buf[:n])
			}
			buf = make([]byte, 8096) // new fresh buffer
			continue
		}

		client.recvBuffer = append(client.recvBuffer, buf[:n]...)

		message := strings.TrimSpace(string(client.recvBuffer))

		logrus.Println("Got message:", hex.EncodeToString(client.recvBuffer))

		if strings.Index(message, `\final\`) == -1 {
			if len(client.recvBuffer) > 1024 {
				// Split message into 2
				client.recvBuffer = make([]byte, 0)
			}
			continue
		}

		client.eventChan <- ClientEvent{Name: "data", Data: message}

		commands := strings.Split(message, `\final\`)
		for _, command := range commands {
			if len(command) == 0 {
				continue
			}

			client.processCommand(command)
		}

		// Add unprocessed commands back into recvBuffer
		client.recvBuffer = []byte(commands[(len(commands) - 1)])
	}
}

func (client *Client) readBuf(buf []byte) (int, error) {
	n, err := client.conn.Read(buf)
	if err != nil {
		if err != io.EOF {
			client.eventChan <- client.FireError(err)
			client.eventChan <- client.FireClose()
			return 0, err
		}
		client.eventChan <- client.FireClose()
		return 0, err
	}
	return n, nil
}

func (c *Client) Key() ClientKey {
	return ClientKey{c.name, c.IpAddr.String()}
}

func (c *Client) Close() {
	logrus.Printf("%s:Client Closing.", c.name)
	c.eventChan <- ClientEvent{Name: "close", Data: c}
	c.conn.Close()
	c.IsActive = false
	c.FireClose()
}

type ClientState struct {
	ServerChallenge string
	ClientChallenge string
	ClientResponse  string
	IpAddress       net.Addr
	HasLogin        bool
	ProfileSent     bool
	LoggedOut       bool
	HeartTicker     *time.Ticker
}
