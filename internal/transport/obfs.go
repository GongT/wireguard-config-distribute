package transport

type obfsBuffer [2048]byte

func newBuffer() obfsBuffer {
	return obfsBuffer{}
}

func (buff *obfsBuffer) encode(len int) []byte {
	for i := 0; i < len; i++ {
		buff[i]++
	}
	return buff[:len]
}

func (buff *obfsBuffer) decode(len int) []byte {
	for i := 0; i < len; i++ {
		buff[i]--
	}
	return buff[:len]
}
