package network

import (
	"crypto/tls"
	"github.com/Synaxis/bfheroesFesl/config"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
	"time"
)

// Socket is a basic event-based TCP-Server
type Socket struct {
	Clients   *Clients
	name      string
	bind      string
	listen    net.Listener
	EventChan chan SocketEvent
	fesl      bool
}

func newSocket(name, bind string, fesl bool) *Socket {
	return &Socket{
		name:      name,
		bind:      bind,
		fesl:      fesl,
		Clients:   newClients(),
		EventChan: make(chan SocketEvent, 1000),
	}
}

// New starts to listen on a new Socket
func NewSocketTCP(name, bind string, fesl bool) (*Socket, error) {
	socket := newSocket(name, bind, fesl)
	listener, err := socket.listenTCP()
	if err != nil {
		return nil, err
	}
	socket.listen = listener
	go socket.run(socket.createClientTCP)

	return socket, nil
}

func NewSocketTLS(name, bind string) (*Socket, error) {
	socket := newSocket(name, bind, true)
	listener, err := socket.listenTLS()
	if err != nil {
		return nil, err
	}
	socket.listen = listener
	go socket.run(socket.createClientTLS)

	return socket, nil
}

func (socket *Socket) listenTCP() (net.Listener, error) {
	listener, err := net.Listen("tcp", socket.bind)
	if err != nil {
		logrus.WithError(err).Errorf("Listening on %s threw an error", socket.bind)
		return nil, err
	}
	return listener, nil
}


func (socket *Socket) listenTLS() (net.Listener, error) {
	cert, err := config.ParseCertificate()
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ClientAuth:         tls.NoClientCert,
		MinVersion:         tls.VersionSSL30,
		InsecureSkipVerify: true,
		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_RC4_128_SHA,
		},
	}

	listener, err := tls.Listen("tcp", socket.bind, config)
	if err != nil {
		logrus.WithError(err).Errorf("Listening on %s threw an error", socket.bind)
		return nil, err
	}

	return listener, nil
}

func (socket *Socket) handleClientEvents(client *Client) {
	defer socket.removeClient(client)

	for client.IsActive {
		select {
		case event := <-client.eventChan:
			switch {
			case event.Name == "close":
				socket.EventChan <- client.FireClientClose(event)
				socket.removeClient(client)
			case strings.Index(event.Name, "command") != -1:
				socket.EventChan <- client.FireClientCommand(event)
			case event.Name == "data":
				socket.EventChan <- client.FireClientData(event)
			default:
				socket.EventChan <- client.FireSomething(event)
			}
		}
	}
}

func (socket *Socket) removeClient(client *Client) {
	logrus.Debugf("Removing client %s", client.name)

	client.IsActive = false
	client.Close()
	socket.Clients.Remove(client)
}

type connAcceptFunc func(conn net.Conn)

func (socket *Socket) run(connect connAcceptFunc) {
	for {
		// Wait and listen for incomming connection
		conn, err := socket.listen.Accept()
		if err != nil {
			logrus.WithError(err).Errorf("A new client connecting to %s threw an error", socket.bind)
			continue
		}

		// Establish connection
		connect(conn)
	}
}

func (socket *Socket) createClientTCP(conn net.Conn) {
	tcpClient := newClientTCP(socket.name, conn, socket.fesl)
	socket.Clients.Add(tcpClient)
	go tcpClient.handleRequest()
	go socket.handleClientEvents(tcpClient)
	socket.EventChan <- socket.FireNewClient(tcpClient)
}

func (socket *Socket) createClientTLS(conn net.Conn) {
	tlscon, ok := conn.(*tls.Conn)
	if !ok {
		conn.Close()
		return
	}

	tlscon.SetDeadline(time.Now().Add(time.Second * 10))
	err := tlscon.Handshake()
	if err != nil {
		logrus.Errorf("%s: A new client connecting threw an error.\n%v\n%v", socket.name, err, tlscon.RemoteAddr())
		socket.EventChan <- socket.FireError(err)
		tlscon.Close()
	}

	state := tlscon.ConnectionState()
	logrus.Debugf("===HANDSHAKE DONE=== %t, %v", state.HandshakeComplete, state)

	// reset deadline after handshake
	tlscon.SetDeadline(time.Time{})

	tlsClient := newClientTLS(socket.name, tlscon)
	go tlsClient.handleRequestTLS()
	go socket.handleClientEvents(tlsClient)

	logrus.Println(socket.name + ":New Client connect")
	socket.Clients.Add(tlsClient)

	socket.EventChan <- socket.FireNewClient(tlsClient)
}

// Close the UDP connection to that client
func (socket *Socket) Close() {
	// Fire closing event
	socket.EventChan <- socket.FireClose()

	// Close socket
	socket.listen.Close()
}

// Socket is a basic event-based TCP-Server
type SocketUDP struct {
	Clients   []*Client
	name      string
	bind      string
	listen    *net.UDPConn
	EventChan chan SocketUDPEvent
	fesl      bool
}

// New starts to listen on a new Socket
func NewSocketUDP(name, bind string, fesl bool) (*SocketUDP, error) {
	socket := &SocketUDP{
		name:    name,
		bind:    bind,
		fesl:    fesl,
		Clients: []*Client{},
	}

	socket.EventChan = make(chan SocketUDPEvent, 1000)

	var err error
	serverAddr, err := net.ResolveUDPAddr("udp", socket.bind)
	if err != nil {
		logrus.Println("%s: Listening on %s threw an error.\n%v", socket.name, socket.bind, err)
		return nil, err
	}

	socket.listen, err = net.ListenUDP("udp", serverAddr)
	if err != nil {
		logrus.Println("%s: Listening on %s threw an error.\n%v", socket.name, socket.bind, err)
		return nil, err
	}

	go socket.run()

	return socket, nil
}

func (socket *SocketUDP) run() {
	buf := make([]byte, 8096)

	for socket.EventChan != nil {
		n, addr, err := socket.listen.ReadFromUDP(buf)
		if err != nil {
			logrus.WithError(err).Error("Error reading from UDP Ln 227 socket.go", err)
			socket.EventChan <- SocketUDPEvent{Name: "error", Addr: addr, Data: err}
			continue
		}

		socket.readFESL(buf[:n], addr)
	}
}


func (socket *SocketUDP) WriteEncode(Packet *codec.Packet, addr *net.UDPAddr) error {	

	// Encode packet
	buf, err := codec.
		NewEncoder().
		EncodePacket(Packet)
	if err != nil {
		logrus.
			WithError(err).
			WithField("type", Packet.Message).
			Error("Cannot encode packet")
		return err
	}

	// Send packet
	_, err = socket.listen.WriteTo(buf.Bytes(), addr)
	if err != nil {
		logrus.
			WithError(err).
			WithField("type", Packet.Message).
			Warn("Cannot send encoded packet")
		return err
	}

	return nil
}

// Close fires a close-event and closes the socket
func (socket *SocketUDP) Close() {
	// Fire closing event
	socket.EventChan <- SocketUDPEvent{Name: "close", Addr: nil, Data: nil}

	// Close socket
	socket.listen.Close()
}

