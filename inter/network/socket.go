package network

import (
	"crypto/tls"
	"net"
	"time"
	"bytes"
	"encoding/binary"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"

	"github.com/Synaxis/bfheroesFesl/config"
)

// Socket is a basic event-based TCP-Server
// TODO: Rename it to broker
type Socket struct {
	bind      string
	listen    net.Listener
	EventChan chan SocketEvent
}

func newSocket(bind string) *Socket {
	return &Socket{
		bind:      bind,
		EventChan: make(chan SocketEvent),
	}
}

func NewSocketTCP(bind string) (*Socket, error) {
	socket := newSocket(bind)
	listener, err := socket.listenTCP()
	if err != nil {
		return nil, err
	}
	socket.listen = listener
	go socket.run(socket.createClientTCP)

	return socket, nil
}

func NewSocketTLS(bind string) (*Socket, error) {
	socket := newSocket(bind)
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
		//MaxVersion:   tls.VersionSSL30,
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
	tcpClient := NewClientTCP(conn)
	go tcpClient.handleRequestTCP()
	go tcpClient.handleClientEvents(socket)
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
		logrus.WithError(err).Errorf("A new client from %s connecting to %s threw an error", tlscon.RemoteAddr(), socket.bind)
		tlscon.Close()
	}
	tlscon.SetDeadline(time.Time{})

	state := tlscon.ConnectionState()
	logrus.Debugf("Connection handshake complete %t, %v", state.HandshakeComplete, state)

	tlsClient := NewClientTLS(tlscon)
	go tlsClient.handleRequestTLS()
	go tlsClient.handleClientEvents(socket)

	logrus.
		WithField("bind", socket.bind).
		WithField("protocol", "tcp").
		Print("A new client connected")

	socket.EventChan <- socket.FireNewClient(tlsClient)
}


// Close fires a close-event and closes the socket
func (socket *Socket) Close() {
	// Fire closing event
	close(socket.EventChan)

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
			logrus.WithError(err).Error("Error reading from UDP", err)
			socket.EventChan <- SocketUDPEvent{Name: "error", Addr: addr, Data: err}
			continue
		}

		socket.readFESL(buf[:n], addr)
	}
}

func (socket *SocketUDP) readFESL(data []byte, addr *net.UDPAddr) {
	p := bytes.NewBuffer(data)
	var payloadID uint32
	var payloadLen uint32

	payloadType := string(data[:4])
	p.Next(4)

	binary.Read(p, binary.BigEndian, &payloadID)
	binary.Read(p, binary.BigEndian, &payloadLen)

	payload := codec.DecodeFESL(data[12:])

	socket.EventChan <- SocketUDPEvent{
		Name: payloadType,
		Addr: addr,
		Data: &codec.Command{
			Query:     payloadType,
			PayloadID: payloadID,
			Message:   payload,
		},
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