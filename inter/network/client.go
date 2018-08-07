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

const (
	FragmentSize = 8096
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
	Conn       net.Conn
	recvBuffer []byte
	eventChan  chan ClientEvent
	receiver   chan ClientEvent
	sender     chan codec.Packet
	IsActive   bool
	reader     *bufio.Reader
	HashState  *level.State
	IpAddr     net.Addr
	State      ClientState
	HeartTicker *time.Ticker
	Type string

	Options ClientOptions
}

type ClientOptions struct {
	FESL bool
}

func newClient(conn net.Conn) *Client {
	return &Client{
		Conn:       conn,
		IpAddr:     conn.RemoteAddr(),
		receiver:   make(chan ClientEvent, 5),
		sender:     make(chan codec.Packet, 5),
		IsActive:   true,
	}
}

func NewClientTCP(conn net.Conn) *Client {
	return newClient(conn)
}

func NewClientTLS(conn *tls.Conn) *Client {
	return newClient(conn)
}


func (client *Client) handleRequestTLS() {
	client.IsActive = true
	buf := make([]byte, FragmentSize)

	for client.IsActive {
		n, err := client.readBuf(buf)
		if err != nil {
			return
		}

		client.readTLSPacket(buf[:n])
		buf = make([]byte, FragmentSize)
	}
}

func (client *Client) handleRequestTCP() {
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
				tempBuf = client.readFESL(buf[:n])
			} else {
				tempBuf = client.readFESL(buf[:n])
			}
			buf = make([]byte, 8096) // new fresh buffer
			continue
		

		client.recvBuffer = append(client.recvBuffer, buf[:n]...)

		message := strings.TrimSpace(string(client.recvBuffer))

		logrus.Debugln("Got message:", hex.EncodeToString(client.recvBuffer))

		if strings.Index(message, `\final\`) == -1 {
			if len(client.recvBuffer) > 1024 {
				// We don't support more than 2048 long messages
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
	n, err := client.Conn.Read(buf)
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

func (client *Client) handleClientEvents(socket *Socket) {
	defer client.Close()

	for client.IsActive {
		select {
		case event := <-client.receiver:
			switch {
			case event.Name == "close":
				return
			case strings.Index(event.Name, "command") != -1:
				socket.EventChan <- client.FireClientCommand(event)
			case event.Name == "data":
				logrus.Warnf("Not implemented: Client send client.data: %s", event.Data)
			default:
				logrus.Warn("Not implemented client.%s for %s", event.Name, event.Data)
			}
		}
	}
}


func (c *Client) Key() ClientKey {
	return ClientKey{c.name, c.IpAddr.String()}
}

func (c *Client) Close() {
	logrus.Printf("%s:Client Closing.", c.name)
	c.eventChan <- ClientEvent{Name: "close", Data: c}
	c.Conn.Close()
	c.IsActive = false
	//c.FireClose()
}

type ClientState struct {
	ServerChallenge string
	ClientChallenge string
	ClientResponse  string
	Username        string
	PlyName         string
	PlyEmail        string
	PlyCountry      string
	PlyPid          int
	Sessionkey      int
	Confirmed       bool
	IpAddress       net.Addr
	HasLogin        bool
	ProfileSent     bool
	LoggedOut       bool
	HeartTicker     *time.Ticker
}
