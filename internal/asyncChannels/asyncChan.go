//go:generate go run ../../tools/create.go SidType "github.com/gongt/wireguard-config-distribute/internal/types"

package asyncChannels

import (
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"golang.org/x/net/context"
	t "github.com/gongt/wireguard-config-distribute/internal/types"
)

type AsyncChanSidType struct {
	ch     chan t.SidType
	closed bool
}

func NewChanSidType() *AsyncChanSidType {
	ch := make(chan t.SidType)
	return &AsyncChanSidType{
		ch:     ch,
		closed: false,
	}
}
func (cc *AsyncChanSidType) Close() {
	cc.closed = true
	close(cc.ch)
}

func (cc *AsyncChanSidType) Read() <-chan t.SidType {
	if cc.closed {
		tools.Error("Read channel after closed!")
	}
	return cc.ch
}
func (cc *AsyncChanSidType) Write(data t.SidType) {
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
