package asyncChan

import (
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"golang.org/x/net/context"
)

type AsyncChan struct {
	ch     chan string
	closed bool
}

func NewChan() *AsyncChan {
	ch := make(chan string)
	return &AsyncChan{
		ch:     ch,
		closed: false,
	}
}
func (cc *AsyncChan) Close() {
	cc.closed = true
	close(cc.ch)
}

func (cc *AsyncChan) Read() <-chan string {
	if cc.closed {
		tools.Error("Read channel after closed!")
	}
	return cc.ch
}
func (cc *AsyncChan) Write(data string) {
	if cc.closed {
		tools.Error("Write channel after closed!")
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		select {
		case <-ctx.Done():
			tools.Error("Data not consumed within 1s! (data is %v)", data)
		case cc.ch <- data:
		}
	}()
}
