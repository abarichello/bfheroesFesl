package network

import (
	"bufio"
	"crypto/tls"

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
		receiver:   make(chan ClientEvent, 5),
		sender:     make(chan codec.Packet, 5),
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

func (client *Client) handleRequest() {
	client.IsActive = true
	buf := make([]byte, FragmentSize)

	for client.IsActive {
		n, err := client.readBuf(buf)
		if err != nil {
			return
		}
		client.readFESL(buf[:n])
		buf = make([]byte, FragmentSize)
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
	c.conn.Close()
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
