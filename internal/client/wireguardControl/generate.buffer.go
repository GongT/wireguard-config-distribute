package wireguardControl

import (
	"bytes"
	"fmt"
	"os"
)

type Buffer struct {
	*bytes.Buffer
	extended bool
}

func newBuffer(extend bool) *Buffer {
	return &Buffer{
		Buffer:   bytes.NewBuffer(make([]byte, 0, 2048)),
		extended: extend,
	}
}

func (b *Buffer) appendLine(line string, args ...interface{}) {
	b.WriteString(fmt.Sprintf(line, args...))
	b.WriteByte('\n')
}

func (b *Buffer) appendLineExtened(line string, args ...interface{}) {
	if b.extended {
		b.appendLine(line, args...)
	}
}

func saveBuffersTo(filename string, datas ...[]byte) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0600))
	if err != nil {
		return err
	}
	for _, data := range datas {
		if _, err = f.Write(data); err != nil {
			return err
		}
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}
