package codec

type Packet struct {
	Type    string
	Step    uint32
	Payload interface{}
}
