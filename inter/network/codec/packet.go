package codec

type Packet struct {
	Message string
	Send    uint32
	Content interface{}
}
