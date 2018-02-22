package codec

type Pkt struct {
	Type    string
	Send    uint32
	Content interface{}
}
