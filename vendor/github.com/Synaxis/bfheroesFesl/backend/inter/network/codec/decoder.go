package codec

func DecodeFESL(data []byte) map[string]string {
	out := map[string]string{}
	key := []byte{}
	pos := 0
	for i, c := range data {
		switch c {
		case charEqual:
			key = data[pos:i]
			pos = i + 1
		case charNewLine:
			out[string(key)] = string(data[pos:i])
			pos = i + 1
		}
	}
	return out
}
