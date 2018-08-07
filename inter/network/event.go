package network

import (
 "net"
 "github.com/Synaxis/bfheroesFesl/inter/network/codec"

)

type SocketEvent struct {
	Name string
	Data interface{}
}

type SocketUDPEvent struct {
	Name string
	Addr *net.UDPAddr
	Data interface{}
}

// ClientEvent is the generic struct for events
type ClientEvent struct {
	Name string
	Data interface{}
}

type EventClientCommand struct {
	Client *Client
	// If TLS (theater then we ignore payloadID - it is always 0x0)
	Command *codec.Packet
}


type EventError struct {
	Error error
}
type EventNewClient struct {
	Client *Client
}

type EventClientClose struct {
	Client *Client
}
type EventClientError struct {
	Client *Client
	Error  error
}
type EvProcess struct {
	Client  *Client
	Process *ProcessFESL // If TLS (theater then we ignore HEX - it is always 0x0)
}

type EventClientData struct {
	Client *Client
	Data   string
}

func (c *Client) FireClientClose(event ClientEvent) SocketEvent {
	return SocketEvent{
		Name: "client.close",
		Data: EventClientClose{Client: c},
	}
}

func (c *Client) FireClose() ClientEvent {
	return ClientEvent{
		Name: "close",
		Data: c,
	}
}

func (c *Client) FireError(err error) ClientEvent {
	return ClientEvent{
		Name: "error",
		Data: err,
	}
}

func (c *Client) FireClientData(event ClientEvent) SocketEvent {
	return SocketEvent{
		Name: "client.data",
		Data: EventClientData{
			Client: c,
			Data:   event.Data.(string),
		},
	}
}

func (c *Client) FireClientCommand(event ClientEvent) SocketEvent {
	return SocketEvent{
		Name: "client." + event.Name,
		Data: EvProcess{
			Client:  c,
			Process: event.Data.(*ProcessFESL),
		},
	}
}


func (c *Client) FireSomething(event ClientEvent) SocketEvent {
	var interfaceSlice = make([]interface{}, 2)
	interfaceSlice[0] = c
	interfaceSlice[1] = event.Data
	return SocketEvent{
		Name: "client." + event.Name,
		Data: interfaceSlice,
	}
}

func (s *Socket) FireError(err error) SocketEvent {
	return SocketEvent{
		Name: "error",
		Data: EventError{Error: err},
	}
}

func (s *Socket) FireClose() SocketEvent {
	return SocketEvent{
		Name: "close",
		Data: nil,
	}
}

func (s *Socket) FireNewClient(client *Client) SocketEvent {
	return SocketEvent{
		Name: "newClient",
		Data: EventNewClient{Client: client},
	}
}
