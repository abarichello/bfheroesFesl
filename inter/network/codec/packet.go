package codec

type Pkt struct {
	Message string
	Send    uint32
	Content interface{}
}
